package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"pick_v2/dao"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/request"
	"pick_v2/utils/xsq_net"
)

func Order(ctx context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	var (
		form req.PurchaseOrderForm
		Id   int
	)

	for i := range messages {
		err := json.Unmarshal(messages[i].Body, &Id)
		if err != nil {
			global.Logger["err"].Infof("解析json失败:%s", err.Error())
			return consumer.ConsumeRetryLater, nil
		}
		form.OrderId = append(form.OrderId, Id)
	}

	orderRsp, err := GetOrderInfo(form)

	if err != nil {
		global.Logger["err"].Infof("获取订单信息失败,错误信息:%s", err.Error())
		return consumer.ConsumeRetryLater, err
	}

	db := global.DB

	for _, info := range orderRsp.Data {

		if info.DistributionType == 6 { //无需出库
			return dao.NoShipping(db, form, info)
		} else {
			return dao.Shipping(db, form, info)
		}

	}

	return consumer.ConsumeRetryLater, errors.New("异常")
}

func PullOrder(c *gin.Context) {
	var form req.PurchaseOrderForm

	orderRsp, err := GetOrderInfo(form)

	if err != nil {
		global.Logger["err"].Infof("获取订单信息失败,错误信息:%s", err.Error())
		xsq_net.ErrorJSON(c, err)
	}

	db := global.DB

	for _, info := range orderRsp.Data {

		if info.DistributionType == 6 { //无需出库
			_, err = dao.NoShipping(db, form, info)
			if err != nil {
				xsq_net.ErrorJSON(c, err)
				return
			}
		} else {
			_, err = dao.Shipping(db, form, info)
			if err != nil {
				xsq_net.ErrorJSON(c, err)
				return
			}
		}

	}

	xsq_net.Success(c)
}

func GetOrderInfo(responseData interface{}) (result rsp.OrderRsp, err error) {

	err = request.Call("api/v1/remote/get/goods/by/id", responseData, &result)

	return
}

func NewCloseOrder(ctx context.Context, messages ...*primitive.MessageExt) (consumeRes consumer.ConsumeResult, err error) {

	var (
		form req.CloseOrderInfo
	)

	for i := range messages {
		err = json.Unmarshal(messages[i].Body, &form.OccId)
		if err != nil {
			global.Logger["err"].Infof("解析json失败:%s", err.Error())
			return consumer.ConsumeRetryLater, nil
		}
	}

	logic, err := NewCloseOrderLogic(form)
	if err != nil {
		return logic, err
	}

	return
}

func CloseOrder(c *gin.Context) {
	var form req.CloseOrderInfo

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	_, err := NewCloseOrderLogic(form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

// 关单逻辑
func NewCloseOrderLogic(form req.CloseOrderInfo) (consumeRes consumer.ConsumeResult, err error) {
	var (
		result     rsp.CloseOrderRsp
		isCommit   bool
		closeOrder model.CloseOrder
	)

	db := global.DB

	err, closeOrder = model.GetCloseOrderByPk(db, form.OccId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.Logger["err"].Infof(err.Error())
		return consumer.ConsumeRetryLater, err
	}

	if closeOrder.Id > 0 {
		global.Logger["err"].Infof("对应的关单任务已存在")
		return consumer.ConsumeSuccess, errors.New("对应的关单任务已存在")
	}

	//调用订货系统接口
	err = request.Call("api/v1/remote/close/info/v2", form, &result)

	if err != nil {
		global.Logger["err"].Infof("接口调用失败:%s", err.Error())
		return consumer.ConsumeRetryLater, nil
	}

	closeOrderRes := result.Data

	goodsInfo := closeOrderRes.GoodsInfo

	if closeOrderRes.Number == "" {
		global.Logger["err"].Infof("订货系统查询订单数据不存在")
		return consumer.ConsumeSuccess, errors.New("订货系统查询订单数据不存在")
	}

	if len(goodsInfo) == 0 {
		global.Logger["err"].Infof("订货系统查询订单商品数据不存在")
		return consumer.ConsumeSuccess, errors.New("订货系统查询订单商品数据不存在")
	}

	var (
		order model.Order
		tips  rsp.CloseTips
	)

	tx := db.Begin()

	//订单状态
	err, order = model.GetOrderByNumber(tx, closeOrderRes.Number)

	if err != nil {
		global.Logger["err"].Infof(err.Error())
		return
	}

	var status = 1

	//如果是新订单
	if order.OrderType == model.NewOrderType {
		//完成状态
		status = 2

		var (
			closeGoodsMp = make(map[int]int, 0)
		)

		for _, info := range goodsInfo {
			closeGoodsMp[info.ID] = info.NeedCloseCount
		}

		//关闭订单逻辑处理
		err, isCommit, _, tips = dao.OrderDataHandle(tx, closeOrderRes.OccId, closeOrderRes.Number, closeOrderRes.Typ, closeGoodsMp)
		if err != nil {
			tx.Rollback()
			global.Logger["err"].Infof(err.Error())
			return consumer.ConsumeRetryLater, err
		}

	} else {
		//不是新订单,查询是否在任务中是新订单
		var outboundOrder model.OutboundOrder

		err, outboundOrder = model.GetOutboundOrderByNumberFirstSortByTaskId(tx, closeOrderRes.Number)

		if err != nil {
			global.Logger["err"].Infof(err.Error())
			return
		}

		//在任务中是新订单
		if outboundOrder.OrderType == model.OutboundOrderTypeNew {
			//完成状态
			status = 2

			var (
				closeGoodsMp = make(map[int]int, 0)
			)

			for _, info := range goodsInfo {
				closeGoodsMp[info.ID] = info.NeedCloseCount
			}

			//关闭订单逻辑处理
			err, isCommit, _, tips = dao.OrderDataHandle(tx, closeOrderRes.OccId, closeOrderRes.Number, closeOrderRes.Typ, closeGoodsMp)
			if err != nil {
				tx.Rollback()
				global.Logger["err"].Infof(err.Error())
				return consumer.ConsumeRetryLater, err
			}

		}
	}

	closeTotal := GetCloseTotal(goodsInfo)

	modelCloseOrder := GetModelCloseOrderData(closeOrderRes, closeTotal, status, tips)

	//需关闭总数
	err = model.CloseOrderReplaceSave(tx, modelCloseOrder, []string{"shop_name", "shop_type"})

	if err != nil {
		tx.Rollback()
		global.Logger["err"].Infof(err.Error())
		return consumer.ConsumeRetryLater, err
	}

	closeGoodsInfo := GetModelCloseOrderGoodsData(goodsInfo, modelCloseOrder.Id, modelCloseOrder.Number)

	err = model.BatchSaveCloseGoods(tx, &closeGoodsInfo)

	if err != nil {
		tx.Rollback()
		global.Logger["err"].Infof(err.Error())
		return consumer.ConsumeRetryLater, err
	}

	if isCommit {
		err = dao.CloseOrderResult(modelCloseOrder.Id, 1)

		if err != nil {
			tx.Rollback()
			global.Logger["err"].Infof(err.Error())
			return consumer.ConsumeRetryLater, err
		}
	}

	tx.Commit()

	return
}

func GetModelCloseOrderData(closeOrder rsp.CloseOrder, closeTotal, status int, tips rsp.CloseTips) (modelCloseOrder *model.CloseOrder) {
	modelCloseOrder = &model.CloseOrder{
		Id:               closeOrder.OccId,
		Number:           closeOrder.Number,
		ShopName:         closeOrder.ShopName,
		ShopType:         closeOrder.ShopType,
		DistributionType: closeOrder.DistributionType,
		OrderRemark:      closeOrder.OrderRemark,
		PayAt:            closeOrder.PayAt,
		PayTotal:         closeOrder.PayTotal,
		NeedCloseTotal:   closeTotal,
		Province:         closeOrder.Province,
		City:             closeOrder.City,
		District:         closeOrder.District,
		Status:           status,
		Applicant:        closeOrder.Applicant,
		ApplyTime:        closeOrder.ApplyTime,
		Colour:           tips.Colour,
		Tips:             tips.Tips,
	}

	return
}

func GetCloseTotal(closeGoodsInfo []rsp.CloseGoodsInfo) (closeTotal int) {
	for _, info := range closeGoodsInfo {
		//需关闭总数
		closeTotal += info.NeedCloseCount
	}
	return
}

func GetModelCloseOrderGoodsData(closeGoodsInfo []rsp.CloseGoodsInfo, closeOrderId int, number string) (modelCloseOrderGoods []model.CloseGoods) {
	for _, info := range closeGoodsInfo {
		modelCloseOrderGoods = append(modelCloseOrderGoods, model.CloseGoods{
			CloseOrderId:   closeOrderId,
			OrderGoodsId:   info.ID,
			Number:         number,
			GoodsName:      info.Name,
			Sku:            info.Sku,
			GoodsSpe:       info.GoodsSpe,
			PayCount:       info.PayCount,
			CloseCount:     info.CloseCount,
			NeedCloseCount: info.NeedCloseCount,
			GoodsRemark:    info.GoodsRemark,
		})
	}

	return
}

func OrderRepair(c *gin.Context) {
	var (
		responseData req.OrderRepairForm
		orderRsp     rsp.OrderRsp
	)

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&responseData, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := request.NewCall("http://127.0.0.1:20020/temp/api/order/number", responseData, &orderRsp)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	db := global.DB

	var (
		compOrders []model.CompleteOrder
	)

	length := len(responseData.Numbers)

	err, compOrders = model.GetCompleteOrderListByNumbers(db, responseData.Numbers)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	if len(compOrders) != length {
		xsq_net.ErrorJSON(c, errors.New("完成订单数量与传递数量不一致"))
		return
	}

	var orders []model.Order

	for _, info := range orderRsp.Data {
		hasRemark := 1

		if info.OrderRemark != "" {
			hasRemark = 2
		}

		orders = append(orders, model.Order{
			Id:               info.OrderID,
			ShopId:           info.ShopID,
			ShopName:         info.ShopName,
			ShopType:         info.ShopType,
			ShopCode:         info.ShopCode,
			Number:           info.Number,
			HouseCode:        info.HouseCode,
			Line:             info.Line,
			DistributionType: info.DistributionType,
			OrderRemark:      info.OrderRemark,
			PayAt:            info.PayAt,
			DeliveryAt:       info.DeliveryAt,
			Province:         info.Province,
			City:             info.City,
			District:         info.District,
			Address:          info.Address,
			ConsigneeName:    info.ConsigneeName,
			ConsigneeTel:     info.ConsigneeTel,
			HasRemark:        hasRemark,
		})
	}

	tx := db.Begin()

	err = model.OrderBatchSave(tx, orders)
	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	err = model.DeleteCompOrderByNumbers(tx, responseData.Numbers)
	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

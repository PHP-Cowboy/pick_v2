package dao

import (
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/utils/slice"
	"pick_v2/utils/timeutil"
	"time"

	"pick_v2/forms/rsp"
	"pick_v2/model"
)

// MQ 订单数据处理 -- 正常拣货
func Shipping(db *gorm.DB, form req.PurchaseOrderForm, info rsp.OrderInfo) (consumer.ConsumeResult, error) {
	var (
		hasRemark  int //是否备注
		orderGoods []model.OrderGoods
	)

	err, exist := OrderOrCompleteOrderExist(db, form.OrderId, info.Number)
	//出现错误，消息消费重试
	if err != nil {
		return consumer.ConsumeRetryLater, err
	}

	//已存在，消息消费成功
	if exist {
		return consumer.ConsumeSuccess, errors.New("此订单在订单或完成订单中已存在")
	}

	for _, goods := range info.GoodsInfo {

		orderGoods = append(orderGoods, model.OrderGoods{
			Id:            goods.ID,
			Number:        info.Number,
			GoodsName:     goods.Name,
			Sku:           goods.Sku,
			GoodsType:     goods.GoodsType,
			GoodsSpe:      goods.GoodsSpe,
			Shelves:       goods.Shelves,
			DiscountPrice: goods.DiscountPrice,
			GoodsUnit:     goods.GoodsUnit,
			SaleUnit:      goods.SaleUnit,
			SaleCode:      goods.SaleCode,
			PayCount:      goods.PayCount,
			LackCount:     goods.PayCount, //欠货数 默认等于 下单数
			GoodsRemark:   goods.GoodsRemark,
			Status:        model.OrderGoodsUnhandledStatus,
		})

		//商品有备注 - 即为订单是有备注的
		if hasRemark != 1 && goods.GoodsRemark != "" {
			hasRemark = 1
		}
	}

	//订单有备注 - 即为订单是有备注的
	if hasRemark != 1 && info.OrderRemark != "" {
		hasRemark = 1
	}

	tx := db.Begin()

	err = model.OrderCreate(tx, &model.Order{
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

	if err != nil {
		tx.Rollback()
		return consumer.ConsumeRetryLater, err
	}

	err = model.OrderGoodsBatchSave(tx, &orderGoods)
	if err != nil {
		tx.Rollback()
		return consumer.ConsumeRetryLater, err
	}

	tx.Commit()

	return consumer.ConsumeSuccess, nil
}

// MQ 订单数据处理 -- 无需发货
func NoShipping(db *gorm.DB, form req.PurchaseOrderForm, info rsp.OrderInfo) (consumer.ConsumeResult, error) {

	var (
		completeOrderDetail []model.CompleteOrderDetail
	)

	//查询订单和完成订单中是否已存在
	//订单中也查询，避免因拣货的错误导致拣货系统也出现错误(极端情况，正常不可能出现,也处理一下，以防万一)
	err, exist := OrderOrCompleteOrderExist(db, form.OrderId, info.Number)
	//出现错误，消息消费重试
	if err != nil {
		return consumer.ConsumeRetryLater, err
	}

	//以存在，消息消费成功
	if exist {
		return consumer.ConsumeSuccess, errors.New("此订单在订单或完成订单中已存在")
	}

	for _, goods := range info.GoodsInfo {

		completeOrderDetail = append(completeOrderDetail, model.CompleteOrderDetail{
			Number:      info.Number,
			GoodsName:   goods.Name,
			Sku:         goods.Sku,
			GoodsType:   goods.GoodsType,
			GoodsSpe:    goods.GoodsSpe,
			Shelves:     goods.Shelves,
			PayCount:    goods.PayCount,
			CloseCount:  0,
			ReviewCount: goods.PayCount,
			GoodsRemark: goods.GoodsRemark,
		})
	}

	tx := db.Begin()

	err = model.CompleteOrderSave(tx, &model.CompleteOrder{
		Number:         info.Number,
		OrderRemark:    info.OrderRemark,
		ShopId:         info.ShopID,
		ShopName:       info.ShopName,
		ShopType:       info.ShopType,
		ShopCode:       info.ShopCode,
		Line:           info.Line,
		DeliveryMethod: info.DistributionType,
		Province:       info.Province,
		City:           info.City,
		District:       info.District,
		PayAt:          info.PayAt,
	})

	if err != nil {
		tx.Rollback()
		return consumer.ConsumeRetryLater, err
	}

	err = model.CompleteOrderDetailBatchSave(tx, &completeOrderDetail)

	if err != nil {
		tx.Rollback()
		return consumer.ConsumeRetryLater, err
	}

	tx.Commit()

	return consumer.ConsumeSuccess, nil
}

// 完成订单
func CompleteOrder(db *gorm.DB, form req.CompleteOrderForm) (err error, res rsp.CompleteOrderRsp) {
	var completeOrder []model.CompleteOrder

	numbers := []string{}

	if form.Sku != "" {
		var completeOrderDetail []model.CompleteOrderDetail

		err, completeOrderDetail = model.GetCompleteOrderDetailBySku(db, form.Sku)

		if err != nil {
			return
		}

		for _, detail := range completeOrderDetail {
			numbers = append(numbers, detail.Number)
		}

		numbers = slice.UniqueSlice(numbers)
	}

	//商品
	local := db.
		Model(&model.CompleteOrder{}).
		Where(&model.CompleteOrder{
			ShopId:         form.ShopId,
			Number:         form.Number,
			Line:           form.Line,
			DeliveryMethod: form.DeliveryMethod,
			ShopType:       form.ShopType,
			Province:       form.Province,
			City:           form.City,
			District:       form.District,
		})

	if len(numbers) > 0 {
		local.Where("number in (?)", numbers)
	}

	if form.IsRemark == 1 { //没有备注
		local.Where("order_remark == ''")
	} else if form.IsRemark == 2 { //有备注
		local.Where("order_remark != ''")
	}

	if form.PayAt != "" {

		var payAt time.Time

		//将支付时间字符串格式时间转换成当天的结束时间
		err, payAt = timeutil.StrToLastTime(form.PayAt)

		if err != nil {
			return
		}

		local.Where("pay_at <= ", payAt)
	}

	result := local.Find(&completeOrder)

	if result.Error != nil {
		err = result.Error
		return
	}

	if len(numbers) == 0 {
		for _, order := range completeOrder {
			numbers = append(numbers, order.Number)
		}
	}

	res.Total = result.RowsAffected

	result = local.Scopes(model.Paginate(form.Page, form.Size)).Find(&completeOrder)

	if result.Error != nil {
		err = result.Error
		return
	}

	err, mpNums := model.CountCompleteOrderNumsByNumber(db, numbers)

	if err != nil {
		return
	}

	list := make([]rsp.CompleteOrder, 0, form.Size)

	for _, o := range completeOrder {

		nums, mpNumsOk := mpNums[o.Number]

		if !mpNumsOk {
			err = errors.New("订单相关数量统计有误")
			return
		}

		list = append(list, rsp.CompleteOrder{
			Number:         o.Number,
			PayAt:          o.PayAt,
			ShopCode:       o.ShopCode,
			ShopName:       o.ShopName,
			ShopType:       o.ShopType,
			PayCount:       nums.SumPayCount,
			OutCount:       nums.SumReviewCount,
			CloseCount:     nums.SumCloseCount,
			Line:           o.Line,
			DeliveryMethod: o.DeliveryMethod,
			Region:         fmt.Sprintf("%s-%s-%s", o.Province, o.City, o.District),
			PickTime:       o.PickTime,
			OrderRemark:    o.OrderRemark,
		})
	}

	res.List = list
	return
}

// 查询订单是否在拣货系统已存在[订单表、完成订单表]
func OrderOrCompleteOrderExist(db *gorm.DB, ids []int, number string) (err error, exist bool) {

	// 根据id查询订单表中是否已存在
	err, exist = model.FindOrderExistByIds(db, ids)

	//有错误或者已存在，直接返回
	if err != nil || exist {
		return
	}

	//查询完成订单里是否存在
	err, exist = model.FindCompleteOrderExist(db, number)

	return
}

// 订单出货记录
func OrderShippingRecord(db *gorm.DB, form req.OrderShippingRecordReq) (err error, res rsp.OrderShippingRecordRsp) {
	var (
		picks   []model.Pick
		pickIds []int
	)

	//根据出库单号获取出库记录数据
	err, picks = model.GetPickListByDeliveryNo(db, form.DeliveryOrderNo)

	if err != nil {
		return
	}

	list := make([]rsp.OrderShippingRecord, 0, len(picks))

	for _, pk := range picks {
		pickIds = append(pickIds, pk.Id)
	}

	query := "pick_id,sum(review_num) as review_num"

	err, numsMp := model.CountPickPoolNumsByPickIds(db, pickIds, query)
	if err != nil {
		return
	}

	for _, p := range picks {
		nums, numsMpOk := numsMp[p.Id]

		if !numsMpOk {
			err = errors.New("拣货相关数量统计有误")
		}

		list = append(list, rsp.OrderShippingRecord{
			Id:              p.Id,
			TakeOrdersTime:  p.TakeOrdersTime,
			PickUser:        p.PickUser,
			ReviewUser:      p.ReviewUser,
			ReviewTime:      p.ReviewTime,
			ReviewNum:       nums.ReviewNum,
			DeliveryOrderNo: p.DeliveryOrderNo,
		})
	}

	res.List = list

	return
}

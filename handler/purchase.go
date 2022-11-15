package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/request"
)

type Form struct {
	OrderId []int `json:"order_id"`
}

func Order(ctx context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	var (
		form Form
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
		global.Logger["err"].Infof(err.Error())
		return consumer.ConsumeRetryLater, err
	}

	for _, info := range orderRsp.Data {

		if info.DistributionType == 6 { //无需拣货
			return NoShipping(form, info)
		} else {
			return Shipping(form, info)
		}
	}

	return consumer.ConsumeRetryLater, errors.New("异常")
}

func GetOrderInfo(responseData interface{}) (rsp.OrderRsp, error) {
	var result rsp.OrderRsp

	body, err := request.Post("api/v1/remote/get/goods/by/id", responseData)

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return result, err
	}

	if result.Code != 200 {
		return result, errors.New(result.Msg)
	}

	return result, nil
}

// 正常拣货
func Shipping(form Form, info rsp.OrderInfo) (consumer.ConsumeResult, error) {
	//是否备注
	var (
		hasRemark     int
		order         []model.Order
		orderGoods    []model.OrderGoods
		completeOrder []model.CompleteOrder
	)

	payTotal := 0

	db := global.DB

	// 查询是否已存在 存在的过滤掉
	result := db.Where("id in (?)", form.OrderId).Find(&order)

	if result.Error != nil {
		return consumer.ConsumeRetryLater, result.Error
	}

	if len(order) > 0 {
		return consumer.ConsumeSuccess, errors.New("订单已存在")
	}

	//查看完成订单里有没有
	result = db.Where("number = ?", info.Number).Find(&completeOrder)

	if result.Error != nil {
		return consumer.ConsumeRetryLater, result.Error
	}

	if len(completeOrder) > 0 {
		return consumer.ConsumeSuccess, errors.New("订单在完成订单中已存在")
	}

	for _, goods := range info.GoodsInfo {
		//订单商品总数量
		payTotal += goods.PayCount

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
		})

		//商品有备注 - 即为订单是有备注的
		if goods.GoodsRemark != "" {
			hasRemark = 1
		}
	}

	//订单有备注 - 即为订单是有备注的
	if info.OrderRemark != "" {
		hasRemark = 1
	}

	order = append(order, model.Order{
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
		PayTotal:         payTotal,
		UnPicked:         payTotal,
		DeliveryAt:       info.DeliveryAt,
		Province:         info.Province,
		City:             info.City,
		District:         info.District,
		Address:          info.Address,
		ConsigneeName:    info.ConsigneeName,
		ConsigneeTel:     info.ConsigneeTel,
		HasRemark:        hasRemark,
	})

	tx := db.Begin()

	result = tx.Save(&order)

	if result.Error != nil {
		tx.Rollback()
		return consumer.ConsumeRetryLater, result.Error
	}

	result = tx.Save(&orderGoods)

	if result.Error != nil {
		tx.Rollback()
		return consumer.ConsumeRetryLater, result.Error
	}

	tx.Commit()

	return consumer.ConsumeSuccess, nil
}

// 无需发货
func NoShipping(form Form, info rsp.OrderInfo) (consumer.ConsumeResult, error) {

	var (
		completeOrder       []model.CompleteOrder
		completeOrderDetail []model.CompleteOrderDetail
	)

	payTotal := 0

	// 查询是否已存在 存在的过滤掉？
	db := global.DB

	result := db.Where("id in (?)", form.OrderId).Find(&completeOrder)

	if result.Error != nil {
		return consumer.ConsumeRetryLater, result.Error
	}

	if len(completeOrder) > 0 {
		return consumer.ConsumeSuccess, errors.New("订单已存在")
	}

	for _, goods := range info.GoodsInfo {
		//订单商品总数量
		payTotal += goods.PayCount

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

	completeOrder = append(completeOrder, model.CompleteOrder{
		Number:         info.Number,
		OrderRemark:    info.OrderRemark,
		ShopId:         info.ShopID,
		ShopName:       info.ShopName,
		ShopType:       info.ShopType,
		ShopCode:       info.ShopCode,
		Line:           info.Line,
		DeliveryMethod: info.DistributionType,
		PayCount:       payTotal,
		CloseCount:     0,
		OutCount:       payTotal,
		Province:       info.Province,
		City:           info.City,
		District:       info.District,
		PayAt:          info.PayAt,
	})

	tx := db.Begin()

	result = tx.Save(&completeOrder)

	if result.Error != nil {
		tx.Rollback()
		return consumer.ConsumeRetryLater, result.Error
	}

	result = tx.Save(&completeOrderDetail)

	if result.Error != nil {
		tx.Rollback()
		return consumer.ConsumeRetryLater, result.Error
	}

	tx.Commit()

	return consumer.ConsumeSuccess, nil
}

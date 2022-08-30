package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	order2 "pick_v2/model/order"
	"pick_v2/utils/request"
)

func Order(ctx context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	type Form struct {
		OrderId []int `json:"order_id"`
	}

	var (
		order      []order2.Order
		orderGoods []order2.OrderGoods
		form       Form
		Id         int
	)

	for i := range messages {

		err := json.Unmarshal(messages[i].Body, &Id)
		if err != nil {
			global.SugarLogger.Errorf("解析json失败:%s", err.Error())
			return consumer.ConsumeSuccess, nil
		}
		form.OrderId = append(form.OrderId, Id)
	}

	// todo 查询是否已存在 存在的过滤掉？

	orderRsp, err := GetOrderInfo(form)

	if err != nil {
		return consumer.ConsumeRetryLater, err
	}

	for _, info := range orderRsp.Data {

		//是否备注
		var isRemark int

		for _, goods := range info.GoodsInfo {
			orderGoods = append(orderGoods, order2.OrderGoods{
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
				isRemark = 1
			}
		}

		//订单有备注 - 即为订单是有备注的
		if info.OrderRemark != "" {
			isRemark = 1
		}

		order = append(order, order2.Order{
			Id:               info.OrderID,
			ShopId:           info.ShopID,
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
			IsRemark:         isRemark,
		})
	}

	tx := global.DB.Begin()

	result := tx.Save(order)

	if result.Error != nil {
		tx.Rollback()
		return consumer.ConsumeRetryLater, result.Error
	}

	result = tx.Save(orderGoods)

	if result.Error != nil {
		tx.Rollback()
		return consumer.ConsumeRetryLater, result.Error
	}

	tx.Commit()

	return consumer.ConsumeSuccess, nil
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

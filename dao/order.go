package dao

import (
	"errors"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"gorm.io/gorm"
	"pick_v2/forms/req"

	"pick_v2/forms/rsp"
	"pick_v2/model"
)

// MQ 订单数据处理 -- 正常拣货
func Shipping(db *gorm.DB, form req.PurchaseOrderForm, info rsp.OrderInfo) (consumer.ConsumeResult, error) {
	var (
		hasRemark  int //是否备注
		orderGoods []model.OrderGoods
	)

	err, exist := model.OrderOrCompleteOrderExist(db, form.OrderId, info.Number)
	//出现错误，消息消费重试
	if err != nil {
		return consumer.ConsumeRetryLater, err
	}

	//以存在，消息消费成功
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

	err = model.OrderSave(tx, &model.Order{
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
	err, exist := model.OrderOrCompleteOrderExist(db, form.OrderId, info.Number)
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

// 变更订单类型&&商品状态
func UpdateOrderAndGoods(db *gorm.DB, orderIds []int, orderGoodsIds []int) error {

	err := model.UpdateOrderByIds(db, orderIds, map[string]interface{}{"order_type": model.PickingOrderType})

	if err != nil {
		return err
	}

	err = model.UpdateOrderGoodsByIds(db, orderGoodsIds, map[string]interface{}{"status": model.OrderGoodsProcessingStatus})

	if err != nil {
		return err
	}

	return nil
}

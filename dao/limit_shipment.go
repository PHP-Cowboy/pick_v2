package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"pick_v2/forms/req"
	"pick_v2/model"
	"pick_v2/utils/ecode"
)

// 订单批量限发
func OrderLimit(db *gorm.DB, form req.OrderLimitForm) error {

	length := len(form.OrderLimit)

	var (
		outboundGoodsJoinOrder []model.OutboundGoodsJoinOrder
		limitShipment          = make([]model.LimitShipment, 0, length)
		sku                    = make([]string, 0, length)
		limitMp                = make(map[string]int, 0)
		outboundGoods          = make([]model.OutboundGoods, 0, length)
	)

	for _, limit := range form.OrderLimit {
		sku = append(sku, limit.Sku)
		limitMp[limit.Sku] = limit.LimitNum
	}

	result := db.Model(&model.OutboundGoodsJoinOrder{}).
		Select("sku,shop_name,goods_name,goods_spe").
		Where(&model.OutboundGoodsJoinOrder{
			TaskId: form.TaskId,
			Number: form.Number,
		}).
		Where("sku in (?)", sku).
		Find(&outboundGoodsJoinOrder)

	if result.Error != nil {
		return result.Error
	}

	for _, orderGoods := range outboundGoodsJoinOrder {

		limitNum, limitOk := limitMp[orderGoods.Sku]

		if !limitOk {
			continue
		}

		limitShipment = append(limitShipment, model.LimitShipment{
			TaskId:    form.TaskId,
			Number:    form.Number,
			Sku:       orderGoods.Sku,
			ShopName:  orderGoods.ShopName,
			GoodsName: orderGoods.GoodsName,
			GoodsSpe:  orderGoods.GoodsSpe,
			LimitNum:  limitNum,
		})

		outboundGoods = append(outboundGoods, model.OutboundGoods{
			TaskId:   form.TaskId,
			Number:   form.Number,
			Sku:      orderGoods.Sku,
			LimitNum: limitNum,
		})
	}

	if len(limitShipment) == 0 {
		return ecode.DataNotExist
	}

	tx := db.Begin()

	err := model.LimitShipmentSave(tx, limitShipment)

	if err != nil {
		tx.Rollback()
		return err
	}

	//更新出库单商品的限发数量
	err = UpdateOutboundGoodsLimit(tx, outboundGoods)

	if err != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()

	return nil
}

// 更新出库单商品的限发数量
func UpdateOutboundGoodsLimit(db *gorm.DB, list []model.OutboundGoods) error {
	result := db.Model(&model.OutboundGoods{}).
		Select("limit_num").
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "task_id,number,sku"}},
			DoUpdates: clause.AssignmentColumns([]string{"limit_num"}),
		}).
		Save(&list)

	return result.Error
}

func TaskLimit(db *gorm.DB, form req.TaskLimitForm) error {

	var (
		outboundGoodsJoinOrder []model.OutboundGoodsJoinOrder
	)

	result := db.Table("t_outbound_goods og").
		Select("og.number,shop_name,goods_name,goods_spe").
		Joins("left join t_outbound_order o on og.task_id = o.task_id and og.number = o.number").
		Where(&model.OutboundGoodsJoinOrder{
			TaskId: form.TaskId,
			Sku:    form.Sku,
		}).
		Find(&outboundGoodsJoinOrder)

	if result.Error != nil {
		return result.Error
	}

	var (
		limitShipment = make([]model.LimitShipment, 0, len(outboundGoodsJoinOrder))
		outboundGoods = make([]model.OutboundGoods, 0, len(outboundGoodsJoinOrder))
	)

	for _, order := range outboundGoodsJoinOrder {
		limitShipment = append(limitShipment, model.LimitShipment{
			TaskId:    form.TaskId,
			Number:    order.Number,
			Sku:       form.Sku,
			ShopName:  order.ShopName,
			GoodsName: order.GoodsName,
			GoodsSpe:  order.GoodsSpe,
			LimitNum:  form.LimitNum,
		})

		outboundGoods = append(outboundGoods, model.OutboundGoods{
			TaskId:   form.TaskId,
			Number:   order.Number,
			Sku:      form.Sku,
			LimitNum: form.LimitNum,
		})
	}

	tx := db.Begin()

	err := model.LimitShipmentSave(tx, limitShipment)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = UpdateOutboundGoodsLimit(tx, outboundGoods)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

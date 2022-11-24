package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
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
			Status:    model.LimitShipmentStatusNormal,
			Typ:       model.LimitShipmentTypOrder,
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
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "task_id,number,sku"}},
			DoUpdates: clause.AssignmentColumns([]string{"limit_num"}),
		}).
		Save(&list)

	return result.Error
}

// 任务批量限发
func TaskLimit(db *gorm.DB, form req.TaskLimitForm) error {

	var (
		outboundGoodsJoinOrder []model.OutboundGoodsJoinOrder
	)

	result := db.Table("t_outbound_goods og").
		Select("og.number,shop_name,goods_name,goods_spe").
		Joins("left join t_outbound_order o on og.task_id = o.task_id and og.number = o.number").
		Where("og.task_id = ? and og.sku = ?", form.TaskId, form.Sku).
		Find(&outboundGoodsJoinOrder)

	if result.Error != nil {
		return result.Error
	}

	if len(outboundGoodsJoinOrder) == 0 {
		return errors.New("任务中没有所选sku请重试")
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
			Status:    model.LimitShipmentStatusNormal,
			Typ:       model.LimitShipmentTypTask,
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

// 撤销限发
func RevokeLimit(db *gorm.DB, form req.RevokeLimitForm) error {

	var outboundGoods model.OutboundGoods

	result := db.Model(&model.OutboundGoods{}).
		Where("task_id = ? and number = ? and sku = ?", form.TaskId, form.Number, form.Sku).
		First(&outboundGoods)

	if result.Error != nil {
		return result.Error
	}

	if outboundGoods.Status != model.OutboundGoodsStatusUnhandled {
		return errors.New("当前订单不允许撤销限发")
	}

	tx := db.Begin()

	result = tx.Model(&model.LimitShipment{}).
		Where("task_id = ? and number = ? and sku = ?", form.TaskId, form.Number, form.Sku).
		Update("status", model.LimitShipmentStatusRevoke)

	if result.Error != nil {
		return result.Error
	}

	//更新限发数量为原始值(欠货数)
	outboundGoods.LimitNum = outboundGoods.LackCount

	list := []model.OutboundGoods{outboundGoods}

	err := model.OutboundGoodsReplaceSave(tx, list, []string{"limit_num"})

	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

// 限发列表
func LimitShipmentList(db *gorm.DB, form req.LimitShipmentListForm) (error, rsp.LimitShipmentListRsp) {
	var res rsp.LimitShipmentListRsp

	err, total, limitShipmentList := model.GetLimitShipmentPageListByTaskIdAndNumber(db, form.TaskId, form.Number, form.Page, form.Size)
	if err != nil {
		return err, res
	}

	res.Total = total

	list := make([]rsp.LimitShipmentList, 0, len(limitShipmentList))

	for _, shipment := range limitShipmentList {
		list = append(list, rsp.LimitShipmentList{
			OutboundNumber: model.GetOutboundNumber(shipment.TaskId, shipment.Number),
			Number:         shipment.Number,
			Sku:            shipment.Sku,
			ShopName:       shipment.ShopName,
			GoodsName:      shipment.GoodsName,
			GoodsSpe:       shipment.GoodsSpe,
			LimitNum:       shipment.LimitNum,
			Status:         shipment.Status,
		})
	}

	res.List = list

	return nil, res
}

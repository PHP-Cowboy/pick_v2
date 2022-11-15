package dao

import (
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/utils/timeutil"
	"strings"
	"time"
)

type NumberMp struct {
	LimitNum  int
	HasRemark int
}

func OutboundTaskList(db *gorm.DB, form req.OutboundTaskListForm) (err error, res rsp.OutboundTaskListRsp) {

	var taskIds []int

	if form.ShopId > 0 || form.Number != "" || form.Sku != "" {
		var outboundOrderAndGoods []model.OutboundGoodsJoinOrder

		orderDb := db.Table("t_outbound_goods og").
			Select("task_id").
			Joins("left join t_outbound_order oo on og.task_id = oo.task_id and og.number = oo.number")

		if form.ShopId > 0 {
			orderDb = orderDb.Where("oo.shop_id = ?", form.ShopId)
		}

		if form.Number != "" {
			orderDb = orderDb.Where("oo.number = ?", form.Number)
		}

		if form.Sku != "" {
			orderDb = orderDb.Where("og.sku = ?", form.Sku)
		}

		orderRes := orderDb.Find(&outboundOrderAndGoods)

		if orderRes.Error != nil {
			return orderRes.Error, res
		}

		for _, good := range outboundOrderAndGoods {
			taskIds = append(taskIds, good.TaskId)
		}
	}

	localDb := db.Model(&model.OutboundTask{})

	if len(taskIds) > 0 {
		localDb.Where("id in (?)", taskIds)
	}

	localDb.Where(&model.OutboundTask{
		Line:             form.Line,
		DistributionType: form.DistributionType,
		Status:           form.Status,
	})

	if form.Status == model.OutboundTaskStatusClosed {
		if form.StartTime != "" {
			localDb.Where("update_time >= ?", form.StartTime)
		}

		if form.EndTime != "" {
			localDb.Where("update_time <= ?", form.EndTime)
		}
	}

	var outboundTask []model.OutboundTask

	result := localDb.Find(&outboundTask)

	if result.Error != nil {
		return result.Error, res
	}

	res.Total = result.RowsAffected

	result = localDb.Scopes(model.Paginate(form.Page, form.Size)).Find(&outboundTask)

	if result.Error != nil {
		return result.Error, res
	}

	list := make([]rsp.OutboundTaskList, 0, len(outboundTask))

	for _, task := range outboundTask {
		list = append(list, rsp.OutboundTaskList{
			Id:                task.Id,
			TaskName:          task.TaskName,
			DeliveryStartTime: task.DeliveryStartTime,
			DeliveryEndTime:   task.DeliveryEndTime,
			Line:              task.Line,
			DistributionType:  task.DistributionType,
			EndTime:           task.EndTime,
			Status:            task.Status,
			IsPush:            task.IsPush,
		})
	}

	res.List = list

	return nil, res
}

// 出库任务相关保存
func OutboundTaskSave(db *gorm.DB, form req.CreateOutboundForm, userInfo *middlewares.CustomClaims) (taskId int, err error) {
	deliveryStartTime, deliveryStartTimeErr := time.ParseInLocation(timeutil.TimeFormat, form.DeliveryStartTime, time.Local)

	if deliveryStartTimeErr != nil {
		return taskId, deliveryStartTimeErr
	}

	deliveryEndTime, deliveryEndTimeErr := time.ParseInLocation(timeutil.TimeFormat, form.DeliveryEndTime, time.Local)

	if deliveryEndTimeErr != nil {
		return taskId, deliveryEndTimeErr
	}

	payTime, payTimeErr := time.ParseInLocation(timeutil.TimeFormat, form.PayTime, time.Local)

	if payTimeErr != nil {
		return taskId, payTimeErr
	}

	task := model.OutboundTask{
		TaskName:          form.OutboundName,
		DeliveryStartTime: (*model.MyTime)(&deliveryStartTime),
		DeliveryEndTime:   (*model.MyTime)(&deliveryEndTime),
		Line:              strings.Join(form.Lines, ""),
		DistributionType:  form.DistributionType,
		EndTime:           (*model.MyTime)(&payTime),
		Creator: model.Creator{
			CreatorId: userInfo.ID,
			Creator:   userInfo.Name,
		},
	}

	err = model.OutboundTaskSave(db, &task)

	if err != nil {
		return taskId, err
	}

	return task.Id, nil
}

// 出库订单相关保存
func OutboundOrderBatchSave(db *gorm.DB, form req.CreateOutboundForm, taskId int) error {

	var (
		orderJoinGoods []model.OrderJoinGoods
	)

	localDb := db.Table("t_order_goods og").
		Joins("left join t_order o on og.number = o.number").
		Select("og.*,o.shop_id,o.shop_name,o.shop_code,o.line,o.distribution_type,o.order_remark").
		Where("o.distribution_type = ? and o.pay_at <= ? and o.delivery_at <= ? ", form.DistributionType, form.PayTime, form.DeliveryEndTime)

	if form.DeliveryStartTime != "" {
		localDb = localDb.Where("o.delivery_at >= ?", form.DeliveryStartTime)
	}

	if len(form.Sku) > 0 {
		localDb = localDb.Where("og.sku in (?)", form.Sku)
	}

	if len(form.Lines) > 0 {
		localDb = localDb.Where("o.line in (?) ", form.Lines)
	}

	//新订单或欠货的订单商品数据
	result := localDb.Where("og.status = ？", model.OrderGoodsUnhandledStatus).
		Find(&orderJoinGoods)

	if result.Error != nil {
		return result.Error
	}

	var (
		outboundOrderMp = make(map[string]model.OutboundOrder, 0)
		numberMp        = make(map[string]NumberMp, 0)
		outboundOrders  = make([]model.OutboundOrder, 0)
		outboundGoods   = make([]model.OutboundGoods, 0, len(orderJoinGoods))
	)

	for _, goods := range orderJoinGoods {
		_, ok := outboundOrderMp[goods.Number]

		if !ok {

			outboundOrderMp[goods.Number] = model.OutboundOrder{
				TaskId:            taskId,
				Number:            goods.Number,
				PayAt:             &goods.PayAt,
				ShopId:            goods.ShopId,
				ShopName:          goods.ShopName,
				ShopType:          goods.ShopType,
				ShopCode:          goods.ShopCode,
				HouseCode:         goods.HouseCode,
				DistributionType:  goods.DistributionType,
				GoodsNum:          goods.PayCount,
				LimitNum:          0,
				CloseNum:          goods.CloseNum,
				Line:              goods.Line,
				Province:          goods.Province,
				City:              goods.City,
				District:          goods.District,
				Address:           goods.Address,
				ConsigneeName:     goods.ConsigneeName,
				ConsigneeTel:      goods.ConsigneeTel,
				OrderType:         model.OutboundOrderTypeNew,
				LatestPickingTime: nil,
				HasRemark:         0,
				OrderRemark:       goods.OrderRemark,
			}
		}

		mp, _ := numberMp[goods.Number]

		mp.LimitNum += goods.LackCount

		if goods.OrderRemark != "" || goods.GoodsRemark != "" {
			mp.HasRemark = 1
		}

		numberMp[goods.Number] = mp

		outboundGoods = append(outboundGoods, model.OutboundGoods{
			TaskId:          taskId,
			Number:          goods.Number,
			Sku:             goods.Sku,
			OrderGoodsId:    goods.Id, //order_goods表ID
			GoodsName:       goods.GoodsName,
			GoodsType:       goods.GoodsType,
			GoodsSpe:        goods.GoodsSpe,
			Shelves:         goods.Shelves,
			DiscountPrice:   goods.DiscountPrice,
			GoodsUnit:       goods.GoodsUnit,
			SaleUnit:        goods.SaleUnit,
			SaleCode:        goods.SaleCode,
			PayCount:        goods.PayCount,
			CloseCount:      goods.CloseCount,
			LackCount:       goods.LackCount,
			OutCount:        goods.OutCount,
			LimitNum:        goods.LackCount,
			GoodsRemark:     goods.GoodsRemark,
			BatchId:         0,
			DeliveryOrderNo: nil,
		})
	}

	for s, order := range outboundOrderMp {

		val, ok := numberMp[s]

		if ok {
			order.HasRemark = val.HasRemark
			order.LimitNum = val.LimitNum
		}

		outboundOrders = append(outboundOrders, order)
	}

	err := model.OutboundOrderBatchSave(db, outboundOrders)
	if err != nil {
		return err
	}

	err = model.OutboundGoodsBatchSave(db, outboundGoods)

	if err != nil {
		return err
	}

	return nil
}

// 出库订单列表
func OutboundOrderList(db *gorm.DB, form req.OutboundOrderListForm) (err error, res rsp.OutboundTaskListRsp) {

	if form.Sku != "" {

	}

	localDb := db.Model(&model.OutboundOrder{}).Where(&model.OutboundOrder{
		TaskId:           form.TaskId,
		Number:           form.Number,
		ShopId:           form.ShopId,
		ShopType:         form.ShopType,
		DistributionType: form.DistributionType,
		Line:             form.Line,
		Province:         form.Province,
		City:             form.City,
		District:         form.District,
		OrderType:        form.OrderType,
	})

	if form.HasRemark != nil {
		localDb.Where("has_remark = ?", *form.HasRemark)
	}

	return nil, res
}

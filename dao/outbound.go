package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/slice"
	"pick_v2/utils/timeutil"
	"strings"
	"time"
)

type NumberMp struct {
	LimitNum  int
	HasRemark int
}

// 出库单任务列表
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
			PayEndTime:        task.PayEndTime,
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
		PayEndTime:        (*model.MyTime)(&payTime),
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
		outboundOrderMp = make(map[string]model.OutboundOrder, 0)             //出库订单map 以订单号为key，用于更新订单备注以及限发总数，然后存储
		numberMp        = make(map[string]NumberMp, 0)                        //订单map，用于处理 订单备注以及限发总数
		outboundOrders  = make([]model.OutboundOrder, 0)                      //出库订单
		outboundGoods   = make([]model.OutboundGoods, 0, len(orderJoinGoods)) //出库订单商品
		order           []model.Order                                         //订单数据，用于更新订单表数据
		orderGoods      []model.OrderGoods                                    //订单商品数据，用于更新订单商品表数据
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

			//更新订单表数据
			order = append(order, model.Order{
				ShopId:    goods.ShopId,
				ShopName:  goods.ShopName,
				ShopType:  goods.ShopType,
				ShopCode:  goods.ShopCode,
				Number:    goods.Number,
				HouseCode: goods.HouseCode,
				Line:      goods.Line,
				OrderType: model.NewOrderType,
			})
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

		//更新订单商品表数据
		orderGoods = append(orderGoods, model.OrderGoods{
			Id:     goods.Id,
			Status: model.OrderGoodsUnhandledStatus,
		})
	}

	for s, oo := range outboundOrderMp {

		val, ok := numberMp[s]

		if ok {
			oo.HasRemark = val.HasRemark
			oo.LimitNum = val.LimitNum
		}

		outboundOrders = append(outboundOrders, oo)
	}

	//出库订单保存
	err := model.OutboundOrderBatchSave(db, outboundOrders)
	if err != nil {
		return err
	}

	//出库商品保存
	err = model.OutboundGoodsBatchSave(db, outboundGoods)

	if err != nil {
		return err
	}

	//更新订单数据
	err = UpdateOrderAndGoods(db, order, orderGoods)

	if err != nil {
		return err
	}

	return nil
}

// 出库订单列表
func OutboundOrderList(db *gorm.DB, form req.OutboundOrderListForm) (err error, res rsp.OutboundOrderListRsp) {

	var (
		numbers        []string
		outboundOrders []model.OutboundOrder
	)

	if form.Sku != "" {
		var outboundGoods []model.OutboundGoods
		goodsRes := db.Model(&model.OutboundGoods{}).
			Select("number").
			Where(&model.OutboundGoods{TaskId: form.TaskId, Sku: form.Sku}).
			Find(&outboundGoods)

		if goodsRes.Error != nil {
			return goodsRes.Error, res
		}

		for _, good := range outboundGoods {
			numbers = append(numbers, good.Number)
		}
	}

	localDb := db.Model(&model.OutboundOrder{}).Where("task_id = ?", form.TaskId)

	if len(numbers) > 0 {
		localDb.Where("number in (?)", numbers)
	}

	localDb.Where(&model.OutboundOrder{
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

	result := localDb.Find(&outboundOrders)

	if result.Error != nil {
		return result.Error, res
	}

	res.Total = result.RowsAffected

	result = localDb.Scopes(model.Paginate(form.Page, form.Size)).Find(&outboundOrders)

	if result.Error != nil {
		return result.Error, res
	}

	list := make([]rsp.OutboundOrderList, 0, len(outboundOrders))

	for _, order := range outboundOrders {
		list = append(list, rsp.OutboundOrderList{
			Number:            order.Number,
			PayAt:             order.PayAt,
			ShopName:          order.ShopName,
			ShopType:          order.ShopType,
			DistributionType:  order.DistributionType,
			GoodsNum:          order.GoodsNum,
			LimitNum:          order.LimitNum,
			CloseNum:          order.CloseNum,
			Line:              order.Line,
			Region:            fmt.Sprintf("%s-%s-%s", order.Province, order.City, order.District),
			LatestPickingTime: order.LatestPickingTime,
			OrderRemark:       order.OrderRemark,
			OrderType:         order.OrderType,
		})
	}

	res.List = list

	return nil, res
}

func OutboundOrderDetail(db *gorm.DB, form req.OutboundOrderDetailForm) (err error, res rsp.OrderDetail) {

	var (
		outboundOrder model.OutboundOrder
		outboundGoods []model.OutboundGoods
	)

	result := db.Model(&model.OutboundOrder{}).
		Where(&model.OutboundOrder{
			TaskId: form.TaskId,
			Number: form.Number,
		}).
		First(&outboundOrder)

	if result.Error != nil {
		return result.Error, res
	}

	result = db.Model(&model.OutboundGoods{}).
		Where(&model.OutboundGoods{
			TaskId: form.TaskId,
			Number: form.Number,
		}).
		Find(&outboundGoods)

	if result.Error != nil {
		return result.Error, res
	}

	mp, err := cache.GetClassification()

	if err != nil {
		return err, res
	}

	detailMap := make(map[string]*rsp.Detail, 0)

	deliveryOrderNoArr := make(model.GormList, 0)

	for _, goods := range outboundGoods {
		goodsType, ok := mp[goods.GoodsType]

		if !ok {
			return errors.New("商品类型:" + goods.GoodsType + "数据未同步"), res
		}

		deliveryOrderNoArr = append(deliveryOrderNoArr, goods.DeliveryOrderNo...)

		if _, detailOk := detailMap[goodsType]; !detailOk {
			detailMap[goodsType] = &rsp.Detail{
				Total: 0,
				List:  make([]*rsp.GoodsDetail, 0),
			}
		}

		detailMap[goodsType].Total += goods.PayCount

		detailMap[goodsType].List = append(detailMap[goodsType].List, &rsp.GoodsDetail{
			Name:        goods.GoodsName,
			GoodsSpe:    goods.GoodsSpe,
			Shelves:     goods.Shelves,
			PayCount:    goods.PayCount,
			CloseCount:  goods.CloseCount,
			LackCount:   goods.LimitNum, //需拣数 以限发数为准
			OutCount:    goods.OutCount,
			GoodsRemark: goods.GoodsRemark,
		})
	}

	res.Number = outboundOrder.Number
	res.PayAt = *outboundOrder.PayAt
	res.ShopCode = outboundOrder.ShopCode
	res.ShopName = outboundOrder.ShopName
	res.Line = outboundOrder.Line
	res.Region = outboundOrder.Province + outboundOrder.City + outboundOrder.District
	res.ShopType = outboundOrder.ShopType
	res.OrderRemark = outboundOrder.OrderRemark

	res.Detail = detailMap

	deliveryOrderNoArr = slice.UniqueStringSlice(deliveryOrderNoArr)
	//历史出库单号
	res.DeliveryOrderNo = deliveryOrderNoArr

	return nil, res
}

// 结束任务
func EndOutboundTask(db *gorm.DB, form req.EndOutboundTaskForm) error {
	return nil
}

// 关闭订单
func OutboundTaskCloseOrder(db *gorm.DB, form req.OutboundTaskCloseOrderForm) (err error) {

	//关闭预拣池订单
	err = ClosePrePickOrder(db, form)
	if err != nil {
		return err
	}

	//关闭出库任务
	err = CloseOutboundOrder(db, form)

	if err != nil {
		return err
	}

	//订单
	err = CloseOrder(db, form)

	if err != nil {
		return err
	}

	return nil
}

// 关闭预拣池订单
func ClosePrePickOrder(db *gorm.DB, form req.OutboundTaskCloseOrderForm) error {

	err, pickGoodsList := model.GetPickGoodsByNumber(db, form.Number)

	if err != nil {
		return err
	}

	if len(pickGoodsList) > 0 {
		var numbers []string
		for _, goods := range pickGoodsList {
			numbers = append(numbers, goods.Number)
		}
		return errors.New("订单号:" + strings.Join(numbers, ",") + "已在拣货池中，请先取消拣货")
	}

	//预拣池
	err, prePickGoodsJoinPrePickList := model.GetPrePickGoodsJoinPrePickListByNumber(db, form.Number)
	if err != nil {
		return err
	}

	var (
		prePickIds      []int
		prePickGoodsIds []int
	)

	for _, pick := range prePickGoodsJoinPrePickList {
		prePickIds = append(prePickIds, pick.PrePickId)
		prePickGoodsIds = append(prePickGoodsIds, pick.PrePickGoodsId)
	}

	prePickIds = slice.UniqueIntSlice(prePickIds)

	//更新预拣池
	err = model.UpdatePrePickStatusByIds(db, prePickIds, model.PrePickStatusClose)

	if err != nil {
		return err
	}

	//更新预拣池商品
	err = model.UpdatePrePickGoodsStatusByIds(db, prePickGoodsIds, model.PrePickGoodsStatusClose)

	if err != nil {
		return err
	}

	return nil
}

// 关闭出库任务
func CloseOutboundOrder(db *gorm.DB, form req.OutboundTaskCloseOrderForm) error {
	//出库任务
	err, outboundGoodsJoinOrderList := model.GetOutboundGoodsJoinOrderListByNumbers(db, form.Number)
	if err != nil {
		return err
	}

	var (
		outboundOrder []model.OutboundOrder
		outboundGoods []model.OutboundGoods
	)

	for _, outbound := range outboundGoodsJoinOrderList {
		outboundOrder = append(outboundOrder, model.OutboundOrder{
			TaskId:    outbound.TaskId,
			Number:    outbound.Number,
			OrderType: model.OutboundOrderTypeClose,
		})

		outboundGoods = append(outboundGoods, model.OutboundGoods{
			TaskId: outbound.TaskId,
			Number: outbound.Number,
			Sku:    outbound.Sku,
			Status: model.OutboundGoodsStatusOutboundClose,
		})
	}

	err = model.OutboundOrderReplaceSave(db, outboundOrder, []string{"order_type"})

	if err != nil {
		return err
	}

	err = model.OutboundGoodsReplaceSave(db, outboundGoods, []string{"status"})

	if err != nil {
		return err
	}

	return nil
}

// 关闭订单
func CloseOrder(db *gorm.DB, form req.OutboundTaskCloseOrderForm) error {
	return nil
}

// 临时加单
func OutboundTaskAddOrder(db *gorm.DB, form req.OutboundTaskAddOrderForm) (err error) {

	var (
		outboundGoodsMp = make(map[string]struct{}, 0)
		outboundOrderMp = make(map[string]struct{}, 0)
		limitShipmentMp = make(map[string]int, 0)
		outboundOrder   []model.OutboundOrder
		outboundGoods   []model.OutboundGoods
	)

	err, limitShipmentList := model.GetLimitShipmentListByTaskIdAndNumbers(db, form.TaskId, form.Number)
	if err != nil {
		return err
	}

	//加单时，如果已经对任务中sku设置了批量限发，则取限发数量
	for _, shipment := range limitShipmentList {
		//不是任务批量限发的，不处理
		if shipment.Typ != model.LimitShipmentTypTask {
			continue
		}
		limitShipmentMp[shipment.Sku] = shipment.LimitNum
	}

	//出库单信息
	err, outboundList := model.GetOutboundGoodsJoinOrderList(db, form.TaskId, form.Number)
	if err != nil {
		return err
	}

	for _, outbound := range outboundList {
		key := fmt.Sprintf("%s%s", outbound.Number, outbound.Sku)
		outboundGoodsMp[key] = struct{}{}
		outboundOrderMp[outbound.Number] = struct{}{}
	}

	//订单信息
	err, orderList := model.GetOrderJoinGoodsList(db, form.Number)
	if err != nil {
		return err
	}

	for _, order := range orderList {
		key := fmt.Sprintf("%s%s", order.Number, order.Sku)

		_, goodsOk := outboundGoodsMp[key]

		//如果当前订单的sku已经在任务中，跳过
		if goodsOk {
			continue
		}

		//key 中包含了number 如果 goodsOk == true 则订单本身已存在了
		_, orderOk := outboundOrderMp[order.Number]

		if !orderOk {
			hasRemark := 0

			if order.OrderRemark != "" {
				hasRemark = 1
			}

			outboundOrder = append(outboundOrder, model.OutboundOrder{
				TaskId:            form.TaskId,
				Number:            order.Number,
				PayAt:             &order.PayAt,
				ShopId:            order.ShopId,
				ShopName:          order.ShopName,
				ShopType:          order.ShopType,
				ShopCode:          order.ShopCode,
				HouseCode:         order.HouseCode,
				DistributionType:  order.DistributionType,
				Line:              order.Line,
				Province:          order.Province,
				City:              order.City,
				District:          order.District,
				Address:           order.Address,
				ConsigneeName:     order.ConsigneeName,
				ConsigneeTel:      order.ConsigneeTel,
				OrderType:         model.OutboundOrderTypeNew,
				LatestPickingTime: nil,
				HasRemark:         hasRemark,
				OrderRemark:       order.OrderRemark,
			})
		}

		//如果任务设置了限发，则取限发数量，否则取欠货数量
		limitNum, limitOk := limitShipmentMp[order.Sku]

		if !limitOk {
			limitNum = order.LackCount
		}

		outboundGoods = append(outboundGoods, model.OutboundGoods{
			TaskId:          form.TaskId,
			Number:          order.Number,
			Sku:             order.Sku,
			OrderGoodsId:    order.Id,
			BatchId:         0,
			GoodsName:       order.GoodsName,
			GoodsType:       order.GoodsType,
			GoodsSpe:        order.GoodsSpe,
			Shelves:         order.Shelves,
			DiscountPrice:   order.DiscountPrice,
			GoodsUnit:       order.GoodsUnit,
			SaleUnit:        order.SaleUnit,
			SaleCode:        order.SaleCode,
			PayCount:        order.PayCount,
			CloseCount:      order.CloseCount,
			LackCount:       order.LackCount,
			OutCount:        order.OutCount,
			LimitNum:        limitNum,
			GoodsRemark:     order.GoodsRemark,
			Status:          model.OutboundGoodsStatusUnhandled,
			DeliveryOrderNo: order.DeliveryOrderNo,
		})
	}

	tx := db.Begin()

	err = model.OutboundOrderBatchSave(tx, outboundOrder)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = model.OutboundGoodsBatchSave(tx, outboundGoods)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// 订单出库记录
func OrderOutboundRecord(db *gorm.DB, form req.OrderOutboundRecordForm) (err error, list rsp.OrderOutboundRecordList) {
	return
}

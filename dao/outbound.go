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
)

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

	//已结束的任务，可以进行结束时间范围搜索
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
			Creator:           task.Creator.Creator,
			UpdateTime:        timeutil.FormatToDateTime(task.UpdateTime),
		})
	}

	res.List = list

	return nil, res
}

// 简化版任务列表
func OutboundTaskListSimple(db *gorm.DB) (error, []rsp.OutboundTaskList) {
	err, dataList := model.GetOutboundTaskStatusOngoingList(db)

	if err != nil {
		return err, nil
	}

	list := make([]rsp.OutboundTaskList, 0, len(dataList))

	for _, d := range dataList {
		list = append(list, rsp.OutboundTaskList{
			Id:       d.Id,
			TaskName: d.TaskName,
		})
	}

	return nil, list
}

// 出库任务状态数量
func OutboundTaskCount(db *gorm.DB) (error, rsp.OutboundTaskCountRsp) {

	err, count := model.OutboundTaskCountGroupStatus(db)
	if err != nil {
		return err, rsp.OutboundTaskCountRsp{}
	}

	var res rsp.OutboundTaskCountRsp

	for _, ct := range count {
		switch ct.Status {
		case model.OutboundTaskStatusOngoing:
			res.Ongoing += ct.Count
			break
		case model.OutboundTaskStatusClosed:
			res.Closed += ct.Count
			break
		}
	}

	return nil, res
}

// 出库任务相关保存
func OutboundTaskSave(db *gorm.DB, form req.CreateOutboundForm, userInfo *middlewares.CustomClaims) (taskId int, err error) {

	deliveryStartTime, deliveryStartTimeErr := timeutil.DateStrToTime(form.DeliveryStartTime)

	if deliveryStartTimeErr != nil {
		return taskId, deliveryStartTimeErr
	}

	deliveryEndTime, deliveryEndTimeErr := timeutil.DateStrToTime(form.DeliveryEndTime)

	if deliveryEndTimeErr != nil {
		return taskId, deliveryEndTimeErr
	}

	payTime, payTimeErr := timeutil.DateStrToTime(form.PayTime)

	if payTimeErr != nil {
		return taskId, payTimeErr
	}

	lines := "全部线路"

	if len(form.Lines) > 0 {
		lines = strings.Join(form.Lines, ",")
	}

	task := model.OutboundTask{
		TaskName:          form.OutboundName,
		DeliveryStartTime: (*model.MyTime)(deliveryStartTime),
		DeliveryEndTime:   (*model.MyTime)(deliveryEndTime),
		Line:              lines,
		DistributionType:  form.DistributionType,
		PayEndTime:        (*model.MyTime)(payTime),
		Creator: model.Creator{
			CreatorId: userInfo.ID,
			Creator:   userInfo.Name,
		},
		Sku:       strings.Join(form.Sku, ","),
		GoodsName: strings.Join(form.GoodsName, ","),
	}

	err = model.OutboundTaskSave(db, &task)

	if err != nil {
		return taskId, err
	}

	return task.Id, nil
}

// 出库订单数量
func OutboundOrderCount(db *gorm.DB, form req.OutboundOrderCountForm) (error, rsp.OutboundOrderCountRsp) {

	err, count := model.OutboundOrderOrderTypeCount(db, form.TaskId)
	if err != nil {
		return err, rsp.OutboundOrderCountRsp{}
	}

	var res rsp.OutboundOrderCountRsp

	for _, ct := range count {
		switch ct.OrderType {
		case model.OutboundOrderTypeNew:
			res.New += ct.Count
			break
		case model.OutboundOrderTypePicking:
			res.Picking += ct.Count
			break
		case model.OutboundOrderTypeComplete:
			res.Complete += ct.Count
			break
		case model.OutboundOrderTypeClose:
			res.Close += ct.Count
			break
		}
		res.Total += ct.Count
	}

	return nil, res
}

// 出库订单相关保存
func OutboundOrderBatchSave(db *gorm.DB, form req.CreateOutboundForm, taskId int) error {

	var (
		orderJoinGoods []model.OrderJoinGoods
		mp             = make(map[string]int, 0) //空map，共用逻辑时使用，这里没什么用途
	)

	localDb := db.Table("t_order_goods og").
		Joins("left join t_order o on og.number = o.number").
		Select("og.*,og.id as order_goods_id,o.*,o.id as order_id").
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

	//订单商品中的 新订单或欠货的订单商品数据
	result := localDb.Where("og.`status` = ?", model.OrderGoodsUnhandledStatus).
		Find(&orderJoinGoods)

	if result.Error != nil {
		return result.Error
	}

	err := OutboundOrderBatchSaveLogic(db, taskId, orderJoinGoods, mp)

	if err != nil {
		return err
	}

	return nil
}

// 出库订单相关保存逻辑
func OutboundOrderBatchSaveLogic(db *gorm.DB, taskId int, orderJoinGoods []model.OrderJoinGoods, mp map[string]int) error {
	var (
		outboundOrderMp = make(map[string]model.OutboundOrder, 0)             //出库订单map 以订单号为key，用于更新订单备注以及限发总数，然后存储
		remarkMp        = make(map[string]struct{}, 0)                        //订单备注map，用于处理 订单是否有备注
		outboundOrders  = make([]model.OutboundOrder, 0)                      //出库订单
		outboundGoods   = make([]model.OutboundGoods, 0, len(orderJoinGoods)) //出库订单商品
		orderIds        []int                                                 //订单id数据，用于更新订单表数据
		orderGoodsIds   []int                                                 //订单商品id数据，用于更新订单商品表数据
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
			orderIds = append(orderIds, goods.OrderId)
		}

		if goods.OrderRemark != "" || goods.GoodsRemark != "" {
			remarkMp[goods.Number] = struct{}{}
		}

		//如果任务设置了限发，则取限发数量，否则取欠货数量
		limitNum, limitOk := mp[goods.Sku]

		if !limitOk {
			limitNum = goods.LackCount
		}

		outboundGoods = append(outboundGoods, model.OutboundGoods{
			TaskId:          taskId,
			Number:          goods.Number,
			Sku:             goods.Sku,
			OrderGoodsId:    goods.OrderGoodsId, //order_goods表ID
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
			LimitNum:        limitNum,
			GoodsRemark:     goods.GoodsRemark,
			BatchId:         0,
			DeliveryOrderNo: nil,
		})

		//更新订单商品表数据
		orderGoodsIds = append(orderGoodsIds, goods.OrderGoodsId)
	}

	for s, oo := range outboundOrderMp {

		//如果 订单号 在 remarkMp 中则为有备注
		_, ok := remarkMp[s]

		if ok {
			oo.HasRemark = 1
		} else {
			oo.HasRemark = 0
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
	err = model.UpdateOrderAndGoodsByIds(db, orderIds, orderGoodsIds, model.PickingOrderType, model.OrderGoodsProcessingStatus)

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
	})

	if form.OrderType > 0 {
		localDb.Where("order_type = ?", form.OrderType)
	} else {
		//全部订单不显示已关闭的
		localDb.Where("order_type != ?", model.OutboundOrderTypeClose)
	}

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

	//筛选没有number时
	if len(numbers) == 0 {
		for _, order := range outboundOrders {
			numbers = append(numbers, order.Number)
		}
	}

	err, numsMp := model.OutboundGoodsNumsStatisticalByTaskIdAndNumbers(db, form.TaskId, numbers)
	if err != nil {
		return err, res
	}

	list := make([]rsp.OutboundOrderList, 0, len(outboundOrders))

	for _, order := range outboundOrders {

		nums, numsOk := numsMp[order.Number]

		if !numsOk {
			return errors.New("订单统计数量不存在"), res
		}

		list = append(list, rsp.OutboundOrderList{
			Number:            order.Number,
			OutboundNumber:    model.GetOutboundNumber(order.TaskId, order.Number),
			PayAt:             order.PayAt,
			ShopName:          order.ShopName,
			ShopType:          order.ShopType,
			DistributionType:  order.DistributionType,
			GoodsNum:          nums.PayCount,
			LimitNum:          nums.LimitNum,
			CloseNum:          nums.CloseCount,
			OutCount:          nums.OutCount,
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

// 出库单订单详情
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

// 出库任务商品列表
func OutboundOrderGoodsList(db *gorm.DB, form req.OutboundOrderGoodsListForm) (error, []rsp.OutboundOrderGoodsList) {
	err, goodsList := model.OutboundOrderDistinctGoodsList(db, form.TaskId)
	if err != nil {
		return err, nil
	}

	list := make([]rsp.OutboundOrderGoodsList, 0, len(goodsList))

	for _, goods := range goodsList {
		list = append(list, rsp.OutboundOrderGoodsList{
			Sku:  goods.Sku,
			Name: goods.GoodsName,
		})
	}
	return nil, list
}

// 结束任务
func EndOutboundTask(db *gorm.DB, form req.EndOutboundTaskForm) error {

	err, batchList := model.GetBatchListByTaskId(db, form.TaskId)

	if err != nil {
		return err
	}

	for _, batch := range batchList {
		if batch.Status != model.BatchClosedStatus {
			return errors.New("批次任务:" + batch.BatchName + "没有结束")
		}
	}

	tx := db.Begin()

	err = model.UpdateOutboundTaskStatusById(tx, form.TaskId)
	if err != nil {
		return err
	}

	//todo 更新订单数据

	tx.Commit()

	return nil
}

// 关闭订单
func OutboundTaskCloseOrder(db *gorm.DB, form req.OutboundTaskCloseOrderForm) (err error) {

	//关闭预拣池订单
	err = ClosePrePickOrder(db, form)
	if err != nil {
		return err
	}

	//关闭出库任务订单
	err = CloseOutboundOrder(db, form)

	if err != nil {
		return err
	}

	//订单状态变更
	err = model.UpdateOrderAndGoodsByNumbers(db, form.Number, model.LackOrderType, model.OrderGoodsUnhandledStatus)

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

// 关闭出库任务订单
func CloseOutboundOrder(db *gorm.DB, form req.OutboundTaskCloseOrderForm) error {
	//出库任务订单
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

// 临时加单
func OutboundTaskAddOrder(db *gorm.DB, form req.OutboundTaskAddOrderForm) error {

	var (
		limitShipmentMp = make(map[string]int, 0)
		orderList       []model.OrderJoinGoods
	)

	//限发
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

	//订单信息
	err, orderList = model.GetOrderJoinGoodsList(db, form.Number)

	if err != nil {
		return err
	}

	err = OutboundOrderBatchSaveLogic(db, form.TaskId, orderList, limitShipmentMp)

	if err != nil {
		return err
	}

	return nil
}

// 订单出库记录
func OrderOutboundRecord(db *gorm.DB, form req.OrderOutboundRecordForm) (err error, list rsp.OrderOutboundRecordList) {
	return
}

// 获取任务某个商品的发货数量
func GetTaskSkuNum(db *gorm.DB, form req.GetTaskSkuNumForm) (err error, num int) {
	err, num = model.GetOutboundGoods(db, form.TaskId, form.Sku)

	return
}

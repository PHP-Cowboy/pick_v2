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

// 出库单任务列表
func OutboundTaskList(db *gorm.DB, form req.OutboundTaskListForm) (err error, res rsp.OutboundTaskListRsp) {

	var taskIds []int

	if form.ShopId > 0 || form.Number != "" || form.Sku != "" {
		var outboundOrderAndGoods []model.GoodsJoinOrder

		orderDb := db.Table("t_outbound_goods og").
			Select("og.task_id").
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

	result = localDb.Scopes(model.Paginate(form.Page, form.Size)).Order("id desc").Find(&outboundTask)

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

	err = model.OutboundTaskCreate(db, &task)

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
		orderJoinGoods []model.GoodsJoinOrder
		mp             = make(map[string]int, 0) //空map，共用逻辑时使用，这里没什么用途
	)

	localDb := db.Table("t_order_goods og").
		Joins("left join t_order o on og.number = o.number").
		Select("og.*,og.id as order_goods_id,o.*,o.id as order_id").
		Where("o.order_type in (?) and o.distribution_type = ? and o.pay_at <= ? and o.delivery_at <= ? ",
			[]int{model.NewOrderType, model.LackOrderType},
			form.DistributionType,
			form.PayTime,
			form.DeliveryEndTime,
		)

	if form.DeliveryStartTime != "" {
		localDb = localDb.Where("o.delivery_at >= ?", form.DeliveryStartTime)
	}

	if len(form.Sku) > 0 {
		localDb = localDb.Where("og.sku in (?)", form.Sku)
	}

	if len(form.Lines) > 0 {
		localDb = localDb.Where("o.line in (?) ", form.Lines)
	}

	if len(form.GoodsType) > 0 {
		localDb = localDb.Where("goods_type in (?) ", form.GoodsType)
	}

	if len(form.ShopIds) > 0 {
		localDb = localDb.Where("shop_id in (?) ", form.ShopIds)
	}

	//订单商品中的 新订单或欠货的订单商品数据
	result := localDb.Where("og.`status` = ?", model.OrderGoodsUnhandledStatus).
		Find(&orderJoinGoods)

	if result.Error != nil {
		return result.Error
	}

	if len(orderJoinGoods) == 0 {
		return errors.New("暂无相关订单")
	}

	err := OutboundOrderBatchSaveLogic(db, taskId, orderJoinGoods, mp)

	if err != nil {
		return err
	}

	return nil
}

// 出库订单相关保存逻辑
func OutboundOrderBatchSaveLogic(db *gorm.DB, taskId int, orderJoinGoods []model.GoodsJoinOrder, mp map[string]int) error {
	var (
		outboundOrderMp = make(map[string]model.OutboundOrder, 0)             //出库订单map 以订单号为key，用于更新订单备注以及限发总数，然后存储
		remarkMp        = make(map[string]struct{}, 0)                        //订单备注map，用于处理 订单是否有备注
		outboundOrders  = make([]model.OutboundOrder, 0)                      //出库订单
		outboundGoods   = make([]model.OutboundGoods, 0, len(orderJoinGoods)) //出库订单商品
		orderIds        []int                                                 //订单id数据，用于更新订单表数据
		orderGoodsIds   []int                                                 //订单商品id数据，用于更新订单商品表数据
	)

	for _, goods := range orderJoinGoods {
		//欠货数量小于0的不进入任务中，订单下无商品的也不会创建订单
		if goods.LackCount <= 0 {
			continue
		}

		_, ok := outboundOrderMp[goods.Number]

		if !ok {
			outboundOrderMp[goods.Number] = model.OutboundOrder{
				TaskId:            taskId,
				Number:            goods.Number,
				PayAt:             goods.PayAt,
				OrderId:           goods.OrderId,
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
				HasRemark:         goods.HasRemark,
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
			OutCount:        0,
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
			oo.HasRemark = 2
		} else {
			oo.HasRemark = 1
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

	if form.Sku != "" || form.GoodsType != "" {
		var outboundGoods []model.OutboundGoods
		goodsRes := db.Model(&model.OutboundGoods{}).
			Select("number").
			Where(&model.OutboundGoods{TaskId: form.TaskId, Sku: form.Sku, GoodsType: form.GoodsType}).
			Find(&outboundGoods)

		if goodsRes.Error != nil {
			return goodsRes.Error, res
		}

		for _, good := range outboundGoods {
			numbers = append(numbers, good.Number)
		}

		if len(numbers) == 0 {
			res.List = make([]rsp.OutboundOrderList, 0)
			return
		}
	}

	if form.Number != "" {
		numbers = append(numbers, form.Number)
	}

	localDb := db.Model(&model.OutboundOrder{}).Where("task_id = ?", form.TaskId)

	if len(numbers) > 0 {
		numbers = slice.UniqueSlice(numbers)
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
		HasRemark:        form.HasRemark,
	})

	if form.OrderType > 0 {
		localDb.Where("order_type = ?", form.OrderType)
	} else {
		//全部订单不显示已关闭的
		localDb.Where("order_type != ?", model.OutboundOrderTypeClose)
	}

	result := localDb.Find(&outboundOrders)

	if result.Error != nil {
		return result.Error, res
	}

	res.Total = result.RowsAffected

	result = localDb.Order("shop_code ASC").Scopes(model.Paginate(form.Page, form.Size)).Find(&outboundOrders)

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
			ShopCode:          order.ShopCode,
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

	classMp, err := cache.GetClassification()

	if err != nil {
		return err, res
	}

	detailMap := make(map[string]*rsp.Detail, 0)

	deliveryOrderNoArr := make(model.GormList, 0)

	for _, goods := range outboundGoods {
		goodsType, ok := classMp[goods.GoodsType]

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
	res.PayAt = outboundOrder.PayAt
	res.ShopCode = outboundOrder.ShopCode
	res.ShopName = outboundOrder.ShopName
	res.Line = outboundOrder.Line
	res.Region = outboundOrder.Province + outboundOrder.City + outboundOrder.District
	res.ShopType = outboundOrder.ShopType
	res.OrderRemark = outboundOrder.OrderRemark

	res.Detail = detailMap

	deliveryOrderNoArr = slice.UniqueSlice(deliveryOrderNoArr)
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
func EndOutboundTask(db *gorm.DB, form req.EndOutboundTaskForm) (err error) {

	var batchList []model.Batch

	//根据出库任务获取批次列表
	err, batchList = model.GetBatchListByTaskId(db, form.TaskId)

	if err != nil {
		return
	}

	for _, batch := range batchList {
		if batch.Status != model.BatchClosedStatus {
			err = errors.New("批次任务:" + batch.BatchName + "没有结束")
			return
		}
	}

	tx := db.Begin()

	err = model.UpdateOutboundTaskStatusById(tx, form.TaskId, model.OutboundTaskStatusClosed)
	if err != nil {
		tx.Rollback()
		return
	}

	//更新订单数据
	err = EndOutboundTaskUpdateOrder(tx, form.TaskId)
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return
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

	prePickIds = slice.UniqueSlice(prePickIds)

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

	err = model.OutboundGoodsReplaceSave(db, &outboundGoods, []string{"status"})

	if err != nil {
		return err
	}

	return nil
}

// 临时加单
func OutboundTaskAddOrder(db *gorm.DB, form req.OutboundTaskAddOrderForm) error {

	var (
		limitShipmentMp = make(map[string]int, 0)
		orderList       []model.GoodsJoinOrder
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
	err, orderList = model.GetOrderJoinGoodsListByNumbers(db, form.Number)

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

// 结束批次更新订单，出库任务相关
func EndOutboundTaskUpdateOrder(tx *gorm.DB, taskId int) (err error) {

	var (
		outboundGoodsJoinOrder []model.GoodsJoinOrder
		outboundNumbers        []string //出库任务订单
		orderGoods             []model.OrderGoods
	)

	//查询任务全部订单&&商品;[新订单和已完成]
	err, outboundGoodsJoinOrder = model.GetOutboundGoodsJoinOrderListByTaskId(tx, taskId)

	//订单最近拣货时间map
	orderPickTimeMp := make(map[string]*model.MyTime, 0)
	//出库任务订单商品表id
	orderGoodsIds := make([]int, 0, len(outboundGoodsJoinOrder))

	//欠货订单编号map
	lackNumbersMap := make(map[string]struct{}, 0)

	for _, o := range outboundGoodsJoinOrder {
		//订单最近拣货时间map
		orderPickTimeMp[o.Number] = o.LatestPickingTime
		//订单商品
		orderGoodsIds = append(orderGoodsIds, o.OrderGoodsId)

		//商品还有欠货，即订单为欠货单
		if o.LackCount > 0 {
			lackNumbersMap[o.Number] = struct{}{}
		}

		outboundNumbers = append(outboundNumbers, o.Number)
	}

	outboundNumbers = slice.UniqueSlice(outboundNumbers)

	//订单历史出库数据map
	orderGoodsHistoryMp := make(map[int]model.HistoryOrderGoods, 0)

	err, orderGoods = model.GetOrderGoodsListByIds(tx, orderGoodsIds)

	if err != nil {
		return
	}

	//获取订单商品历史出库单数据
	for _, good := range orderGoods {
		orderGoodsHistoryMp[good.Id] = model.HistoryOrderGoods{
			DeliveryNumber: good.DeliveryOrderNo,
			OutCount:       good.OutCount,
		}
	}

	//订单更新处理
	err = OrderUpdateHandle(tx, outboundGoodsJoinOrder, lackNumbersMap, orderGoodsHistoryMp, orderPickTimeMp, nil)

	if err != nil {
		return
	}

	//更新出库任务订单为已完成
	err = model.UpdateOutboundOrderByTaskIdAndNumbers(tx, taskId, outboundNumbers, map[string]interface{}{"order_type": model.OutboundOrderTypeComplete})
	if err != nil {
		return
	}

	return
}

// 订单更新处理
func OrderUpdateHandle(
	tx *gorm.DB,
	goodsJoinOrders []model.GoodsJoinOrder,
	lackNumbersMap map[string]struct{},
	orderGoodsHistoryMp map[int]model.HistoryOrderGoods,
	orderPickTimeMp map[string]*model.MyTime,
	deliveryOrderNo model.GormList,
) (
	err error,
) {
	var (
		lackOrder             []model.Order                  //欠货订单
		lackGoods             []model.OrderGoods             //欠货订单商品出库数据更新
		completeOrder         []model.CompleteOrder          //完成订单
		completeOrderDetail   []model.CompleteOrderDetail    //完成订单商品数据
		completeIds           []int                          //完成订单id，删除订单使用
		completeOrderGoodsIds []int                          //完成订单商品id，删订单商品数据使用
		completeOrderMp       = make(map[string]struct{}, 0) //完成订单编号去重
		lackOrderMp           = make(map[int]struct{}, 0)    //欠货订单ID去重
	)

	var (
		orderPickTimeMpOk bool
		pickTime          *model.MyTime
	)

	for _, goodsJoinOrder := range goodsJoinOrders {

		//最近拣货时间 [确认出库时orderPickTimeMp传nil]
		// 结束批次时，最近拣货时间取的任务中订单表的最近拣货时间
		// 确认出库时，订单最近拣货时间为当前时间
		if orderPickTimeMp != nil {
			//结束批次时
			//最近拣货时间
			pickTime, orderPickTimeMpOk = orderPickTimeMp[goodsJoinOrder.Number]

			if !orderPickTimeMpOk {
				pickTime = nil
			}
		} else {
			//确认出库时
			now := time.Now()
			//最近拣货时间
			pickTime = (*model.MyTime)(&now)
		}

		//出库单号
		deliveryOrderNoArr := make(model.GormList, 0)
		// 最近拣货时间 [确认出库时orderGoodsHistoryMp传nil]
		// 结束批次时，读取的是订单商品表数据，这里可以取到历史出库单号
		// 确认出库时，没有去读取历史出库单号，只是生成了本次的出库单号
		if orderGoodsHistoryMp != nil {
			//结束批次时
			//历史出库订单
			orderGoodsHistory, orderGoodsHistoryMpOk := orderGoodsHistoryMp[goodsJoinOrder.OrderGoodsId]

			historyDeliveryOrderNoArr := []string{}

			if orderGoodsHistoryMpOk {
				historyDeliveryOrderNoArr = orderGoodsHistory.DeliveryNumber
			}

			deliveryOrderNoArr = append(deliveryOrderNoArr, goodsJoinOrder.DeliveryOrderNo...)
			deliveryOrderNoArr = append(deliveryOrderNoArr, historyDeliveryOrderNoArr...)

			/*
				出库任务中订单数据只有本次出库数量
				历史出库数量在订单商品表数据中，这里先给他加进来
				欠货数量不用处理，订单第一次进入任务之后再进入时，欠货数量已经是最新欠货数量了
			*/

			//订单商品出货数量 = 历史出货数量 + 本次出货数量
			goodsJoinOrder.OutCount += orderGoodsHistory.OutCount
		} else {
			//确认出库时
			deliveryOrderNoArr = append(deliveryOrderNoArr, goodsJoinOrder.DeliveryOrderNo...)
			deliveryOrderNoArr = append(deliveryOrderNoArr, deliveryOrderNo...)
		}

		//是否欠货单
		_, lackNumbersMapOk := lackNumbersMap[goodsJoinOrder.Number]

		if lackNumbersMapOk {

			//欠货订单商品相关更新
			lackGoods = append(lackGoods, model.OrderGoods{
				Id:              goodsJoinOrder.OrderGoodsId,
				Number:          goodsJoinOrder.Number,
				LackCount:       goodsJoinOrder.LackCount,
				OutCount:        goodsJoinOrder.OutCount,
				Status:          model.OrderGoodsUnhandledStatus, //欠货，更新成未处理
				DeliveryOrderNo: deliveryOrderNoArr,
			})

			_, lackOrderMpOk := lackOrderMp[goodsJoinOrder.OrderId]

			if lackOrderMpOk {
				continue
			}

			//欠货订单
			lackOrder = append(lackOrder, model.Order{
				Id:                goodsJoinOrder.OrderId,
				ShopId:            goodsJoinOrder.ShopId,
				ShopName:          goodsJoinOrder.ShopName,
				ShopType:          goodsJoinOrder.ShopType,
				ShopCode:          goodsJoinOrder.ShopCode,
				Number:            goodsJoinOrder.Number,
				HouseCode:         goodsJoinOrder.HouseCode,
				Line:              goodsJoinOrder.Line,
				OrderType:         model.LackOrderType,
				LatestPickingTime: pickTime,
			})

			lackOrderMp[goodsJoinOrder.OrderId] = struct{}{}
		} else {

			//完成订单详情
			completeOrderDetail = append(completeOrderDetail, model.CompleteOrderDetail{
				OrderGoodsId:    goodsJoinOrder.OrderGoodsId,
				Number:          goodsJoinOrder.Number,
				GoodsName:       goodsJoinOrder.GoodsName,
				Sku:             goodsJoinOrder.Sku,
				GoodsSpe:        goodsJoinOrder.GoodsSpe,
				GoodsType:       goodsJoinOrder.GoodsType,
				Shelves:         goodsJoinOrder.Shelves,
				PayCount:        goodsJoinOrder.PayCount,
				CloseCount:      goodsJoinOrder.CloseCount,
				ReviewCount:     goodsJoinOrder.OutCount,
				GoodsRemark:     goodsJoinOrder.GoodsRemark,
				DiscountPrice:   goodsJoinOrder.DiscountPrice,
				DeliveryOrderNo: deliveryOrderNoArr,
			})

			//完成订单商品id，删除完成订单商品数据
			completeOrderGoodsIds = append(completeOrderGoodsIds, goodsJoinOrder.OrderGoodsId)

			_, completeOrderMpOk := completeOrderMp[goodsJoinOrder.Number]

			if completeOrderMpOk {
				continue
			}

			//完成订单id，删除订单表完成订单
			completeIds = append(completeIds, goodsJoinOrder.OrderId)

			//完成订单
			completeOrder = append(completeOrder, model.CompleteOrder{
				Number:         goodsJoinOrder.Number,
				OrderRemark:    goodsJoinOrder.OrderRemark,
				ShopId:         goodsJoinOrder.ShopId,
				ShopName:       goodsJoinOrder.ShopName,
				ShopType:       goodsJoinOrder.ShopType,
				ShopCode:       goodsJoinOrder.ShopCode,
				Line:           goodsJoinOrder.Line,
				DeliveryMethod: goodsJoinOrder.DistributionType,
				HouseCode:      goodsJoinOrder.HouseCode,
				Province:       goodsJoinOrder.Province,
				City:           goodsJoinOrder.City,
				District:       goodsJoinOrder.District,
				PickTime:       pickTime,
				PayAt:          goodsJoinOrder.PayAt,
			})

			completeOrderMp[goodsJoinOrder.Number] = struct{}{}
		}
	}

	//欠货订单更新字段
	lackOrderValues := []string{"shop_id", "shop_name", "shop_type", "shop_code", "house_code", "line", "order_type", "latest_picking_time"}

	//欠货商品更新字段
	lackGoodsValues := []string{"lack_count", "out_count", "delivery_order_no", "status"}

	//结束任务更新订单数据
	err = UpdateOrders(tx, &lackOrder, &lackGoods, lackOrderValues, lackGoodsValues, &completeOrder, &completeOrderDetail, completeIds, completeOrderGoodsIds)

	if err != nil {
		return
	}

	return
}

// 结束任务或复核完成更新订单数据
func UpdateOrders(
	tx *gorm.DB,
	lackOrder *[]model.Order,
	lackGoods *[]model.OrderGoods,
	lackOrderValues,
	lackGoodsValues []string,
	completeOrder *[]model.CompleteOrder,
	completeOrderDetail *[]model.CompleteOrderDetail,
	completeIds,
	completeOrderGoodsIds []int,
) (err error) {
	// 在model里判断了是否为空才更新

	//更新订单欠货数据
	err = model.OrderReplaceSave(tx, lackOrder, lackOrderValues)

	if err != nil {
		return
	}

	//更新订单商品数据
	err = model.OrderGoodsReplaceSave(tx, lackGoods, lackGoodsValues)
	if err != nil {
		return
	}

	//完成订单保存
	err = model.CompleteOrderBatchSave(tx, completeOrder)

	if err != nil {
		return
	}

	//完成订单商品数据保存
	err = model.CompleteOrderDetailBatchSave(tx, completeOrderDetail)

	if err != nil {
		return
	}

	//删除订单表已完成数据
	err = model.DeleteOrderByIds(tx, completeIds)
	if err != nil {
		return
	}

	//删除订单商品表数据
	err = model.DeleteOrderGoodsByIds(tx, completeOrderGoodsIds)

	if err != nil {
		return
	}

	return
}

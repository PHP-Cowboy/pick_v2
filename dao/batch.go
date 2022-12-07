package dao

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"gorm.io/gorm"

	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/helper"
	"pick_v2/utils/slice"
	"pick_v2/utils/timeutil"
)

// 全量拣货 -按任务创建批次
func CreateBatchByTask(db *gorm.DB, form req.CreateBatchByTaskForm, claims *middlewares.CustomClaims) (err error) {

	//根据任务ID查询出库任务订单表数据，获取任务的全部订单号
	err, outboundOrderList := model.GetOutboundOrderByTaskId(db, form.TaskId)
	if err != nil {
		return err
	}

	//构造订单号数据
	numbers := make([]string, 0, len(outboundOrderList))

	//根据订单数据获取订单号
	for _, order := range outboundOrderList {
		numbers = append(numbers, order.Number)
	}

	if len(numbers) == 0 {
		err = errors.New("任务订单已全部生成批次")
		return
	}

	newCreateBatchForm := req.NewCreateBatchForm{
		TaskId:    form.TaskId,
		Number:    numbers,
		BatchName: form.TaskName,
		Typ:       form.Typ,
	}

	err = CreateBatch(db, newCreateBatchForm, claims)
	if err != nil {
		return err
	}

	return
}

// 创建批次
func CreateBatch(db *gorm.DB, form req.NewCreateBatchForm, claims *middlewares.CustomClaims) error {

	var (
		orderGoodsIds []int
		outboundGoods []model.OutboundGoods
	)

	tx := db.Begin()

	//批次
	err, batch := BatchSaveLogic(tx, form, claims)

	if err != nil {
		tx.Rollback()
		return err
	}

	//预拣池逻辑
	err, orderGoodsIds, outboundGoods, _, _, _, _, _ = CreatePrePickLogic(tx, form, claims, batch.Id)

	if err != nil {
		tx.Rollback()
		return err
	}

	//批量更新 t_order_goods 表 batch_id
	err = model.UpdateOrderGoodsByIds(tx, orderGoodsIds, map[string]interface{}{"batch_id": batch.Id})

	if err != nil {
		tx.Rollback()
		return err
	}

	//构造更新 t_outbound_order 表 order_type 数据
	outboundOrder := make([]model.OutboundOrder, 0, len(form.Number))

	//根据传递过来的任务id和订单编号生成需要被变更order_type状态为拣货中的数据
	for _, number := range form.Number {
		outboundOrder = append(outboundOrder, model.OutboundOrder{
			TaskId:    form.TaskId,
			Number:    number,
			OrderType: model.OutboundOrderTypePicking,
		})
	}

	//批量更新 t_outbound_order 表 order_type 状态为拣货中
	err = model.OutboundOrderReplaceSave(db, outboundOrder, []string{"order_type"})

	if err != nil {
		return err
	}

	//批量更新 t_outbound_goods 状态 为拣货中
	if err = model.OutboundGoodsReplaceSave(tx, outboundGoods, []string{"status"}); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// 生成批次数据逻辑
func BatchSaveLogic(db *gorm.DB, form req.NewCreateBatchForm, claims *middlewares.CustomClaims) (err error, batch model.Batch) {

	var (
		outboundTask model.OutboundTask
	)
	//获取出库任务信息
	err, outboundTask = model.GetOutboundTaskById(db, form.TaskId)

	if err != nil {
		return
	}

	batch = model.Batch{
		TaskId:            form.TaskId,
		WarehouseId:       claims.WarehouseId,
		BatchName:         form.BatchName,
		DeliveryStartTime: outboundTask.DeliveryStartTime,
		DeliveryEndTime:   outboundTask.DeliveryEndTime,
		UserName:          claims.Name,
		Line:              outboundTask.Line,
		DeliveryMethod:    outboundTask.DistributionType,
		EndTime:           nil,
		Status:            0,
		Sort:              0,
		PayEndTime:        outboundTask.PayEndTime,
		Version:           0,
		Typ:               form.Typ,
	}

	//t_batch
	err = model.BatchSave(db, &batch)

	if err != nil {
		return
	}

	return
}

// 快递批次
func CourierBatch(db *gorm.DB, form req.NewCreateBatchForm, claims *middlewares.CustomClaims) error {

	var (
		orderGoodsIds          []int
		outboundGoods          []model.OutboundGoods
		outboundGoodsJoinOrder []model.OutboundGoodsJoinOrder
		prePickIds             []int
		prePicks               []model.PrePick
		prePickGoods           []model.PrePickGoods
		prePickRemarks         []model.PrePickRemark
	)

	tx := db.Begin()

	//生成批次
	err, batch := BatchSaveLogic(tx, form, claims)

	if err != nil {
		tx.Rollback()
		return err
	}

	//生成预拣池
	//todo 在快递批次时直接把状态设置成已进入拣货池？
	//todo 拣货池逻辑中就可以不修改状态了，但是后续是否会快递批次被改成先集中拣货完成再到二次分拣？
	//TODO 个人觉得集中拣货和二次分拣同时进行在实际业务中是有问题的
	err, orderGoodsIds, outboundGoods, outboundGoodsJoinOrder, prePickIds, prePicks, prePickGoods, prePickRemarks = CreatePrePickLogic(tx, form, claims, batch.Id)

	if err != nil {
		tx.Rollback()
		return err
	}

	//批量更新 t_order_goods batch_id
	err = model.UpdateOrderGoodsByIds(tx, orderGoodsIds, map[string]interface{}{"batch_id": batch.Id})

	if err != nil {
		tx.Rollback()
		return err
	}

	//构造更新 t_outbound_order 表 order_type 数据
	outboundOrder := make([]model.OutboundOrder, 0, len(form.Number))

	//根据传递过来的任务id和订单编号生成需要被变更order_type状态为拣货中的数据
	for _, number := range form.Number {
		outboundOrder = append(outboundOrder, model.OutboundOrder{
			TaskId:    form.TaskId,
			Number:    number,
			OrderType: model.OutboundOrderTypePicking,
		})
	}

	//批量更新 t_outbound_order 表 order_type 状态为拣货中
	err = model.OutboundOrderReplaceSave(db, outboundOrder, []string{"order_type"})

	if err != nil {
		return err
	}

	//批量更新 t_outbound_goods 状态
	if err = model.OutboundGoodsReplaceSave(tx, outboundGoods, []string{"status"}); err != nil {
		tx.Rollback()
		return err
	}

	//生成集中拣货
	err = CreateCentralizedPick(db, outboundGoodsJoinOrder, batch.Id)
	if err != nil {
		return err
	}

	pick := req.BatchPickForm{
		BatchId:     batch.Id,
		Ids:         prePickIds,
		Type:        1,
		TypeParam:   []string{},
		WarehouseId: claims.WarehouseId,
		Typ:         2,
	}

	//生成拣货池
	err = BatchPickByParams(db, pick, prePicks, prePickGoods, prePickRemarks)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// 根据sku,number,shop_id查批次id
func GetBatchIdsFromPrePickGoods(db *gorm.DB, form req.GetBatchListForm) (err error, batchIds []int) {
	var prePickGoods []model.PrePickGoods

	err = db.Model(&model.PrePickGoods{}).
		Where(model.PrePickGoods{
			Sku:    form.Sku,
			Number: form.Number,
			ShopId: form.ShopId,
		}).
		Select("batch_id").
		Find(&prePickGoods).
		Error

	if err != nil {
		return
	}

	//未找到，直接返回
	if len(prePickGoods) == 0 {
		return
	}

	//batchIds
	for _, b := range prePickGoods {
		batchIds = append(batchIds, b.BatchId)
	}

	batchIds = slice.UniqueSlice(batchIds)

	return
}

// 批次列表
func BatchList(db *gorm.DB, form req.GetBatchListForm) (err error, res rsp.GetBatchListRsp) {
	var (
		batches  []model.Batch
		batchIds []int
	)

	local := db.Model(&model.Batch{})

	//子表数据
	if form.Sku != "" || form.Number != "" || form.ShopId > 0 {
		err, batchIds = GetBatchIdsFromPrePickGoods(db, form)
		if err != nil {
			return
		}
		local.Where("id in (?)", batchIds)
	}

	if form.Line != "" {
		local.Where("line like ?", form.Line+"%")
	}

	if form.CreateTime != "" {
		local.Where("create_time <= ?", form.CreateTime)
	}

	if form.EndTime != "" {
		local.Where("end_time <= ?", form.EndTime)
	}

	if *form.Status == 0 {
		local.Where("status in (?)", []int{model.BatchOngoingStatus, model.BatchSuspendStatus})
	} else {
		local.Where("status = ?", model.BatchClosedStatus)
	}

	local.Where(&model.Batch{DeliveryMethod: form.DeliveryMethod, Typ: form.Typ})

	result := local.Find(&batches)

	err = result.Error

	if err != nil {
		return
	}

	res.Total = result.RowsAffected

	result = local.Scopes(model.Paginate(form.Page, form.Size)).
		Order("sort desc, id desc").
		Find(&batches)

	err = result.Error

	if err != nil {
		return
	}

	//batchIds
	if len(batchIds) == 0 {
		for _, batch := range batches {
			batchIds = append(batchIds, batch.Id)
		}
	}

	// 统计批次 门店数量、预拣单数等
	err, numsMp := model.CountBatchNumsByBatchIds(db, batchIds)
	if err != nil {
		return
	}

	list := make([]*rsp.Batch, 0, len(batches))
	for _, b := range batches {

		nums, numsOk := numsMp[b.Id]

		if !numsOk {
			err = errors.New("拣货池数据异常")
			return
		}

		list = append(list, &rsp.Batch{
			Id:                b.Id,
			CreateTime:        b.CreateTime.Format(timeutil.MinuteFormat),
			UpdateTime:        b.UpdateTime.Format(timeutil.MinuteFormat),
			BatchName:         b.BatchName + helper.GetDeliveryMethod(b.DeliveryMethod),
			DeliveryStartTime: b.DeliveryStartTime,
			DeliveryEndTime:   b.DeliveryEndTime,
			ShopNum:           nums.ShopNum,
			OrderNum:          nums.OrderNum,
			GoodsNum:          nums.GoodsNum,
			UserName:          b.UserName,
			Line:              b.Line,
			DeliveryMethod:    b.DeliveryMethod,
			EndTime:           b.EndTime,
			Status:            b.Status,
			PrePickNum:        nums.PrePickNum,
			PickNum:           nums.PickNum,
			RecheckSheetNum:   nums.ReviewNum,
		})
	}

	res.List = list

	return
}

// 结束批次
func EndBatch(db *gorm.DB, form req.EndBatchForm) (err error) {
	var (
		batch          model.Batch
		pickGoods      []model.PickGoods
		orderJoinGoods []model.OrderJoinGoods
	)

	err, batch = model.GetBatchByPk(db, form.Id)

	if err != nil {
		return
	}

	if batch.Status != model.BatchSuspendStatus {
		err = errors.New("请先停止拣货")
		return
	}

	tx := db.Begin()

	//修改批次状态为已结束
	err = model.UpdateBatchByPk(tx, batch.Id, map[string]interface{}{"status": model.BatchClosedStatus})

	if err != nil {
		tx.Rollback()
		return
	}

	// 根据批次id查询订单&&订单商品数据
	err, orderJoinGoods = model.GetOrderGoodsJoinOrderByBatchId(db, batch.Id)

	if err != nil {
		tx.Rollback()
		return
	}

	//查询批次下全部订单
	err, pickGoods = model.GetPickGoodsByBatchId(db, batch.Id)

	if err != nil {
		tx.Rollback()
		err = errors.New("批次结束成功，但推送u8拣货数据查询失败:" + err.Error())
		return
	}

	err = YongYouLog(tx, pickGoods, orderJoinGoods, batch.Id)

	if err != nil {
		tx.Rollback()
		return
	}

	var orderNumbers = make([]string, 0, len(orderJoinGoods))

	for _, good := range pickGoods {
		orderNumbers = append(orderNumbers, good.Number)
	}

	orderNumbers = slice.UniqueSlice(orderNumbers)

	//更新 OutboundOrder 为已完成
	err = model.UpdateOutboundOrderByTaskIdAndNumbers(db, batch.TaskId, orderNumbers, map[string]interface{}{"order_type": model.OutboundOrderTypeComplete})

	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return
}

// 批次出库订单和商品明细
func GetBatchOrderAndGoods(db *gorm.DB, form req.GetBatchOrderAndGoodsForm) (err error, res rsp.GetBatchOrderAndGoodsRsp) {
	var (
		batch         model.Batch
		outboundOrder []model.OutboundOrder
		outboundGoods []model.OutboundGoods
		mp            = make(map[string][]rsp.OutGoods)
		numbers       = make([]string, 0)
	)

	err, batch = model.GetBatchByPk(db, form.Id)

	if err != nil {
		return
	}

	//状态:0:进行中,1:已结束,2:暂停
	if batch.Status != 1 {
		err = errors.New("批次未结束")
		return
	}

	err, outboundGoods = model.GetOutboundGoodsList(db, model.OutboundGoods{BatchId: batch.Id})
	if err != nil {
		return
	}

	totalGoodsNum := 0

	for _, good := range outboundGoods {
		//出库为0的不推送
		if good.OutCount == 0 {
			continue
		}

		totalGoodsNum++

		//编号 ，查询订单
		numbers = append(numbers, good.Number)

		_, ok := mp[good.Number]

		if !ok {
			mp[good.Number] = make([]rsp.OutGoods, 0)
		}

		mp[good.Number] = append(mp[good.Number], rsp.OutGoods{
			Id:            good.OrderGoodsId,
			Name:          good.GoodsName,
			Sku:           good.Sku,
			GoodsType:     good.GoodsType,
			GoodsSpe:      good.GoodsSpe,
			DiscountPrice: good.DiscountPrice,
			GoodsUnit:     good.GoodsUnit,
			SaleUnit:      good.SaleUnit,
			SaleCode:      good.SaleCode,
			OutCount:      good.OutCount,
			OutAt:         good.UpdateTime.Format(timeutil.TimeFormat),
			Number:        good.Number,
			CkNumber:      strings.Join(good.DeliveryOrderNo, ","),
		})

	}

	numbers = slice.UniqueSlice(numbers)

	err, outboundOrder = model.GetOutboundOrderByTaskIdAndNumbers(db, batch.TaskId, numbers)

	if err != nil {
		return
	}

	list := make([]rsp.OutOrder, 0, len(outboundOrder))

	for _, order := range outboundOrder {
		goodsInfo, ok := mp[order.Number]

		if !ok {
			err = ecode.DataQueryError
			return
		}

		list = append(list, rsp.OutOrder{
			DistributionType: order.DistributionType,
			PayAt:            *order.PayAt,
			OrderId:          order.OrderId,
			GoodsInfo:        goodsInfo,
		})
	}

	res.Count = totalGoodsNum

	res.List = list
	return
}

// 推送批次信息到消息队列
func SyncBatch(batchId int) error {
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{global.ServerConfig.RocketMQ})),
		producer.WithRetry(2),
	)

	err := p.Start()

	if err != nil {
		global.Logger["err"].Infof("start producer error: %s", err.Error())
		return err
	}

	topic := "pick_batch"

	msg := &primitive.Message{
		Topic: topic,
		Body:  []byte(strconv.Itoa(batchId)),
	}

	res, err := p.SendSync(context.Background(), msg)

	if err != nil {
		global.Logger["err"].Infof("send message error: %s", err.Error())
		return err
	} else {
		global.Logger["info"].Infof("send message success: result=%s", res.String())
	}

	err = p.Shutdown()

	if err != nil {
		global.Logger["err"].Infof("shutdown producer error: %s", err.Error())
		return err
	}

	return nil
}

// 结束批次更新订单，出库任务相关
func UpdateCompleteOrder(db *gorm.DB, batchId, taskId int, orderJoinGoods []model.OrderJoinGoods) (err error) {

	var (
		orderGoods []model.OrderGoods
		isSendMQ   = true
	)

	var (
		//完成订单map
		//key => number
		completeMp          = make(map[string]interface{}, 0)
		completeNumbers     []string
		deleteIds           []int    //待删除订单表id
		deleteNumbers       []string //删除订单商品表
		completeOrder       = make([]model.CompleteOrder, 0)
		completeOrderDetail = make([]model.CompleteOrderDetail, 0)
		lackNumbers         []string //待更新为欠货订单表number
		numsMp              = make(map[string]model.OrderGoodsNumsStatistical, 0)
		orderNumbers        []string
	)

	for _, good := range orderJoinGoods {
		orderNumbers = append(orderNumbers, good.Number)
	}

	query := "number,sum(lack_count) as lack_count"

	err, numsMp = model.OrderGoodsNumsStatisticalByNumbers(db, query, orderNumbers)

	if err != nil {
		return
	}

	for _, o := range orderJoinGoods {

		nums, numsOk := numsMp[o.Number]

		if !numsOk {
			err = errors.New("订单欠货数量统计异常")
			return
		}

		//还有欠货 [这里需要保证每次出库(复核完成)时，都更新了商品订单欠货数]
		if nums.LackCount > 0 {
			lackNumbers = append(lackNumbers, o.Number)
			continue
		}

		deleteIds = append(deleteIds, o.Id)

		completeMp[o.Number] = struct{}{}

		//完成订单订单号
		completeNumbers = append(completeNumbers, o.Number)

		//完成的订单在订单表中删除
		deleteNumbers = append(deleteNumbers, o.Number)

		//完成订单
		completeOrder = append(completeOrder, model.CompleteOrder{
			Number:         o.Number,
			OrderRemark:    o.OrderRemark,
			ShopId:         o.ShopId,
			ShopName:       o.ShopName,
			ShopType:       o.ShopType,
			ShopCode:       o.ShopCode,
			Line:           o.Line,
			DeliveryMethod: o.DistributionType,
			Province:       o.Province,
			City:           o.City,
			District:       o.District,
			PickTime:       o.LatestPickingTime,
			PayAt:          o.PayAt,
		})
	}

	if len(completeNumbers) > 0 {
		//如果有完成订单重新查询订单商品，因为这批商品可能是多个批次拣的，根据批次查的商品不全
		err, orderGoods = model.GetOrderGoodsListByNumbers(db, completeNumbers)

		if err != nil {
			return
		}

		for _, og := range orderGoods {
			//完成订单map中不存在订单的跳过
			_, ok := completeMp[og.Number]

			if !ok {
				continue
			}
			//完成订单详情
			completeOrderDetail = append(completeOrderDetail, model.CompleteOrderDetail{
				Number:          og.Number,
				GoodsName:       og.GoodsName,
				Sku:             og.Sku,
				GoodsSpe:        og.GoodsSpe,
				GoodsType:       og.GoodsType,
				Shelves:         og.Shelves,
				PayCount:        og.PayCount,
				CloseCount:      og.CloseCount,
				ReviewCount:     og.OutCount,
				GoodsRemark:     og.GoodsRemark,
				DeliveryOrderNo: og.DeliveryOrderNo,
			})
		}
	}

	if len(deleteIds) > 0 {
		err = model.DeleteOrderByIds(db, deleteIds)

		if err != nil {
			return
		}
	}

	if len(deleteNumbers) > 0 {
		deleteNumbers = slice.UniqueSlice(deleteNumbers)

		err = model.DeleteOrderGoodsByNumbers(db, deleteNumbers)

		if err != nil {
			return
		}
	}

	// 欠货的单 拣货池还有未完成的，不更新为欠货
	if len(lackNumbers) > 0 {
		var (
			pickAndGoods   []model.PickAndGoods
			pendingNumbers []string
			diffSlice      []string
		)

		//获取欠货的订单number是否有在拣货池中未复核完成的数据，如果有，过滤掉欠货的订单number
		err, pickAndGoods = model.GetPickGoodsJoinPickByNumbers(db, lackNumbers)

		if err != nil {
			return
		}

		//获取拣货id，根据拣货id查出 拣货单中 未复核完成的订单，不更新为欠货，
		//且 有未复核完成的订单 不发送到mq中，完成后再发送到mq中
		for _, p := range pickAndGoods {
			//已经被接单，且未完成复核
			if p.Status < model.ReviewCompletedStatus && p.PickUser != "" {
				pendingNumbers = append(pendingNumbers, p.Number)
				isSendMQ = false
			}
		}

		pendingNumbers = slice.UniqueSlice(pendingNumbers)

		diffSlice = slice.StrDiff(lackNumbers, pendingNumbers) // 在 lackNumbers 不在 pendingNumbers 中的

		if len(diffSlice) > 0 {
			//更新为欠货
			err = model.UpdateOrderByNumbers(db, diffSlice, map[string]interface{}{"order_type": model.LackOrderType})
			if err != nil {
				return
			}
		}

	}

	//保存完成订单
	if len(completeOrder) > 0 {
		err = model.CompleteOrderBatchSave(db, &completeOrder)
		if err != nil {
			return
		}
	}

	//保存完成订单详情
	if len(completeOrderDetail) > 0 {
		err = model.CompleteOrderDetailBatchSave(db, &completeOrderDetail)

		if err != nil {
			return
		}
	}

	//更新 OutboundOrder 为已完成
	err = model.UpdateOutboundOrderByTaskIdAndNumbers(db, taskId, orderNumbers, map[string]interface{}{"order_type": model.OutboundOrderTypeComplete})

	if err != nil {

		return
	}

	if isSendMQ {
		//mq 存入 批次id
		err = SyncBatch(batchId)
		if err != nil {
			return errors.New("写入mq失败")
		}
	}

	return nil
}

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

	//只能从暂停状态变更为关闭
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

	//查询批次下全部订单商品数据
	err, pickGoods = model.GetPickGoodsByBatchId(db, batch.Id)

	if err != nil {
		tx.Rollback()
		err = errors.New("推送u8拣货数据查询失败:" + err.Error())
		return
	}

	err = YongYouLog(tx, pickGoods, orderJoinGoods, batch.Id)

	if err != nil {
		tx.Rollback()
		return
	}

	var (
		orderNumbers     = make([]string, 0, len(orderJoinGoods))
		notReviewNumbers = make([]string, 0, len(orderJoinGoods))
	)

	// 批次结束时如果有已接单但未复核完成的，不变更任务订单状态为完成
	for _, good := range pickGoods {
		if good.ReviewNum == 0 {
			notReviewNumbers = append(notReviewNumbers, good.Number)
		} else {
			orderNumbers = append(orderNumbers, good.Number)
		}
	}

	orderNumbers = slice.UniqueSlice(orderNumbers)
	notReviewNumbers = slice.UniqueSlice(notReviewNumbers)

	//去掉已结但未复核完成的
	orderNumbers = slice.Diff(orderNumbers, notReviewNumbers)

	//更新 OutboundOrder 为已完成
	err = model.UpdateOutboundOrderByTaskIdAndNumbers(db, batch.TaskId, orderNumbers, map[string]interface{}{"order_type": model.OutboundOrderTypeComplete})

	if err != nil {
		tx.Rollback()
		return
	}

	err = SendBatchMsgToPurchase(db, batch.Id, 0)

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

// 批次结束 || 确认出库 订货系统MQ交互逻辑
func SendBatchMsgToPurchase(tx *gorm.DB, batchId int, pickId int) (err error) {
	var picks []model.Pick

	err, picks = model.GetPickList(tx, &model.Pick{BatchId: batchId})

	if err != nil {
		return
	}

	var isSend = true

	for _, ps := range picks {
		//确认出库中刚更新的数据立即查询可能数据还在缓存中，那边传递拣货ID过来，直接跳过，认为时拣货复核完成的。
		//批次中传0
		if ps.Id == pickId {
			continue
		}

		//拣货池有被接单且状态不是复核完成状态，则不发送msg，确认出库时再验证是否发送
		if ps.PickUser != "" && ps.Status < model.ReviewCompletedStatus {
			isSend = false
			break
		}
	}

	if isSend {
		err = SyncBatch(batchId)
	}

	return
}

// 批量变更批次状态为 暂停||进行中
func BatchCloseBatch(tx *gorm.DB, status int) (err error) {
	var (
		batchList []model.Batch
		batchIds  []int
		statusMp  = make(map[int]int, 0)
	)

	err, batchList = model.GetBatchList(tx, &model.Batch{Status: status})

	if err != nil {
		return
	}

	if len(batchList) == 0 {
		err = errors.New("没有可以被暂停或开启的批次数据")
		return
	}

	//状态相互转换
	statusMp = map[int]int{
		model.BatchOngoingStatus: model.BatchSuspendStatus,
		model.BatchSuspendStatus: model.BatchOngoingStatus,
	}

	statusVal, statusOk := statusMp[status]

	if !statusOk {
		err = errors.New("状态异常")
		return
	}

	//批次id
	for _, batch := range batchList {
		batchIds = append(batchIds, batch.Id)
	}

	//变更状态
	err = model.UpdateBatchByIds(tx, batchIds, map[string]interface{}{"status": statusVal})

	if err != nil {
		return
	}

	return
}

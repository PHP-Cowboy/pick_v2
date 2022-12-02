package dao

import (
	"context"
	"errors"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"gorm.io/gorm"

	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
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
	err, orderGoodsIds, outboundGoods, _, _ = CreatePrePickLogic(tx, form, claims, batch.Id)

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
	result := db.Model(&model.OutboundTask{}).First(&outboundTask, form.TaskId)

	if result.Error != nil {
		return result.Error, batch
	}

	//t_batch
	err, batch = model.BatchSave(db, model.Batch{
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
	})

	if err != nil {
		return err, batch
	}

	return
}

// 快递批次
func CourierBatch(db *gorm.DB, form req.NewCreateBatchForm, claims *middlewares.CustomClaims) error {

	var (
		orderGoodsIds          []int
		outboundGoods          []model.OutboundGoods
		outboundGoodsJoinOrder []model.OutboundGoodsJoinOrder
		//prePickIds             []int
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
	err, orderGoodsIds, outboundGoods, outboundGoodsJoinOrder, _ = CreatePrePickLogic(tx, form, claims, batch.Id)

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

	//pick := req.BatchPickForm{
	//	BatchId:     batch.Id,
	//	Ids:         prePickIds,
	//	Type:        1,
	//	TypeParam:   []string{},
	//	WarehouseId: claims.WarehouseId,
	//	Typ:         2,
	//}
	//
	////生成拣货池
	//err = BatchPickByParams(db, pick, 2)
	//if err != nil {
	//	tx.Rollback()
	//	return err
	//}

	tx.Commit()

	return nil
}

// 根据sku,number,shop_id查批次id
func GetBatchIdsFromPrePickGoods(db *gorm.DB, form req.GetBatchListForm) (err error, batchIds []int) {
	var prePickGoods []model.PrePickGoods

	result := global.DB.Model(&model.PrePickGoods{}).
		Where(model.PrePickGoods{
			Sku:    form.Sku,
			Number: form.Number,
			ShopId: form.ShopId,
		}).
		Select("batch_id").
		Find(&prePickGoods)

	err = result.Error

	if err != nil {
		return
	}

	//未找到，直接返回
	if result.RowsAffected == 0 {
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

	local.Where(&model.Batch{DeliveryMethod: form.DeliveryMethod})

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

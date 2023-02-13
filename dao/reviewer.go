package dao

import (
	"errors"
	"pick_v2/forms/rsp"
	"time"

	"gorm.io/gorm"
	"pick_v2/common/constant"
	"pick_v2/forms/req"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/slice"
)

func ReviewList(db *gorm.DB, form req.ReviewListReq) (err error, res rsp.ReviewListRsp) {
	var (
		pick       []model.Pick
		pickRemark []model.PickRemark
	)

	localDb := db.Model(&model.Pick{}).Where("status = ?", form.Status)

	//1:待复核,2:复核完成
	if form.Status == 1 {
		localDb.Where("review_user = ? or review_user = ''", form.UserName)
	} else {
		localDb.Where("review_user = ? ", form.UserName).Order("review_time desc")
	}

	result := localDb.Where(model.Pick{PickUser: form.Name}).Find(&pick)

	err = result.Error

	if err != nil {
		return
	}

	res.Total = result.RowsAffected

	size := len(pick)

	list := make([]rsp.Pick, 0, size)

	if size == 0 {
		return
	}

	//拣货ids
	pickIds := make([]int, 0, len(pick))

	for _, p := range pick {
		pickIds = append(pickIds, p.Id)
	}

	//拣货ids 的订单备注
	result = db.Where("pick_id in (?)", pickIds).Find(&pickRemark)

	err = result.Error

	if err != nil {
		return
	}

	query := "pick_id,count(distinct(shop_id)) as shop_num,count(distinct(number)) as order_num,sum(need_num) as need_num,sum(complete_num) as complete_num,sum(review_num) as review_num"

	err, numsMp := model.CountPickPoolNumsByPickIds(db, pickIds, query)
	if err != nil {
		return
	}

	//构建pickId 对应的订单 是否有备注map
	remarkMp := make(map[int]struct{}, 0) //key 存在即为有
	for _, remark := range pickRemark {
		remarkMp[remark.PickId] = struct{}{}
	}

	isRemark := false

	for _, p := range pick {

		_, remarkMpOk := remarkMp[p.Id]
		if remarkMpOk { //拣货id在拣货备注中存在，即为有备注
			isRemark = true
		}

		nums, numsOk := numsMp[p.Id]

		if !numsOk {
			err = errors.New("拣货相关数量统计有误")
			return
		}

		list = append(list, rsp.Pick{
			Id:             p.Id,
			TaskName:       p.TaskName,
			ShopCode:       p.ShopCode,
			ShopName:       p.ShopName,
			ShopNum:        nums.ShopNum,
			OrderNum:       nums.OrderNum,
			NeedNum:        nums.NeedNum,
			PickUser:       p.PickUser,
			TakeOrdersTime: p.TakeOrdersTime,
			IsRemark:       isRemark,
			PickNum:        nums.CompleteNum,
			ReviewNum:      nums.ReviewNum,
			Num:            p.Num,
			ReviewUser:     p.ReviewUser,
			ReviewTime:     p.ReviewTime,
		})
	}

	res.List = list

	return
}

// 确认出库
func ConfirmDelivery(db *gorm.DB, form req.ConfirmDeliveryReq) (err error) {

	var (
		pick           model.Pick
		pickGoods      []model.PickGoods
		batch          model.Batch
		orderJoinGoods []model.GoodsJoinOrder
	)

	//根据id获取拣货数据
	err, pick = model.GetPickByPk(db, form.Id)

	if err != nil {
		return
	}

	//不是待复核状态
	if pick.Status != model.ToBeReviewedStatus {
		err = ecode.OrderNotToBeReviewed
		return
	}

	//复核接单的人不是当前用户
	if pick.ReviewUser != form.UserName {
		err = ecode.DataNotExist
		return
	}

	//获取批次数据
	err, batch = model.GetBatchByPk(db, pick.BatchId)

	if err != nil {
		return
	}

	//生成出库单号
	deliveryOrderNo, err := cache.GetIncrNumberByKey(constant.DELIVERY_ORDER_NO, 3)

	if err != nil {
		return
	}

	var (
		orderGoodsIds      []int
		orderPickGoodsIdMp = make(map[int]int, 0)
		skuCompleteNumMp   = make(map[string]int, 0)
	)

	for _, cp := range form.CompleteReview {

		//全部订单数据id
		for _, ids := range cp.ParamsId {
			orderGoodsIds = append(orderGoodsIds, ids.OrderGoodsId)
			//map[订单表id]拣货商品表id
			orderPickGoodsIdMp[ids.OrderGoodsId] = ids.PickGoodsId
		}
		//sku完成数量
		skuCompleteNumMp[cp.Sku] = cp.ReviewNum
	}

	//step: 根据 订单表id切片 查出订单数据 根据支付时间升序
	err, orderJoinGoods = model.GetOrderGoodsJoinOrderByIds(db, orderGoodsIds)

	if err != nil {
		return
	}

	//拣货商品表 id 和 拣货复核数量
	pickGoodsReviewNumMp := make(map[int]int, 0)

	//出库订单商品
	outboundGoods := make([]model.OutboundGoods, 0, len(orderJoinGoods))

	//step: 构造 拣货商品表 id, 完成数量 并扣减 sku 完成数量
	for _, info := range orderJoinGoods {
		//完成数量
		completeNum, completeOk := skuCompleteNumMp[info.Sku]

		if !completeOk {
			continue
		}

		pickGoodsId, mpOk := orderPickGoodsIdMp[info.OrderGoodsId]

		if !mpOk {
			continue
		}

		reviewCompleteNum := 0

		if completeNum >= info.LackCount { //完成数量大于等于需拣数量
			reviewCompleteNum = info.LackCount
			skuCompleteNumMp[info.Sku] = completeNum - info.LackCount //减
		} else {
			//按下单时间拣货少于需拣时
			reviewCompleteNum = completeNum
			skuCompleteNumMp[info.Sku] = 0
		}

		pickGoodsReviewNumMp[pickGoodsId] = reviewCompleteNum

		deliveryOrderNoArr := make(model.GormList, 0)
		//一个任务下一个商品只会有一个出库单
		deliveryOrderNoArr = append(deliveryOrderNoArr, deliveryOrderNo)

		//构造更新出库单商品表数据
		outboundGoods = append(outboundGoods, model.OutboundGoods{
			TaskId:          pick.TaskId,
			Number:          info.Number,
			Sku:             info.Sku,
			LackCount:       info.LackCount - reviewCompleteNum,
			OutCount:        reviewCompleteNum,
			Status:          model.OutboundGoodsStatusOutboundDelivery,
			DeliveryOrderNo: deliveryOrderNoArr,
		})
	}

	//获取拣货商品数据
	err, pickGoods = model.GetPickGoodsByPickIds(db, []int{form.Id})

	if err != nil {
		return
	}

	//构造打印 chan 结构体数据
	printChMp := make(map[int]struct{}, 0)

	//构造更新 出库任务订单表数据 使用
	var orderNumbers []string

	for k, pg := range pickGoods {
		_, printChOk := printChMp[pg.ShopId]

		if !printChOk {
			printChMp[pg.ShopId] = struct{}{}
		}

		completeNum, pickGoodsReviewNumMpOk := pickGoodsReviewNumMp[pg.Id]

		if !pickGoodsReviewNumMpOk {
			continue
		}

		pickGoods[k].ReviewNum = completeNum

		//更新订单表
		orderNumbers = append(orderNumbers, pg.Number)

	}

	orderNumbers = slice.UniqueSlice(orderNumbers)

	//当前时间
	now := time.Now()

	//出库订单
	outboundOrder := make([]model.OutboundOrder, 0, len(orderNumbers))

	//构造更新出库订单数据
	for _, number := range orderNumbers {
		outboundOrder = append(outboundOrder, model.OutboundOrder{
			TaskId:            pick.TaskId,
			Number:            number,
			LatestPickingTime: (*model.MyTime)(&now),
			OrderType:         model.PickingOrderType,
		})
	}

	tx := db.Begin()

	//更新出库任务商品表数据
	err = model.OutboundGoodsReplaceSave(tx, &outboundGoods, []string{"lack_count", "out_count", "status", "delivery_order_no"})
	if err != nil {
		tx.Rollback()
		return
	}

	no := model.GormList{deliveryOrderNo}

	val, err := no.Value()

	if err != nil {
		tx.Rollback()
		return
	}

	outboundType := model.OutboundTypeNormal

	if *form.Num == 0 {
		outboundType = model.OutboundTypeNoNeedToIssue
	}

	//更新拣货池表
	err = model.UpdatePickByPk(tx, pick.Id, map[string]interface{}{
		"status":            model.ReviewCompletedStatus,
		"review_time":       &now,
		"num":               *form.Num,
		"delivery_order_no": val,
		"delivery_no":       deliveryOrderNo,
		"print_num":         gorm.Expr("print_num + ?", 1),
		"outbound_type":     outboundType,
	})

	if err != nil {
		tx.Rollback()
		return
	}

	//更新拣货商品数据
	err = model.PickGoodsReplaceSave(tx, &pickGoods, []string{"review_num"})

	if err != nil {
		tx.Rollback()
		return
	}

	//变更出库任务订单表最近拣货时间&&订单类型
	err = model.OutboundOrderReplaceSave(tx, outboundOrder, []string{"latest_picking_time", "order_type"})
	if err != nil {
		tx.Rollback()
		return
	}

	//批次已结束的复核出库要单独推u8
	if batch.Status == model.BatchClosedStatus {
		//如果出库任务已结束，则需要更新订单和订单商品表&&完成订单和完成订单表状态&&推送订货系统【前面已经更新了出库单相关数据】
		var (
			outboundTask    model.OutboundTask
			picks           []model.Pick
			reviewedNumbers []string
		)

		err, picks = model.GetPickList(tx, &model.Pick{BatchId: batch.Id})

		if err != nil {
			return
		}

		//批次结束时，要更新出库任务订单状态
		err, reviewedNumbers = UpdateOutboundOrderLogic(tx, orderNumbers, batch.TaskId, pick.Id, picks)

		if err != nil {
			tx.Rollback()
			return
		}

		//如果出库任务也结束了 则需要更新 订单和订单商品表&&完成订单和完成订单表状态
		//出库任务的关闭必须要所有批次被关闭，所以这个可以写在批次结束判断内部
		err, outboundTask = model.GetOutboundTaskById(db, pick.TaskId)

		if err != nil {
			tx.Rollback()
			return
		}

		//任务是否结束
		if outboundTask.Status == model.OutboundTaskStatusClosed {
			//任务结束时出库更新订单
			err = TaskEndDeliveryUpdateOrders(tx, orderNumbers, reviewedNumbers, pickGoods, no)
			if err != nil {
				tx.Rollback()
				return
			}
		}

		err = YongYouLog(tx, pickGoods, orderJoinGoods, pick.BatchId)
		if err != nil {
			tx.Rollback()
			return
		}

		//消息中处理了是否发送消息逻辑
		err = SendBatchMsgToPurchase(batch.Id, pick.Id, picks)

		if err != nil {
			tx.Rollback()
			return
		}

	}

	tx.Commit()

	//拆单 -打印
	for shopId := range printChMp {
		AddPrintJobMap(constant.JH_HUOSE_CODE, &global.PrintCh{
			DeliveryOrderNo: deliveryOrderNo,
			ShopId:          shopId,
			Type:            1, // 1-全部打印 2-打印箱单 3-打印出库单 第一次全打，后边的前段选
		})
	}

	return
}

func UpdateOutboundOrderLogic(tx *gorm.DB, numbers []string, taskId, pickId int, picks []model.Pick) (err error, reviewedNumbers []string) {
	var (
		pickIds            []int
		pickGoods          []model.PickGoods
		notReviewedNumbers []string
	)
	//step1:查询当前批次已接单未复核完成的任务
	for _, ps := range picks {
		//确认出库中刚更新的数据立即查询可能数据还在缓存中，那边传递拣货ID过来，直接跳过，认为时拣货复核完成的。
		if ps.Id == pickId {
			continue
		}

		//拣货池有被接单且状态不是复核完成状态
		if ps.PickUser != "" && ps.Status < model.ReviewCompletedStatus {
			pickIds = append(pickIds, ps.Id)
			break
		}
	}

	//step2:查询任务中订单
	if len(pickIds) > 0 {
		err, pickGoods = model.GetPickGoodsByPickIds(tx, pickIds)

		if err != nil {
			return
		}

		for _, good := range pickGoods {
			notReviewedNumbers = append(notReviewedNumbers, good.Number)
		}

		//去重
		notReviewedNumbers = slice.UniqueSlice(notReviewedNumbers)

	}

	if len(notReviewedNumbers) > 0 {
		//本次出库订单号，且不在其他已接单未复核出库任务中的单号
		numbers = slice.Diff(numbers, notReviewedNumbers)
	}

	//step3:取出在本次拣货中但不在step2中查询到的订单编号里的单号，更新成已完成

	err = model.UpdateOutboundOrderByTaskIdAndNumbers(tx, taskId, numbers, map[string]interface{}{"order_type": model.OutboundOrderTypeComplete})
	if err != nil {
		return
	}

	reviewedNumbers = numbers

	return
}

// 任务结束时出库更新订单
func TaskEndDeliveryUpdateOrders(
	tx *gorm.DB,
	orderNumbers,
	reviewedNumbers []string,
	pickGoods []model.PickGoods,
	deliveryOrderNo model.GormList,
) (err error) {
	var (
		orderJoinGoods []model.GoodsJoinOrder
	)

	//根据拣货订单编号查询 订单&&商品信息
	err, orderJoinGoods = model.GetOrderJoinGoodsListByNumbers(tx, orderNumbers)

	if err != nil {
		return
	}

	//欠货订单编号map
	lackNumbersMap := make(map[string]struct{}, 0)
	//订单复核数map
	orderGoodsReviewMp := make(map[int]int, 0)
	//已经完成复核的订单map,不在其他已接未复核完成的单里
	reviewedNumbersMp := slice.SliceToMap(reviewedNumbers)

	for _, pg := range pickGoods {
		orderGoodsReviewMp[pg.OrderGoodsId] = pg.ReviewNum
	}

	//构造欠货订单map数据
	for i, good := range orderJoinGoods {
		orderGoodsReviewNum, orderGoodsReviewMpOk := orderGoodsReviewMp[good.OrderGoodsId]

		lackCount := good.LackCount
		//本次拣货，可能没有订单中的这个品[非全单拣货]
		if orderGoodsReviewMpOk {
			//本次拣了这个品，欠货数量 -= 本次单品复核数量
			lackCount -= orderGoodsReviewNum
		}

		//订单商品欠货数量
		orderJoinGoods[i].LackCount = lackCount

		//订单商品出货数量 = 历史出货数量 + 本次出货数量
		orderJoinGoods[i].OutCount = good.OutCount + orderGoodsReviewNum

		_, reviewedNumbersMpOk := reviewedNumbersMp[good.Number]

		//商品还有欠货,且当前单不在其他已接未复核完成的单里，即订单为欠货单
		if lackCount > 0 && reviewedNumbersMpOk {
			lackNumbersMap[good.Number] = struct{}{}
		}
	}

	err = OrderUpdateHandle(tx, orderJoinGoods, lackNumbersMap, nil, nil, deliveryOrderNo)

	if err != nil {
		return
	}

	return
}

// 快捷出库
func QuickDelivery(db *gorm.DB, form req.QuickDeliveryReq) (err error) {

	var (
		picks                   []model.Pick
		pickGoods               []model.PickGoods
		batch                   model.Batch
		pickIdDeliveryOrderNoMp = make(map[int]string, 0)
	)

	//根据id获取拣货数据
	err, picks = model.GetPickListByIds(db, form.Ids)

	if err != nil {
		return
	}

	//当前时间
	now := time.Now()

	for i, p := range picks {
		//不是待复核状态
		if p.Status != model.ToBeReviewedStatus {
			err = ecode.OrderNotToBeReviewed
			return
		}

		var deliveryOrderNo string

		//生成出库单号
		deliveryOrderNo, err = cache.GetIncrNumberByKey(constant.DELIVERY_ORDER_NO, 3)

		if err != nil {
			return
		}

		pickIdDeliveryOrderNoMp[p.Id] = deliveryOrderNo

		no := model.GormList{deliveryOrderNo}

		picks[i].Status = model.ReviewCompletedStatus
		picks[i].ReviewTime = (*model.MyTime)(&now)
		picks[i].DeliveryOrderNo = no
		picks[i].DeliveryNo = deliveryOrderNo
		picks[i].PrintNum += 1
		picks[i].OutboundType = model.OutboundTypeNormal

	}

	//获取批次数据
	err, batch = model.GetBatchByPk(db, form.BatchId)

	if err != nil {
		return
	}

	err, pickGoods = model.GetPickGoodsByPickIds(db, form.Ids)

	if err != nil {
		return
	}

	var (
		printChMp      = make(map[int]int, 0)               //构造打印 chan 结构体数据
		orderNumbers   = make([]string, 0)                  //构造更新 出库任务订单表数据 使用
		skuCompleteNum = make(map[string]map[string]int, 0) //订单号对应的sku的拣货数量 map[number][sku] = CompleteNum
		numberPickIdMp = make(map[string]int, 0)            //订单号对应的拣货id，在后台拣货中一个订单只会对应一个拣货id，其他的拣货任务中不一定【合并拣货一个单可能对应多个拣货ID】
	)

	for k, pg := range pickGoods {

		numberSkuMp, skuCompleteNumOk := skuCompleteNum[pg.Number]

		if !skuCompleteNumOk {
			numberSkuMp = make(map[string]int, 0)
		}

		numberSkuMp[pg.Sku] = pg.CompleteNum

		skuCompleteNum[pg.Number] = numberSkuMp

		_, printChOk := printChMp[pg.ShopId]

		if !printChOk {
			printChMp[pg.ShopId] = pg.PickId
		}

		pickGoods[k].ReviewNum = pg.CompleteNum

		numberPickIdMp[pg.Number] = pg.PickId

		//更新订单表
		orderNumbers = append(orderNumbers, pg.Number)

	}

	orderNumbers = slice.UniqueSlice(orderNumbers)

	//出库订单
	outboundOrder := make([]model.OutboundOrder, 0, len(orderNumbers))

	//构造更新出库订单数据
	for _, number := range orderNumbers {
		outboundOrder = append(outboundOrder, model.OutboundOrder{
			TaskId:            batch.TaskId,
			Number:            number,
			LatestPickingTime: (*model.MyTime)(&now),
			OrderType:         model.PickingOrderType,
		})
	}

	//出库订单商品
	outboundGoods := make([]model.OutboundGoods, 0)

	err, outboundGoods = model.GetOutboundGoodsListByTaskIdAndNumbers(db, batch.TaskId, orderNumbers)

	if err != nil {
		return
	}

	for _, og := range outboundGoods {
		pickId, numberPickIdMpOk := numberPickIdMp[og.Number]

		if !numberPickIdMpOk {
			err = errors.New("订单编号和拣货任务ID对应关系出错")
			return
		}

		deliveryOrderNo, pickIdDeliveryOrderNoMpOK := pickIdDeliveryOrderNoMp[pickId]

		if !pickIdDeliveryOrderNoMpOK {
			err = errors.New("拣货任务对应的出库单号出错")
			return
		}

		numberSkuMp, skuCompleteNumOk := skuCompleteNum[og.Number]

		if !skuCompleteNumOk {
			err = errors.New("任务商品订单数据有误")
			return
		}

		completeNum, numberSkuMpOk := numberSkuMp[og.Sku]

		if !numberSkuMpOk {
			continue
		}

		deliveryOrderNoArr := make(model.GormList, 0)
		//一个任务下一个商品只会有一个出库单
		deliveryOrderNoArr = append(deliveryOrderNoArr, deliveryOrderNo)

		//构造更新出库单商品表数据
		outboundGoods = append(outboundGoods, model.OutboundGoods{
			TaskId:          batch.TaskId,
			Number:          og.Number,
			Sku:             og.Sku,
			LackCount:       og.LackCount - completeNum,
			OutCount:        completeNum,
			Status:          model.OutboundGoodsStatusOutboundDelivery,
			DeliveryOrderNo: deliveryOrderNoArr,
		})
	}

	tx := db.Begin()

	//更新出库任务商品表数据
	err = model.OutboundGoodsReplaceSave(tx, &outboundGoods, []string{"lack_count", "out_count", "status", "delivery_order_no"})
	if err != nil {
		tx.Rollback()
		return
	}

	//拣货任务数据更新
	err = model.PickReplaceSave(tx, &picks, []string{"status", "review_time", "delivery_order_no", "delivery_no", "print_num", "outbound_type"})

	if err != nil {
		tx.Rollback()
		return
	}

	//更新拣货商品数据
	err = model.PickGoodsReplaceSave(tx, &pickGoods, []string{"review_num"})

	if err != nil {
		tx.Rollback()
		return
	}

	//变更出库任务订单表最近拣货时间&&订单类型
	err = model.OutboundOrderReplaceSave(tx, outboundOrder, []string{"latest_picking_time", "order_type"})
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	//拆单 -打印
	for shopId, pickId := range printChMp {
		deliveryOrderNo, pickIdDeliveryOrderNoMpOk := pickIdDeliveryOrderNoMp[pickId]

		if !pickIdDeliveryOrderNoMpOk {
			continue
		}

		AddPrintJobMap(constant.JH_HUOSE_CODE, &global.PrintCh{
			DeliveryOrderNo: deliveryOrderNo,
			ShopId:          shopId,
			Type:            3, // 1-全部打印 2-打印箱单 3-打印出库单
		})
	}

	return
}

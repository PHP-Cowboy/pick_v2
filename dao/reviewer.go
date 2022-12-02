package dao

import (
	"errors"
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

func ConfirmDelivery(db *gorm.DB, form req.ConfirmDeliveryReq) (err error) {

	var (
		pick           model.Pick
		pickGoods      []model.PickGoods
		batch          model.Batch
		orderGoods     []model.OrderGoods
		orderJoinGoods []model.OrderJoinGoods
	)

	//根据id获取拣货数据
	err, pick = model.GetPickByPk(db, form.Id)

	if err != nil {
		return
	}

	if pick.Status == 2 {
		err = ecode.OrderHasBeenReviewedAndCompleted
		return
	}

	if pick.ReviewUser != form.UserName {
		err = ecode.DataNotExist
		return
	}

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

	//拣货表 id 和 拣货数量
	mp := make(map[int]int, 0)

	type OrderGoods struct {
		OutCount           int
		deliveryOrderNoArr model.GormList
	}

	orderGoodsMp := make(map[int]OrderGoods, 0)

	var (
		pickGoodsIds []int
		//出库订单商品
		outboundGoods = make([]model.OutboundGoods, 0, len(orderJoinGoods))
		//订单商品id
		orderGoodsId []int
	)

	//step: 构造 拣货商品表 id, 完成数量 并扣减 sku 完成数量
	for _, info := range orderJoinGoods {
		//完成数量
		completeNum, completeOk := skuCompleteNumMp[info.Sku]

		if !completeOk {
			continue
		}

		pickGoodsId, mpOk := orderPickGoodsIdMp[info.Id]

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
		pickGoodsIds = append(pickGoodsIds, pickGoodsId)
		mp[pickGoodsId] = reviewCompleteNum

		deliveryOrderNoArr := make(model.GormList, 0)

		deliveryOrderNoArr = append(deliveryOrderNoArr, info.DeliveryOrderNo...)
		deliveryOrderNoArr = append(deliveryOrderNoArr, deliveryOrderNo)

		orderGoodsMp[info.Id] = OrderGoods{
			OutCount:           reviewCompleteNum,
			deliveryOrderNoArr: deliveryOrderNoArr,
		}

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

		orderGoodsId = append(orderGoodsId, info.Id)
	}

	//获取拣货商品数据
	err, pickGoods = model.GetPickGoodsByPickIds(db, []int{form.Id})

	if err != nil {
		return
	}

	//构造打印 chan 结构体数据
	printChMp := make(map[int]struct{}, 0)

	//构造更新 订单表 订单商品 表完成出库数据
	var orderNumbers []string

	for k, pg := range pickGoods {
		_, printChOk := printChMp[pg.ShopId]

		if !printChOk {
			printChMp[pg.ShopId] = struct{}{}
		}

		completeNum, mpOk := mp[pg.Id]

		if !mpOk {
			continue
		}

		pickGoods[k].ReviewNum = completeNum

		//更新订单表
		orderNumbers = append(orderNumbers, pg.Number)

	}

	orderNumbers = slice.UniqueSlice(orderNumbers)

	//order_goods 这里会被mysql排序
	err, orderGoods = model.GetOrderGoodsListByIds(db, orderGoodsIds)
	if err != nil {
		return
	}

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
		})
	}

	orderPickMp := make(map[string]int)

	for i, good := range orderGoods {

		val, ok := orderGoodsMp[good.Id]

		if !ok {
			continue
		}

		_, ogMpOk := orderPickMp[good.Number]

		if !ogMpOk {
			orderPickMp[good.Number] = val.OutCount
		} else {
			orderPickMp[good.Number] += val.OutCount
		}

		orderGoods[i].LackCount = good.LackCount - val.OutCount
		orderGoods[i].OutCount = val.OutCount
		orderGoods[i].DeliveryOrderNo = val.deliveryOrderNoArr
	}

	var order []model.Order
	//更新订单表 已拣 未拣
	err, order = model.GetOrderListByNumbers(db, orderNumbers)

	if err != nil {
		return
	}

	tx := db.Begin()

	//更新拣货单商品表数据
	err = model.OutboundGoodsReplaceSave(tx, outboundGoods, []string{"lack_count", "out_count", "status", "delivery_order_no"})
	if err != nil {
		tx.Rollback()
		return
	}

	//更新订单商品数据
	err = model.OrderGoodsReplaceSave(db, &orderGoods, []string{"lack_count", "out_count", "delivery_order_no"})
	if err != nil {
		tx.Rollback()
		return
	}

	//变更出库任务订单表最近拣货时间
	err = model.OutboundOrderReplaceSave(tx, outboundOrder, []string{"latest_picking_time"})
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

	//更新拣货池表
	err = model.UpdatePickByPk(tx, pick.Id, map[string]interface{}{
		"status":            model.ReviewCompletedStatus,
		"review_time":       &now,
		"num":               form.Num,
		"delivery_order_no": val,
		"delivery_no":       deliveryOrderNo,
		"print_num":         gorm.Expr("print_num + ?", 1),
	})

	if err != nil {
		tx.Rollback()
		return
	}

	//前边更新的 pick.Id的 status 后面查询还未生效
	pickStatusMp := make(map[int]int, 1)
	pickStatusMp[pick.Id] = model.ReviewCompletedStatus

	lackNumberMp := make(map[string]struct{}, 0) //要被更新为欠货的 number map

	//批次已结束的
	if batch.Status == model.BatchClosedStatus {
		// 欠货逻辑 查出当前出库的所有商品 number ，这些number如果还有未复核完成状态的就不更新为欠货
		var (
			picks          []model.Pick
			isSendMQ       = true
			pickNumbers    = make([]string, 0) //当前出库的所有商品 number
			pickAndGoods   []model.PickAndGoods
			pendingNumbers []string
			diffSlice      []string
		)
		//查出当前批次拣货池数据，如果有未完成复核的，就不发送mq消息
		err, picks = model.GetPickList(db, model.Pick{BatchId: batch.Id})

		if err != nil {
			tx.Rollback()
			return
		}

		for _, ps := range picks {

			status, ok := pickStatusMp[ps.Id]

			if ok {
				ps.Status = status //刚更新的状态立即查询，mysql数据可能还在缓存中差不到，更新成map中保存的状态
			}

			if ps.Status < model.ReviewCompletedStatus {
				//批次有未复核完成的拣货单
				isSendMQ = false
				continue
			}
		}

		//获取当前出库的所有商品 number
		for _, good := range pickGoods {
			pickNumbers = append(pickNumbers, good.Number)
		}

		pickNumbers = slice.UniqueStringSlice(pickNumbers)

		//获取欠货的订单number是否有在拣货池中未复核完成的数据，如果有，过滤掉欠货的订单number
		err, pickAndGoods = model.GetPickGoodsJoinPickByNumbers(db, pickNumbers)
		if err != nil {
			tx.Rollback()
			return
		}

		//获取拣货id，根据拣货id查出 拣货单中 未复核完成状态的订单，不更新为欠货，
		//且 有未复核完成的订单 不发送到mq中，完成后再发送到mq中
		for _, p := range pickAndGoods {
			status, ok := pickStatusMp[p.PickId]

			if ok {
				p.Status = status //刚更新的状态立即查询，mysql数据可能还在缓存中差不到，更新成map中保存的状态
			}

			if p.Status < model.ReviewCompletedStatus {
				pendingNumbers = append(pendingNumbers, p.Number)
				isSendMQ = false
			}
		}

		pendingNumbers = slice.UniqueStringSlice(pendingNumbers)

		diffSlice = slice.StrDiff(pickNumbers, pendingNumbers) // 在 pickNumbers 不在 pendingNumbers 中的

		if len(diffSlice) > 0 {
			//构造 更新为欠货 number map
			for _, s := range diffSlice {
				lackNumberMp[s] = struct{}{}
			}
		}

		if isSendMQ {
			//mq 存入 批次id
			err = SyncBatch(batch.Id)
			if err != nil {
				tx.Rollback()
				return
			}
		}
	}

	//更新完成订单
	var (
		completeOrder  []model.CompleteOrder
		completeNumber []string // 查询&&删除 完成订单详情使用
		numsMp         = make(map[string]model.OrderGoodsNumsStatistical, 0)
	)

	query := "number,sum(lack_count) as lack_count"

	err, numsMp = model.OrderGoodsNumsStatisticalByNumbers(db, query, orderNumbers)

	if err != nil {
		tx.Rollback()
		return
	}

	//如果是批次结束时 还在拣货池中的单的出库，且拣货池中没有未复核的商品，要更新 order_type
	for i, o := range order {
		picked, ogMpOk := orderPickMp[o.Number]

		if !ogMpOk {
			continue
		}

		nums, numsOk := numsMp[o.Number]

		if !numsOk {
			err = errors.New("订单欠货数量统计异常")
			tx.Rollback()
			return
		}

		order[i].LatestPickingTime = (*model.MyTime)(&now)

		_, ok := lackNumberMp[o.Number]

		if ok {
			order[i].OrderType = model.LackOrderType //更新为欠货
		}

		//之前的订单欠货数，减去本次订单拣货数 为0的 且 批次结束 改成 完成订单
		if nums.LackCount-picked == 0 && batch.Status == model.BatchClosedStatus {

			completeNumber = append(completeNumber, o.Number)

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
	}

	//更新订单数据
	err = model.OrderReplaceSave(tx, order, []string{"shop_id", "shop_name", "shop_type", "shop_code", "house_code", "line", "order_type", "latest_picking_time"})

	if err != nil {
		tx.Rollback()
		return
	}

	if len(completeNumber) > 0 {
		//保存完成订单
		err = model.CompleteOrderBatchSave(tx, &completeOrder)

		if err != nil {
			tx.Rollback()
			return
		}

		//根据条件重新查询完成订单详情
		err, orderGoods = model.GetOrderGoodsListByNumbers(db, completeNumber)

		if err != nil {
			tx.Rollback()
			return
		}

		var completeOrderDetail []model.CompleteOrderDetail

		//保存完成订单详情
		for _, good := range orderGoods {
			completeOrderDetail = append(completeOrderDetail, model.CompleteOrderDetail{
				Number:          good.Number,
				GoodsName:       good.GoodsName,
				Sku:             good.Sku,
				GoodsSpe:        good.GoodsSpe,
				GoodsType:       good.GoodsType,
				Shelves:         good.Shelves,
				PayCount:        good.PayCount,
				CloseCount:      good.CloseCount,
				ReviewCount:     good.OutCount,
				GoodsRemark:     good.GoodsRemark,
				DeliveryOrderNo: good.DeliveryOrderNo,
			})
		}

		err = model.CompleteOrderDetailBatchSave(tx, &completeOrderDetail)

		if err != nil {
			tx.Rollback()
			return
		}

		//删除完成订单
		err = model.DeleteOrderByNumbers(tx, completeNumber)

		if err != nil {
			tx.Rollback()
			return
		}

		//删除完成订单详情
		err = model.DeleteOrderGoodsByNumbers(tx, completeNumber)

		if err != nil {
			tx.Rollback()
			return
		}
	}

	//更新拣货商品数据
	result := tx.Save(&pickGoods)

	if result.Error != nil {
		tx.Rollback()
		return
	}

	//拆单 -打印
	for shopId, _ := range printChMp {
		AddPrintJobMap(constant.JH_HUOSE_CODE, &global.PrintCh{
			DeliveryOrderNo: deliveryOrderNo,
			ShopId:          shopId,
			Type:            1, // 1-全部打印 2-打印箱单 3-打印出库单 第一次全打，后边的前段选
		})
	}

	result = db.First(&batch, pick.BatchId)
	if result.Error != nil {
		tx.Rollback()
		return
	}

	//批次已结束,这个不能往前移，里面有commit，移到前面去如果进入commit，后面的又有失败的，事务无法保证一致性了
	if batch.Status == model.BatchClosedStatus {
		err = YongYouLog(tx, pickGoods, orderJoinGoods, pick.BatchId)
		if err != nil {
			tx.Commit() // u8推送失败不能影响仓库出货，只提示，业务继续
			return
		}
	}

	tx.Commit()

	return
}

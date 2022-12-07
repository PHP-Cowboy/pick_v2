package dao

import (
	"errors"
	"pick_v2/forms/rsp"
	"strconv"
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
		orderJoinGoods []model.OrderJoinGoods
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

		pickGoodsReviewNumMp[pickGoodsId] = reviewCompleteNum

		deliveryOrderNoArr := make(model.GormList, 0)

		deliveryOrderNoArr = append(deliveryOrderNoArr, info.DeliveryOrderNo...)
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

	//构造更新 订单表 订单商品 表完成出库数据
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
		})
	}

	tx := db.Begin()

	//更新出库任务商品表数据
	err = model.OutboundGoodsReplaceSave(tx, outboundGoods, []string{"lack_count", "out_count", "status", "delivery_order_no"})
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

	//更新拣货商品数据
	err = model.PickGoodsReplaceSave(tx, &pickGoods, []string{"review_num"})

	if err != nil {
		tx.Rollback()
		return
	}

	//拆单 -打印
	for shopId := range printChMp {
		AddPrintJobMap(constant.JH_HUOSE_CODE, &global.PrintCh{
			DeliveryOrderNo: deliveryOrderNo,
			ShopId:          shopId,
			Type:            1, // 1-全部打印 2-打印箱单 3-打印出库单 第一次全打，后边的前段选
		})
	}

	//批次已结束的复核出库要单独推u8
	if batch.Status == model.BatchClosedStatus {
		err = YongYouLog(tx, pickGoods, orderJoinGoods, pick.BatchId)
		if err != nil {
			err = errors.New("出库成功，但是推送u8失败:" + strconv.Itoa(form.Id))
			tx.Commit() // u8推送失败不能影响仓库出货，只提示，业务继续
			return
		}
	}

	tx.Commit()

	return
}

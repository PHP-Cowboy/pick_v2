package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"pick_v2/common/constant"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/slice"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"sort"
	"time"
)

// 复核列表 通过状态区分是否已完成
func ReviewList(c *gin.Context) {
	var (
		form       req.ReviewListReq
		res        rsp.ReviewListRsp
		pick       []model.Pick
		pickRemark []model.PickRemark
		//pickListModel []rsp.PickListModel
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	//result := db.Model(&batch.Pick{}).Where("status = ?", form.Status).Where(batch.Pick{PickUser: form.Name}).Find(&pick)

	localDb := db.Model(&model.Pick{}).Where("status = ?", form.Status)

	claims, ok := c.Get("claims")

	if !ok {
		xsq_net.ErrorJSON(c, errors.New("claims 获取失败"))
		return
	}

	userInfo := claims.(*middlewares.CustomClaims)

	//1:待复核,2:复核完成
	if form.Status == 1 {
		localDb.Where("review_user = ? or review_user = ''", userInfo.Name)
	} else {
		localDb.Where("review_user = ? ", userInfo.Name)
	}

	result := localDb.Where(model.Pick{PickUser: form.Name}).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	list := make([]rsp.Pick, 0, form.Size)

	if len(pick) > 0 {
		//拣货ids
		pickIds := make([]int, 0, len(pick))
		for _, p := range pick {
			pickIds = append(pickIds, p.Id)
		}

		//拣货ids 的订单备注
		result = db.Where("pick_id in (?)", pickIds).Find(&pickRemark)
		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
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

			reviewTime := ""

			if p.ReviewTime != nil {
				reviewTime = p.ReviewTime.Format(timeutil.TimeFormat)
			}

			list = append(list, rsp.Pick{
				Id:             p.Id,
				TaskName:       p.TaskName,
				ShopCode:       p.ShopCode,
				ShopName:       p.ShopName,
				ShopNum:        p.ShopNum,
				OrderNum:       p.OrderNum,
				NeedNum:        p.NeedNum,
				PickUser:       p.PickUser,
				TakeOrdersTime: p.TakeOrdersTime.Format(timeutil.TimeFormat),
				IsRemark:       isRemark,
				PickNum:        p.PickNum,
				ReviewNum:      p.ReviewNum,
				Num:            p.Num,
				ReviewUser:     p.ReviewUser,
				ReviewTime:     reviewTime,
			})
		}
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// 复核明细
func ReviewDetail(c *gin.Context) {
	var (
		form req.ReviewDetailReq
		res  rsp.ReviewDetailRsp
		pick model.Pick
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	result := db.First(&pick, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//复核员为空，接单复核
	if pick.ReviewUser == "" {
		claims, ok := c.Get("claims")

		if !ok {
			xsq_net.ErrorJSON(c, errors.New("claims 获取失败"))
			return
		}

		userInfo := claims.(*middlewares.CustomClaims)

		result = db.Model(&model.Pick{}).
			Where("id = ? and version = ?", pick.Id, pick.Version).
			Updates(map[string]interface{}{"review_user": userInfo.Name, "version": pick.Version + 1})

		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	}

	var (
		pickGoods  []model.PickGoods
		pickRemark []model.PickRemark
	)

	res.TaskName = pick.TaskName
	res.ShopCode = pick.ShopCode
	res.PickUser = pick.PickUser
	res.TakeOrdersTime = pick.TakeOrdersTime.Format(timeutil.TimeFormat)
	res.ReviewUser = pick.ReviewUser

	var reviewTime string

	if pick.ReviewTime != nil {
		reviewTime = pick.ReviewTime.Format(timeutil.TimeFormat)
	}
	res.ReviewTime = reviewTime

	result = db.Where("pick_id = ?", form.Id).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	pickGoodsSkuMp := make(map[string]rsp.MergePickGoods, 0)

	goodsMap := make(map[string][]rsp.MergePickGoods, 0)

	for _, goods := range pickGoods {

		val, ok := pickGoodsSkuMp[goods.Sku]

		paramsId := rsp.ParamsId{
			PickGoodsId:  goods.Id,
			OrderGoodsId: goods.OrderGoodsId,
		}

		if !ok {

			pickGoodsSkuMp[goods.Sku] = rsp.MergePickGoods{
				Id:          goods.Id,
				Sku:         goods.Sku,
				GoodsName:   goods.GoodsName,
				GoodsType:   goods.GoodsType,
				GoodsSpe:    goods.GoodsSpe,
				Shelves:     goods.Shelves,
				NeedNum:     goods.NeedNum,
				CompleteNum: goods.CompleteNum,
				ReviewNum:   goods.ReviewNum,
				Unit:        goods.Unit,
				ParamsId:    []rsp.ParamsId{paramsId},
			}
		} else {
			val.NeedNum += goods.NeedNum
			val.CompleteNum += goods.CompleteNum
			val.ReviewNum += goods.ReviewNum
			val.ParamsId = append(val.ParamsId, paramsId)
			pickGoodsSkuMp[goods.Sku] = val
		}
	}

	needTotal := 0
	completeTotal := 0
	reviewTotal := 0

	for _, goods := range pickGoodsSkuMp {
		completeTotal += goods.CompleteNum
		needTotal += goods.NeedNum
		reviewTotal += goods.ReviewNum
		goodsMap[goods.GoodsType] = append(goodsMap[goods.GoodsType], rsp.MergePickGoods{
			Id:          goods.Id,
			Sku:         goods.Sku,
			GoodsName:   goods.GoodsName,
			GoodsType:   goods.GoodsType,
			GoodsSpe:    goods.GoodsSpe,
			Shelves:     goods.Shelves,
			NeedNum:     goods.NeedNum,
			CompleteNum: goods.CompleteNum,
			ReviewNum:   goods.ReviewNum,
			Unit:        goods.Unit,
			ParamsId:    goods.ParamsId,
		})
	}

	res.ShopCode = pick.ShopCode
	res.OutTotal = completeTotal
	res.UnselectedTotal = needTotal - completeTotal
	res.ReviewTotal = reviewTotal

	//按货架号排序
	for s, goods := range goodsMap {

		ret := rsp.MyMergePickGoods(goods)

		sort.Sort(ret)

		goodsMap[s] = ret
	}

	res.Goods = goodsMap

	result = db.Where("pick_id = ?", form.Id).Find(&pickRemark)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := []rsp.PickRemark{}
	for _, remark := range pickRemark {
		list = append(list, rsp.PickRemark{
			Number:      remark.Number,
			OrderRemark: remark.OrderRemark,
			GoodsRemark: remark.GoodsRemark,
		})
	}

	res.RemarkList = list

	xsq_net.SucJson(c, res)
}

// 确认出库
func ConfirmDelivery(c *gin.Context) {
	var (
		form req.ConfirmDeliveryReq
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		pick          model.Pick
		pickGoods     []model.PickGoods
		batch         model.Batch
		orderAndGoods []rsp.OrderAndGoods
	)

	db := global.DB

	//根据id获取拣货数据
	result := db.First(&pick, form.Id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			xsq_net.ErrorJSON(c, ecode.DataNotExist)
			return
		}
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if pick.Status == 2 {
		xsq_net.ErrorJSON(c, ecode.OrderHasBeenReviewedAndCompleted)
		return
	}

	result = db.First(&batch, pick.BatchId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, errors.New("获取用户数据失败"))
		return
	}

	if pick.ReviewUser != userInfo.Name {
		xsq_net.ErrorJSON(c, ecode.DataNotExist)
		return
	}

	//生成出库单号
	deliveryOrderNo, err := cache.GetIncrNumberByKey(constant.DELIVERY_ORDER_NO, 3)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var (
		orderGoodsIds      []int
		orderPickGoodsIdMp = make(map[int]int, 0)
		skuCompleteNumMp   = make(map[string]int, 0)
		totalNum           int //更新拣货池复核数量
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
		totalNum += cp.ReviewNum //总复核数量
	}

	//step: 根据 订单表id切片 查出订单数据 根据支付时间升序
	result = db.Table("t_pick_order_goods og").
		Select("og.*").
		Joins("left join t_pick_order o on og.pick_order_id = o.id").
		Where("og.id in (?)", orderGoodsIds).
		Order("pay_at ASC").
		Find(&orderAndGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
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
		//拣订单商品
		pickOrderGoods = make([]model.PickOrderGoods, 0, len(orderAndGoods))
		//订单商品
		orderGoods []model.OrderGoods
		//
		orderGoodsId []int
	)

	//step: 构造 拣货商品表 id, 完成数量 并扣减 sku 完成数量
	for _, info := range orderAndGoods {
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

		orderGoodsMp[info.OrderGoodsId] = OrderGoods{
			OutCount:           reviewCompleteNum,
			deliveryOrderNoArr: deliveryOrderNoArr,
		}
		//构造更新拣货单商品表数据
		pickOrderGoods = append(pickOrderGoods, model.PickOrderGoods{
			Base: model.Base{
				Id: info.Id,
			},
			LackCount:       info.LackCount - reviewCompleteNum,
			OutCount:        reviewCompleteNum,
			Status:          2,
			DeliveryOrderNo: deliveryOrderNoArr,
		})

		orderGoodsId = append(orderGoodsId, info.OrderGoodsId)
	}

	//获取拣货商品数据
	result = db.Where("pick_id = ?", form.Id).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//构造打印 chan 结构体数据
	printChMp := make(map[int]struct{}, 0)

	//构造更新 订单表 订单商品 表完成出库数据
	orderNumbers := []string{}
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

	//order_goods 这里会被mysql排序
	result = db.Model(&model.OrderGoods{}).Where("id in (?)", orderGoodsId).Find(&orderGoods)

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
	result = db.Model(&model.Order{}).Where("number in (?)", orderNumbers).Find(&order)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	now := time.Now()

	tx := db.Begin()

	//更新拣货单商品表数据
	result = tx.Select("id", "update_time", "lack_count", "out_count", "status", "delivery_order_no").Save(&pickOrderGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//更新订单商品数据
	result = tx.Select("id", "update_time", "lack_count", "out_count", "delivery_order_no").Save(&orderGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//更新 pick_order 表 最近拣货时间
	result = tx.Model(&model.PickOrder{}).
		Where("number in (?)", orderNumbers).
		Updates(map[string]interface{}{
			"latest_picking_time": now.Format(timeutil.TimeFormat),
		})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	no := model.GormList{deliveryOrderNo}

	val, err := no.Value()

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	//更新拣货池表
	result = tx.Model(&model.Pick{}).
		Where("id = ?", pick.Id).
		Updates(map[string]interface{}{
			"status":            model.ReviewCompletedStatus,
			"review_time":       &now,
			"num":               form.Num,
			"review_num":        totalNum,
			"delivery_order_no": val,
			"delivery_no":       deliveryOrderNo,
			"print_num":         gorm.Expr("print_num + ?", 1),
		})

	//前边更新的 pick.Id的 status 后面查询还未生效
	pickStatusMp := make(map[int]int, 1)
	pickStatusMp[pick.Id] = model.ReviewCompletedStatus

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

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
		result = db.Model(&model.Pick{}).Where("batch_id = ?", batch.Id).Find(&picks)
		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
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
		result = db.Table("t_pick_goods pg").
			Select("p.id as pick_id,p.status,pg.number").
			Joins("left join t_pick p on pg.pick_id = p.id").
			Where("number in (?)", pickNumbers).
			Find(&pickAndGoods)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
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
				xsq_net.ErrorJSON(c, errors.New("写入mq失败"))
				return
			}
		}
	}

	//更新完成订单
	var (
		completeOrder  []model.CompleteOrder
		completeNumber []string // 查询&&删除 完成订单详情使用
	)

	//如果是批次结束时 还在拣货池中的单的出库，且拣货池中没有未复核的商品，要更新 order_type
	for i, o := range order {
		picked, ogMpOk := orderPickMp[o.Number]

		if !ogMpOk {
			continue
		}

		order[i].Picked = picked + o.Picked //订单被拆分成多个出，已拣累加
		order[i].UnPicked = o.UnPicked - picked

		order[i].LatestPickingTime = &now

		_, ok := lackNumberMp[o.Number]

		if ok {
			order[i].OrderType = model.LackOrderType //更新为欠货
		}

		//之前的欠货数，减去本次订单拣货数 为0的 且 批次结束 改成 完成订单
		if o.UnPicked-picked == 0 && batch.Status == model.BatchClosedStatus {

			completeNumber = append(completeNumber, o.Number)

			payAt, payAtErr := time.ParseInLocation(timeutil.TimeZoneFormat, o.PayAt, time.Local)

			if payAtErr != nil {
				xsq_net.ErrorJSON(c, ecode.DataTransformationError)
				return
			}

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
				PayCount:       o.PayTotal,
				CloseCount:     o.CloseNum,
				OutCount:       o.Picked,
				Province:       o.Province,
				City:           o.City,
				District:       o.District,
				PickTime:       o.LatestPickingTime,
				PayAt:          payAt.Format(timeutil.TimeFormat),
			})
		}
	}

	//更新订单数据
	result = tx.
		Omit("number", "pay_at", "delivery_at", "has_remark", "address", "order_remark").
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"latest_picking_time", "picked", "un_picked", "order_type"}),
		}).
		Save(&order)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if len(completeNumber) > 0 {
		//保存完成订单
		result = tx.Save(&completeOrder)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		//根据条件重新查询完成订单详情
		result = db.Model(&model.OrderGoods{}).Where("number in (?)", completeNumber).Find(&orderGoods)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
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

		result = tx.Save(&completeOrderDetail)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		//删除完成订单
		result = tx.Delete(&model.Order{}, "number in (?)", completeNumber)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		//删除完成订单详情
		result = tx.Delete(&model.OrderGoods{}, "number in (?)", completeNumber)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	}

	//更新拣货商品数据
	result = tx.Save(&pickGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//拆单 -打印
	for shopId, _ := range printChMp {
		AddPrintJobMap(constant.JH_HUOSE_CODE, &global.PrintCh{
			DeliveryOrderNo: deliveryOrderNo,
			ShopId:          shopId,
		})
	}

	result = db.First(&batch, pick.BatchId)
	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	err = UpdateBatchPickNums(tx, pick.BatchId)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//批次已结束,这个不能往前移，里面有commit，移到前面去如果进入commit，后面的又有失败的，事务无法保证一致性了
	if batch.Status == model.BatchClosedStatus {
		err = YongYouLog(tx, pickGoods, orderAndGoods, pick.BatchId)
		if err != nil {
			tx.Commit() // u8推送失败不能影响仓库出货，只提示，业务继续
			xsq_net.ErrorJSON(c, err)
			return
		}
	}

	tx.Commit()

	xsq_net.Success(c)
}

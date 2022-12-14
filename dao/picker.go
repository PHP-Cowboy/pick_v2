package dao

import (
	"errors"

	"fmt"
	"pick_v2/forms/rsp"
	"time"

	"gorm.io/gorm"

	"pick_v2/forms/req"
	"pick_v2/model"
	"pick_v2/utils/ecode"
)

func GetPick(db *gorm.DB, pick []model.Pick) (res rsp.ReceivingOrdersRsp, err error) {

	if len(pick) == 1 { //只查到一条
		res.Id = pick[0].Id
		res.BatchId = pick[0].BatchId
		res.Version = pick[0].Version
		res.TakeOrdersTime = pick[0].TakeOrdersTime
	} else { //查到多条
		//排序
		var (
			batchIds []int
			batchMp  = make(map[int]struct{}, 0)
			pickMp   = make(map[int][]model.Pick, 0)
		)

		//去重，构造批次id切片
		for _, b := range pick {
			//构造批次下的拣货池数据map
			//批次排序后，直接获取某个批次的全部拣货池数据。
			//然后对这部分数据排序
			pickMp[b.BatchId] = append(pickMp[b.BatchId], b)
			//已经存入了批次map的，跳过
			_, bMpOk := batchMp[b.BatchId]
			if bMpOk {
				continue
			}
			//写入批次mp
			batchMp[b.BatchId] = struct{}{}
			//存入批次id切片
			batchIds = append(batchIds, b.BatchId)
		}

		var (
			batch model.Batch
		)

		if len(batchIds) == 0 { //只有一个批次
			batch.Id = batchIds[0]
		} else {
			//多个批次
			err, batch = model.GetBatchListByBatchIdsAndSort(db, batchIds, "sort desc")

			if err != nil {
				return
			}
		}

		maxSort := 0

		res.BatchId = batch.Id

		//循环排序最大的批次下的拣货数据，并取出sort最大的那个的id
		for _, pm := range pickMp[batch.Id] {
			if pm.Sort >= maxSort {
				res.Id = pm.Id
				res.Version = pm.Version
				res.TakeOrdersTime = pm.TakeOrdersTime
			}
		}
	}

	return
}

func ReceivingOrders(db *gorm.DB, form req.ReceivingOrdersForm) (err error, res rsp.ReceivingOrdersRsp) {
	var (
		picks   []model.Pick
		batches []model.Batch
	)

	// 先查询是否有当前拣货员被分配的任务或已经接单且未完成拣货的数据,如果被分配多条，第一按批次优先级，第二按拣货池优先级 优先拣货
	err, picks = model.GetPickListByPickUserAndStatusAndTyp(db, form.UserName, model.ToBePickedStatus, form.Typ)

	if err != nil {
		return
	}

	now := time.Now()

	//有分配的拣货任务
	if len(picks) > 0 {
		res, err = GetPick(db, picks)
		if err != nil {
			return
		}
		//后台分配的单没有接单时间,更新接单时间
		if res.TakeOrdersTime == nil {
			err = model.UpdatePickByIds(db, []int{res.Id}, map[string]interface{}{"take_orders_time": &now})

			if err != nil {
				return
			}
		}
		return
	}

	//进行中的批次
	err, batches = model.GetBatchListByStatusAndTyp(db, model.BatchOngoingStatus, form.Typ)

	if err != nil {
		return
	}

	batchIds := make([]int, 0)

	for _, b := range batches {
		batchIds = append(batchIds, b.Id)
	}

	if len(batchIds) == 0 {
		err = errors.New("没有进行中的批次,无法接单")
		return
	}

	//查询未被接单的拣货池数据
	err, picks = model.GetPickListNoOrderReceived(db, batchIds, form.Typ)

	if err != nil {
		return
	}

	//拣货池有未接单的数据
	if len(picks) > 0 {

		res, err = GetPick(db, picks)
		if err != nil {
			return
		}

		//更新拣货池 + version 防并发
		err = model.UpdatePickByPkAndVersion(db, res.Id, res.Version, map[string]interface{}{
			"pick_user":        form.UserName,
			"take_orders_time": &now,
			"version":          res.Version + 1,
		})

		if err != nil {
			return
		}
		return
	} else {
		err = errors.New("暂无拣货单")
		return
	}
}

func CompletePick(db *gorm.DB, form req.CompletePickForm) (err error) {
	// 这里是否需要做并发处理
	var (
		pick           model.Pick
		pickGoods      []model.PickGoods
		orderJoinGoods []model.OrderJoinGoods
	)

	err, pick = model.GetPickByPk(db, form.PickId)

	if err != nil {
		return
	}

	if pick.Status == 1 {
		err = ecode.OrderPickingCompleted
		return
	}

	if pick.PickUser != form.UserName {
		err = errors.New("请确认拣货单是否被分配给其他拣货员")
		return
	}

	tx := db.Begin()

	//****************************** 无需拣货 ******************************//
	if form.Type == 2 {
		//更新主表 无需拣货直接更新为复核完成
		err = model.UpdatePickByIds(tx, []int{pick.Id}, map[string]interface{}{"status": model.ReviewCompletedStatus})

		if err != nil {
			tx.Rollback()
			return
		}

		// 更新拣货数量(PickGoods.CompleteNum)为0
		err = model.UpdatePickGoodsByPickId(tx, pick.Id, map[string]interface{}{"complete_num": 0, "review_num": 0})

		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()

		return
	}
	//****************************** 无需拣货逻辑完成 ******************************//

	//****************************** 正常拣货逻辑 ******************************//
	//step:处理前端传递的拣货数据，构造[订单表id切片,订单表id和拣货商品表id map,sku完成数量 map]
	//step: 根据 订单表id切片 查出订单数据 根据支付时间升序
	//step: 构造 拣货商品表 id, 完成数量 并扣减 sku 完成数量
	//step: 更新拣货商品表

	var (
		orderGoodsIds      []int
		orderPickGoodsIdMp = make(map[int]int, 0)
		skuCompleteNumMp   = make(map[string]int, 0)
	)

	//step:处理前端传递的拣货数据，构造[订单表id切片,订单表id和拣货商品表id映射,sku完成数量映射]
	for _, cp := range form.CompletePick {
		//全部订单数据id
		for _, ids := range cp.ParamsId {
			orderGoodsIds = append(orderGoodsIds, ids.OrderGoodsId)
			//map[订单表id]拣货商品表id
			orderPickGoodsIdMp[ids.OrderGoodsId] = ids.PickGoodsId
		}
		//sku完成数量
		skuCompleteNumMp[cp.Sku] = cp.CompleteNum
	}

	//step: 根据 订单表id切片 查出订单数据 根据支付时间升序
	err, orderJoinGoods = model.GetOrderGoodsJoinOrderByIds(db, orderGoodsIds)

	if err != nil {
		return
	}

	type CloseGoodsInfo struct {
		Name string `json:"name"`
		Num  int    `json:"num"`
	}

	//订单商品需拣数量按sku累计
	//1.防止拣货数量大于需拣数量
	//2.关单或关品后拣货员未刷新提示拣货数量大于需拣数量
	skuNeedPickGoodsInfoMp := make(map[string]CloseGoodsInfo, 0)

	for _, good := range orderJoinGoods {
		//需拣数量
		needPickGoodsInfo, skuNeedPickGoodsInfoMpOk := skuNeedPickGoodsInfoMp[good.Sku]

		if !skuNeedPickGoodsInfoMpOk {
			needPickGoodsInfo = CloseGoodsInfo{
				Num:  0,
				Name: good.GoodsName,
			}
		}

		needPickGoodsInfo.Num += good.LackCount

		skuNeedPickGoodsInfoMp[good.Sku] = needPickGoodsInfo
	}

	//关闭商品map
	closeMp := make(map[string]CloseGoodsInfo, 0)

	for sku, needPickGoodsInfo := range skuNeedPickGoodsInfoMp {
		completeNum, skuCompleteNumMpOk := skuCompleteNumMp[sku]

		if !skuCompleteNumMpOk {
			err = errors.New("需拣数量完成数量对比异常")
			return
		}

		needPickNum := needPickGoodsInfo.Num

		if completeNum > needPickNum {
			closeMp[sku] = CloseGoodsInfo{
				Name: needPickGoodsInfo.Name,
				Num:  completeNum - needPickNum,
			}
			continue
		}
	}

	if len(closeMp) > 0 {
		errStr := "请注意，有商品关闭，"

		for _, cl := range closeMp {
			errStr += fmt.Sprintf("%s，关闭%d件", cl.Name, cl.Num)
		}

		err = errors.New(errStr)

		return
	}

	//拣货表 id 和 拣货数量
	mp := make(map[int]int, 0)

	var pickGoodsIds []int

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

		pickCompleteNum := 0

		if completeNum >= info.LackCount { //完成数量大于等于需拣数量
			pickCompleteNum = info.LackCount
			skuCompleteNumMp[info.Sku] = completeNum - info.LackCount //减
		} else {
			//按下单时间拣货少于需拣时
			pickCompleteNum = completeNum
			skuCompleteNumMp[info.Sku] = 0
		}
		pickGoodsIds = append(pickGoodsIds, pickGoodsId)
		mp[pickGoodsId] = pickCompleteNum

	}

	//查出拣货商品数据
	err, pickGoods = model.GetPickGoodsByIds(db, pickGoodsIds)

	if err != nil {
		return
	}

	//更新拣货数量数据
	for i, good := range pickGoods {
		completeNum, mpOk := mp[good.Id]

		if !mpOk {
			continue
		}

		pickGoods[i].CompleteNum = completeNum
	}

	//正常拣货 更新拣货数量
	err = model.PickGoodsReplaceSave(tx, &pickGoods, []string{"need_num", "complete_num"})

	if err != nil {
		tx.Rollback()
		return
	}

	//更新主表
	err = model.UpdatePickByIds(db, []int{pick.Id}, map[string]interface{}{"status": model.ToBeReviewedStatus})

	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}

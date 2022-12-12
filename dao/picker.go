package dao

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"pick_v2/forms/req"
	"pick_v2/model"
	"pick_v2/utils/ecode"
)

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
		err = model.UpdatePickGoodsByPickId(tx, pick.Id, map[string]interface{}{"complete_num": 0})

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

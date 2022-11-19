package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
)

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
	err, orderGoodsIds, outboundGoods = PrePickLogic(tx, form, claims, batch.Id)

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
	})

	if err != nil {
		return err, batch
	}

	return
}

// 生成预拣池数据逻辑
func PrePickLogic(db *gorm.DB, form req.NewCreateBatchForm, claims *middlewares.CustomClaims, batchId int) (err error, orderGoodsIds []int, outboundGoods []model.OutboundGoods) {
	var (
		outboundGoodsJoinOrder []model.OutboundGoodsJoinOrder
	)

	mp, err := cache.GetClassification()

	if err != nil {
		return err, nil, nil
	}

	//缓存中的线路数据
	lineCacheMp, errCache := cache.GetShopLine()

	if errCache != nil {
		return errCache, nil, nil
	}

	err, outboundGoodsJoinOrder = model.GetOutboundGoodsJoinOrderList(db, form.TaskId, form.Number)

	if err != nil {
		return err, nil, nil
	}

	var (
		prePicks      []model.PrePick
		prePickGoods  []model.PrePickGoods
		prePickRemark []model.PrePickRemark
		shopNumMp     = make(map[int]struct{}, 0) //店铺
	)

	for _, og := range outboundGoodsJoinOrder {
		//线路
		cacheMpLine, cacheMpOk := lineCacheMp[og.ShopId]

		if !cacheMpOk {
			return errors.New("店铺:" + og.ShopName + "线路未同步，请先同步"), nil, nil
		}

		//拣货池根据门店合并订单数据
		_, shopMpOk := shopNumMp[og.ShopId]

		if shopMpOk {
			continue
		}

		prePicks = append(prePicks, model.PrePick{
			WarehouseId: claims.WarehouseId,
			BatchId:     batchId,
			ShopId:      og.ShopId,
			ShopCode:    og.ShopCode,
			ShopName:    og.ShopName,
			Line:        cacheMpLine,
			Status:      0,
		})

		shopNumMp[og.ShopId] = struct{}{}
	}

	//预拣池数量

	if len(prePicks) == 0 {
		return ecode.NoOrderFound, nil, nil
	}

	//预拣池
	err, prePicksRes := model.PrePickBatchSave(db, prePicks)
	if err != nil {
		return err, nil, nil
	}

	//单次创建批次时，预拣池表数据根据shop_id唯一
	//构造 map[shop_id]pre_pick_id 更新 t_pre_pick_goods t_pre_pick_remark 表 pre_pick_id
	prePickMp := make(map[int]int, 0)
	for _, pp := range prePicksRes {
		prePickMp[pp.ShopId] = pp.Id
	}

	for _, og := range outboundGoodsJoinOrder {
		//用于更新 t_order_goods 的 batch_id
		orderGoodsIds = append(orderGoodsIds, og.OrderGoodsId)

		prePickId, ppMpOk := prePickMp[og.ShopId]

		if !ppMpOk {
			return errors.New(fmt.Sprintf("店铺ID: %v 无预拣池数据", og.ShopId)), nil, nil
		}

		//商品类型
		goodsType, mpOk := mp[og.GoodsType]

		if !mpOk {
			return errors.New("商品类型:" + og.GoodsType + "数据未同步"), nil, nil
		}

		//线路
		cacheMpLine, cacheMpOk := lineCacheMp[og.ShopId]

		if !cacheMpOk {
			return errors.New("店铺:" + og.ShopName + "线路未同步，请先同步"), nil, nil
		}

		needNum := og.LackCount

		//如果欠货数量大于限发数量，需拣货数量为限货数
		if og.LackCount > og.LimitNum {
			needNum = og.LimitNum
		}

		prePickGoods = append(prePickGoods, model.PrePickGoods{
			WarehouseId:      claims.WarehouseId,
			BatchId:          batchId,
			OrderGoodsId:     og.OrderGoodsId,
			Number:           og.Number,
			PrePickId:        prePickId,
			ShopId:           og.ShopId,
			DistributionType: og.DistributionType,
			Sku:              og.Sku,
			GoodsName:        og.GoodsName,
			GoodsType:        goodsType,
			GoodsSpe:         og.GoodsSpe,
			Shelves:          og.Shelves,
			DiscountPrice:    og.DiscountPrice,
			Unit:             og.GoodsUnit,
			NeedNum:          needNum,
			CloseNum:         og.CloseCount,
			OutCount:         0,
			NeedOutNum:       og.LackCount,
		})

		if og.GoodsRemark != "" || og.OrderRemark != "" {
			prePickRemark = append(prePickRemark, model.PrePickRemark{
				WarehouseId:  claims.WarehouseId,
				BatchId:      batchId,
				OrderGoodsId: og.OrderGoodsId,
				ShopId:       og.ShopId,
				PrePickId:    prePickId,
				Number:       og.Number,
				OrderRemark:  og.OrderRemark,
				GoodsRemark:  og.GoodsRemark,
				ShopName:     og.ShopName,
				Line:         cacheMpLine,
			})
		}

		outboundGoods = append(outboundGoods, model.OutboundGoods{
			TaskId: og.TaskId,
			Number: og.Number,
			Sku:    og.Sku,
			Status: model.OutboundGoodsStatusPicking,
		})
	}

	//预拣池商品
	err = model.PrePickGoodsBatchSave(db, prePickGoods)
	if err != nil {
		return err, nil, nil
	}

	//预拣池商品备注，可能没有备注
	if len(prePickRemark) > 0 {
		err = model.PrePickRemarkBatchSave(db, prePickRemark)
		if err != nil {
			return err, nil, nil
		}
	}

	return nil, orderGoodsIds, outboundGoods
}

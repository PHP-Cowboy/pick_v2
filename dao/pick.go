package dao

import (
	"errors"
	"pick_v2/utils/slice"
	"strconv"

	"gorm.io/gorm"

	"pick_v2/forms/req"
	"pick_v2/model"
)

// 批量拣货 - 根据参数类型
func BatchPickByParams(db *gorm.DB, form req.BatchPickForm, batchType int) error {

	var (
		prePick        []model.PrePick
		prePickRemarks []model.PrePickRemark
	)

	//0:未处理,1:已进入拣货池
	result := db.Model(&model.PrePick{}).Where("id in (?) and status = 0", form.Ids).Find(&prePick)

	if result.Error != nil {
		return result.Error
	}

	//按分类或商品获取未进入拣货池的商品数据
	err, prePickGoods := model.GetPrePickGoodsByTypeParam(db, form.Ids, form.Type, form.TypeParam)

	if err != nil {
		return err
	}

	if len(prePickGoods) == 0 {
		return errors.New("对应的拣货池商品不存在")
	}

	var (
		prePickGoodsIds   []int
		prePickRemarksIds []int
		pickGoods         []model.PickGoods
		pickRemark        []model.PickRemark
	)

	var picks = make([]model.Pick, 0, len(prePick))

	//拣货池数据处理
	for _, pre := range prePick {

		picks = append(picks, model.Pick{
			WarehouseId:    form.WarehouseId,
			BatchId:        pre.BatchId,
			PrePickIds:     strconv.Itoa(pre.Id),
			TaskName:       pre.ShopName,
			ShopCode:       pre.ShopCode,
			ShopName:       pre.ShopName,
			Line:           pre.Line,
			PickUser:       "",
			ReviewUser:     "",
			TakeOrdersTime: nil,
			Sort:           0,
			Version:        0,
			Typ:            batchType,
		})
	}

	err = model.PickSave(db, &picks)

	if err != nil {
		return err
	}

	//prePick 和 pick id 关系映射
	prePickIdsMp := make(map[int]int, 0)

	for _, p := range picks {
		//合并拣货时PrePickIds才会有多个，这里是只有一个的
		prePickId, atoiErr := strconv.Atoi(p.PrePickIds)

		if atoiErr != nil {
			return atoiErr
		}

		prePickIdsMp[prePickId] = p.Id
	}

	var orderGoodsIds []int

	for _, goods := range prePickGoods {

		pickId, prePickIdsMpOk := prePickIdsMp[goods.PrePickId]

		if !prePickIdsMpOk {
			return errors.New("pick_id不存在")
		}

		orderGoodsIds = append(orderGoodsIds, goods.OrderGoodsId)

		//更新 prePickGoods 使用
		prePickGoodsIds = append(prePickGoodsIds, goods.Id)

		pickGoods = append(pickGoods, model.PickGoods{
			WarehouseId:      form.WarehouseId,
			PickId:           pickId,
			BatchId:          goods.BatchId,
			PrePickGoodsId:   goods.Id,
			OrderGoodsId:     goods.OrderGoodsId,
			Number:           goods.Number,
			ShopId:           goods.ShopId,
			DistributionType: goods.DistributionType,
			Sku:              goods.Sku,
			GoodsName:        goods.GoodsName,
			GoodsType:        goods.GoodsType,
			GoodsSpe:         goods.GoodsSpe,
			Shelves:          goods.Shelves,
			DiscountPrice:    goods.DiscountPrice,
			NeedNum:          goods.NeedNum,
			Unit:             goods.Unit,
		})
	}

	result = db.Where("order_goods_id in (?)", orderGoodsIds).Find(&prePickRemarks)

	if result.Error != nil {
		return result.Error
	}

	for _, remark := range prePickRemarks {
		pickId, prePickIdsMpOk := prePickIdsMp[remark.PrePickId]

		if !prePickIdsMpOk {
			return errors.New("pick_id不存在")
		}

		//更新 prePickRemarks 使用
		prePickRemarksIds = append(prePickRemarksIds, remark.Id)

		pickRemark = append(pickRemark, model.PickRemark{
			WarehouseId:     form.WarehouseId,
			BatchId:         remark.BatchId,
			PickId:          pickId,
			PrePickRemarkId: remark.Id,
			OrderGoodsId:    remark.OrderGoodsId,
			Number:          remark.Number,
			OrderRemark:     remark.OrderRemark,
			GoodsRemark:     remark.GoodsRemark,
			ShopName:        remark.ShopName,
			Line:            remark.Line,
		})
	}

	//商品数据保存
	result = db.Save(&pickGoods)

	if result.Error != nil {
		return result.Error
	}

	//订单备注数据
	if len(pickRemark) > 0 {
		result = db.Save(&pickRemark)

		if result.Error != nil {
			return result.Error
		}
	}

	//更新预拣池商品表的商品数据状态
	if len(prePickGoodsIds) > 0 {
		err = model.UpdatePrePickGoodsByIds(db, prePickGoodsIds, map[string]interface{}{"status": model.PrePickGoodsStatusProcessing})

		if err != nil {
			return err
		}
	}

	var prePickIds []int
	//预拣池内商品全部进入拣货池时 更新 对应的 预拣池状态
	if form.Type == 1 { //全单拣货
		prePickIds = form.Ids
	} else {
		//0:未处理,1:已进入拣货池
		result = db.Model(&model.PrePickGoods{}).Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickGoods)
		if result.Error != nil {
			return result.Error
		}

		//将传过来的id转换成map
		idsMp := slice.SliceToMap(form.Ids)

		//去除未处理的预拣池id
		for _, good := range prePickGoods {
			delete(idsMp, good.PrePickId)
		}

		//将map转回切片
		prePickIds = slice.MapToSlice(idsMp)
	}

	if len(prePickIds) > 0 {
		err = model.UpdatePrePickByIds(db, prePickIds, map[string]interface{}{"status": model.PrePickStatusProcessing})

		if err != nil {
			return err
		}
	}

	//更新预拣池商品备注表的数据状态
	if len(prePickRemarksIds) > 0 {
		err = model.UpdatePrePickRemarkByIds(db, prePickRemarksIds, map[string]interface{}{"status": model.PrePickRemarkStatusProcessing})
		if err != nil {
			return err
		}
	}

	return nil
}

// 合并拣货
func MergePickByParams(db *gorm.DB, form req.MergePickForm) error {
	var (
		prePickGoods   []model.PrePickGoods
		prePickRemarks []model.PrePickRemark
		pickGoods      []model.PickGoods
		pickRemarks    []model.PickRemark
	)

	var (
		prePickIds string
		prePickGoodsIds,
		orderGoodsIds,
		prePickRemarksIds []int
	)

	result := db.Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickGoods)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("商品数据未找到")
	}

	tx := db.Begin()

	pick := model.Pick{
		WarehouseId:    form.WarehouseId,
		BatchId:        form.BatchId,
		PrePickIds:     prePickIds,
		TaskName:       form.TaskName,
		ShopCode:       "",
		ShopName:       form.TaskName,
		Line:           "",
		PickUser:       "",
		ReviewUser:     "",
		TakeOrdersTime: nil,
		Sort:           0,
		Version:        0,
	}

	result = tx.Save(&pick)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	for _, goods := range prePickGoods {

		prePickGoodsIds = append(prePickGoodsIds, goods.Id)

		orderGoodsIds = append(orderGoodsIds, goods.OrderGoodsId)

		pickGoods = append(pickGoods, model.PickGoods{
			WarehouseId:      form.WarehouseId,
			PickId:           pick.Id,
			BatchId:          goods.BatchId,
			PrePickGoodsId:   goods.Id,
			OrderGoodsId:     goods.OrderGoodsId,
			Number:           goods.Number,
			ShopId:           goods.ShopId,
			DistributionType: goods.DistributionType,
			Sku:              goods.Sku,
			GoodsName:        goods.GoodsName,
			GoodsType:        goods.GoodsType,
			GoodsSpe:         goods.GoodsSpe,
			Shelves:          goods.Shelves,
			DiscountPrice:    goods.DiscountPrice,
			NeedNum:          goods.NeedNum,
			Unit:             goods.Unit,
		})
	}

	result = tx.Save(&pickGoods)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	//更新预拣货池商品相关数据状态
	result = tx.Model(model.PrePickGoods{}).Where("id in (?)", prePickGoodsIds).Updates(map[string]interface{}{"status": 1})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	//预拣池内商品全部进入拣货池时 更新 对应的 预拣池状态
	if form.Type == 1 { //全单拣货
		result = tx.Model(model.PrePick{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"status": 1})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	} else {
		//0:未处理,1:已进入拣货池
		result = tx.Model(model.PrePickGoods{}).Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickGoods)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		//将传过来的id转换成map
		idsMp := make(map[int]struct{}, 0)

		for _, id := range form.Ids {
			idsMp[id] = struct{}{}
		}

		//去除未处理的预拣池id
		for _, good := range prePickGoods {
			delete(idsMp, good.PrePickId)
		}

		//将map转回切片
		prePickIdSlice := []int{}
		for id, _ := range idsMp {
			prePickIdSlice = append(prePickIdSlice, id)
		}

		if len(prePickIdSlice) > 0 {
			result = tx.Model(model.PrePick{}).Where("id in (?)", prePickIdSlice).Updates(map[string]interface{}{"status": 1})
			if result.Error != nil {
				tx.Rollback()
				return result.Error
			}
		}
	}

	result = db.Where("order_goods_id in (?)", orderGoodsIds).Find(&prePickRemarks)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if len(prePickRemarks) > 0 {
		for _, remark := range prePickRemarks {

			prePickRemarksIds = append(prePickRemarksIds, remark.Id)

			pickRemarks = append(pickRemarks, model.PickRemark{
				WarehouseId:     form.WarehouseId,
				BatchId:         form.BatchId,
				PickId:          pick.Id,
				PrePickRemarkId: remark.Id,
				OrderGoodsId:    remark.OrderGoodsId,
				Number:          remark.Number,
				OrderRemark:     remark.OrderRemark,
				GoodsRemark:     remark.GoodsRemark,
				ShopName:        remark.ShopName,
				Line:            remark.Line,
			})
		}

		result = tx.Save(&pickRemarks)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		//更新预拣货池备注相关数据状态
		result = tx.Model(model.PrePickRemark{}).Where("id in (?)", prePickRemarksIds).Updates(map[string]interface{}{"status": 1})

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	tx.Commit()

	return nil
}

package dao

import (
	"errors"
	"pick_v2/utils/slice"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"pick_v2/forms/req"
	"pick_v2/model"
)

func BatchPick(db *gorm.DB, form req.BatchPickForm) (err error) {
	var batch model.Batch

	err, batch = model.GetBatchByPk(db, form.BatchId)

	if err != nil {
		return
	}

	if batch.Status == 1 { //状态:0:进行中,1:已结束,2:暂停
		err = errors.New("请先开启拣货")
		return
	}

	err = BatchPickByParams(db, form, nil, nil, nil)

	return
}

// 批量拣货 - 根据参数类型
func BatchPickByParams(db *gorm.DB, form req.BatchPickForm, prePicks []model.PrePick, prePickGoods []model.PrePickGoods, prePickRemarks []model.PrePickRemark) (err error) {
	if form.Typ == 1 { //常规批次拣货池数据要查询，快递批次用传递过来的数据
		err, prePicks = model.GetPrePickByIdsAndStatus(db, form.Ids, model.PrePickStatusUnhandled)

		if err != nil {
			return
		}

		if len(prePicks) == 0 {
			err = errors.New("对应的预拣池数据不存在")
			return
		}

		//按分类或商品获取未进入拣货池的商品数据
		err, prePickGoods = model.GetPrePickGoodsByTypeParam(db, form.Ids, form.Type, form.TypeParam)

		if err != nil {
			return
		}

		if len(prePickGoods) == 0 {
			err = errors.New("对应的预拣池商品不存在")
			return
		}
	}

	var (
		prePickGoodsIds   []int
		prePickRemarksIds []int
		pickGoods         []model.PickGoods
		pickRemark        []model.PickRemark
	)

	var picks = make([]model.Pick, 0, len(prePicks))

	//拣货池数据处理
	for _, pre := range prePicks {

		picks = append(picks, model.Pick{
			WarehouseId:    form.WarehouseId,
			TaskId:         pre.TaskId,
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
			Typ:            form.Typ,
		})
	}

	err = model.PickBatchSave(db, &picks)

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

	if form.Typ == 1 {
		result := db.Where("order_goods_id in (?)", orderGoodsIds).Find(&prePickRemarks)

		if result.Error != nil {
			return result.Error
		}
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
	result := db.Save(&pickGoods)

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

	if form.Typ == 1 { //常规批次更新预拣池相关数据状态，快递批次不需要更新
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
	}

	return nil
}

// 合并拣货
func MergePick(db *gorm.DB, form req.MergePickForm) (err error) {
	var batch model.Batch

	err, batch = model.GetBatchByPk(db, form.BatchId)

	if err != nil {
		return
	}

	if batch.Status == 1 {
		err = errors.New("请先开启拣货")
		return
	}

	err = MergePickByParams(db, form, batch.Id)

	return
}

// 合并拣货 - 根据参数类型
func MergePickByParams(db *gorm.DB, form req.MergePickForm, taskId int) (err error) {
	var (
		prePickGoods   []model.PrePickGoods
		prePickRemarks []model.PrePickRemark
		pickGoods      []model.PickGoods
		pickRemarks    []model.PickRemark
		prePickIds     []string
		prePickGoodsIds,
		orderGoodsIds,
		prePickRemarksIds []int
	)

	err, prePickGoods = model.GetPrePickGoodsByPrePickIdAndStatus(db, form.Ids, model.PrePickGoodsStatusUnhandled)

	if err != nil {
		return err
	}

	if len(prePickGoods) == 0 {
		err = errors.New("商品数据未找到")
		return
	}

	//构造 prePickIds
	for _, good := range prePickGoods {
		prePickIds = append(prePickIds, strconv.Itoa(good.PrePickId))
	}

	//去重
	prePickIds = slice.UniqueSlice(prePickIds)

	tx := db.Begin()

	pick := model.Pick{
		WarehouseId:    form.WarehouseId,
		TaskId:         taskId,
		BatchId:        form.BatchId,
		PrePickIds:     strings.Join(prePickIds, ","),
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

	err = model.PickSave(tx, &pick)

	if err != nil {
		tx.Rollback()
		return
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

	err = model.PickGoodsSave(tx, &pickGoods)

	if err != nil {
		tx.Rollback()
		return
	}

	//更新预拣货池商品相关数据状态
	err = model.UpdatePrePickGoodsStatusByIds(tx, prePickGoodsIds, model.PrePickGoodsStatusProcessing)

	if err != nil {
		tx.Rollback()
		return
	}

	prePickIdSlice := []int{}

	//预拣池内商品全部进入拣货池时 更新 对应的 预拣池状态
	if form.Type == 1 { //全单拣货
		prePickIdSlice = form.Ids
	} else {
		//0:未处理,1:已进入拣货池 -前面修改过pre_pick_goods 状态了，重新查询
		err, prePickGoods = model.GetPrePickGoodsByPrePickIdAndStatus(db, form.Ids, model.PrePickGoodsStatusUnhandled)

		if err != nil {
			tx.Rollback()
			return
		}

		//将传过来的id转换成map
		idsMp := slice.SliceToMap(form.Ids)

		//去除未处理的预拣池id
		for _, good := range prePickGoods {
			delete(idsMp, good.PrePickId)
		}

		//将map转回切片
		prePickIdSlice = slice.MapToSlice(idsMp)
	}

	if len(prePickIdSlice) > 0 {
		err = model.UpdatePrePickStatusByIds(tx, prePickIdSlice, model.PrePickStatusProcessing)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	err, prePickRemarks = model.GetPrePickRemarkByOrderGoodsIds(db, orderGoodsIds)

	if err != nil {
		tx.Rollback()
		return
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

		err = model.PickRemarkBatchSave(tx, &pickRemarks)

		if err != nil {
			tx.Rollback()
			return
		}

		//更新预拣货池备注相关数据状态
		err = model.UpdatePrePickRemarkByIds(tx, prePickRemarksIds, map[string]interface{}{"status": model.PrePickRemarkStatusProcessing})
		if err != nil {
			return err
		}
	}

	tx.Commit()

	return nil
}

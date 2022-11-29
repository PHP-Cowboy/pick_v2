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

func CheckTyp(typ, distributionType int) bool {
	return typ == model.ExpressDeliveryBatchTyp && distributionType == model.DistributionTypeCourier
}

// 生成预拣池数据逻辑
func CreatePrePickLogic(
	db *gorm.DB,
	form req.NewCreateBatchForm,
	claims *middlewares.CustomClaims,
	batchId int,
) (
	err error,
	orderGoodsIds []int,
	outboundGoods []model.OutboundGoods,
	outboundGoodsJoinOrder []model.OutboundGoodsJoinOrder,
	prePickIds []int,
) {

	mp, err := cache.GetClassification()

	if err != nil {
		return
	}

	//缓存中的线路数据
	lineCacheMp, err := cache.GetShopLine()

	if err != nil {
		return
	}

	err, outboundGoodsJoinOrder = model.GetOutboundGoodsJoinOrderList(db, form.TaskId, form.Number)

	if err != nil {
		return
	}

	var (
		prePicks      []model.PrePick
		prePickGoods  []model.PrePickGoods
		prePickRemark []model.PrePickRemark
		shopNumMp     = make(map[int]struct{}, 0) //店铺
	)

	for _, og := range outboundGoodsJoinOrder {
		//todo 快递批次上线后删除 form.Typ == model.ExpressDeliveryBatchTyp &&
		if form.Typ == model.ExpressDeliveryBatchTyp && !CheckTyp(form.Typ, og.DistributionType) {
			err = ecode.RegularOrExpressDeliveryBatchInvalid
			return
		}
		//线路
		cacheMpLine, cacheMpOk := lineCacheMp[og.ShopId]

		if !cacheMpOk {
			err = errors.New("店铺:" + og.ShopName + "线路未同步，请先同步")
			return
		}

		//拣货池根据门店合并订单数据
		_, shopMpOk := shopNumMp[og.ShopId]

		if shopMpOk {
			continue
		}

		//生成预拣池数据
		prePicks = append(prePicks, model.PrePick{
			WarehouseId: claims.WarehouseId,
			BatchId:     batchId,
			ShopId:      og.ShopId,
			ShopCode:    og.ShopCode,
			ShopName:    og.ShopName,
			Line:        cacheMpLine,
			Status:      0,
			Typ:         form.Typ,
		})

		//门店进入生成预拣池数据后，map赋值，下次不再生成预拣池数据
		shopNumMp[og.ShopId] = struct{}{}
	}

	//预拣池数量

	if len(prePicks) == 0 {
		err = ecode.NoOrderFound
		return
	}

	//预拣池任务数据保存
	err, prePicksRes := model.PrePickBatchSave(db, prePicks)
	if err != nil {
		return
	}

	//单次创建批次时，预拣池表数据根据shop_id唯一
	//构造 map[shop_id]pre_pick_id 更新 t_pre_pick_goods t_pre_pick_remark 表 pre_pick_id
	prePickMp := make(map[int]int, 0)
	for _, pp := range prePicksRes {
		prePickMp[pp.ShopId] = pp.Id
		//预拣池ID
		prePickIds = append(prePickIds, pp.Id)
	}

	for _, og := range outboundGoodsJoinOrder {
		//用于更新 t_order_goods 的 batch_id
		orderGoodsIds = append(orderGoodsIds, og.OrderGoodsId)

		prePickId, ppMpOk := prePickMp[og.ShopId]

		if !ppMpOk {
			err = errors.New(fmt.Sprintf("店铺ID: %v 无预拣池数据", og.ShopId))
			return
		}

		//商品类型
		goodsType, mpOk := mp[og.GoodsType]

		if !mpOk {
			err = errors.New("商品类型:" + og.GoodsType + "数据未同步")
			return
		}

		//线路
		cacheMpLine, cacheMpOk := lineCacheMp[og.ShopId]

		if !cacheMpOk {
			err = errors.New("店铺:" + og.ShopName + "线路未同步，请先同步")
			return
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
			Typ:              form.Typ,
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
				Typ:          form.Typ,
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
		return
	}

	//预拣池商品备注，可能没有备注
	if len(prePickRemark) > 0 {
		err = model.PrePickRemarkBatchSave(db, prePickRemark)
		if err != nil {
			return
		}
	}

	return
}

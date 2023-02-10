package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"strconv"
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
	outboundGoodsJoinOrder []model.GoodsJoinOrder,
	prePickIds []int,
	prePicks []model.PrePick,
	prePickGoods []model.PrePickGoods,
	prePickRemarks []model.PrePickRemark,
) {

	classMp, err := cache.GetClassification()

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
		shopNumMp = make(map[int]struct{}, 0) //店铺
		prePickStatus,
		prePickGoodsStatus,
		prePickRemarkStatus int
		picks       []model.Pick
		pickGoods   []model.PickGoods
		pickRemarks []model.PickRemark
	)

	if form.Typ == 1 { // 常规批次，所有的状态都是待处理
		prePickStatus = model.PrePickStatusUnhandled
		prePickGoodsStatus = model.PrePickGoodsStatusUnhandled
		prePickRemarkStatus = model.PrePickRemarkStatusUnhandled
	} else { //快递批次，所有的状态都是已处理，直接进入预拣池和拣货池
		prePickStatus = model.PrePickStatusProcessing
		prePickGoodsStatus = model.PrePickGoodsStatusProcessing
		prePickRemarkStatus = model.PrePickRemarkStatusProcessing
	}

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
			TaskId:      form.TaskId,
			BatchId:     batchId,
			ShopId:      og.ShopId,
			ShopCode:    og.ShopCode,
			ShopName:    og.ShopName,
			Line:        cacheMpLine,
			Status:      prePickStatus,
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
	err = model.PrePickBatchSave(db, &prePicks)
	if err != nil {
		return
	}

	//单次创建批次时，预拣池表数据根据shop_id唯一
	//构造 map[shop_id]pre_pick_id 更新 t_pre_pick_goods t_pre_pick_remark 表 pre_pick_id
	prePickMp := make(map[int]int, 0)
	for _, pp := range prePicks {
		prePickMp[pp.ShopId] = pp.Id
		//预拣池ID
		prePickIds = append(prePickIds, pp.Id)
		if form.Typ == 2 {
			picks = append(picks, model.Pick{
				WarehouseId:     pp.WarehouseId,
				TaskId:          pp.TaskId,
				BatchId:         pp.BatchId,
				PrePickIds:      strconv.Itoa(pp.Id),
				TaskName:        "",
				ShopCode:        pp.ShopCode,
				ShopName:        pp.ShopName,
				Line:            pp.Line,
				Num:             0,
				PrintNum:        0,
				PickUser:        "",
				TakeOrdersTime:  nil,
				ReviewUser:      "",
				ReviewTime:      nil,
				Sort:            0,
				Version:         0,
				Status:          1,
				OutboundType:    1,
				DeliveryNo:      "",
				Typ:             pp.Typ,
				DeliveryOrderNo: nil,
			})
		}
	}

	var prePickIdAndPickIdMp = make(map[string]int, 0)

	//拣货池任务数据保存
	if form.Typ == 2 {
		err = model.PickBatchSave(db, &picks)
		if err != nil {
			return
		}

		//构造预拣池id和拣货池id对应关系map
		for _, pk := range picks {
			prePickIdAndPickIdMp[pk.PrePickIds] = pk.Id
		}
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
		goodsType, classMpOk := classMp[og.GoodsType]

		if !classMpOk {
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

		//构造预拣池数据
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
			Status:           prePickGoodsStatus,
			Typ:              form.Typ,
		})

		//如果有备注，构造预拣池备注数据
		if og.GoodsRemark != "" || og.OrderRemark != "" {
			prePickRemarks = append(prePickRemarks, model.PrePickRemark{
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
				Status:       prePickRemarkStatus,
				Typ:          form.Typ,
			})
		}

		//构造 更新 t_outbound_goods 表 status 状态 为拣货中
		outboundGoods = append(outboundGoods, model.OutboundGoods{
			TaskId:  og.TaskId,
			Number:  og.Number,
			Sku:     og.Sku,
			BatchId: batchId,
			Status:  model.OutboundGoodsStatusPicking,
		})

	}

	//预拣池商品
	err = model.PrePickGoodsBatchSave(db, &prePickGoods)
	if err != nil {
		return
	}

	if form.Typ == 2 {

		for _, og := range prePickGoods {
			prePickId, ppMpOk := prePickMp[og.ShopId]

			if !ppMpOk {
				err = errors.New(fmt.Sprintf("店铺ID: %v 无预拣池数据", og.ShopId))
				return
			}
			pickId, prePickIdAndPickIdMpOk := prePickIdAndPickIdMp[strconv.Itoa(prePickId)]

			if !prePickIdAndPickIdMpOk {
				err = errors.New("拣货池id有误")
				return
			}

			needNum := og.NeedNum

			pickGoods = append(pickGoods, model.PickGoods{
				WarehouseId:      claims.WarehouseId,
				PickId:           pickId,
				BatchId:          og.BatchId,
				PrePickGoodsId:   og.Id,
				OrderGoodsId:     og.OrderGoodsId,
				Number:           og.Number,
				ShopId:           og.ShopId,
				DistributionType: og.DistributionType,
				Sku:              og.Sku,
				GoodsName:        og.GoodsName,
				GoodsType:        og.GoodsType,
				GoodsSpe:         og.GoodsSpe,
				Shelves:          og.Shelves,
				DiscountPrice:    og.DiscountPrice,
				NeedNum:          needNum,
				CompleteNum:      needNum,
				ReviewNum:        needNum,
				CloseNum:         og.CloseNum,
				Status:           1,
				Unit:             og.Unit,
			})

		}

		//拣货池商品
		err = model.PickGoodsSave(db, &pickGoods)
		if err != nil {
			return
		}

	}

	//预拣池商品备注，可能没有备注
	if len(prePickRemarks) > 0 {
		err = model.PrePickRemarkBatchSave(db, &prePickRemarks)
		if err != nil {
			return
		}

		if form.Typ == 2 {
			for _, remark := range prePickRemarks {

				prePickId, ppMpOk := prePickMp[remark.ShopId]

				if !ppMpOk {
					err = errors.New(fmt.Sprintf("店铺ID: %v 无预拣池数据", remark.ShopId))
					return
				}

				pickId, prePickIdAndPickIdMpOk := prePickIdAndPickIdMp[strconv.Itoa(prePickId)]

				if !prePickIdAndPickIdMpOk {
					err = errors.New("拣货池id有误")
					return
				}

				pickRemarks = append(pickRemarks, model.PickRemark{
					WarehouseId:     claims.WarehouseId,
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

			err = model.PickRemarkBatchSave(db, &pickRemarks)
			if err != nil {
				return
			}
		}
	}

	return
}

// 预拣货明细
func GetPrePickDetail(db *gorm.DB, form req.GetPrePickDetailForm) (err error, res rsp.GetPrePickDetailRsp) {

	var (
		prePick       model.PrePick
		prePickGoods  []model.PrePickGoods
		prePickRemark []model.PrePickRemark
	)

	err, prePick = model.GetPrePickByPk(db, form.PrePickId)

	if err != nil {
		return
	}

	res.TaskName = prePick.ShopName
	res.Line = prePick.Line

	err, prePickGoods = model.GetPrePickGoodsByPrePickIdAndStatus(db, []int{form.PrePickId}, model.PrePickGoodsStatusUnhandled)

	if err != nil {
		return
	}

	prePickGoodsSkuMp := make(map[string]rsp.MergePrePickGoods, 0)

	goodsNum := 0

	orderNumMp := make(map[string]struct{}, 0)

	//相同sku合并处理
	for _, goods := range prePickGoods {

		orderNumMp[goods.Number] = struct{}{}

		goodsNum += goods.NeedNum

		val, ok := prePickGoodsSkuMp[goods.Sku]

		paramsId := rsp.ParamsId{
			PickGoodsId:  goods.Id,
			OrderGoodsId: goods.OrderGoodsId,
		}

		if !ok {

			prePickGoodsSkuMp[goods.Sku] = rsp.MergePrePickGoods{
				Id:        goods.Id,
				Sku:       goods.Sku,
				GoodsName: goods.GoodsName,
				GoodsType: goods.GoodsType,
				GoodsSpe:  goods.GoodsSpe,
				Shelves:   goods.Shelves,
				NeedNum:   goods.NeedNum,
				CloseNum:  goods.CloseNum,
				Unit:      goods.Unit,
				ParamsId:  []rsp.ParamsId{paramsId},
			}

		} else {
			val.CloseNum += goods.CloseNum
			val.NeedNum += goods.NeedNum
			val.ParamsId = append(val.ParamsId, paramsId)
			prePickGoodsSkuMp[goods.Sku] = val
		}
	}

	//订单数
	res.OrderNum = len(orderNumMp)

	//商品数
	res.GoodsNum = goodsNum

	goodsMap := make(map[string][]rsp.MergePrePickGoods, 0)

	for _, goods := range prePickGoodsSkuMp {

		goodsMap[goods.GoodsType] = append(goodsMap[goods.GoodsType], rsp.MergePrePickGoods{
			Id:        goods.Id,
			Sku:       goods.Sku,
			GoodsName: goods.GoodsName,
			GoodsType: goods.GoodsType,
			GoodsSpe:  goods.GoodsSpe,
			Shelves:   goods.Shelves,
			NeedNum:   goods.NeedNum,
			CloseNum:  goods.CloseNum,
			Unit:      goods.Unit,
			ParamsId:  goods.ParamsId,
		})
	}

	res.Goods = goodsMap

	err, prePickRemark = model.GetPrePickRemarkByPrePickId(db, form.PrePickId)

	if err != nil {
		return
	}

	remarkMap := make(map[string]rsp.Remark, 0)

	list := []rsp.Remark{}
	for _, remark := range prePickRemark {
		remarkMap[remark.Number] = rsp.Remark{
			Number:      remark.Number,
			OrderRemark: remark.OrderRemark,
			GoodsRemark: remark.GoodsRemark,
		}
	}

	for n := range orderNumMp {
		remark, remarkMapOk := remarkMap[n]

		if !remarkMapOk {
			list = append(list, rsp.Remark{
				Number:      n,
				OrderRemark: "",
				GoodsRemark: "",
			})
		} else {
			list = append(list, remark)
		}
	}

	res.RemarkList = list

	return
}

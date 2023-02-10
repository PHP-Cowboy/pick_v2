package dao

import (
	"errors"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/slice"
	"sort"
	"time"
)

// 后台拣货数据列表
func AdminPickList(form req.AdminPickListReq) (err error, res rsp.AdminPickListRsp) {
	var (
		pickGoodsList []model.PickGoods

		orderMp = make(map[int]rsp.AdminPickList)

		shopIds   []int
		shopMp    = make(map[int]model.Shop, 0)
		shopList  []model.Shop
		shopNumMp = make(map[int][]string, 0)
	)

	db := global.DB

	err, pickGoodsList = model.GetPickGoodsList(db, &model.PickGoods{BatchId: form.BatchId, GoodsName: form.GoodsName, GoodsType: form.GoodsType})

	if err != nil {
		return
	}

	for _, goods := range pickGoodsList {
		shopIds = append(shopIds, goods.ShopId)

		_, shopNumMpOk := shopNumMp[goods.ShopId]

		if !shopNumMpOk {
			shopNumMp[goods.ShopId] = make([]string, 0)
		}

		shopNumMp[goods.ShopId] = append(shopNumMp[goods.ShopId], goods.Number)
	}

	err, shopList = model.GetShopListByIds(db, "shop_name,shop_code", shopIds)

	if err != nil {
		return
	}

	for _, shop := range shopList {
		shopMp[shop.ShopId] = shop
	}

	for _, goods := range pickGoodsList {
		adminPick, orderMpOk := orderMp[goods.ShopId]

		needNum := goods.NeedNum

		shop, shopMpOk := shopMp[goods.ShopId]

		if !shopMpOk {
			shop = model.Shop{}
		}

		if !orderMpOk {
			shopNum, shopNumMpOk := shopNumMp[goods.ShopId]

			if !shopNumMpOk {
				err = errors.New("门店信息有误")
				return
			}

			orderNum := len(slice.UniqueSlice(shopNum))

			adminPick = rsp.AdminPickList{
				ShopId:   goods.ShopId,
				ShopCode: shop.ShopCode,
				ShopName: shop.ShopName,
				OrderNum: orderNum,
				NeedNum:  needNum,
				Remark:   "",
				Num:      needNum,
			}
		} else {
			adminPick.NeedNum += needNum
		}

		orderMp[goods.ShopId] = adminPick
	}

	list := make([]rsp.AdminPickList, 0)

	for _, goods := range pickGoodsList {
		adminPick, orderMpOk := orderMp[goods.ShopId]

		if !orderMpOk {
			err = errors.New("")
			return
		}

		list = append(list, adminPick)
	}

	res.Total = len(list)

	res.List = list

	return
}

// 后台拣货数据详情
func AdminPickDetail(form req.AdminPickDetailReq) (err error, res rsp.AdminPickDetailRsp) {
	db := global.DB

	var (
		batch      model.Batch
		picks      []model.Pick
		pickGoods  []model.PickGoods
		pickRemark []model.PickRemark
		pickIds    []int
	)

	err, batch = model.GetBatchByPk(db, form.BatchId)
	if err != nil {
		return
	}

	err, picks = model.GetPickList(db, &model.Pick{BatchId: form.BatchId, ShopName: form.ShopName})

	if err != nil {
		return
	}

	for _, pick := range picks {
		pickIds = append(pickIds, pick.Id)
	}

	query := "pick_id,count(distinct(shop_id)) as shop_num,count(distinct(number)) as order_num,sum(need_num) as need_num"

	err, numsMp := model.CountPickPoolNumsByPickIds(db, pickIds, query)

	if err != nil {
		return
	}

	for _, nums := range numsMp {
		res.ShopNum += nums.ShopNum
		res.OrderNum += nums.OrderNum
		res.GoodsNum = +nums.NeedNum
	}

	res.BatchId = form.BatchId
	res.TaskName = batch.BatchName
	res.ShopCode = form.ShopCode

	now := time.Now()

	res.PickUser = "新时沏1号"
	res.TakeOrdersTime = (*model.MyTime)(&now)

	err, pickGoods = model.GetPickGoodsByPickIds(db, pickIds)

	if err != nil {
		return
	}

	orderNumbers := make([]string, 0, res.OrderNum)

	pickGoodsSkuMp := make(map[string]rsp.MergePickGoods, 0)
	//相同sku合并处理
	for _, goods := range pickGoods {
		orderNumbers = append(orderNumbers, goods.Number)

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

	//去重
	orderNumbers = slice.UniqueSlice(orderNumbers)

	goodsMap := make(map[string][]rsp.MergePickGoods, 0)

	for _, goods := range pickGoodsSkuMp {

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

	//按货架号排序
	for s, goods := range goodsMap {

		ret := rsp.MyMergePickGoods(goods)

		sort.Sort(ret)

		goodsMap[s] = ret
	}

	res.Goods = goodsMap

	err, pickRemark = model.GetPickRemarkListByPickIds(db, pickIds)

	if err != nil {
		return
	}

	remarkMap := make(map[string]rsp.PickRemark, 0)

	list := []rsp.PickRemark{}
	for _, remark := range pickRemark {
		remarkMap[remark.Number] = rsp.PickRemark{
			Number:      remark.Number,
			OrderRemark: remark.OrderRemark,
			GoodsRemark: remark.GoodsRemark,
		}
	}

	for _, n := range orderNumbers {
		remark, remarkMapOk := remarkMap[n]

		if !remarkMapOk {
			list = append(list, rsp.PickRemark{
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

// 批次门店商品列表
func BatchShopGoodsList(form req.BatchShopGoodsListReq) (err error, list []rsp.BatchShopGoodsList) {
	db := global.DB

	var (
		pickGoodsList []model.PickGoods
		goodsMp       = make(map[string]model.PickGoods, 0)
	)

	err, pickGoodsList = model.GetPickGoodsList(db, &model.PickGoods{BatchId: form.BatchId, ShopId: form.ShopId, Status: 1})
	if err != nil {
		return
	}

	list = make([]rsp.BatchShopGoodsList, 0)

	skus := make([]string, 0, len(pickGoodsList))

	for _, pg := range pickGoodsList {
		skus = append(skus, pg.Sku)

		goods, goodsMpOk := goodsMp[pg.Sku]

		if !goodsMpOk {
			goods = pg
		} else {
			goods.CompleteNum += pg.CompleteNum
			goods.ReviewNum += pg.ReviewNum
		}

		goodsMp[pg.Sku] = goods
	}

	skus = slice.UniqueSlice(skus)

	for _, s := range skus {
		goods, goodsMpOk := goodsMp[s]

		if !goodsMpOk {
			err = errors.New("sku数据有误")
			return
		}

		list = append(list, rsp.BatchShopGoodsList{
			Sku:         goods.Sku,
			GoodsName:   goods.GoodsName,
			CompleteNum: goods.CompleteNum,
			ReviewNum:   goods.ReviewNum,
		})
	}

	return
}

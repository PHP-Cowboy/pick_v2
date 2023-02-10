package dao

import (
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/slice"
	"strconv"
)

// 货品汇总单
func GoodsSummaryList(db *gorm.DB, form req.GoodsSummaryListForm) (err error, mp map[string]map[string]string, column, shopCodes []string) {

	var list []model.PrePickGoodsJoinPrePick

	err, list = model.GetPrePickGoodsJoinPrePickListByBatchId(db, form.BatchId, []string{})

	for _, l := range list {
		shopCodes = append(shopCodes, l.ShopCode)
	}

	shopCodes = slice.UniqueSlice(shopCodes)

	mp = make(map[string]map[string]string, 0)

	mpSum := make(map[string]int, 0)

	for _, pg := range list {

		subMp, ok := mp[pg.Sku]

		if !ok {
			subMp = make(map[string]string, 0)
			mp[pg.Sku] = subMp
		}

		for _, code := range shopCodes {
			_, has := subMp[code]
			if !has {
				subMp[code] = "0"
			}

			if code == pg.ShopCode {
				subMp[code] = strconv.Itoa(pg.NeedNum)
			}
		}

		subMp["商品名称"] = pg.GoodsName
		subMp["规格"] = pg.GoodsSpe
		subMp["单位"] = pg.Unit

		_, msOk := mpSum[pg.Sku]
		if !msOk {
			mpSum[pg.Sku] = pg.NeedNum
		} else {
			mpSum[pg.Sku] += pg.NeedNum
		}

		subMp["总计"] = strconv.Itoa(mpSum[pg.Sku])

		mp[pg.Sku] = subMp
	}

	column = []string{"商品名称", "规格", "单位", "总计"}

	shopCodes = slice.UniqueSlice(shopCodes)

	column = append(column, shopCodes...)

	return
}

func ShopAddress(form req.ShopAddressReq) (err error, mp map[int]map[string]string) {
	db := global.DB

	var (
		picks     []model.Pick
		pickGoods []model.PickGoods
		orders    []model.Order
		pickIds   []int
		numbers   []string
	)
	mp = make(map[int]map[string]string, 0)

	err, picks = model.GetPickList(db, &model.Pick{BatchId: form.BatchId})
	if err != nil {
		return
	}

	for _, pick := range picks {
		pickIds = append(pickIds, pick.Id)
	}

	err, pickGoods = model.GetPickGoodsByPickIds(db, pickIds)
	if err != nil {
		return
	}

	for _, good := range pickGoods {
		numbers = append(numbers, good.Number)
	}

	numbers = slice.UniqueSlice(numbers)

	err, orders = model.GetOrderListByNumbers(db, numbers)
	if err != nil {
		return
	}

	for _, o := range orders {
		val, mpOk := mp[o.ShopId]

		if !mpOk {
			val = make(map[string]string, 0)

			val["consignee_name"] = o.ConsigneeName
			val["consignee_tel"] = o.ConsigneeTel
			val["province"] = o.Province
			val["city"] = o.City
			val["district"] = o.District
			val["address"] = o.Address
		}

		mp[o.ShopId] = val
	}

	return
}

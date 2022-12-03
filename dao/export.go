package dao

import (
	"gorm.io/gorm"
	"pick_v2/model"
	"pick_v2/utils/slice"
	"strconv"
)

// 货品汇总单
func GoodsSummaryList(db *gorm.DB, batchId int) (err error, mp map[string]map[string]string, column, shopCodes []string) {

	var list []model.PrePickGoodsJoinPrePick

	err, list = model.GetPrePickGoodsJoinPrePickListByBatchId(db, batchId)

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

package dao

import (
	"pick_v2/common/constant"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"strconv"
	"time"
)

func Print(form req.PrintCallGetReq) (err error, ret []rsp.PrintCallGetRsp) {

	printCh := GetPrintJobMap(form.HouseCode, form.Typ)

	//通道中没有任务
	if printCh == nil {
		return
	}

	var (
		pick      model.Pick
		pickGoods []model.PickGoods
	)

	db := global.DB

	result := db.Model(&model.Pick{}).Where("delivery_no = ?", printCh.DeliveryOrderNo).Find(&pick)

	if result.Error != nil {
		err = result.Error
		return
	}

	result = db.Model(&model.PickGoods{}).Where("pick_id = ? and shop_id = ? and review_num > 0", pick.Id, printCh.ShopId).Find(&pickGoods)

	if result.Error != nil {
		err = result.Error
		return
	}

	length := len(pickGoods) //有多少条pickGoods就有多少条OrderInfo数据，map数也是

	orderGoodsIds := make([]int, 0, length)

	goodsMp := make(map[int]model.PickGoods, length)

	for _, good := range pickGoods {
		orderGoodsIds = append(orderGoodsIds, good.OrderGoodsId)

		goodsMp[good.OrderGoodsId] = good
	}

	err, orderJoinGoods := model.GetOrderGoodsJoinOrderByIdsNoSort(db, orderGoodsIds)
	if err != nil {
		return
	}

	var (
		compOrderJoinGoods []model.GoodsJoinOrder
	)

	err, compOrderJoinGoods = model.GetCompleteOrderJoinGoodsByOrderGoodsId(db, orderGoodsIds)

	if err != nil {
		return
	}

	for _, good := range compOrderJoinGoods {
		orderJoinGoods = append(orderJoinGoods, good)
	}

	if len(orderJoinGoods) <= 0 {
		err = ecode.OrderDataNotFound
		return
	}

	packages := pick.Num

	if pick.ShopCode == "" {
		pick.ShopCode = "未设置店编"
	}

	//合并单不打印，ShopCode = "MergePickingTasks" 说明是合并单，合并单不打印箱单
	if pick.ShopCode == "MergePickingTasks" {
		packages = 0
	}

	item := rsp.PrintCallGetRsp{
		ShopName:    pick.ShopName,
		JHNumber:    strconv.Itoa(pick.Id),
		PickName:    pick.PickUser, //拣货人
		ShopType:    orderJoinGoods[0].ShopType,
		CheckName:   pick.ReviewUser,                                              //复核员
		HouseName:   TransferHouse(orderJoinGoods[0].HouseCode),                   //TransferHouse(info.HouseCode)
		Delivery:    TransferDistributionType(orderJoinGoods[0].DistributionType), //TransferDistributionType(info.DistributionType),
		OrderRemark: orderJoinGoods[0].OrderRemark,
		Consignee:   orderJoinGoods[0].ConsigneeName, //info.ConsigneeName
		Shop_code:   pick.ShopCode,
		Packages:    packages,
		Phone:       orderJoinGoods[0].ConsigneeTel, //info.ConsigneeTel,
		PriType:     printCh.Type,
	}

	if orderJoinGoods[0].ShopCode != "" {
		item.ShopName = orderJoinGoods[0].ShopCode + "--" + orderJoinGoods[0].ShopName
	}

	item2Mp := make(map[string]rsp.CallGetGoodsView, 0)

	for _, info := range orderJoinGoods {

		pgs, ok := goodsMp[info.OrderGoodsId]

		if !ok {
			continue
		}

		item2val, item2ok := item2Mp[info.Number]

		if !item2ok {
			item2val = rsp.CallGetGoodsView{
				SaleNumber:  info.Number,
				Date:        timeutil.FormatToDateTime(time.Time(*info.PayAt)),
				OrderRemark: info.OrderRemark,
			}
		}

		var lackCount int

		if info.LackCount > 0 {
			lackCount = info.LackCount - pgs.ReviewNum
		}

		item3 := rsp.CallGetGoods{
			GoodsName:    info.GoodsName,
			GoodsSpe:     info.GoodsSpe,
			GoodsCount:   info.PayCount,
			RealOutCount: pgs.ReviewNum,
			GoodsUnit:    info.GoodsUnit,
			Price:        int64(info.DiscountPrice) * int64(pgs.ReviewNum),
			LackCount:    lackCount,
		}
		item2val.List = append(item2val.List, item3)

		item2Mp[info.Number] = item2val
	}

	for _, item2 := range item2Mp {
		item.GoodsList = append(item.GoodsList, item2)
	}

	ret = make([]rsp.PrintCallGetRsp, 0, 1)

	ret = append(ret, item)

	return
}

func TransferHouse(s string) string {
	switch s {
	case constant.JH_HUOSE_CODE:
		return constant.JH_HUOSE_NAME
	default:
		return constant.OT_HUOSE_NAME
	}
}

func TransferDistributionType(t int) (method string) {
	switch t {
	case 1:
		method = "公司配送"
		break
	case 2:
		method = "用户自提"
		break
	case 3:
		method = "三方物流"
		break
	case 4:
		method = "快递配送"
		break
	case 5:
		method = "首批物料|设备单"
		break
	default:
		method = "其他"
		break
	}

	return method
}

package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/slice"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
)

// 拣货单列表
func PickOrderList(c *gin.Context) {
	var form req.PickOrderListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var (
		pickOrder      []model.PickOrder
		pickOrderGoods []model.PickOrderGoods
		numbers        []string
	)

	db := global.DB

	if form.Sku != "" {
		result := db.Where("sku = ?", form.Sku).Find(&pickOrderGoods)

		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, good := range pickOrderGoods {
			numbers = append(numbers, good.Number)
		}
	}

	localDb := db.Model(&model.PickOrder{})

	if len(numbers) > 0 {
		numbers = slice.UniqueStringSlice(numbers)
		localDb = localDb.Where("number in (?)", numbers)
	}

	localDb.Where(&model.PickOrder{
		OrderType:        form.OrderType,
		ShopId:           form.ShopId,
		Number:           form.Number,
		Line:             form.Lines,
		DistributionType: form.DistributionType,
		ShopType:         form.ShopType,
		Province:         form.Province,
		City:             form.City,
		District:         form.District,
		HasRemark:        form.HasRemark,
	})

	if form.PayEndTime != "" {
		localDb = localDb.Where("pay_at >= ?", form.PayEndTime)
	}

	var (
		total int64
		res   rsp.PickOrderGoodsListRsp
	)

	result := localDb.Count(&total)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = localDb.Scopes(model.Paginate(form.Page, form.Size)).Find(&pickOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.PickOrder, 0, form.Size)

	for _, o := range pickOrder {
		latestPickingTime := ""

		if o.LatestPickingTime != nil {
			latestPickingTime = o.LatestPickingTime.Format(timeutil.TimeFormat)
		}

		list = append(list, rsp.PickOrder{
			Number:            o.Number,
			PayAt:             o.PayAt,
			ShopCode:          o.ShopCode,
			ShopName:          o.ShopName,
			ShopType:          o.ShopType,
			DistributionType:  o.DistributionType,
			PayCount:          o.PayTotal,
			LimitNum:          o.LimitNum,
			Line:              o.Line,
			Region:            o.Province + o.City + o.District,
			OrderRemark:       o.OrderRemark,
			OrderType:         o.OrderType,
			LatestPickingTime: latestPickingTime,
		})
	}

	res.Total = total
	res.List = list

	xsq_net.SucJson(c, res)
}

// 拣货单明细
func GetPickOrderDetail(c *gin.Context) {

	var form req.GetOrderDetailForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		orders     model.PickOrder
		orderGoods []model.PickOrderGoods
	)

	db := global.DB

	result := db.Model(&model.PickOrder{}).Where("number = ?", form.Number).First(&orders)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.PickOrderGoods{}).Where("number = ?", form.Number).Find(&orderGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	mp, err := cache.GetClassification()

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var res rsp.OrderDetail

	detailMap := make(map[string]*rsp.Detail, 0)

	deliveryOrderNoArr := make([]string, 0)

	for _, og := range orderGoods {

		deliveryOrderNoArr = append(deliveryOrderNoArr, og.DeliveryOrderNo...)

		goodsType, ok := mp[og.GoodsType]

		if !ok {
			xsq_net.ErrorJSON(c, errors.New("商品类型:"+og.GoodsType+"数据未同步"))
			return
		}

		if _, detailOk := detailMap[goodsType]; !detailOk {
			detailMap[goodsType] = &rsp.Detail{
				Total: 0,
				List:  make([]*rsp.GoodsDetail, 0),
			}
		}

		detailMap[goodsType].Total += og.PayCount

		detailMap[goodsType].List = append(detailMap[goodsType].List, &rsp.GoodsDetail{
			Name:        og.GoodsName,
			GoodsSpe:    og.GoodsSpe,
			Shelves:     og.Shelves,
			PayCount:    og.PayCount,
			CloseCount:  og.CloseCount,
			LackCount:   og.LimitNum, //需拣数 以限发数为准
			GoodsRemark: og.GoodsRemark,
		})
	}

	res.Number = orders.Number
	res.PayAt = orders.PayAt
	res.ShopCode = orders.ShopCode
	res.ShopName = orders.ShopName
	res.Line = orders.Line
	res.Region = orders.Province + orders.City + orders.District
	res.ShopType = orders.ShopType
	res.OrderRemark = orders.OrderRemark

	res.Detail = detailMap

	deliveryOrderNoArr = slice.UniqueStringSlice(deliveryOrderNoArr)
	//历史出库单号
	res.DeliveryOrderNo = deliveryOrderNoArr

	xsq_net.SucJson(c, res)
}

package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"pick_v2/dao"

	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/request"
	"pick_v2/utils/slice"
	"pick_v2/utils/xsq_net"
)

// 订单列表
func GetGoodsList(c *gin.Context) {
	var form req.GoodsListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		orders     []model.Order
		orderGoods []model.OrderGoods
		numbers    []string
	)

	db := global.DB

	if form.Sku != "" {
		result := db.Where("sku = ?", form.Sku).Find(&orderGoods)

		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, good := range orderGoods {
			numbers = append(numbers, good.Number)
		}
	}

	localDb := db.Model(&model.Order{})

	if len(numbers) > 0 {
		localDb = localDb.Where("number in (?)", numbers)
	}

	localDb.Where(&model.Order{
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

	localDb.Where("order_type != 2") //不要拣货中的

	var (
		res   rsp.GoodsListRsp
		dbRes []OrderNum
		total int64
	)

	result := localDb.Count(&total)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	localDb.Order(fmt.Sprintf("pay_at %s", form.PayAtSort))

	if form.ShopCodeSort != "" {
		localDb.Order(fmt.Sprintf("shop_code  %s", form.ShopCodeSort))
	}

	result = localDb.Scopes(model.Paginate(form.Page, form.Size)).Find(&orders)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if len(numbers) == 0 {
		for _, o := range orders {
			numbers = append(numbers, o.Number)
		}
	}

	query := "number,sum(pay_count) as pay_count,sum(close_count) as close_count,sum(out_count) as out_count,sum(lack_count) as lack_count"

	err, numsMp := model.OrderGoodsNumsStatisticalByNumbers(db, query, numbers)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	list := make([]rsp.Order, 0, form.Size)

	for _, o := range orders {

		nums, numsOk := numsMp[o.Number]

		if !numsOk {
			xsq_net.ErrorJSON(c, errors.New("订单统计数量不存在"))
			return
		}

		list = append(list, rsp.Order{
			Number:            o.Number,
			PayAt:             o.PayAt,
			ShopCode:          o.ShopCode,
			ShopName:          o.ShopName,
			ShopType:          o.ShopType,
			DistributionType:  o.DistributionType,
			PayCount:          nums.PayCount,
			Picked:            nums.OutCount,
			UnPicked:          nums.LackCount,
			CloseNum:          nums.CloseCount,
			Line:              o.Line,
			Region:            o.Province + o.City + o.District,
			OrderRemark:       o.OrderRemark,
			OrderType:         o.OrderType,
			LatestPickingTime: o.LatestPickingTime,
		})
	}

	result = localDb.Model(&model.Order{}).Select("count(id) as count, order_type").Group("order_type").Find(&dbRes)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = total

	res.List = list

	xsq_net.SucJson(c, res)
}

// 订单明细
func GetOrderDetail(c *gin.Context) {

	var form req.GetOrderDetailForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		orders     model.Order
		orderGoods []model.OrderGoods
	)

	db := global.DB

	result := db.Model(&model.Order{}).Where("number = ?", form.Number).First(&orders)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.OrderGoods{}).Where("number = ?", form.Number).Find(&orderGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	classMp, err := cache.GetClassification()

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var res rsp.OrderDetail

	detailMap := make(map[string]*rsp.Detail, 0)

	deliveryOrderNoArr := make(model.GormList, 0)

	for _, og := range orderGoods {

		if og.LackCount == 0 {
			//是欠货单时只显示欠货的
			continue
		}

		deliveryOrderNoArr = append(deliveryOrderNoArr, og.DeliveryOrderNo...)

		goodsType, classMpOk := classMp[og.GoodsType]

		if !classMpOk {
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
			LackCount:   og.LackCount,
			OutCount:    og.OutCount,
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

	deliveryOrderNoArr = slice.UniqueSlice(deliveryOrderNoArr)
	//历史出库单号
	res.DeliveryOrderNo = deliveryOrderNoArr

	xsq_net.SucJson(c, res)
}

// 商品列表
func CommodityList(c *gin.Context) {
	var (
		form   req.CommodityListForm
		result rsp.CommodityListRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	err := request.Call("api/v1/remote/shop/sku", form, &result)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, result.Data)
}

// 订单出货记录
func OrderShippingRecord(c *gin.Context) {
	var (
		form req.OrderShippingRecordReq
	)

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, res := dao.OrderShippingRecord(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, res)
}

// 订单出货记录明细
func ShippingRecordDetail(c *gin.Context) {
	var (
		form req.ShippingRecordDetailReq
		//res  rsp.ShippingRecordDetailRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var pickGoods []model.PickGoods

	result := global.DB.Where("pick_id in (?) and review_num > 0", form.Id).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	type Goods struct {
		GoodsName string `json:"goods_name"`
		GoodsSpe  string `json:"goods_spe"`
		ReviewNum int    `json:"review_num"`
	}

	mp := make(map[string][]Goods, 0)

	for _, good := range pickGoods {

		mp[good.GoodsType] = append(mp[good.GoodsType], Goods{
			GoodsName: good.GoodsName,
			GoodsSpe:  good.GoodsSpe,
			ReviewNum: good.ReviewNum,
		})
	}

	xsq_net.SucJson(c, mp)
}

// 完成订单
func CompleteOrder(c *gin.Context) {
	var (
		form req.CompleteOrderForm
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, res := dao.CompleteOrder(global.DB, form)
	if err != nil {
		return
	}

	xsq_net.SucJson(c, res)
}

// 完成订单详情
func CompleteOrderDetail(c *gin.Context) {
	var (
		form req.CompleteOrderDetailReq
		res  rsp.CompleteOrderDetailRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		completeOrder       model.CompleteOrder
		completeOrderDetail []model.CompleteOrderDetail
	)

	db := global.DB

	result := db.Model(&model.CompleteOrder{}).Where("number = ?", form.Number).First(&completeOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Where("number = ?", form.Number).Find(&completeOrderDetail)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	classMp, err := cache.GetClassification()

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	detailMap := make(map[string]*rsp.Detail, 0)

	deliveryOrderNoArr := make(model.GormList, 0)

	for _, goods := range completeOrderDetail {

		goodsType, classMpOk := classMp[goods.GoodsType]

		if !classMpOk {
			xsq_net.ErrorJSON(c, errors.New("商品类型:"+goods.GoodsType+"数据未同步"))
			return
		}

		if _, detailOk := detailMap[goodsType]; !detailOk {
			detailMap[goodsType] = &rsp.Detail{
				Total: 0,
				List:  make([]*rsp.GoodsDetail, 0),
			}
		}

		detailMap[goodsType].Total += goods.PayCount

		detailMap[goodsType].List = append(detailMap[goodsType].List, &rsp.GoodsDetail{
			Name:        goods.GoodsName,
			GoodsSpe:    goods.GoodsSpe,
			Shelves:     goods.Shelves,
			PayCount:    goods.PayCount,
			CloseCount:  goods.CloseCount,
			OutCount:    goods.ReviewCount,
			GoodsRemark: goods.GoodsRemark,
		})

		deliveryOrderNoArr = append(deliveryOrderNoArr, goods.DeliveryOrderNo...)
	}

	deliveryOrderNoArr = slice.UniqueSlice(deliveryOrderNoArr)

	res.Goods = detailMap

	res.ShopName = completeOrder.ShopName
	res.ShopCode = completeOrder.ShopCode
	res.Line = completeOrder.Line
	res.Region = fmt.Sprintf("%s-%s-%s", completeOrder.Province, completeOrder.City, completeOrder.District)
	res.ShopType = completeOrder.ShopType
	res.Number = completeOrder.Number
	res.OrderRemark = completeOrder.OrderRemark
	res.DeliveryOrderNo = deliveryOrderNoArr

	xsq_net.SucJson(c, res)
}

type OrderNum struct {
	Count     int `json:"count"`
	OrderType int `json:"order_type"`
}

func Count(c *gin.Context) {

	var form req.CountFrom

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		orderGoods []model.OrderGoods
		numbers    []string
	)

	db := global.DB

	if form.Sku != "" {
		result := db.Where("sku = ?", form.Sku).Find(&orderGoods)

		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, good := range orderGoods {
			numbers = append(numbers, good.Number)
		}
	}

	localDb := db.Model(&model.Order{})

	if len(numbers) > 0 {
		localDb = localDb.Where("number in (?)", numbers)
	}

	localDb.Where(&model.Order{
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

	localDb.Where("order_type != 2") //不要拣货中的

	type OrderNum struct {
		Count     int `json:"count"`
		OrderType int `json:"order_type"`
	}

	var (
		dbRes []OrderNum
		res   rsp.CountRes
	)

	result := localDb.Model(&model.Order{}).Select("count(id) as count, order_type").Group("order_type").Find(&dbRes)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, r := range dbRes {
		switch r.OrderType {
		case 1: //1:新订单
			res.NewCount = r.Count
			break
		case 2: //2:拣货中
			continue //不统计拣货中的
			//res.PickCount = r.Count
			//break
		case 3: //3:欠货单
			res.OldCount = r.Count
			break
		case 4: //已关闭
			res.CloseCount = r.Count
			break
		}
		res.AllCount += r.Count
	}

	xsq_net.SucJson(c, res)
}

package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/model/batch"
	"pick_v2/model/order"
	"pick_v2/utils/cache"
	"pick_v2/utils/slice"
	"pick_v2/utils/timeutil"
	"time"

	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/utils/ecode"
	"pick_v2/utils/request"
	"pick_v2/utils/xsq_net"
)

// 订单列表
func GetGoodsList(c *gin.Context) {
	var form req.GoodsListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var (
		orders     []order.Order
		orderGoods []order.OrderGoods
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

	localDb := db.Model(&order.Order{})

	if len(numbers) > 0 {
		localDb = localDb.Where("number in (?)", numbers)
	}

	localDb.Where(&order.Order{
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
		res   rsp.GoodsListRsp
	)

	result := localDb.Count(&total)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = localDb.Scopes(model.Paginate(form.Page, form.Size)).Find(&orders)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.Order, 0, form.Size)

	for _, o := range orders {
		latestPickingTime := ""

		if o.LatestPickingTime != nil {
			latestPickingTime = o.LatestPickingTime.Format(timeutil.TimeFormat)
		}

		list = append(list, rsp.Order{
			Number:            o.Number,
			PayAt:             o.PayAt,
			ShopCode:          o.ShopCode,
			ShopName:          o.ShopName,
			ShopType:          o.ShopType,
			DistributionType:  o.DistributionType,
			PayCount:          o.PayTotal,
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

// 订单明细
func GetOrderDetail(c *gin.Context) {

	var form req.GetOrderDetailForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		orders     order.Order
		orderGoods []order.OrderGoods
	)

	db := global.DB

	result := db.Model(&order.Order{}).Where("number = ?", form.Number).First(&orders)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&order.OrderGoods{}).Where("number = ?", form.Number).Find(&orderGoods)

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

	for _, og := range orderGoods {
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
			LackCount:   og.LackCount,
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

	//欠货单 需要查询 历史出库单号
	if orders.OrderType == 3 { //订单类型:1:新订单,2:拣货中,3:欠货单
		var (
			pickGoods          []batch.PickGoods
			pick               []batch.Pick
			pickIds            []int
			deliveryOrderNoArr []string
		)

		dbRes := db.Where("number = ?", form.Number).Find(&pickGoods)

		if dbRes.Error != nil {
			xsq_net.ErrorJSON(c, dbRes.Error)
			return
		}

		//拣货id map 去重
		pickIdsMp := make(map[int]struct{}, len(pickGoods))

		for _, pg := range pickGoods {

			_, ok := pickIdsMp[pg.PickId]

			if ok {
				continue
			}

			//获取复核数量和需拣货数量不一致的
			if pg.NeedNum != pg.ReviewNum {
				pickIds = append(pickIds, pg.PickId)
				pickIdsMp[pg.PickId] = struct{}{}
			}
		}

		if len(pickIds) > 0 {
			dbRes = db.Where("id in (?)", pickIds).Find(&pick)
			if dbRes.Error != nil {
				xsq_net.ErrorJSON(c, dbRes.Error)
				return
			}

			for _, p := range pick {
				deliveryOrderNoArr = append(deliveryOrderNoArr, p.DeliveryOrderNo)
			}
		}

		res.DeliveryOrderNo = deliveryOrderNoArr
	}

	xsq_net.SucJson(c, res)
}

func RequestGoodsList(responseData interface{}) (rsp.ApiGoodsListRsp, error) {

	var result rsp.ApiGoodsListRsp

	body, err := request.Post("api/v1/remote/lack/list", responseData)

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return result, err
	}
	if result.Code != 200 {
		return result, errors.New(result.Msg)
	}

	return result, nil
}

// 商品列表
func CommodityList(c *gin.Context) {
	var result rsp.CommodityListRsp

	body, err := request.Post("api/v1/remote/shop/sku", nil)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	if result.Code != 200 {
		xsq_net.ErrorJSON(c, errors.New(result.Msg))
		return
	}

	xsq_net.SucJson(c, result.Data)
}

// 订单出货记录
func OrderShippingRecord(c *gin.Context) {
	var (
		form req.OrderShippingRecordReq
		res  rsp.OrderShippingRecordRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var pick []batch.Pick

	result := global.DB.Where("delivery_order_no in (?)", form.DeliveryOrderNo).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.OrderShippingRecord, 0, result.RowsAffected)

	for _, p := range pick {
		list = append(list, rsp.OrderShippingRecord{
			Id:              p.Id,
			TakeOrdersTime:  p.TakeOrdersTime.Format(timeutil.TimeFormat),
			PickUser:        p.PickUser,
			ReviewUser:      p.ReviewUser,
			ReviewTime:      p.ReviewTime.Format(timeutil.TimeFormat),
			ReviewNum:       p.ReviewNum,
			DeliveryOrderNo: p.DeliveryOrderNo,
		})
	}

	res.List = list

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

	var pickGoods []batch.PickGoods

	result := global.DB.Where("pick_id in (?)", form.Id).Find(&pickGoods)

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
		form req.CompleteOrderReq
		res  rsp.CompleteOrderRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var completeOrder []order.CompleteOrder

	db := global.DB

	numbers := []string{}

	if form.Sku != "" {
		var completeOrderDetail []order.CompleteOrderDetail
		result := db.Model(&order.CompleteOrderDetail{}).Where("sku = ?", form.Sku).Find(&completeOrderDetail)
		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, detail := range completeOrderDetail {
			numbers = append(numbers, detail.Number)
		}

		numbers = slice.UniqueStringSlice(numbers)
	}

	//商品
	local := db.
		Model(&order.CompleteOrder{}).
		Where(&order.CompleteOrder{
			ShopId:         form.ShopId,
			Number:         form.Number,
			Line:           form.Line,
			DeliveryMethod: form.DeliveryMethod,
			ShopType:       form.ShopType,
			Province:       form.Province,
			City:           form.City,
			District:       form.District,
		})

	if len(numbers) > 0 {
		local.Where("number in (?)", numbers)
	}

	if form.IsRemark == 1 { //没有备注
		local.Where("order_remark == ''")
	} else if form.IsRemark == 2 { //有备注
		local.Where("order_remark != ''")
	}

	if form.PayAt != "" {
		t, err := time.ParseInLocation(timeutil.DateFormat, form.PayAt, time.Local)

		if err != nil {
			xsq_net.ErrorJSON(c, ecode.DataTransformationError)
			return
		}

		local.Where("pay_at <= ", timeutil.GetLastTime(t))
	}

	result := local.Find(&completeOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	result = local.Scopes(model.Paginate(form.Page, form.Size)).Find(&completeOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.CompleteOrder, 0, form.Size)

	for _, o := range completeOrder {
		list = append(list, rsp.CompleteOrder{
			Number:      o.Number,
			PayAt:       o.PayAt,
			ShopCode:    o.ShopCode,
			ShopName:    o.ShopName,
			ShopType:    o.ShopType,
			PayCount:    o.PayCount,
			OutCount:    o.OutCount,
			CloseCount:  o.CloseCount,
			Line:        o.Line,
			Region:      fmt.Sprintf("%s-%s-%s", o.Province, o.City, o.District),
			PickTime:    o.PickTime.Format(timeutil.TimeFormat),
			OrderRemark: o.OrderRemark,
		})
	}

	res.List = list

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
		completeOrder       order.CompleteOrder
		completeOrderDetail []order.CompleteOrderDetail
	)

	db := global.DB

	result := db.Model(&order.CompleteOrder{}).Where("number = ?", form.Number).First(&completeOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Where("number = ?", form.Number).Find(&completeOrderDetail)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	goodsMap := make(map[string][]rsp.PrePickGoods, 0)

	for _, goods := range completeOrderDetail {
		goodsMap[goods.GoodsType] = append(goodsMap[goods.GoodsType], rsp.PrePickGoods{
			GoodsName:   goods.Name,
			GoodsSpe:    goods.GoodsSpe,
			Shelves:     goods.Shelves,
			NeedNum:     goods.PayCount,
			CloseNum:    goods.CloseCount,
			OutCount:    goods.ReviewCount,
			NeedOutNum:  goods.PayCount,
			GoodsRemark: goods.GoodsRemark,
		})
	}

	res.Goods = goodsMap

	res.ShopName = completeOrder.ShopName
	res.ShopCode = completeOrder.ShopCode
	res.Line = completeOrder.Line
	res.Region = fmt.Sprintf("%s-%s-%s", completeOrder.Province, completeOrder.City, completeOrder.District)
	res.ShopType = completeOrder.ShopType
	res.Number = completeOrder.Number
	res.OrderRemark = completeOrder.OrderRemark

	xsq_net.SucJson(c, res)
}

func Count(c *gin.Context) {

	type OrderNum struct {
		Count     int `json:"count"`
		OrderType int `json:"order_type"`
	}

	var (
		dbRes []OrderNum
		res   rsp.CountRes
	)

	result := global.DB.Model(&order.Order{}).Select("count(id) as count, order_type").Group("order_type").Find(&dbRes)

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
			res.PickCount = r.Count
			break
		case 3: //3:欠货单
			res.OldCount = r.Count
			break
		}
		res.AllCount += r.Count
	}

	xsq_net.SucJson(c, res)
}

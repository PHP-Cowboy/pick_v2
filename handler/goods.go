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

// 获取待拣货订单商品列表
func GetGoodsList(c *gin.Context) {

	var form req.GetGoodsListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}
	//商品系统接口页数为index
	form.Index = form.Page

	fmt.Printf("%+v", form)

	result, err := RequestGoodsList(form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	orderList := OrderList(result.Data.List)
	xsq_net.SucJson(c, gin.H{"list": orderList, "total": result.Data.Count})
}

func OrderList(goodsList []*rsp.ApiGoods) []*rsp.OrderList {

	res := make(map[string]struct{}, 0)

	list := make([]*rsp.OrderList, 0, 16)

	for _, goods := range goodsList {
		if _, ok := res[goods.Number]; ok {
			continue
		}
		res[goods.Number] = struct{}{}

		list = append(list, &rsp.OrderList{
			Number:            goods.Number,
			PayAt:             goods.PayAt,
			ShopCode:          goods.ShopCode,
			ShopName:          goods.ShopName,
			ShopType:          goods.ShopType,
			DistributionType:  goods.DistributionType,
			SaleUnit:          goods.SaleUnit,
			PayCount:          goods.PayCount,
			OutCount:          goods.OutCount,
			LackCount:         goods.LackCount,
			Line:              goods.Line,
			Region:            goods.Province + goods.City + goods.District,
			OrderRemark:       goods.OrderRemark,
			LatestPickingTime: goods.LatestPickingTime,
		})
	}

	return list
}

// 订单明细
func GetOrderDetail(c *gin.Context) {
	var form req.GetOrderDetailForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result, err := RequestGoodsList(form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var (
		mp   = make(map[string]string, 0)
		list []*rsp.ApiGoods
	)

	mp, err = cache.GetClassification()

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	for _, l := range result.Data.List {
		goodsType, ok := mp[l.SecondType]

		if !ok {
			xsq_net.ErrorJSON(c, errors.New("商品类型:"+l.SecondType+"数据未同步"))
			return
		}

		l.GoodsType = goodsType
		list = append(list, l)
	}

	r := OrderDetail(list)
	xsq_net.SucJson(c, r)

}

func OrderDetail(goodsList []*rsp.ApiGoods) rsp.OrderDetail {

	var result rsp.OrderDetail

	detailMap := make(map[string]*rsp.Detail, 0)

	for _, list := range goodsList {

		if _, ok := detailMap[list.GoodsType]; !ok {
			detailMap[list.GoodsType] = &rsp.Detail{
				Total: 0,
				List:  make([]*rsp.GoodsDetail, 0),
			}
		}

		detailMap[list.GoodsType].Total += list.PayCount
		detailMap[list.GoodsType].List = append(detailMap[list.GoodsType].List, &rsp.GoodsDetail{
			Name:        list.Name,
			GoodsSpe:    list.GoodsSpe,
			Shelves:     list.Shelves,
			PayCount:    list.PayCount,
			CloseCount:  list.CloseCount,
			LackCount:   list.LackCount,
			GoodsRemark: list.GoodsRemark,
		})

		result.Number = list.Number
		result.PayAt = list.PayAt
		result.ShopCode = list.ShopCode
		result.ShopName = list.ShopName
		result.Line = list.Line
		result.Region = list.Province + list.City + list.District
		result.ShopType = list.ShopType
		result.OrderRemark = list.OrderRemark
	}

	result.Detail = detailMap

	return result
}

func RequestGoodsList(responseData interface{}) (rsp.ApiGoodsListRsp, error) {

	var result rsp.ApiGoodsListRsp

	body, err := request.Post("api/v1/remote/lack/list", responseData)

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)

	fmt.Println(result)
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
	var result rsp.CountRsp

	body, err := request.Get("api/v1/remote/order/node/count")

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	err = json.Unmarshal(body, &result)

	fmt.Println(result)
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

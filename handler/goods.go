package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model/other"
	"pick_v2/utils/ecode"
	"pick_v2/utils/request"
	"pick_v2/utils/xsq_net"
	"time"
)

//获取待拣货订单商品列表
func GetGoodsList(c *gin.Context) {

	var form req.GetGoodsListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}
	//商品系统接口页数为index
	form.Index = form.Page

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
			Number:           goods.Number,
			PayAt:            goods.PayAt,
			ShopCode:         goods.ShopCode,
			ShopName:         goods.ShopName,
			ShopType:         goods.ShopType,
			DistributionType: goods.DistributionType,
			SaleUnit:         goods.SaleUnit,
			PayCount:         goods.PayCount,
			Line:             goods.Line,
			Region:           goods.Province + goods.City + goods.District,
			OrderRemark:      goods.OrderRemark,
		})
	}

	return list
}

//订单明细
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

	//OrderDetail(result.Data.List)

	var list []*rsp.ApiGoods

	for _, l := range result.Data.List {
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

	body, err := request.Post("api/v1/remote/pick/lack/list", responseData)

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

//商品列表
func CommodityList(c *gin.Context) {
	var result rsp.CommodityListRsp

	body, err := request.Post("api/v1/remote/pick/shop/sku", nil)

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

//订单出货记录
func OrderShippingRecord(c *gin.Context) {

}

//订单出货记录明细
func ShippingRecordDetail(c *gin.Context) {

}

//商品分类转仓库分类
func ShopClassToWarehouseClass(shopClass string) (val string, err error) {

	goodsClassMap, ok := goCache.Get("goodsClassMap")

	mp := make(map[string]string, 0)

	if !ok {
		mp, err = Cache()
		if err != nil {
			return "", err
		}
	} else {
		mp = goodsClassMap.(map[string]string)
	}

	val, ok = mp[shopClass]

	if !ok {
		return "", errors.New("key:" + shopClass + "未找到")
	}
	return val, nil
}

var goCache *cache.Cache

func init() {
	//创建一个默认过期时间为24小时的缓存
	//每12小时清洗一次过期物品
	goCache = cache.New(24*time.Hour, 12*time.Hour)
}

func Cache() (map[string]string, error) {

	var class []other.Classification

	result := global.DB.Find(&class)

	if result.Error != nil {
		return nil, result.Error
	}

	mp := make(map[string]string, 0)

	for _, cl := range class {
		mp[cl.GoodsClass] = cl.WarehouseClass
	}

	goCache.Set("goodsClassMap", mp, cache.DefaultExpiration)

	return mp, nil
}

package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/utils/ecode"
	"pick_v2/utils/request"
	"pick_v2/utils/xsq_net"
)

//获取待拣货订单商品列表
func GetGoodsList(c *gin.Context) {

	var form req.GetGoodsListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result, err := RequestGoodsList(nil)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	orderList := OrderList(result.Data.List)
	xsq_net.SucJson(c, orderList)
}

func OrderList(goodsList []*rsp.ApiGoodsList) []*rsp.OrderList {

	res := make(map[string]*rsp.OrderList, 0)

	for _, list := range goodsList {
		if _, ok := res[list.Number]; ok {
			continue
		}
		res[list.Number] = &rsp.OrderList{
			Number:           list.Number,
			PayAt:            list.PayAt,
			ShopCode:         list.ShopCode,
			ShopName:         list.ShopName,
			ShopType:         list.ShopType,
			DistributionType: list.DistributionType,
			SaleUnit:         list.SaleUnit,
			PayCount:         list.PayCount,
			Line:             list.Line,
			Region:           list.Province + list.City + list.District,
			OrderRemark:      list.OrderRemark,
		}
	}

	list := make([]*rsp.OrderList, 0, len(res))

	for _, r := range res {
		list = append(list, r)
	}

	return list
}

func GetOrderDetail(c *gin.Context) {
	var form req.GetOrderDetailForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result, err := RequestGoodsList(nil)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	//OrderDetail(result.Data.List)

	var list []*rsp.ApiGoodsList

	for _, l := range result.Data.List {
		if l.Number != "QZG2207230005" {
			continue
		}
		list = append(list, l)
	}

	r := OrderDetail(list)
	xsq_net.SucJson(c, r)

}

func OrderDetail(goodsList []*rsp.ApiGoodsList) rsp.OrderDetail {

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

func RequestGoodsList(responseData map[string]interface{}) (rsp.ApiGoodsListRsp, error) {
	var result rsp.ApiGoodsListRsp

	body, err := request.Post("api/v1/remote/pick/lack/list", responseData)

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return result, err
	}

	return result, nil
}

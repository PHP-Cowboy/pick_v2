package rsp

import "pick_v2/model"

type OrderRsp struct {
	Code  int         `json:"code"`
	Data  []OrderInfo `json:"data"`
	Count int         `json:"count"`
	Msg   string      `json:"msg"`
}

type GoodsInfo struct {
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	Sku             string         `json:"sku"`
	GoodsType       string         `json:"goods_type"`
	GoodsSpe        string         `json:"goods_spe"`
	Shelves         string         `json:"shelves"`
	DiscountPrice   int            `json:"discount_price"`
	GoodsUnit       string         `json:"goods_unit"`
	SaleUnit        string         `json:"sale_unit"`
	SaleCode        string         `json:"sale_code"`
	PayCount        int            `json:"pay_count"`
	GoodsRemark     string         `json:"goods_remark"`
	Number          string         `json:"number"`
	DeliveryOrderNo model.GormList `json:"delivery_order_no"`
}

type OrderInfo struct {
	ShopID           int         `json:"shop_id"`
	ShopName         string      `json:"shop_name"`
	ShopType         string      `json:"shop_type"`
	Number           string      `json:"number"`
	HouseCode        string      `json:"house_code"`
	Line             string      `json:"line"`
	ShopCode         string      `json:"shop_code"`
	DistributionType int         `json:"distribution_type"`
	OrderRemark      string      `json:"order_remark"`
	PayAt            string      `json:"pay_at"`
	DeliveryAt       string      `json:"delivery_at"`
	Province         string      `json:"province"`
	City             string      `json:"city"`
	District         string      `json:"district"`
	OrderID          int         `json:"order_id"`
	Address          string      `json:"address"`
	ConsigneeName    string      `json:"consignee_name"`
	ConsigneeTel     string      `json:"consignee_tel"`
	GoodsInfo        []GoodsInfo `json:"goods_info"`
}

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
	ShopID           int           `json:"shop_id"`
	ShopName         string        `json:"shop_name"`
	ShopType         string        `json:"shop_type"`
	Number           string        `json:"number"`
	HouseCode        string        `json:"house_code"`
	Line             string        `json:"line"`
	ShopCode         string        `json:"shop_code"`
	DistributionType int           `json:"distribution_type"`
	OrderRemark      string        `json:"order_remark"`
	PayAt            *model.MyTime `json:"pay_at"`
	DeliveryAt       model.MyTime  `json:"delivery_at"`
	Province         string        `json:"province"`
	City             string        `json:"city"`
	District         string        `json:"district"`
	OrderID          int           `json:"order_id"`
	Address          string        `json:"address"`
	ConsigneeName    string        `json:"consignee_name"`
	ConsigneeTel     string        `json:"consignee_tel"`
	GoodsInfo        []GoodsInfo   `json:"goods_info"`
}

// 返回给订货系统

type GetBatchOrderAndGoodsRsp struct {
	Count int        `json:"count"`
	List  []OutOrder `json:"list"`
}

type OutOrder struct {
	DistributionType int          `json:"distribution_type"`
	PayAt            model.MyTime `json:"pay_at"`
	OrderId          int          `json:"order_id"`
	GoodsInfo        []OutGoods   `json:"goods_info"`
}

type OutGoods struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Sku           string `json:"sku"`
	GoodsType     string `json:"goods_type"`
	GoodsSpe      string `json:"goods_spe"`
	DiscountPrice int    `json:"discount_price"`
	GoodsUnit     string `json:"goods_unit"`
	SaleUnit      string `json:"sale_unit"`
	SaleCode      string `json:"sale_code"`
	OutCount      int    `json:"out_count"`
	OutAt         string `json:"out_at"`
	Number        string `json:"number"`
	CkNumber      string `json:"ck_number"`
}

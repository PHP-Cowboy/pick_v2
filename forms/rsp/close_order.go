package rsp

import "pick_v2/model"

type CloseOrderCountRsp struct {
	PendingNum   int `json:"pending_num"`
	CompleteNum  int `json:"complete_num"`
	ExceptionNum int `json:"exception_num"`
}

type CloseOrderListRsp struct {
	Total int64            `json:"total"`
	List  []CloseOrderList `json:"list"`
}

type CloseOrderList struct {
	Id               int           `json:"id"`
	Number           string        `json:"number"`
	PayAt            *model.MyTime `json:"pay_at"`
	PayTotal         int           `json:"pay_total"`
	NeedCloseTotal   int           `json:"need_close_total"`
	ShopName         string        `json:"shop_name"`
	ShopType         string        `json:"shop_type"`
	DistributionType int           `json:"distribution_type"`
	Province         string        `json:"province"`
	City             string        `json:"city"`
	District         string        `json:"district"`
	OrderRemark      string        `json:"order_remark"`
	Status           int           `json:"status"`
}

type CloseOrderDetailRsp struct {
	Number           string           `json:"number"`
	ShopName         string           `json:"shop_name"`
	DistributionType int              `json:"distribution_type"`
	District         string           `json:"district"`
	Status           int              `json:"status"`
	OrderRemark      string           `json:"order_remark"`
	List             []CloseGoodsList `json:"list"`
}

type CloseGoodsList struct {
	GoodsName      string `json:"goods_name"`
	GoodsSpe       string `json:"goods_spe"`
	PayCount       int    `json:"pay_count"`
	CloseCount     int    `json:"close_count"`      //已关闭数量
	OutCount       int    `json:"out_count"`        //已出库数量
	NeedCloseCount int    `json:"need_close_count"` //需关闭数量
	GoodsRemark    string `json:"goods_remark"`
}

// 关闭订单 订单号和类型
type CloseOrderNumberTyp struct {
	Number string
	Typ    int
}

type CloseOrderAndGoodsList struct {
	Number      string           `json:"number"`
	ShopName    string           `json:"shop_name"`
	OrderRemark string           `json:"order_remark"`
	List        []CloseGoodsList `json:"list"`
}

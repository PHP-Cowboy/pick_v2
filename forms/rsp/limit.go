package rsp

type LimitShipmentList struct {
	OutboundNumber string `json:"outbound_number"`
	ShopName       string `json:"shop_name"`
	GoodsName      string `json:"goods_name"`
	GoodsSpe       string `json:"goods_spe"`
	LimitNum       int    `json:"limit_num"`
	Status         int    `json:"status"`
}

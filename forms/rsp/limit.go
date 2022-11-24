package rsp

type LimitShipmentListRsp struct {
	Total int64               `json:"total"`
	List  []LimitShipmentList `json:"list"`
}

type LimitShipmentList struct {
	OutboundNumber string `json:"outbound_number"`
	Number         string `json:"number"`
	Sku            string `json:"sku"`
	ShopName       string `json:"shop_name"`
	GoodsName      string `json:"goods_name"`
	GoodsSpe       string `json:"goods_spe"`
	LimitNum       int    `json:"limit_num"`
	Status         int    `json:"status"`
}

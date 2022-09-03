package rsp

type PickOrderGoodsListRsp struct {
	Total int64       `json:"total"`
	List  []PickOrder `json:"list"`
}

type PickOrder struct {
	Number            string `json:"number"`
	PayAt             string `json:"pay_time"`
	ShopCode          string `json:"shop_code"`
	ShopName          string `json:"shop_name"`
	ShopType          string `json:"shop_type"`
	DistributionType  int    `json:"distribution_type"` //配送方式
	PayCount          int    `json:"pay_count"`         //下单数量
	LimitNum          int    `json:"limit_num"`
	Line              string `json:"line"`
	Region            string `json:"region"`
	OrderRemark       string `json:"order_remark"` //订单备注
	OrderType         int    `json:"order_type"`   //1:新订单,2:拣货中,3:欠货单
	LatestPickingTime string `json:"latest_picking_time"`
}

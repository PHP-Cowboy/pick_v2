package rsp

type PickOrderGoodsListRsp struct {
	Total int64       `json:"total"`
	List  []PickOrder `json:"list"`
}

type PickOrder struct {
	Id                int    `json:"id"`
	Number            string `json:"number"`
	PickNumber        string `json:"pick_number"`
	PayAt             string `json:"pay_time"`
	ShopCode          string `json:"shop_code"`
	ShopName          string `json:"shop_name"`
	ShopType          string `json:"shop_type"`
	DistributionType  int    `json:"distribution_type"` //配送方式
	ShipmentsNum      int    `json:"shipments_num"`     //发货数量
	LimitNum          int    `json:"limit_num"`
	CloseNum          int    `json:"close_num"`
	Line              string `json:"line"`
	Region            string `json:"region"`
	OrderRemark       string `json:"order_remark"` //订单备注
	OrderType         int    `json:"order_type"`   //1:新订单,2:拣货中,3:欠货单
	LatestPickingTime string `json:"latest_picking_time"`
}

type PickOrderGoods struct {
	Id        int    `json:"id"`
	GoodsName string `json:"goods_name"`
	LackCount int    `json:"lack_count"`
	LimitNum  int    `json:"limit_num"`
}

type PickOrderCount struct {
	AllCount      int `json:"all_count"`
	NewCount      int `json:"new_count"`
	PickCount     int `json:"pick_count"`
	CloseCount    int `json:"close_count"`
	CompleteCount int `json:"complete_count"`
}

type DeliveryMethodInfoRsp struct {
	UserName         string `json:"user_name"`
	Tel              string `json:"tel"`
	Province         string `json:"province"`
	City             string `json:"city"`
	District         string `json:"district"`
	Address          string `json:"address"`
	DistributionType int    `json:"distribution_type"`
}

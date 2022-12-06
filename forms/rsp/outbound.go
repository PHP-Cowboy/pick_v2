package rsp

import "pick_v2/model"

type OutboundTaskListRsp struct {
	Total int64              `json:"total"`
	List  []OutboundTaskList `json:"list"`
}

type OutboundTaskList struct {
	Id                int           `json:"id"`
	TaskName          string        `json:"task_name"`
	DeliveryStartTime *model.MyTime `json:"delivery_start_time"`
	DeliveryEndTime   *model.MyTime `json:"delivery_end_time"`
	Line              string        `json:"line"`
	DistributionType  int           `json:"distribution_type"`
	PayEndTime        *model.MyTime `json:"pay_end_time"`
	Status            int           `json:"status"`
	IsPush            int           `json:"is_push"`
	Creator           string        `json:"creator"`
	UpdateTime        string        `json:"update_time"`
}

type OutboundTaskCountRsp struct {
	Ongoing int `json:"ongoing"`
	Closed  int `json:"closed"`
}

type OutboundOrderCountRsp struct {
	Total    int `json:"total"`
	New      int `json:"new"`
	Picking  int `json:"picking"`
	Complete int `json:"complete"`
	Close    int `json:"close"`
}

type OutboundOrderListRsp struct {
	Total int64               `json:"total"`
	List  []OutboundOrderList `json:"list"`
}

type OutboundOrderList struct {
	Number            string        `json:"number"`
	OutboundNumber    string        `json:"outbound_number"`
	PayAt             *model.MyTime `json:"pay_at"`
	ShopName          string        `json:"shop_name"`
	ShopType          string        `json:"shop_type"`
	DistributionType  int           `json:"distribution_type"`
	Sku               string        `json:"sku"`
	GoodsNum          int           `json:"goods_num"`
	LimitNum          int           `json:"limit_num"`
	CloseNum          int           `json:"close_num"`
	OutCount          int           `json:"out_count"`
	Line              string        `json:"line"`
	Region            string        `json:"region"`
	LatestPickingTime *model.MyTime `json:"latest_picking_time"`
	OrderRemark       string        `json:"order_remark"`
	OrderType         int           `json:"order_type"`
}

type OutboundOrderGoodsList struct {
	Sku  string `json:"sku"`
	Name string `json:"name"`
}

type OutboundOrderDetailRsp struct {
}

type OrderOutboundRecordList struct {
}

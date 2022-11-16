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
	EndTime           *model.MyTime `json:"end_time"`
	Status            int           `json:"status"`
	IsPush            int           `json:"is_push"`
}

type OutboundOrderListRsp struct {
	Total int64               `json:"total"`
	List  []OutboundOrderList `json:"list"`
}

type OutboundOrderList struct {
	Number            string        `json:"number"`
	PayAt             *model.MyTime `json:"pay_at"`
	ShopName          string        `json:"shop_name"`
	ShopType          string        `json:"shop_type"`
	DistributionType  int           `json:"distribution_type"`
	GoodsNum          int           `json:"goods_num"`
	LimitNum          int           `json:"limit_num"`
	CloseNum          int           `json:"close_num"`
	Line              string        `json:"line"`
	Region            string        `json:"region"`
	LatestPickingTime *model.MyTime `json:"latest_picking_time"`
	OrderRemark       string        `json:"order_remark"`
	OrderType         int           `json:"order_type"`
}

type OutboundOrderDetailRsp struct {
}

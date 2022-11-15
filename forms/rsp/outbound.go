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
	DistributionType  int           `json:"delivery_method"`
	EndTime           *model.MyTime `json:"end_time"`
	Status            int           `json:"status"`
	IsPush            int           `json:"is_push"`
}

type OutboundOrderListRsp struct {
	Total int64               `json:"total"`
	List  []OutboundOrderList `json:"list"`
}

type OutboundOrderList struct {
}

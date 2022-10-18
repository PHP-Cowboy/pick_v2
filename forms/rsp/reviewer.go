package rsp

import "pick_v2/model"

type ReviewListRsp struct {
	Total int64  `json:"total"`
	List  []Pick `json:"list"`
}

type ReviewList struct {
}

type ReviewDetailRsp struct {
	TaskName        string                      `json:"task_name"`
	ShopCode        string                      `json:"shop_code"`
	OutTotal        int                         `json:"out_total"`
	UnselectedTotal int                         `json:"unselected_total"`
	ReviewTotal     int                         `json:"review_total"`
	PickUser        string                      `json:"pick_user"`
	TakeOrdersTime  *model.MyTime               `json:"take_orders_time"`
	ReviewUser      string                      `json:"review_user"`
	ReviewTime      *model.MyTime               `json:"review_time"`
	Goods           map[string][]MergePickGoods `json:"goods"`
	RemarkList      []PickRemark                `json:"remark_list"`
}

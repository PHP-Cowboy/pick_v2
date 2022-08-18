package rsp

type ReviewListRsp struct {
	Total int64  `json:"total"`
	List  []Pick `json:"list"`
}

type ReviewList struct {
}

type ReviewDetailRsp struct {
	TaskName        string                 `json:"task_name"`
	OutTotal        int                    `json:"out_total"`
	UnselectedTotal int                    `json:"unselected_total"`
	ReviewTotal     int                    `json:"review_total"`
	PickUser        string                 `json:"pick_user"`
	TakeOrdersTime  string                 `json:"take_orders_time"`
	ReviewUser      string                 `json:"review_user"`
	ReviewTime      string                 `json:"review_time"`
	Goods           map[string][]PickGoods `json:"goods"`
	RemarkList      []PickRemark           `json:"remark_list"`
}

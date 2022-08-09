package rsp

type ReceivingOrdersRsp struct {
	Id int `json:"id"`
}

type PickingRecordDetailRsp struct {
	TaskName        string                 `json:"task_name"`
	OutTotal        int                    `json:"out_total"`
	UnselectedTotal int                    `json:"unselected_total"`
	PickUser        string                 `json:"pick_user"`
	TakeOrdersTime  string                 `json:"take_orders_time"`
	ReviewUser      string                 `json:"review_user"`
	ReviewTime      string                 `json:"review_time"`
	Goods           map[string][]PickGoods `json:"goods"`
	RemarkList      []PickRemark           `json:"remark_list"`
}

type PickingRecordRsp struct {
	Total int64           `json:"total"`
	List  []PickingRecord `json:"list"`
}

type PickingRecord struct {
	Id             int    `json:"id"`
	TaskName       string `json:"task_name"`
	ShopCode       string `json:"shop_code"`
	ShopNum        int    `json:"shop_num"`
	OrderNum       int    `json:"order_num"`
	NeedNum        int    `json:"need_num"`
	TakeOrdersTime string `json:"take_orders_time"`
	ReviewUser     string `json:"review_user"`
	OutNum         int    `json:"out_num"`
	ReviewStatus   string `json:"review_status"`
}

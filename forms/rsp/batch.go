package rsp

type GetBatchListRsp struct {
	Total int64    `json:"total"`
	List  []*Batch `json:"list"`
}

type Batch struct {
	BatchName         string `json:"batch_name"`
	DeliveryStartTime string `json:"delivery_start_time"`
	DeliveryEndTime   string `json:"delivery_end_time"`
	ShopNum           int    `json:"shop_num"`
	OrderNum          int    `json:"order_num"`
	UserName          string `json:"user_name"`
	Line              string `json:"line"`
	DeliveryMethod    int    `json:"delivery_method"`
	EndTime           string `json:"end_time"`
	Status            int    `json:"status"`
	PickNum           int    `json:"pick_num"`
	RecheckSheetNum   int    `json:"recheck_sheet_num"`
}

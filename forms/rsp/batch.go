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

type GetPrePickListRsp struct {
	Total int64      `json:"total"`
	List  []*PrePick `json:"list"`
}

type PrePick struct {
	ShopCode     string               `json:"shop_code"`
	ShopName     string               `json:"shop_name"`
	Line         string               `json:"line"`
	CategoryInfo map[string]PickCount `json:"category_info"`
}

type PickCount struct {
	WaitingPick int `json:"waiting_pick"`
	PickedCount int `json:"picked_count"`
}

type GetPrePickDetailRsp struct {
	Goods      map[string][]PrePickGoods `json:"goods"`
	RemarkList []Remark                  `json:"remark_list"`
}

type PrePickGoods struct {
	GoodsName  string `json:"goods_name"`
	GoodsSpe   string `json:"goods_spe"`
	Shelves    string `json:"shelves"`
	NeedNum    int    `json:"need_num"`
	CloseNum   int    `json:"close_num"`
	OutCount   int    `json:"out_count"`
	NeedOutNum int    `json:"need_out_num"`
}

type Remark struct {
	Number      string `json:"number"`
	OrderRemark string `json:"order_remark"`
	GoodsRemark string `json:"goods_remark"`
}

package rsp

import "time"

type ReceivingOrdersRsp struct {
	Id             int        `json:"id"`
	BatchId        int        `json:"batch_id"`
	Version        int        `json:"version"`
	TakeOrdersTime *time.Time `json:"take_orders_time"`
}

type PickingRecordDetailRsp struct {
	TaskName        string                      `json:"task_name"`
	OutTotal        int                         `json:"out_total"`
	UnselectedTotal int                         `json:"unselected_total"`
	PickUser        string                      `json:"pick_user"`
	TakeOrdersTime  string                      `json:"take_orders_time"`
	ReviewUser      string                      `json:"review_user"`
	ReviewTime      string                      `json:"review_time"`
	Goods           map[string][]MergePickGoods `json:"goods"`
	RemarkList      []PickRemark                `json:"remark_list"`
}

type MergePickGoods struct {
	Id          int        `json:"id"`
	Sku         string     `json:"sku"`
	GoodsName   string     `json:"goods_name"`
	GoodsType   string     `json:"goods_type"`
	GoodsSpe    string     `json:"goods_spe"`
	Shelves     string     `json:"shelves"`
	NeedNum     int        `json:"need_num"`
	CompleteNum int        `json:"complete_num"`
	ReviewNum   int        `json:"review_num"`
	Unit        string     `json:"unit"`
	ParamsId    []ParamsId `json:"params_id"`
}

type MyMergePickGoods []MergePickGoods

func (pg MyMergePickGoods) Len() int {
	return len(pg)
}

func (pg MyMergePickGoods) Less(i, j int) bool {
	return pg[i].Shelves < pg[j].Shelves
}

func (pg MyMergePickGoods) Swap(i, j int) {
	pg[i].Shelves, pg[j].Shelves = pg[j].Shelves, pg[i].Shelves
}

type ParamsId struct {
	PickGoodsId  int `json:"pick_goods_id"`
	OrderGoodsId int `json:"order_goods_id"`
}

type PickingRecordRsp struct {
	Total int64           `json:"total"`
	List  []PickingRecord `json:"list"`
}

type PickingRecord struct {
	Id               int    `json:"id"`
	TaskName         string `json:"task_name"`
	ShopCode         string `json:"shop_code"`
	ShopNum          int    `json:"shop_num"`
	OrderNum         int    `json:"order_num"`
	NeedNum          int    `json:"need_num"`
	TakeOrdersTime   string `json:"take_orders_time"`
	ReviewUser       string `json:"review_user"`
	OutNum           int    `json:"out_num"`
	ReviewStatus     string `json:"review_status"`
	DistributionType int    `json:"distribution_type"`
	IsRemark         bool   `json:"is_remark"`
}

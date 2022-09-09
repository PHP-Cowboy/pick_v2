package rsp

import "time"

type PickListModel struct {
	Id             int        `json:"id"`
	TaskName       string     `json:"task_name"`
	ShopCode       string     `json:"shop_code"`
	ShopName       string     `json:"shop_name"`
	ShopNum        int        `json:"shop_num"`
	OrderNum       int        `json:"order_num"`
	NeedNum        int        `json:"need_num"`
	PickUser       string     `json:"pick_user"`
	TakeOrdersTime *time.Time `json:"take_orders_time"`
	UpdateTime     time.Time  `json:"update_time"`
	Status         int        `json:"status"`
	ReviewNum      int        `json:"review_num"`
	Num            int        `json:"num"`
	ReviewUser     string     `json:"review_user"`
	ReviewTime     *time.Time `json:"review_time"`
	OrderRemark    string     `json:"order_remark"`
	GoodsRemark    string     `json:"goods_remark"`
	PickNum        int        `json:"pick_num"`
}

type PickListRsp struct {
	Total int64  `json:"total"`
	List  []Pick `json:"list"`
}

type Pick struct {
	Id             int    `json:"id"`
	TaskName       string `json:"task_name"`
	ShopCode       string `json:"shop_code"`
	ShopName       string `json:"shop_name"`
	ShopNum        int    `json:"shop_num"`
	OrderNum       int    `json:"order_num"`
	NeedNum        int    `json:"need_num"`
	PickUser       string `json:"pick_user"`
	TakeOrdersTime string `json:"take_orders_time"`
	IsRemark       bool   `json:"is_remark"`
	Status         int    `json:"status"`
	UpdateTime     string `json:"update_time"` //结束时间
	PickNum        int    `json:"pick_num"`    //已拣数量
	ReviewNum      int    `json:"review_num"`  //复核数
	Num            int    `json:"num"`         //件数
	ReviewUser     string `json:"review_user"` //复核人
	ReviewTime     string `json:"review_time"` //复核时间
	//Shelves    string `json:"shelves"`     //货位号不要了
}

type GetPickDetailRsp struct {
	BatchId        int                         `json:"batch_id"`
	PickId         int                         `json:"pick_id"`
	TaskName       string                      `json:"task_name"`
	ShopNum        int                         `json:"shop_num"`
	OrderNum       int                         `json:"order_num"`
	GoodsNum       int                         `json:"goods_num"`
	PickUser       string                      `json:"pick_user"`
	TakeOrdersTime string                      `json:"take_orders_time"`
	Goods          map[string][]MergePickGoods `json:"goods"`
	RemarkList     []PickRemark                `json:"remark_list"`
}

type PickGoods struct {
	Id          int    `json:"id"`
	Sku         string `json:"sku"`
	GoodsName   string `json:"goods_name"`
	GoodsSpe    string `json:"goods_spe"`
	Shelves     string `json:"shelves"`
	NeedNum     int    `json:"need_num"`
	CompleteNum int    `json:"complete_num"`
	ReviewNum   int    `json:"review_num"`
	Unit        string `json:"unit"`
}

type PickRemark struct {
	Number      string `json:"number"`
	OrderRemark string `json:"order_remark"`
	GoodsRemark string `json:"goods_remark"`
}

package rsp

import "time"

type PickListModel struct {
	Id             int        `json:"id"`
	ShopCode       string     `json:"shop_code"`
	ShopName       string     `json:"shop_name"`
	ShopNum        int        `json:"shop_num"`
	OrderNum       int        `json:"order_num"`
	NeedNum        int        `json:"need_num"`
	PickUser       string     `json:"pick_user"`
	TakeOrdersTime *time.Time `json:"take_orders_time"`
	OrderRemark    string     `json:"order_remark"`
	GoodsRemark    string     `json:"goods_remark"`
}

type PickListRsp struct {
	Total int64  `json:"total"`
	List  []Pick `json:"list"`
}

type Pick struct {
	Id             int        `json:"id"`
	ShopCode       string     `json:"shop_code"`
	ShopName       string     `json:"shop_name"`
	ShopNum        int        `json:"shop_num"`
	OrderNum       int        `json:"order_num"`
	NeedNum        int        `json:"need_num"`
	PickUser       string     `json:"pick_user"`
	TakeOrdersTime *time.Time `json:"take_orders_time"`
	IsRemark       bool       `json:"is_remark"`
}

type GetPickDetailRsp struct {
	TaskName       string                 `json:"task_name"`
	ShopNum        int                    `json:"shop_num"`
	OrderNum       int                    `json:"order_num"`
	GoodsNum       int                    `json:"goods_num"`
	PickUser       string                 `json:"pick_user"`
	TakeOrdersTime *time.Time             `json:"take_orders_time"`
	Goods          map[string][]PickGoods `json:"goods"`
	RemarkList     []PickRemark           `json:"remark_list"`
}

type PickGoods struct {
	GoodsName   string `json:"goods_name"`
	GoodsSpe    string `json:"goods_spe"`
	Shelves     string `json:"shelves"`
	NeedNum     int    `json:"need_num"`
	CompleteNum int    `json:"complete_num"`
}

type PickRemark struct {
	Number      string `json:"number"`
	OrderRemark string `json:"order_remark"`
	GoodsRemark string `json:"goods_remark"`
}

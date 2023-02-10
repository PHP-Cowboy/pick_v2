package rsp

import (
	"pick_v2/model"
)

type AdminPickListRsp struct {
	Total int             `json:"total"`
	List  []AdminPickList `json:"list"`
}

type AdminPickList struct {
	ShopId   int    `json:"shop_id"`
	ShopCode string `json:"shop_code"`
	ShopName string `json:"shop_name"`
	OrderNum int    `json:"order_num"`
	NeedNum  int    `json:"need_num"`
	Remark   string `json:"remark"`
	Num      int    `json:"num"`
}

type AdminPickDetailRsp struct {
	BatchId        int                         `json:"batch_id"`
	PickId         int                         `json:"pick_id"`
	TaskName       string                      `json:"task_name"`
	ShopCode       string                      `json:"shop_code"`
	ShopNum        int                         `json:"shop_num"`
	OrderNum       int                         `json:"order_num"`
	GoodsNum       int                         `json:"goods_num"`
	PickUser       string                      `json:"pick_user"`
	TakeOrdersTime *model.MyTime               `json:"take_orders_time"`
	Goods          map[string][]MergePickGoods `json:"goods"`
	RemarkList     []PickRemark                `json:"remark_list"`
}

type BatchShopGoodsList struct {
	Sku         string `json:"sku"`
	GoodsName   string `json:"goods_name"`
	CompleteNum int    `json:"complete_num"`
	ReviewNum   int    `json:"review_num"`
}

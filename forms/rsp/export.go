package rsp

type PrePickAndGoods struct {
	PrePickGoodsId int    `json:"pre_pick_goods_id"`
	ShopCode       string `json:"shop_code"`
	ShopName       string `json:"shop_name"`
	GoodsType      string `json:"goods_type"`
	GoodsName      string `json:"goods_name"`
	GoodsSpe       string `json:"goods_spe"`
	Unit           string `json:"unit"`
	NeedNum        int    `json:"need_num"`
	ReviewNum      int    `json:"review_num"`
}

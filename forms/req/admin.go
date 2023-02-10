package req

type AdminPickListReq struct {
	BatchId   int    `json:"batch_id" form:"batch_id" binding:"required"`
	ShopId    int    `json:"shop_id" form:"shop_id"`
	GoodsName string `json:"goods_name" form:"goods_name"`
	GoodsType string `json:"goods_type" form:"goods_type"`
	IsRemark  int    `json:"is_remark" form:"is_remark"`
}

type AdminPickDetailReq struct {
	BatchId  int    `json:"batch_id" form:"batch_id" binding:"required"`
	ShopName string `json:"shop_name" form:"shop_name" binding:"required"`
	ShopCode string `json:"shop_code" form:"shop_code"`
}

type BatchShopGoodsListReq struct {
	BatchId int `json:"batch_id" form:"batch_id" binding:"required"`
	ShopId  int `json:"shop_id" form:"shop_id"`
}

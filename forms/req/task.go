package req

type PickListForm struct {
	Paging
	BatchId int    `json:"batch_id" form:"batch_id" binding:"required"`
	ShopId  int    `json:"shop_id" form:"shop_id"`
	Goods   string `json:"goods" form:"goods"`
	Number  string `json:"number" form:"number"`
	Status  *int   `json:"status" form:"status"`
}

type PickToppingForm struct {
	Id int `json:"id" form:"id"`
}

type GetPickDetailForm struct {
	PickId int `json:"pick_id" form:"pick_id"`
}

type ChangeNumReq struct {
	Id  int  `json:"id" binding:"required"`
	Num *int `json:"num" binding:"required"`
}

type PrintReq struct {
	Ids []int `json:"ids" binding:"required"`
}

type AssignReq struct {
	UserId  int   `json:"user_id" binding:"required"`
	PickIds []int `json:"pick_ids" binding:"required"`
}

type ChangeReviewNumForm struct {
	BatchId   int         `json:"batch_id" binding:"required"`
	PickId    int         `json:"pick_id" binding:"required"`
	SkuReview []SkuReview `json:"sku_review"`
}

type SkuReview struct {
	Sku string `json:"sku" binding:"required"`
	Num *int   `json:"num" binding:"required"`
}

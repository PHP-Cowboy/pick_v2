package req

type PickListForm struct {
	Paging
	BatchId   int    `json:"batch_id" form:"batch_id" binding:"required"`
	ShopId    int    `json:"shop_id" form:"shop_id"`
	GoodsName string `json:"goods_name" form:"goods_name"`
	GoodsType string `json:"goods_type" form:"goods_type"`
	Number    string `json:"number" form:"number"`
	Status    *int   `json:"status" form:"status"`
	HasRemark int    `json:"has_remark" form:"has_remark"` //0:全部;1:有备注;2:无备注
}

type PickToppingForm struct {
	Id int `json:"id" form:"id"`
}

type GetPickDetailForm struct {
	PickId int `json:"pick_id" form:"pick_id"`
}

type VoidExpressBillForm struct {
}

type ReprintExpressBillForm struct {
	CourierNumber string `json:"courier_number"`
}

type ChangeNumReq struct {
	Id  int  `json:"id" binding:"required"`
	Num *int `json:"num" binding:"required"`
}

type PrintReq struct {
	Ids  []int `json:"ids" binding:"required"`
	Type int   `json:"type" binding:"required"` // 1-全部打印 2-打印箱单 3-打印出库单 第一次全打，后边的前段选
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
	Sku       string `json:"sku" binding:"required"`
	ReviewNum *int   `json:"review_num" binding:"required"`
}

type ChangeCompleteNumForm struct {
	BatchId     int           `json:"batch_id" binding:"required"`
	PickId      int           `json:"pick_id" binding:"required"`
	SkuComplete []SkuComplete `json:"sku_complete"`
}

type SkuComplete struct {
	Sku         string `json:"sku" binding:"required"`
	CompleteNum *int   `json:"complete_num" binding:"required"`
}

type CancelPickForm struct {
	Ids []int `json:"ids" form:"ids"`
}

type CentralizedPickDetailPDAForm struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type NoNeedToReview struct {
	Id int `json:"id" binding:"required"`
}

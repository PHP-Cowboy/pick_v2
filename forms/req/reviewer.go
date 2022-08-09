package req

type ReviewListReq struct {
	Paging
	Name   string `json:"name" form:"name"`
	Status int    `json:"status" form:"status" binding:"required,oneof=1 2"`
}

type ReviewDetailReq struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type ConfirmDeliveryReq struct {
	Id             int                    `json:"id"`
	CompleteReview []CompleteReviewDetail `json:"complete_review"`
}

type CompleteReviewDetail struct {
	PickGoodsId int `json:"pick_goods_id"`
	ReviewNum   int `json:"review_num"`
}

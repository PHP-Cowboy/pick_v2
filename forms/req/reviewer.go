package req

type ReviewListReq struct {
	Paging
	Name     string `json:"name" form:"name"`
	Status   int    `json:"status" form:"status" binding:"required,oneof=1 2"`
	UserName string `json:"user_name"`
}

type ReviewDetailReq struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type ConfirmDeliveryReq struct {
	Id             int                    `json:"id"`
	Num            int                    `json:"num" binding:"required"`
	CompleteReview []CompleteReviewDetail `json:"complete_review"`
	UserName       string                 `json:"user_name"`
}

type CompleteReviewDetail struct {
	Sku       string     `json:"sku"`
	ParamsId  []ParamsId `json:"params_id"`
	ReviewNum int        `json:"review_num"`
}

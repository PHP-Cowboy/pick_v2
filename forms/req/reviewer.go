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
	Id             int                  `json:"id"`
	Num            int                  `json:"num" binding:"required"`
	CompleteReview []CompletePickDetail `json:"complete_review"`
}

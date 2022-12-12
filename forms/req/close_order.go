package req

type CloseOrderListForm struct {
	Paging
	Number string `json:"number" form:"number"`
	Status int    `json:"status" form:"status" binding:"required"`
}

type CloseOrderDetailForm struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type CloseOrder struct {
	Number []string `json:"number"`
}

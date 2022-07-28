package req

type GetGoodsListForm struct {
	Paging
}

type GetOrderDetailForm struct {
	Number string `json:"number" form:"number" binding:"required"`
}

package req

type ShopListForm struct {
	Paging
	ShopName   string `json:"shop_name" form:"shop_name"`
	LineStatus int    `json:"line_status" form:"line_status"`
	Line       string `json:"line" form:"line"`
	IsPage     bool   `json:"is_page" form:"is_page"`
}

type BatchSetLineForm struct {
	Ids  []int  `json:"ids"`
	Line string `json:"line" binding:"required"`
}

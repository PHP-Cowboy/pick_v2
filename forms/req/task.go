package req

type PickListForm struct {
	Paging
	ShopId int    `json:"shop_id" form:"shop_id"`
	Goods  string `json:"goods" form:"goods"`
	Number string `json:"number" form:"number"`
}

type PickToppingForm struct {
	Id int `json:"id" form:"id"`
}

type GetPickDetailForm struct {
	PickId int `json:"pick_id" form:"pick_id"`
}

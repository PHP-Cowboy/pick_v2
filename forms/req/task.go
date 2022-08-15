package req

type PickListForm struct {
	Paging
	ShopId int    `json:"shop_id" form:"shop_id"`
	Goods  string `json:"goods" form:"goods"`
	Number string `json:"number" form:"number"`
	Status *int   `json:"status"`
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

type PrintParams struct {
	WarehouseId     int    `json:"warehouse_id" binding:"required"`
	DeliveryOrderNo string `json:"delivery_order_no" binding:"required"`
	ShopCode        string `json:"shop_code" binding:"required"`
}

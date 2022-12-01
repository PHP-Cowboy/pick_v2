package req

type ReceivingOrdersForm struct {
	Typ int `json:"typ" binding:"required"`
}

type PickingRecordForm struct {
	Paging
	Status *int `json:"status" form:"status"`
}

type CompletePickForm struct {
	PickId       int                  `json:"pick_id" form:"pick_id"`
	CompletePick []CompletePickDetail `json:"complete_pick"`
	Type         int                  `json:"type" form:"type" binding:"required,oneof=1 2"` //1正常拣货 2无需拣货
	UserName     string               `json:"user_name"`
}

type CompleteConcentratedPickForm struct {
	Id      int `json:"id" binding:"required"`
	PickNum int `json:"pick_num" binding:"required"`
	Typ     int `json:"typ" binding:"required"`
}

type RemainingQuantityForm struct {
	Typ int `json:"typ"`
}

type CompletePickDetail struct {
	Sku         string     `json:"sku"`
	ParamsId    []ParamsId `json:"params_id"`
	CompleteNum int        `json:"complete_num"`
}

type ParamsId struct {
	PickGoodsId  int `json:"pick_goods_id"`
	OrderGoodsId int `json:"order_goods_id"`
}

type PickingRecordDetailForm struct {
	PickId int `json:"pick_id" form:"pick_id"`
}

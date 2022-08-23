package req

type ReceivingOrdersForm struct {
	Paging
}

type PickingRecordForm struct {
	Paging
	Status *int `json:"status" form:"status"`
}

type CompletePickForm struct {
	PickId       int                  `json:"pick_id" form:"pick_id"`
	CompletePick []CompletePickDetail `json:"complete_pick"`
	Type         int                  `json:"type" form:"type" binding:"required,oneof=1 2"` //1正常拣货 2无需拣货
}

type CompletePickDetail struct {
	Sku         string     `json:"sku"`
	ParamsId    []ParamsId `json:"params_id"`
	CompleteNum int        `json:"complete_num"`
}

type ParamsId struct {
	PickGoodsId int `json:"pick_goods_id"`
	OrderInfoId int `json:"order_info_id"`
}

type PickingRecordDetailForm struct {
	PickId int `json:"pick_id" form:"pick_id"`
}

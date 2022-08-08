package req

type ReceivingOrdersForm struct {
	Paging
}

type PickingRecordForm struct {
	Paging
}

type CompletePickForm struct {
	PickId       int                  `json:"pick_id" form:"pick_id"`
	CompletePick []CompletePickDetail `json:"complete_pick"`
}

type CompletePickDetail struct {
	PickGoodsId int `json:"pick_goods_id"`
	CompleteNum int `json:"complete_num"`
}

type PickingRecordDetailForm struct {
	PickId int `json:"pick_id" form:"pick_id"`
}

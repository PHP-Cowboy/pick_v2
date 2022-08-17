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
	Id          int `json:"id"`
	CompleteNum int `json:"complete_num"`
}

type PickingRecordDetailForm struct {
	PickId int `json:"pick_id" form:"pick_id"`
}

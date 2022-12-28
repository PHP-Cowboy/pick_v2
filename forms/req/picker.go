package req

type ReceivingOrdersForm struct {
	Typ      int    `json:"typ" binding:"required"`
	UserName string `json:"-"`
}

type PickingRecordForm struct {
	Paging
	Status     *int   `json:"status" form:"status"`
	PickUser   string `json:"-"`
	TowDaysAgo string `json:"-"`
}

type CompletePickForm struct {
	PickId       int                  `json:"pick_id" form:"pick_id"`
	CompletePick []CompletePickDetail `json:"complete_pick"`
	Type         int                  `json:"type" form:"type" binding:"required,oneof=1 2"` //1正常拣货 2无需拣货
	UserName     string               `json:"user_name"`
}

type CompleteConcentratedPickForm struct {
	Id          int `json:"id" binding:"required"`
	CompleteNum int `json:"complete_num" binding:"required"`
}

type RemainingQuantityForm struct {
	Typ      int    `json:"typ" form:"typ"`
	PickUser string `json:"-"`
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

type ConcentratedPickReceivingOrdersForm struct {
	BatchId int `json:"batch_id" binding:"required"`
}

type CustomsDeclarationReminderForm struct {
	Id         int          `json:"id"`
	SkuNeedNum []SkuNeedNum `json:"sku_need_num"`
}

type SkuNeedNum struct {
	Sku     string `json:"sku"`
	NeedNum int    `json:"need_num"`
}

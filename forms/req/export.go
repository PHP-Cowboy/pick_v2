package req

type FirstMaterialExportReq struct {
	Id int `json:"id" form:"id"`
}

type OutboundBatchFormReq struct {
	Id int `json:"id" form:"id" form:"id"`
}

type BatchTaskForm struct {
	Id      int `json:"id" form:"id"`             //pick_id
	BatchId int `json:"batch_id" form:"batch_id"` //pick_id
}

type LackForm struct {
}

type BatchShopForm struct {
	Id int `json:"id" form:"id" binding:"required"` //batch_id
}

type BatchShopMaterialForm struct {
	Id int `json:"id" form:"id" binding:"required"` //batch_id
}

type GoodsSummaryListForm struct {
	BatchId    int      `form:"batch_id" binding:"required"`
	Typ        int      `form:"typ"`
	GoodsTypes []string `json:"goods_typs"`
}

type ShopAddressReq struct {
	BatchId int `json:"batch_id" form:"batch_id" binding:"required"`
}

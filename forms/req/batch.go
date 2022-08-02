package req

type CreateBatchForm struct {
	WarehouseId       int    `json:"warehouse_id" form:"warehouse_id"`
	Lines             string `json:"lines" form:"lines"`
	PayEndTime        string `json:"pay_end_time"  form:"pay_end_time"`
	DeliveryStartTime string `json:"delivery_start_time" form:"delivery_start_time"`
	DeliveryEndTime   string `json:"delivery_end_time" form:"delivery_end_time"`
	DeType            int    `json:"de_type" form:"de_type" validate:"oneof=1 2 3 4 5"` //1-公司配送 2-用户自提 3-三方物流 4-快递配送 5-首批物料|设备单
	Sku               string `json:"sku" form:"sku"`
	Goods             string `json:"goods" form:"goods"`
}

type GetBatchListForm struct {
	Paging
}

package req

type CreateBatchForm struct {
	WarehouseId       int    `json:"warehouse_id" form:"warehouse_id"`
	PayEndTime        string `json:"pay_end_time" form:"pay_end_time"`
	DeliveryStartTime string `json:"delivery_start_time" form:"delivery_start_time"`
	DeliveryEndTime   string `json:"delivery_end_time" form:"delivery_end_time"`
	Line              string `json:"line" form:"line"`
	DeliveryMethod    int    `json:"delivery_method" form:"delivery_method"`
	Goods             string `json:"goods" form:"goods"`
}

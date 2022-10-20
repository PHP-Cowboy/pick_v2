package req

type TaskListForm struct {
	Paging
	Status int `json:"status" form:"status"`
}

type ChangeTaskForm struct {
	OrderNo string `json:"order_no" binding:"required"`
	Status  int    `json:"status" binding:"required"`
}

type TaskRecordListForm struct {
	Paging
	OrderNo   string `json:"order_no" form:"order_no" binding:"required"`
	GoodsName string `json:"goods_name" form:"goods_name"`
	GoodsType string `json:"goods_type" form:"goods_type"`
	IsNeed    bool   `json:"is_need" form:"is_need"`
}

type ExportForm struct {
	OrderNo string `json:"order_no" form:"order_no" binding:"required"`
}

type TypeListForm struct {
	OrderNo string `json:"order_no" form:"order_no" binding:"required"`
}

type InventoryRecordListForm struct {
	OrderNo   string `json:"order_no" form:"order_no"`
	Sku       string `json:"sku" form:"sku"`
	SortField string `json:"sort_field" form:"sort_field"`
	SortRule  string `json:"sort_rule" form:"sort_rule" default:"asc"`
}

type InventoryRecordDeleteForm struct {
	Id int `json:"id"`
}

type CountForm struct {
	OrderNo string `json:"order_no" form:"order_no"`
}

type UserInventoryRecordListForm struct {
	Paging
	OrderNo string `json:"order_no" form:"order_no" binding:"required"`
	Sku     string `json:"sku" form:"sku"`
}

type UpdateInventoryRecordForm struct {
	Id           int     `json:"id"`
	InventoryNum float64 `json:"inventory_num"`
}

type BatchCreateForm struct {
	OrderNo string            `json:"order_no" form:"order_no"`
	Records []InventoryRecord `json:"records"`
}

type InventoryRecord struct {
	Sku          string  `json:"sku" form:"sku"`
	InventoryNum float64 `json:"inventory_num" form:"inventory_num"`
}

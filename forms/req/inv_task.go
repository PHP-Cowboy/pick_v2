package req

type ChangeTaskForm struct {
	Id     string `json:"id" binding:"required"`
	Status int    `json:"status" binding:"required"`
}

type TaskRecordListForm struct {
	Paging
	SelfBuiltId int    `json:"self_built_id" form:"self_built_id"`
	GoodsName   string `json:"goods_name" form:"goods_name"`
	GoodsType   string `json:"goods_type" form:"goods_type"`
	IsNeed      bool   `json:"is_need" form:"is_need"`
	SortField   string `json:"sort_field" form:"sort_field"`
	SortRule    string `json:"sort_rule" form:"sort_rule" default:"asc"`
	InvType     int    `json:"inv_type" form:"inv_type" binding:"required,oneof=1 2"`
}

type ExportForm struct {
	OrderNo string `json:"order_no" form:"order_no" binding:"required"`
}

type TypeListForm struct {
	OrderNo string `json:"order_no" form:"order_no" binding:"required"`
}

type InventoryRecordListForm struct {
	SelfBuiltId int    `json:"self_built_id" form:"self_built_id"`
	Sku         string `json:"sku" form:"sku"`
	InvType     int    `json:"inv_type" form:"inv_type" binding:"required,oneof=1 2"`
}

type InventoryRecordDeleteForm struct {
	Id int `json:"id"`
}

type CountForm struct {
	SelfBuiltId int `json:"self_built_id" form:"self_built_id"`
	InvType     int `json:"inv_type" binding:"required,oneof=1 2"`
}

type NotInvCountForm struct {
	OrderNo string `json:"order_no" form:"order_no"`
}

type UserNotInventoryRecordListForm struct {
	SelfBuiltId int    `json:"self_built_id" form:"self_built_id" binding:"required"`
	OrderNo     string `json:"order_no" form:"order_no" binding:"required"`
	Sku         string `json:"sku" form:"sku"`
}

type UserInventoryRecordListForm struct {
	Paging
	SelfBuiltId int    `json:"self_built_id" form:"self_built_id" binding:"required"`
	Sku         string `json:"sku" form:"sku"`
	InvType     int    `json:"inv_type" binding:"required,oneof=1 2"`
}

type UpdateInventoryRecordForm struct {
	Id           int     `json:"id"`
	InventoryNum float64 `json:"inventory_num"`
}

type BatchCreateForm struct {
	SelfBuiltId int               `json:"self_built_id" form:"self_built_id"`
	InvType     int               `json:"inv_type" binding:"required,oneof=1 2"`
	Records     []InventoryRecord `json:"records"`
}

type InventoryRecord struct {
	Sku          string  `json:"sku" form:"sku"`
	InventoryNum float64 `json:"inventory_num" form:"inventory_num"`
}

type InvTaskForm struct {
	TaskName string `json:"task_name" binding:"required"`
	OrderNo  string `json:"order_no"`
}

type ChangeSelfBuiltTaskForm struct {
	Id      int    `json:"id" binding:"required"`
	OrderNo string `json:"order_no" binding:"required"`
}

type SelfBuiltTaskListForm struct {
	Paging
}

type SetSecondInventoryForm struct {
	Id      int      `json:"id" binding:"required"`
	Sku     []string `json:"sku"`
	InvType int      `json:"inv_type" binding:"required,oneof=1 2"`
}

type InvAgainForm struct {
	SelfBuiltId int      `json:"self_built_id" form:"self_built_id"`
	Sku         []string `json:"sku"`
}

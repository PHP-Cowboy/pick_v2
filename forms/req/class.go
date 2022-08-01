package req

type ClassListForm struct {
	Paging
	ClassStatus    int    `json:"class_status" form:"class_status"`
	WarehouseClass string `json:"warehouse_class" form:"warehouse_class"`
}

type BatchSetClassForm struct {
	Ids            []int  `json:"ids"`
	WarehouseClass string `json:"warehouse_class" binding:"required"`
}

type ClassNameListForm struct {
}

package req

type GetWarehouseListForm struct {
	Paging
}

type CreateWarehouseForm struct {
	WarehouseName string `json:"warehouse_name"`
	Abbreviation  string `json:"abbreviation"`
}

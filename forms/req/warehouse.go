package req

type GetWarehouseListForm struct {
}

type CreateWarehouseForm struct {
	WarehouseName string `json:"warehouse_name"`
	Abbreviation  string `json:"abbreviation"`
}

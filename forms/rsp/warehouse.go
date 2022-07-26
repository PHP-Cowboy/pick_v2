package rsp

type GetWarehouseListRsp struct {
	Total int64            `json:"total"`
	List  []*WarehouseList `json:"list"`
}

type WarehouseList struct {
	Id            int    `json:"id"`
	WarehouseName string `json:"warehouse_name"`
	Abbreviation  string `json:"abbreviation"`
}

type CreateWarehouseRsp struct {
	Id            int    `json:"id"`
	WarehouseName string `json:"warehouse_name"`
	Abbreviation  string `json:"abbreviation"`
}

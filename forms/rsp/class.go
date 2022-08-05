package rsp

type ClassListRsp struct {
	Total int64    `json:"total"`
	List  []*Class `json:"list"`
}

type Class struct {
	Id             int    `json:"id"`
	GoodsClass     string `json:"goods_class"`
	WarehouseClass string `json:"warehouse_class"`
}

type ClassNameListRsp struct {
	WarehouseClass string `json:"warehouse_class"`
}

type GoodsClassListRsp struct {
	GoodsClass string `json:"goods_class"`
}

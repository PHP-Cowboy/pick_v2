package rsp

type AddUserRsp struct {
	Id          int    `json:"id"`
	Account     string `json:"account"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	WarehouseId int    `json:"warehouse_id"`
	Status      int    `json:"status"`
	CreateTime  string `json:"create_time"`
}

type UserListRsp struct {
	Total int64   `json:"total"`
	List  []*User `json:"list"`
}

type User struct {
	Id          int    `json:"id"`
	Account     string `json:"account"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	WarehouseId int    `json:"warehouse_id"`
	Status      bool   `json:"status"` //bool
	CreateTime  string `json:"create_time"`
}

type GetWarehouseUserCountListRsp struct {
	Count         int    `json:"count"`
	WarehouseId   int    `json:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"`
}

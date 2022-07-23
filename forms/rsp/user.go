package rsp

type AddUserRsp struct {
	Id          int    `json:"id"`
	Account     string `json:"account"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	WarehouseId int    `json:"warehouse_id"`
}

type UserListRsp struct {
	Total int64        `json:"total"`
	Data  []*AddUserRsp `json:"data"`
}

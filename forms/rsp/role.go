package rsp

type GetRoleList struct {
	Total int64   `json:"total"`
	Data  []*Role `json:"data"`
}

type Role struct {
	Id         int    `json:"id"`
	CreateTime string `json:"create_time"`
	Name       string `json:"name"`
}

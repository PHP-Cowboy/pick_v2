package rsp

type GetRoleList struct {
	Total int64   `json:"total"`
	List  []*Role `json:"list"`
}

type Role struct {
	Id         int    `json:"id"`
	CreateTime string `json:"create_time"`
	Name       string `json:"name"`
}

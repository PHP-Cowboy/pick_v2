package rsp

type DictTypeListRsp struct {
	Total int64       `json:"total"`
	List  []*DictType `json:"list"`
}

type DictType struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	CreateTime string `json:"create_time"`
}

type DictListRsp struct {
	Code       string `json:"code"`
	TypeCode   string `json:"type_code"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	IsEdit     int    `json:"is_edit"`
	CreateTime string `json:"create_time"`
}

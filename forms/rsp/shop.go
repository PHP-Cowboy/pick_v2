package rsp

type ShopListRsp struct {
	Total int64   `json:"total"`
	List  []*Shop `json:"list"`
}

type Shop struct {
	Id       int    `json:"id"`
	ShopId   int    `json:"shop_id"`
	ShopName string `json:"shop_name"`
	ShopCode string `json:"shop_code"`
	Line     string `json:"line"`
}

type LineListRsp struct {
	Line string `json:"line"`
}

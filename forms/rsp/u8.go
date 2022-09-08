package rsp

type LogListRsp struct {
	Total int64     `json:"total"`
	List  []LogList `json:"list"`
}

type LogList struct {
	Id          int    `json:"id"`
	Number      string `json:"number"`
	BatchId     int    `json:"batch_id"`
	Status      int    `json:"status"`
	RequestXml  string `json:"request_xml"`
	ResponseXml string `json:"response_xml"`
	ResponseNo  string `json:"response_no"`
	Msg         string `json:"msg"`
	ShopName    string `json:"shop_name"`
}

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

type LogDetailRsp struct {
	Id               int    `json:"id"`
	PickId           int    `json:"pick_id"`
	BatchId          int    `json:"batch_id"`
	PrePickGoodsId   int    `json:"pre_pick_goods_id"`
	OrderGoodsId     int    `json:"order_goods_id"`
	Number           string `json:"number"`
	ShopId           int    `json:"shop_id"`
	DistributionType int    `json:"distribution_type"`
	Sku              string `json:"sku"`
	GoodsName        string `json:"goods_name"`
	GoodsType        string `json:"goods_type"`
	GoodsSpe         string `json:"goods_spe"`
	Shelves          string `json:"shelves"`
	DiscountPrice    int    `json:"discount_price"`
	NeedNum          int    `json:"need_num"`
	CompleteNum      int    `json:"complete_num"`
	ReviewNum        int    `json:"review_num"`
	Unit             string `json:"unit"`
}

package rsp

type LogListRsp struct {
	Total int64     `json:"total"`
	List  []LogList `json:"list"`
}

type LogList struct {
	Id          int    `json:"id"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	Number      string `json:"number"`
	BatchId     int    `json:"batch_id"`
	PickId      int    `json:"pick_id"`
	Status      int    `json:"status"`
	RequestXml  string `json:"request_xml"`
	ResponseXml string `json:"response_xml"`
	ResponseNo  string `json:"response_no"`
	Msg         string `json:"msg"`
	ShopName    string `json:"shop_name"`
}

type LogDetail struct {
	Id               int     `json:"id"`
	UpdateTime       string  `json:"update_time"`
	PickId           int     `json:"pick_id"`
	BatchId          int     `json:"batch_id"`
	PrePickGoodsId   int     `json:"pre_pick_goods_id"`
	OrderGoodsId     int     `json:"order_goods_id"`
	Number           string  `json:"number"`
	ShopId           int     `json:"shop_id"`
	DistributionType int     `json:"distribution_type"`
	Sku              string  `json:"sku"`
	GoodsName        string  `json:"goods_name"`
	GoodsType        string  `json:"goods_type"`
	GoodsSpe         string  `json:"goods_spe"`
	Shelves          string  `json:"shelves"`
	DiscountPrice    float64 `json:"discount_price"`
	NeedNum          int     `json:"need_num"`
	CompleteNum      int     `json:"complete_num"`
	ReviewNum        int     `json:"review_num"`
	Unit             string  `json:"unit"`
}

type LogDetailRsp struct {
	ShopName       string      `json:"shop_name"`
	PickUser       string      `json:"pick_user"`
	TakeOrdersTime string      `json:"take_orders_time"`
	ReviewUser     string      `json:"review_user"`
	ReviewTime     string      `json:"review_time"`
	PayAt          string      `json:"pay_at"`
	List           []LogDetail `json:"list"`
}

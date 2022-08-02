package rsp

type CommodityListRsp struct {
	Code int              `json:"code"`
	Data []*CommodityList `json:"data"`
	Msg  string           `json:"msg"`
}

type CommodityList struct {
	Name string `json:"name"`
	Sku  string `json:"sku"`
}

type ApiGoodsListRsp struct {
	Code int             `json:"code"`
	Data []*ApiGoodsList `json:"data"`
	Msg  string          `json:"msg"`
}

//type ApiGoodsData struct {
//	List []*ApiGoodsList `json:"list"`
//}

type ApiGoodsList struct {
	Id               int     `json:"id"`
	ShopId           int     `json:"shop_id"`
	ShopName         string  `json:"shop_name"`
	ShopType         string  `json:"shop_type"`
	ShopCode         string  `json:"shop_code"`
	HouseCode        string  `json:"house_code"`
	Line             string  `json:"line"`
	Number           string  `json:"number"`
	Status           int     `json:"status"`
	DeliveryAt       string  `json:"delivery_at"`
	DistributionType int     `json:"distribution_type"`
	OrderRemark      string  `json:"order_remark"`
	Province         string  `json:"province"`
	City             string  `json:"city"`
	District         string  `json:"district"`
	Address          string  `json:"address"`
	ConsigneeName    string  `json:"consignee_name"`
	ConsigneeTel     string  `json:"consignee_tel"`
	Name             string  `json:"name"`
	Sku              string  `json:"sku"`
	GoodsSpe         string  `json:"goods_spe"`
	GoodsType        string  `json:"goods_type"`
	Shelves          string  `json:"shelves"`
	OriginalPrice    int     `json:"original_price"`
	DiscountPrice    float64 `json:"discount_price"`
	GoodsUnit        string  `json:"goods_unit"`
	SaleUnit         string  `json:"sale_unit"`
	SaleCode         string  `json:"sale_code"`
	PayCount         int     `json:"pay_count"`
	CloseCount       int     `json:"close_count"`
	OutCount         int     `json:"out_count"`
	GoodsRemark      string  `json:"goods_remark"`
	PickStatus       int     `json:"pick_status"`
	PayAt            string  `json:"pay_at"`
	LackCount        int     `json:"lack_count"`
}

type OrderList struct {
	Number           string `json:"number"`
	PayAt            string `json:"pay_time"`
	ShopCode         string `json:"shop_code"`
	ShopName         string `json:"shop_name"`
	ShopType         string `json:"shop_type"`
	DistributionType int    `json:"distribution_type"` //配送方式
	SaleUnit         string `json:"sale_unit"`         //销售单位
	PayCount         int    `json:"pay_count"`         //下单数量
	Line             string `json:"line"`
	Region           string `json:"region"`
	OrderRemark      string `json:"order_remark"` //订单备注
}

type OrderDetail struct {
	Number      string             `json:"number"`
	PayAt       string             `json:"pay_time"`
	ShopCode    string             `json:"shop_code"`
	ShopName    string             `json:"shop_name"`
	OrderRemark string             `json:"order_remark"` //订单备注
	Line        string             `json:"line"`
	Region      string             `json:"region"`
	ShopType    string             `json:"shop_type"`
	Detail      map[string]*Detail `json:"detail"`
}

type Detail struct {
	Total int            `json:"total"`
	List  []*GoodsDetail `json:"list"`
}

type GoodsDetail struct {
	Name        string `json:"name"`
	GoodsSpe    string `json:"goods_spe"`
	Shelves     string `json:"shelves"`
	PayCount    int    `json:"pay_count"`
	CloseCount  int    `json:"close_count"`
	LackCount   int    `json:"need_count"` //欠货数量(需出数量)
	GoodsRemark string `json:"goods_remark"`
}

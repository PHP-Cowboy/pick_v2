package rsp

type ApiGoodsListRsp struct {
	Code int           `json:"code"`
	Data *ApiGoodsData `json:"data"`
}

type ApiGoodsData struct {
	List []*ApiGoodsList `json:"list"`
}

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
	LockCount        int     `json:"lock_count"`
}

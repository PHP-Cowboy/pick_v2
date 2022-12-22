package model

type GoodsJoinOrder struct {
	TaskId            int      `json:"task_id"`
	Number            string   `json:"number"`
	OrderId           int      `json:"order_id"`
	PayAt             *MyTime  `json:"pay_at"`
	ShopId            int      `json:"shop_id"`
	ShopName          string   `json:"shop_name"`
	ShopType          string   `json:"shop_type"`
	ShopCode          string   `json:"shop_code"`
	HouseCode         string   `json:"house_code"`
	DistributionType  int      `json:"distribution_type"`
	Line              string   `json:"line"`
	Province          string   `json:"province"`
	City              string   `json:"city"`
	District          string   `json:"district"`
	Address           string   `json:"address"`
	ConsigneeName     string   `json:"consignee_name"`
	ConsigneeTel      string   `json:"consignee_tel"`
	OrderType         int      `json:"order_type"`
	LatestPickingTime *MyTime  `json:"latest_picking_time"`
	HasRemark         int      `json:"has_remark"`
	OrderRemark       string   `json:"order_remark"`
	Sku               string   `json:"sku"`
	OrderGoodsId      int      `json:"order_goods_id"` //订单商品表id
	BatchId           int      `json:"batch_id"`
	GoodsName         string   `json:"goods_name"`
	GoodsType         string   `json:"goods_type"`
	GoodsSpe          string   `json:"goods_spe"`
	Shelves           string   `json:"shelves"`
	DiscountPrice     int      `json:"discount_price"`
	GoodsUnit         string   `json:"goods_unit"`
	SaleUnit          string   `json:"sale_unit"`
	SaleCode          string   `json:"sale_code"`
	PayCount          int      `json:"pay_count"`
	CloseCount        int      `json:"close_count"`
	LackCount         int      `json:"lack_count"`
	OutCount          int      `json:"out_count"`
	LimitNum          int      `json:"limit_num"`
	GoodsRemark       string   `json:"goods_remark"`
	Status            int      `json:"status"`
	DeliveryOrderNo   GormList `json:"delivery_order_no"`
	DeliveryAt        MyTime   `json:"delivery_at"`
}

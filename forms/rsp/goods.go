package rsp

import "pick_v2/model"

type GoodsListRsp struct {
	Total int64   `json:"total"`
	List  []Order `json:"list"`
}

type Order struct {
	Number            string `json:"number"`
	PayAt             string `json:"pay_time"`
	ShopCode          string `json:"shop_code"`
	ShopName          string `json:"shop_name"`
	ShopType          string `json:"shop_type"`
	DistributionType  int    `json:"distribution_type"` //配送方式
	PayCount          int    `json:"pay_count"`         //下单数量
	Picked            int    `json:"picked"`
	UnPicked          int    `json:"un_picked"`
	CloseNum          int    `json:"close_num"`
	Line              string `json:"line"`
	Region            string `json:"region"`
	OrderRemark       string `json:"order_remark"` //订单备注
	OrderType         int    `json:"order_type"`   //1:新订单,2:拣货中,3:欠货单
	LatestPickingTime string `json:"latest_picking_time"`
}

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
	Code int          `json:"code"`
	Data ApiGoodsData `json:"data"`
	Msg  string       `json:"msg"`
}

type ApiGoodsData struct {
	List  []*ApiGoods `json:"list"`
	Count int         `json:"count"`
}

type ApiGoods struct {
	Id                int    `json:"id"`
	ShopId            int    `json:"shop_id"`
	ShopName          string `json:"shop_name"`
	ShopType          string `json:"shop_type"`
	ShopCode          string `json:"shop_code"`
	HouseCode         string `json:"house_code"`
	Line              string `json:"line"`
	Number            string `json:"number"`
	Status            int    `json:"status"`
	DeliveryAt        string `json:"delivery_at"`
	DistributionType  int    `json:"distribution_type"`
	OrderRemark       string `json:"order_remark"`
	Province          string `json:"province"`
	City              string `json:"city"`
	District          string `json:"district"`
	Address           string `json:"address"`
	ConsigneeName     string `json:"consignee_name"`
	ConsigneeTel      string `json:"consignee_tel"`
	Name              string `json:"name"`
	Sku               string `json:"sku"`
	GoodsSpe          string `json:"goods_spe"`
	GoodsType         string `json:"goods_type"`
	SecondType        string `json:"second_type"`
	Shelves           string `json:"shelves"`
	OriginalPrice     int    `json:"original_price"`
	DiscountPrice     int    `json:"discount_price"`
	GoodsUnit         string `json:"goods_unit"`
	SaleUnit          string `json:"sale_unit"`
	SaleCode          string `json:"sale_code"`
	PayCount          int    `json:"pay_count"`
	CloseCount        int    `json:"close_count"`
	OutCount          int    `json:"out_count"`
	GoodsRemark       string `json:"goods_remark"`
	PickStatus        int    `json:"pick_status"`
	PayAt             string `json:"pay_at"`
	LackCount         int    `json:"lack_count"`
	LatestPickingTime string `json:"latest_picking_time"`
}

type OrderList struct {
	Number            string `json:"number"`
	PayAt             string `json:"pay_time"`
	ShopCode          string `json:"shop_code"`
	ShopName          string `json:"shop_name"`
	ShopType          string `json:"shop_type"`
	DistributionType  int    `json:"distribution_type"` //配送方式
	SaleUnit          string `json:"sale_unit"`         //销售单位
	PayCount          int    `json:"pay_count"`         //下单数量
	OutCount          int    `json:"out_count"`
	LackCount         int    `json:"lack_count"`
	Line              string `json:"line"`
	Region            string `json:"region"`
	OrderRemark       string `json:"order_remark"` //订单备注
	LatestPickingTime string `json:"latest_picking_time"`
}

type OrderDetail struct {
	Number          string             `json:"number"`
	PayAt           string             `json:"pay_time"`
	ShopCode        string             `json:"shop_code"`
	ShopName        string             `json:"shop_name"`
	OrderRemark     string             `json:"order_remark"` //订单备注
	Line            string             `json:"line"`
	Region          string             `json:"region"`
	ShopType        string             `json:"shop_type"`
	Detail          map[string]*Detail `json:"detail"`
	DeliveryOrderNo []string           `json:"delivery_order_no"`
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

type OrderShippingRecordRsp struct {
	List []OrderShippingRecord `json:"list"`
}

type OrderShippingRecord struct {
	Id              int            `json:"id"`
	TakeOrdersTime  string         `json:"take_orders_time"`
	PickUser        string         `json:"pick_user"`
	ReviewUser      string         `json:"review_user"`
	ReviewTime      string         `json:"review_time"`
	ReviewNum       int            `json:"review_num"`
	DeliveryOrderNo model.GormList `json:"delivery_order_no"`
}

type ShippingRecordDetailRsp struct {
}

type CompleteOrderRsp struct {
	Total int64           `json:"total"`
	List  []CompleteOrder `json:"list"`
}

type CompleteOrder struct {
	Number      string `json:"number"`
	PayAt       string `json:"pay_at"`
	ShopCode    string `json:"shop_code"`
	ShopName    string `json:"shop_name"`
	ShopType    string `json:"shop_type"`
	PayCount    int    `json:"pay_count"`
	OutCount    int    `json:"out_count"`
	CloseCount  int    `json:"close_count"`
	Line        string `json:"line"`
	Region      string `json:"region"`
	PickTime    string `json:"pick_time"`
	OrderRemark string `json:"order_remark"`
}

type CompleteOrderDetailRsp struct {
	ShopName        string                    `json:"shop_name"`
	ShopCode        string                    `json:"shop_code"`
	Line            string                    `json:"line"`
	Region          string                    `json:"regin"`
	ShopType        string                    `json:"shop_type"`
	Number          string                    `json:"number"`
	OrderRemark     string                    `json:"order_remark"`
	DeliveryOrderNo model.GormList            `json:"delivery_order_no"`
	Goods           map[string][]PrePickGoods `json:"goods"`
}

type CountRes struct {
	AllCount   int `json:"all_count"`
	NewCount   int `json:"new_count"`
	OldCount   int `json:"old_count"`
	PickCount  int `json:"pick_count"`
	CloseCount int `json:"close_count"`
}

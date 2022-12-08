package req

type PickOrderListForm struct {
	Paging
	OrderType        int    `json:"order_type" form:"order_type" binding:"oneof=0 1 2 3 4"`
	ShopId           int    `json:"shop_id" form:"shop_id"`
	Sku              string `json:"sku" form:"sku"`
	Number           string `json:"number" form:"number"`
	Lines            string `json:"lines" form:"lines"`
	DistributionType int    `json:"distribution_type" form:"distribution_type"` //1-公司配送 2-用户自提 3-三方物流 4-快递配送 5-首批物料|设备单
	ShopType         string `json:"shop_type" form:"shop_type"`
	Province         string `json:"province" form:"province"`
	City             string `json:"city" form:"city"`
	District         string `json:"district" form:"district"`
	HasRemark        int    `json:"has_remark" form:"has_remark"` ////是否存在备注  0-默认全部 1-有备注 2-无备注
	PayEndTime       string `json:"pay_end_time"  form:"pay_end_time"`
}

type PickOrderCountForm struct {
	ShopId           int    `json:"shop_id" form:"shop_id"`
	Sku              string `json:"sku" form:"sku"`
	Number           string `json:"number" form:"number"`
	Lines            string `json:"lines" form:"lines"`
	DistributionType int    `json:"distribution_type" form:"distribution_type"` //1-公司配送 2-用户自提 3-三方物流 4-快递配送 5-首批物料|设备单
	ShopType         string `json:"shop_type" form:"shop_type"`
	Province         string `json:"province" form:"province"`
	City             string `json:"city" form:"city"`
	District         string `json:"district" form:"district"`
	HasRemark        int    `json:"has_remark" form:"has_remark"` ////是否存在备注  0-默认全部 1-有备注 2-无备注
	PayEndTime       string `json:"pay_end_time"  form:"pay_end_time"`
}

type OrderGoodsListForm struct {
	PickNumber string `json:"pick_number" form:"pick_number"`
}

type RestrictedShipmentForm struct {
	PickNumber string           `json:"pick_number" form:"pick_number"`
	Params     []ShipmentParams `json:"params" form:"params"`
}

type ShipmentParams struct {
	Id       int `json:"id"`
	LimitNum int `json:"limit_num"`
}

type BatchRestrictedShipmentForm struct {
	Sku      string `json:"sku" form:"sku" binding:"required"`
	LimitNum int    `json:"limit_num"  binding:"required"`
}

type GoodsNumForm struct {
	Sku string `json:"sku" form:"sku"`
}

type RestrictedShipmentListForm struct {
	Paging
}

type RevokeRestrictedShipmentForm struct {
	PickOrderGoodsId int `json:"pick_order_goods_id"`
}

type CloseOrderForm struct {
	OrderId int `json:"order_id"`
}

type CloseOrderGoodsForm struct {
	GoodsId  int `json:"goods_id" binding:"required"`
	CloseNum int `json:"close_num" binding:"required"`
}

type GetPickOrderDetailForm struct {
	Id int `json:"id" form:"id" binding:"required"`
}

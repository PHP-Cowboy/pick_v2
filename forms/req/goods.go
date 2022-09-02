package req

type GoodsListForm struct {
	Paging
	OrderType        int    `json:"order_type" form:"order_type" binding:"oneof=0 1 2 3"`
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

type GetOrderDetailForm struct {
	Number string `json:"number" form:"number" binding:"required"`
}

type OrderShippingRecordReq struct {
	DeliveryOrderNo []string `json:"delivery_order_no" form:"delivery_order_no"`
}

type ShippingRecordDetailReq struct {
	Id int `json:"id" form:"id"`
}

type CompleteOrderReq struct {
	Paging
	Number         string `json:"number"` //订单号
	ShopId         int    `json:"shop_id"`
	Sku            string `json:"sku"`
	Line           string `json:"line"`
	DeliveryMethod int    `json:"delivery_method"`
	ShopType       string `json:"shop_type"` //门店类型
	Province       string `json:"province"`
	City           string `json:"city"`
	District       string `json:"district"`
	IsRemark       int    `json:"is_remark"`
	PayAt          string `json:"pay_at"`
}

type CompleteOrderDetailReq struct {
	Number string `json:"number" form:"number" binding:"required"` //订单号
}

package req

// MQ 采购订单
type PurchaseOrderForm struct {
	OrderId []int `json:"order_id"`
}

type CloseOrderInfo struct {
	OccId int `json:"occ_id"`
}

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
	PayEndTime       string `json:"pay_end_time" form:"pay_end_time"`
	PayAtSort        string `json:"pay_at_sort" form:"pay_at_sort" binding:"required"`
	ShopCodeSort     string `json:"shop_code_sort" form:"shop_code_sort"`
}

type GoodsSort struct {
	Field string `json:"field" form:"field"`
	Rule  string `json:"rule" form:"rule"`
}

type GetOrderDetailForm struct {
	Number string `json:"number" form:"number" binding:"required"`
	IsLack int    `json:"is_lack" form:"is_lack"`
}

type CommodityListForm struct {
	Name string `json:"name" form:"name"`
}

type OrderShippingRecordReq struct {
	DeliveryOrderNo []string `json:"delivery_order_no" form:"delivery_order_no"`
}

type ChangeDistributionForm struct {
	Ids              []int `json:"ids" binding:"required"`
	DistributionType int   `json:"distribution_type" binding:"required"`
}

type ShippingRecordDetailReq struct {
	Id     int    `json:"id" form:"id"`
	Number string `json:"number" form:"number"`
}

type CompleteOrderForm struct {
	Paging
	Number         string `json:"number" form:"number"` //订单号
	ShopId         int    `json:"shop_id" form:"shop_id"`
	Sku            string `json:"sku" form:"sku"`
	Lines          string `json:"lines" form:"lines"`
	DeliveryMethod int    `json:"delivery_method" form:"delivery_method"`
	ShopType       string `json:"shop_type" form:"shop_type"` //门店类型
	Province       string `json:"province" form:"province"`
	City           string `json:"city" form:"city"`
	District       string `json:"district" form:"district"`
	IsRemark       int    `json:"is_remark" form:"is_remark"`
	PayAt          string `json:"pay_at" form:"pay_at"`
}

type CompleteOrderDetailReq struct {
	Number string `json:"number" form:"number" binding:"required"` //订单号
}

type CreatePickOrderForm struct {
	Numbers []string `json:"numbers" form:"numbers" binding:"required"`
}

type CountFrom struct {
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
	PayEndTime       string `json:"pay_end_time" form:"pay_end_time"`
}

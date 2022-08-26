package req

// pay_range,batch_number,de_start,de_end
type GetGoodsListForm struct {
	Paging
	PickSta      *int   `json:"pick_sta" form:"pick_sta" binding:"required,oneof=0 1 2"`
	Number       string `json:"number" form:"number"`
	ShopId       int    `json:"shop_id" form:"shop_id"`
	Sku          string `json:"sku" form:"sku"`
	Lines        string `json:"lines" form:"lines"`
	DeType       int    `json:"de_type" form:"de_type"` //1-公司配送 2-用户自提 3-三方物流 4-快递配送 5-首批物料|设备单
	ShopType     string `json:"shop_type" form:"shop_type"`
	Province     string `json:"province" form:"province"`
	City         string `json:"city" form:"city"`
	District     string `json:"district" form:"district"`
	HasRemark    int    `json:"has_remark" form:"has_remark"` ////是否存在备注  0-默认全部 1-有备注 2-无备注
	PayStartTime string `json:"pay_start_time" form:"pay_start_time"`
	PayEndTime   string `json:"pay_end_time"  form:"pay_end_time"`
	Index        int    `json:"index"`
	GoodsName    string `json:"goods_name" form:"goods_name"` //
	ShopName     string `json:"shop_name" form:"shop_name"`   //
}

type GetOrderDetailForm struct {
	Number string `json:"number" form:"number" binding:"required"`
	IsLack bool   `json:"is_lack" form:"is_lack"`
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

package req

//pay_range,batch_number,de_start,de_end
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
	HasRemark    bool   `json:"has_remark" form:"has_remark"`
	PayStartTime string `json:"pay_start_time" form:"pay_start_time"`
	PayEndTime   string `json:"pay_end_time"  form:"pay_end_time"`
	Index        int    `json:"index"`
}

type GetOrderDetailForm struct {
	Number string `json:"number" form:"number" binding:"required"`
}

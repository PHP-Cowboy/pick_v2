package req

type CreateOutboundForm struct {
	OutboundName      string   `json:"outbound_name" form:"outbound_name" binding:"required"`
	DeliveryStartTime string   `json:"delivery_start_time" form:"delivery_start_time"`                //发货开始
	DeliveryEndTime   string   `json:"delivery_end_time" form:"delivery_end_time" binding:"required"` //发货截止
	Lines             []string `json:"lines" form:"lines"`                                            //线路
	DistributionType  int      `json:"distribution_type" form:"distribution_type" binding:"required"` //配送方式 1-公司配送 2-用户自提 3-三方物流 4-快递配送 5-首批物料|设备单
	PayTime           string   `json:"pay_time"  form:"pay_time" binding:"required"`                  //支付时间
	Sku               []string `json:"sku" form:"sku"`
	GoodsName         []string `json:"goods_name" form:"goods_name"`
	GoodsType         []string `json:"goods_type" form:"goods_type"`
}

type OutboundTaskListForm struct {
	Paging
	ShopId           int    `json:"shop_id" form:"shop_id"`
	Sku              string `json:"sku" form:"sku"`
	Number           string `json:"number" form:"number"`
	Line             string `json:"line" form:"line"`
	DistributionType int    `json:"distribution_type" form:"distribution_type"`
	Status           int    `json:"status" form:"status" binding:"required"`
	StartTime        string `json:"start_time" form:"start_time"`
	EndTime          string `json:"end_time" form:"end_time"`
}

type OutboundOrderListForm struct {
	Paging
	TaskId           int    `json:"task_id" form:"task_id" binding:"required"`
	ShopId           int    `json:"shop_id" form:"shop_id" `
	Number           string `json:"number" form:"number" `
	Sku              string `json:"sku" form:"sku" `
	Line             string `json:"line" form:"line" `
	DistributionType int    `json:"distribution_type" form:"distribution_type" `
	ShopType         string `json:"shop_type" form:"shop_type" `
	Province         string `json:"province" form:"province" `
	City             string `json:"city" form:"city" `
	District         string `json:"district" form:"district" `
	HasRemark        int    `json:"has_remark" form:"has_remark" `
	GoodsType        string `json:"goods_type" form:"goods_type"`
	OrderType        int    `json:"order_type" form:"order_type" `
}

type OutboundOrderCountForm struct {
	TaskId int `json:"task_id" form:"task_id" binding:"required"`
}

type GetTaskSkuNumForm struct {
	TaskId int    `json:"task_id" form:"task_id" binding:"required"`
	Sku    string `json:"sku" form:"sku" binding:"required"`
}

type OutboundOrderDetailForm struct {
	TaskId int    `json:"task_id" form:"task_id" binding:"required"`
	Number string `json:"number" form:"number" binding:"required"`
}

type OutboundOrderGoodsListForm struct {
	TaskId int `json:"task_id" form:"task_id" binding:"required"`
}

type OrderLimit struct {
	Sku      string `json:"sku" binding:"required"`
	LimitNum int    `json:"limit_num" binding:"required"`
}

type EndOutboundTaskForm struct {
	TaskId int `json:"task_id" form:"task_id" binding:"required"`
}

type OutboundTaskCloseOrderForm struct {
	Number []string `json:"number"`
}

type OutboundTaskAddOrderForm struct {
	TaskId int      `json:"task_id" binding:"required"`
	Number []string `json:"number" binding:"required"`
}

type OrderOutboundRecordForm struct {
}

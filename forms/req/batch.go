package req

type CreateBatchForm struct {
	WarehouseId       int     `json:"warehouse_id" form:"warehouse_id"`
	Lines             string  `json:"lines" form:"lines"`
	PayEndTime        string  `json:"pay_end_time"  form:"pay_end_time"`
	DeliveryStartTime string  `json:"delivery_start_time" form:"delivery_start_time"`
	DeliveryEndTime   string  `json:"delivery_end_time" form:"delivery_end_time"`
	DeType            int     `json:"de_type" form:"de_type"` //1-公司配送 2-用户自提 3-三方物流 4-快递配送 5-首批物料|设备单
	Goods             []Goods `json:"goods" form:"goods"`
	Sku               string  `json:"sku" form:"sku"`
	GoodsName         string  `json:"goods_name" form:"goods_name"`
	BatchNumber       string
}

type Goods struct {
	Sku  string `json:"sku" form:"sku"`
	Name string `json:"name" form:"name"`
}

type GetBatchListForm struct {
	Paging
	Status         int    `json:"status" form:"status" validate:"required"` //进行中，已结束
	ShopId         int    `json:"shop_id" form:"shop_id"`                   //店铺
	GoodsName      string `json:"goods_name" form:"goods_name"`
	Number         string `json:"number" form:"number"`
	Line           string `json:"line" form:"line"`
	DeliveryMethod int    `json:"delivery_method" form:"delivery_method"`
	CreateTime     string `json:"create_time" form:"create_time"`
	EndTime        string `json:"end_time" form:"end_time"`
}

type GetPrePickListForm struct {
	Paging
}

type GetBaseForm struct {
	BatchId int `json:"batch_id" form:"batch_id"`
}

type GetPrePickDetailForm struct {
	PrePickId int `json:"pre_pick_id" form:"pre_pick_id"`
}

type ToppingForm struct {
	Id int `json:"id" form:"id"`
}

type BatchPickForm struct {
	Ids []int `json:"ids"`
}

type MergePickForm struct {
}

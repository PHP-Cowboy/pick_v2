package req

type CreateBatchForm struct {
	WarehouseId       int     `json:"warehouse_id" form:"warehouse_id"`
	Number            string  `json:"number"`
	Lines             string  `json:"lines" form:"lines"`
	PayEndTime        string  `json:"pay_end_time"  form:"pay_end_time"`
	DeliveryStartTime string  `json:"delivery_start_time" form:"delivery_start_time"`
	DeliveryEndTime   string  `json:"delivery_end_time" form:"delivery_end_time"`
	DeType            int     `json:"de_type" form:"de_type"` //1-公司配送 2-用户自提 3-三方物流 4-快递配送 5-首批物料|设备单
	Goods             []Goods `json:"goods" form:"goods"`
	Sku               string  `json:"sku" form:"sku"`
	GoodsName         string  `json:"goods_name" form:"goods_name"`
	BatchNumber       string  `json:"batch_number"`
}

type Goods struct {
	Sku  string `json:"sku" form:"sku"`
	Name string `json:"name" form:"name"`
}

type StopPickForm struct {
	Id     int  `json:"id" binding:"required"`
	Status *int `json:"status" binding:"required,oneof=0 2"`
}

type EndBatchForm struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type OutGoods struct {
	BatchNumber string         `json:"batch_number"`
	List        []OutGoodsList `json:"list"`
}

type OutGoodsList struct {
	GoodsLogId   int    `json:"goods_log_id"`
	Number       string `json:"number"`
	OutNumber    string `json:"out_number"`
	CkNumber     string `json:"ck_number"`
	Sku          string `json:"sku"`
	Name         string `json:"name"`
	OutCount     int    `json:"out_count"`
	Price        int    `json:"price"`
	SumPrice     int    `json:"sum_price"`
	OutAt        string `json:"out_at"`
	PayAt        string `json:"pay_at"`
	GoodsSpe     string `json:"goods_spe"`
	GoodsUnit    string `json:"goods_unit"`
	DeliveryType int    `json:"delivery_type"`
}

type GetBatchListForm struct {
	Paging
	Status         *int   `json:"status" form:"status" binding:"required"` //进行中，已结束
	ShopId         int    `json:"shop_id" form:"shop_id"`                  //店铺
	GoodsName      string `json:"goods_name" form:"goods_name"`
	Number         string `json:"number" form:"number"`
	Line           string `json:"line" form:"line"`
	DeliveryMethod int    `json:"delivery_method" form:"delivery_method"`
	CreateTime     string `json:"create_time" form:"create_time"`
	EndTime        string `json:"end_time" form:"end_time"`
}

type GetPrePickListForm struct {
	Paging
	ShopId int    `json:"shop_id"`
	Line   string `json:"line"`
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
	BatchId     int      `json:"batch_id" binding:"required"`
	Ids         []int    `json:"ids" binding:"required"`
	Type        int      `json:"type" binding:"required,oneof=1 2 3"`
	TypeParam   []string `json:"type_param" binding:"required"`
	WarehouseId int      `json:"warehouse_id"`
}

type MergePickForm struct {
	BatchId     int      `json:"batch_id" binding:"required"`
	Ids         []int    `json:"ids" binding:"required"`
	Type        int      `json:"type" binding:"required,oneof=1 2 3"`
	TypeParam   []string `json:"type_param" binding:"required"`
	TaskName    string   `json:"task_name" binding:"required"`
	WarehouseId int      `json:"warehouse_id"`
}

type PrintCallGetReq struct {
	HouseCode string `json:"house_code" binding:"required"`
}

type GetPoolNumReq struct {
	BatchId int `json:"batch_id" form:"batch_id" binding:"required"`
}

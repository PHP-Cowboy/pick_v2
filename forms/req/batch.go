package req

type CreateBatchForm struct {
	BatchName         string   `json:"batch_name" form:"batch_name"`
	DeliveryStartTime string   `json:"delivery_start_time" form:"delivery_start_time"`                //发货开始
	DeliveryEndTime   string   `json:"delivery_end_time" form:"delivery_end_time" binding:"required"` //发货截止
	Lines             []string `json:"lines" form:"lines" binding:"required"`                         //线路
	DistributionType  int      `json:"distribution_type" form:"distribution_type" binding:"required"` //配送方式 1-公司配送 2-用户自提 3-三方物流 4-快递配送 5-首批物料|设备单
	PayTime           string   `json:"pay_time"  form:"pay_time" binding:"required"`                  //支付时间
	Sku               []string `json:"sku" form:"sku"`
	GoodsName         []string `json:"goods_name" form:"goods_name"`
}

type StopPickForm struct {
	Id     int  `json:"id" binding:"required"`
	Status *int `json:"status" binding:"required,oneof=0 2"`
}

type CreateByOrderReq struct {
	PickNumber string `json:"pick_number" form:"pick_number" binding:"required"`
}

type EndBatchForm struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type EditBatchForm struct {
	Id        int    `json:"id" form:"id" binding:"required"`
	BatchName string `json:"batch_name"`
}

type GetBatchOrderAndGoodsForm struct {
	Id int `json:"id" form:"id" binding:"required"`
}

type GetBatchListForm struct {
	Paging
	Status         *int   `json:"status" form:"status" binding:"required"` //进行中，已结束
	ShopId         int    `json:"shop_id" form:"shop_id"`                  //店铺
	Sku            string `json:"sku" form:"sku"`
	Number         string `json:"number" form:"number"`
	Line           string `json:"line" form:"line"`
	DeliveryMethod int    `json:"delivery_method" form:"delivery_method"`
	CreateTime     string `json:"create_time" form:"create_time"`
	EndTime        string `json:"end_time" form:"end_time"`
}

type GetPrePickListForm struct {
	Paging
	BatchId int    `json:"batch_id" form:"batch_id"`
	ShopId  int    `json:"shop_id" form:"shop_id"`
	Line    string `json:"line" form:"line"`
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
	HouseCode string `json:"house_code" form:"house_code" binding:"required"`
}

type GetPoolNumReq struct {
	BatchId int `json:"batch_id" form:"batch_id" binding:"required"`
}

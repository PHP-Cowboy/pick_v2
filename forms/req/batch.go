package req

type CreateBatchByTaskForm struct {
	TaskId   int    `json:"task_id" binding:"required"`
	TaskName string `json:"task_name" binding:"required"`
	Typ      int    `json:"typ"` //1:常规批次,2:快递批次
}

type NewCreateBatchForm struct {
	TaskId    int      `json:"task_id" binding:"required"`
	Number    []string `json:"number" binding:"required"`
	BatchName string   `json:"batch_name" binding:"required"`
	Typ       int      `json:"typ"` //1:常规批次,2:快递批次,3:后台拣货
}

type CentralizedPickListForm struct {
	Paging
	BatchId   int    `json:"batch_id" form:"batch_id" binding:"required"`
	GoodsName string `json:"goods_name" form:"goods_name"`
	GoodsType string `json:"goods_type" form:"goods_type"`
	IsRemark  int    `json:"is_remark" form:"is_remark"`
}

type CentralizedPickDetailForm struct {
	BatchId int    `json:"batch_id" form:"batch_id" binding:"required"`
	Sku     string `json:"sku" form:"sku" binding:"required"`
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
	Typ            int    `json:"typ" form:"typ"`
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
	Typ         int      `json:"typ"` //类型：1:常规批次,2:快递批次;前端不需要传,程序里赋值
}

type MergePickForm struct {
	BatchId     int      `json:"batch_id" binding:"required"`
	Ids         []int    `json:"ids" binding:"required"`
	Type        int      `json:"type" binding:"required,oneof=1 2 3"`
	TypeParam   []string `json:"type_param" binding:"required"`
	TaskName    string   `json:"task_name" binding:"required"`
	WarehouseId int      `json:"warehouse_id"`
	Typ         int      `json:"typ"` //类型：1:常规批次,2:快递批次;前端不需要传,程序里赋值
}

type BatchChangeBatchForm struct {
	Status *int `json:"status" binding:"required,oneof=0 2"`
}

type PrintCallGetReq struct {
	HouseCode string `json:"house_code" form:"house_code" binding:"required"`
	Typ       int    `json:"typ" form:"typ" ` //默认 仓库拣货打印 1:办公室拣货[后台拣货/快递拣货]打印
}

type GetPoolNumReq struct {
	BatchId int `json:"batch_id" form:"batch_id" binding:"required"`
}

type GetBatchPoolNumForm struct {
	Typ int `json:"typ" form:"typ"`
}

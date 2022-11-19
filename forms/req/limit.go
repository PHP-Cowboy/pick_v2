package req

type OrderLimitForm struct {
	TaskId     int          `json:"task_id" binding:"required"`
	Number     string       `json:"number" binding:"required"`
	OrderLimit []OrderLimit `json:"order_limit" binding:"required"`
}

type TaskLimitForm struct {
	TaskId   int    `json:"task_id" form:"task_id" binding:"required"`
	Sku      string `json:"sku" form:"sku" binding:"required"`
	LimitNum int    `json:"limit_num" form:"limit_num" binding:"required"`
}

type RevokeLimitForm struct {
	TaskId int    `json:"task_id" form:"task_id" binding:"required"`
	Sku    string `json:"sku" form:"sku" binding:"required"`
	Number string `json:"number" form:"number" binding:"required"`
}

type LimitShipmentListForm struct {
	TaskId int `json:"task_id" form:"task_id" binding:"required"`
}

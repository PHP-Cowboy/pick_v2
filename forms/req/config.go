package req

type DeliveryMethodInfoForm struct {
	TaskId int    `json:"task_id" form:"task_id" binding:"required"`
	Number string `json:"number" form:"number" binding:"required"`
}

type ChangeDeliveryMethodForm struct {
	TaskId         int    `json:"task_id" binding:"required"`
	Number         string `json:"number" form:"number" binding:"required"`
	DeliveryMethod int    `json:"delivery_method"  binding:"required"`
}

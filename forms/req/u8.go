package req

type LogListForm struct {
	Paging
	Status    int    `json:"status" form:"status"`
	StartTime string `json:"start_time" form:"start_time"`
	EndTime   string `json:"end_time" form:"end_time"`
}

type BatchSupplementForm struct {
	Ids []int `json:"ids"`
}

type LogDetailForm struct {
	Number string `json:"number" form:"number" binding:"required"`
	PickId int    `json:"pick_id" form:"pick_id" binding:"required"`
}

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
	BatchId int    `json:"batch_id"`
	Number  string `json:"number"`
}

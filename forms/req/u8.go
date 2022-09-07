package req

type LogListForm struct {
	Paging
	Status    int    `json:"status"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

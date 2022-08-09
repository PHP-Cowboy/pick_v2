package req

type CreateMenuReq struct {
	Type  int    `json:"type" binding:"required"`
	Title string `json:"title" binding:"required"`
	Path  string `json:"path"`
	PId   int    `json:"p_id"`
	Sort  int    `json:"sort"`
}

type ChangeMenuReq struct {
	Id    int    `json:"id" binding:"required"`
	Type  int    `json:"type" binding:"required"`
	Title string `json:"title" binding:"required"`
	Path  string `json:"path"`
	PId   int    `json:"p_id"`
	Sort  int    `json:"sort"`
}

type BatchDeleteMenuReq struct {
	Ids []int `json:"ids"`
}

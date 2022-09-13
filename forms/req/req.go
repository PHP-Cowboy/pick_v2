package req

type Paging struct {
	Page int `form:"page" json:"page" binding:"required,gt=0"`
	Size int `form:"size" json:"size" binding:"required,gt=0,lte=500"`
}

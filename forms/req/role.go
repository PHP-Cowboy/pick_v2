package req

//创建角色
type CreateRoleForm struct {
	Name string `json:"name"`
}

//更新角色
type ChangeRoleForm struct {
	Id       int  `json:"id" binding:"required"`
	IsDelete bool `json:"is_delete"`
}

type GetRoleListForm struct {
	Paging
	Name string `json:"name"`
}

type BatchDeleteRoleForm struct {
	Ids []int `json:"ids" binding:"required"`
}

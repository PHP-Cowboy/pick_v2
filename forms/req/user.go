package req

type AddUserForm struct {
	WarehouseId int    `form:"warehouse_id" json:"warehouse_id" binding:"required"`
	RoleId      int    `form:"role_id" json:"role_id" binding:"required"`
	Name        string `form:"name" json:"name" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
}

type GetUserListForm struct {
	Paging
	WarehouseId int `form:"warehouse_id" json:"warehouse_id"`
}

type LoginForm struct {
	WarehouseId int    `form:"warehouse_id" json:"warehouse_id" binding:"required"`
	Account     string `json:"account" form:"account"`
	Password    string `json:"password" form:"password"`
}

type CheckPwdForm struct {
	Id          int    `json:"id" binding:"required"`
	NewPassword string `json:"new_password"`
	Name        string `json:"name"`
	Status      int    `json:"status"`
	RoleId      int    `json:"role_id"`
}

type BatchDeleteUserForm struct {
	Ids []int `json:"ids" binding:"required"`
}

type WarehouseUserCountForm struct {
}

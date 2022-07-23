package req

type AddUserForm struct {
	WarehouseId int    `form:"warehouse_id" json:"warehouse_id" binding:"required"`
	RoleId      int    `form:"role_id" json:"role_id" binding:"required"`
	Name        string `form:"name" json:"name" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
}

type GetUserListForm struct {
	Page int `form:"page" json:"page" binding:"required,gt=0"`
	Size int `form:"size" json:"size" binding:"required,gt=0,lte=100"`
}

type LoginForm struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type CheckPwdForm struct {
	Id          int    `json:"id" binding:"required"`
	NewPassword string `json:"new_password"`
	Name        string `json:"name"`
	Status      int    `json:"status"`
	RoleId      int    `json:"role_id"`
}

type WarehouseUserCountForm struct {
	WarehouseId int    `form:"warehouse_id" json:"warehouse_id" binding:"required"`
}

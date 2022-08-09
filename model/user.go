package model

const (
	Other    = iota
	Admin    //管理员
	Picker   //拣货员
	Reviewer //复核员
)

//用户
type User struct {
	Base
	Account     string `gorm:"type:varchar(64);index:unique;not null;comment:账号"`
	Password    string `gorm:"type:varchar(100);not null;comment:密码"`
	Name        string `gorm:"type:varchar(16);not null;comment:姓名"`
	RoleId      int    `gorm:"type:tinyint;not null;comment:角色表id"`
	Role        string `gorm:"type:varchar(16);not null;comment:角色(岗位)"`
	Status      int    `gorm:"type:tinyint;not null;default:1;comment:状态:0:未知,1:正常,2:禁用"`
	WarehouseId int    `gorm:"type:int(11);comment:仓库"`
}

type Role struct {
	Base
	Name string `gorm:"type:varchar(32);unique;not null;comment:角色名"`
}

type UserRole struct {
	Base
	UserId int `gorm:"not null;comment:用户id"`
	RoleId int `gorm:"not null;comment:角色id"`
}

type Menu struct {
	Base
	Type  int    `gorm:"comment:1:菜单,2:功能"`
	Title string `gorm:"comment:菜单名称"`
	Path  string `gorm:"comment:路由地址"`
	PId   int    `gorm:"上级权限id"`
	Sort  int    `gorm:"排序"`
}

type RoleMenu struct {
	Base
	RoleId int `gorm:"not null;comment:角色id"`
	MenuId int `gorm:"not null;comment:菜单id"`
}

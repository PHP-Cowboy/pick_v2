package model

// 用户
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

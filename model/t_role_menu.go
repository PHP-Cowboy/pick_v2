package model

type RoleMenu struct {
	Base
	RoleId int `gorm:"type:int(11);not null;comment:角色id"`
	MenuId int `gorm:"type:int(11);not null;comment:菜单id"`
}

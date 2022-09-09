package model

type RoleMenu struct {
	Base
	RoleId int `gorm:"not null;comment:角色id"`
	MenuId int `gorm:"not null;comment:菜单id"`
}

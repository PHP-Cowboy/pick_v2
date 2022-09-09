package model

type UserRole struct {
	Base
	UserId int `gorm:"not null;comment:用户id"`
	RoleId int `gorm:"not null;comment:角色id"`
}

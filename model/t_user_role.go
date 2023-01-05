package model

type UserRole struct {
	Base
	UserId int `gorm:"type:int(11);not null;comment:用户id"`
	RoleId int `gorm:"type:int(11);not null;comment:角色id"`
}

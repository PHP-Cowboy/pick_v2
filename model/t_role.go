package model

type Role struct {
	Base
	Name string `gorm:"type:varchar(32);unique;not null;comment:角色名"`
}

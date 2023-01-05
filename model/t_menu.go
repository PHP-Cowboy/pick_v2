package model

type Menu struct {
	Base
	Type  int    `gorm:"type:tinyint;comment:1:菜单,2:功能"`
	Title string `gorm:"type:varchar(64);comment:菜单名称"`
	Path  string `gorm:"type:varchar(255);comment:路由地址"`
	PId   int    `gorm:"type:int(11) unsigned;上级权限id"`
	Sort  int    `gorm:"type:int(11) unsigned;comment:排序"`
}

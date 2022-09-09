package model

type Menu struct {
	Base
	Type  int    `gorm:"comment:1:菜单,2:功能"`
	Title string `gorm:"comment:菜单名称"`
	Path  string `gorm:"comment:路由地址"`
	PId   int    `gorm:"上级权限id"`
	Sort  int    `gorm:"排序"`
}

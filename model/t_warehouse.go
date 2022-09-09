package model

// 仓库
type Warehouse struct {
	Base
	WarehouseName string `gorm:"type:varchar(64);not null;comment:仓库名称"`
	Abbreviation  string `gorm:"type:varchar(64);not null;comment:仓库简称"`
}

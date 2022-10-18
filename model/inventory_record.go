package model

type InventoryRecord struct {
	Base
	OrderNo      string `gorm:"index;comment:任务编号"`
	Sku          string `gorm:"index;type:varchar(64);comment:sku"`
	GoodsName    string `gorm:"type:varchar(64);comment:商品名称"`
	GoodsSpe     string `gorm:"type:varchar(64);comment:商品规格"`
	GoodsUnit    string `gorm:"type:varchar(32);comment:商品单位"`
	InventoryNum int    `gorm:"not null;default:0;comment:盘点数量"`
	UserName     string `gorm:"type:varchar(16);not null;default:'';comment:盘点人"`
	IsDelete     int    `gorm:"type:tinyint;default:1;comment:1:正常,2:删除"`
}

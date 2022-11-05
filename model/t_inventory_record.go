package model

// 用户盘点记录表
type InventoryRecord struct {
	Base
	SelfBuiltId  int     `gorm:"type:int(11);comment:自建盘点任务ID"`
	Sku          string  `gorm:"type:varchar(64);comment:sku"`
	UserName     string  `gorm:"type:varchar(16);not null;comment:盘点人"`
	GoodsName    string  `gorm:"type:varchar(64);comment:商品名称"`
	GoodsSpe     string  `gorm:"type:varchar(64);comment:商品规格"`
	GoodsUnit    string  `gorm:"type:varchar(32);comment:商品单位"`
	InventoryNum float64 `gorm:"type:decimal(10,2);not null;default:0;comment:盘点数量"`
	IsDelete     int     `gorm:"type:tinyint;default:1;comment:1:正常,2:删除"`
	InvType      int     `gorm:"type:tinyint;default:1;comment:盘点类型:1:首次,2:复盘"`
}

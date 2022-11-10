package model

type InvRecordSum struct {
	SelfBuiltId  int     `gorm:"type:int(11);comment:自建盘点任务ID"`
	Sku          string  `gorm:"type:varchar(64);comment:sku"`
	InvType      int     `gorm:"type:tinyint;default:1;comment:盘点类型:1:首次,2:复盘"`
	InventoryNum float64 `gorm:"type:decimal(10,2);not null;default:0;comment:盘点数量"`
}

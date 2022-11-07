package model

import "time"

// 盘点任务商品记录 视图
type InvOrderSkuSum struct {
	SelfBuiltId  int       `gorm:"index;type:int(11);comment:自建盘点任务ID"`
	OrderNo      string    `gorm:"primaryKey;type:varchar(64);comment:任务编号"`
	Sku          string    `gorm:"primaryKey;type:varchar(64);comment:sku"`
	CreateTime   time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime   time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime   time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
	GoodsName    string    `gorm:"type:varchar(64);comment:商品名称"`
	GoodsType    string    `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe     string    `gorm:"type:varchar(64);comment:商品规格"`
	GoodsUnit    string    `gorm:"type:varchar(32);comment:商品单位"`
	BookNum      float64   `gorm:"type:decimal(10,2);default:0;comment:账面数量"`
	InventoryNum float64   `gorm:"type:decimal(10,2);default:0;comment:盘点数量"`
	InvType      int       `gorm:"type:tinyint;default:1;comment:盘点类型:1:首次,2:复盘"`
	IsDelete     int       `gorm:"type:tinyint;default:1;comment:1:正常,2:删除"`
}

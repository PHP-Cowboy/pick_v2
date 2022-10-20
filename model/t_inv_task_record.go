package model

import "time"

type InvTaskRecord struct {
	OrderNo    string    `gorm:"primaryKey;type:varchar(64);comment:任务编号"`
	Sku        string    `gorm:"primaryKey;type:varchar(64);comment:sku"`
	CreateTime time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
	GoodsName  string    `gorm:"type:varchar(64);comment:商品名称"`
	GoodsType  string    `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe   string    `gorm:"type:varchar(64);comment:商品规格"`
	GoodsUnit  string    `gorm:"type:varchar(32);comment:商品单位"`
	BookNum    float64   `gorm:"type:decimal(10,2);not null;default:0;comment:账面数量"`
}

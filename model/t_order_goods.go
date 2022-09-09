package model

import "time"

// 订单商品表
type OrderGoods struct {
	Id              int       `gorm:"primaryKey;type:int(11) unsigned;comment:id"`
	CreateTime      time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime      time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime      time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
	Number          string    `gorm:"type:varchar(64);index:number_sku_idx;comment:订单编号"`
	GoodsName       string    `gorm:"type:varchar(64);comment:商品名称"`
	Sku             string    `gorm:"type:varchar(64);index:number_sku_idx;comment:sku"`
	GoodsType       string    `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe        string    `gorm:"type:varchar(128);comment:商品规格"`
	Shelves         string    `gorm:"type:varchar(64);comment:货架"`
	DiscountPrice   int       `gorm:"comment:折扣价"`
	GoodsUnit       string    `gorm:"type:varchar(64);comment:商品单位"`
	SaleUnit        string    `gorm:"type:varchar(64);comment:销售单位"`
	SaleCode        string    `gorm:"comment:销售编码"`
	PayCount        int       `gorm:"comment:下单数量"`
	CloseCount      int       `gorm:"type:int;default:0;comment:关闭数量"`
	LackCount       int       `gorm:"type:int;comment:欠货数量"`
	OutCount        int       `gorm:"type:int;comment:出库数量"`
	GoodsRemark     string    `gorm:"type:varchar(255);comment:商品备注"`
	BatchId         int       `gorm:"type:int(11);index;comment:批次id"`
	DeliveryOrderNo GormList  `gorm:"type:varchar(255);comment:出库单号"`
}

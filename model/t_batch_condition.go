package model

import "time"

type BatchCondition struct {
	Base
	BatchId           int        `gorm:"type:int(11);index;comment:批次id"`
	WarehouseId       int        `gorm:"type:int(11);comment:仓库"`
	PayEndTime        *time.Time `gorm:"type:datetime;comment:支付截止时间"`
	DeliveryStartTime *time.Time `gorm:"type:datetime;default:null;comment:发货起始时间"`
	DeliveryEndTime   *time.Time `gorm:"type:datetime;comment:发货截止时间"`
	Line              string     `gorm:"type:varchar(255);comment:线路;default:null"`
	DeliveryMethod    int        `gorm:"type:tinyint;not null;comment:配送方式"`
	Sku               string     `gorm:"type:varchar(255);default:null;comment:商品sku"`
	Goods             string     `gorm:"type:varchar(255);comment:商品"`
}

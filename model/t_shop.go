package model

import "time"

// 店铺 先同步后勾选批量设置线路
type Shop struct {
	Id               int       `gorm:"primaryKey;type:int(11) unsigned;comment:id" json:"id"`
	ShopId           int       `gorm:"not null;comment:哗啦啦店铺id" json:"shop_id"`
	ShopName         string    `gorm:"type:varchar(64);not null;comment:店铺名称" json:"shop_name"`
	HouseCode        string    `gorm:"type:varchar(64);not null;comment:店铺编码" json:"house_code"`
	Warehouse        string    `gorm:"type:varchar(64);not null;comment:仓库" json:"warehouse"`
	Typ              string    `gorm:"type:varchar(64);not null;comment:类型" json:"typ"`
	Province         string    `gorm:"type:varchar(64);not null;comment:省" json:"province"`
	City             string    `gorm:"type:varchar(64);not null;comment:市" json:"city"`
	District         string    `gorm:"type:varchar(64);not null;comment:地区" json:"district"`
	Line             string    `gorm:"type:varchar(64);not null;comment:线路" json:"line"`
	ShopCode         string    `gorm:"type:varchar(255);not null;comment:店铺编号" json:"shop_code"`
	Status           int       `gorm:"not null;comment:状态" json:"status"`
	DistributionType int       `gorm:"type:tinyint;default:null;comment:配送方式"`
	CreateAt         time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间" json:"create_at"`
	UpdateAt         time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间" json:"update_at"`
}

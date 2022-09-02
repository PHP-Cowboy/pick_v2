package model

import (
	"time"
)

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

// 分类
type Classification struct {
	Base
	GoodsClass     string `gorm:"type:varchar(64);not null;comment:商品分类" json:"types"`
	WarehouseClass string `gorm:"type:varchar(64);default:'';comment:仓库分类" json:"warehouse_class"`
}

// 仓库
type Warehouse struct {
	Base
	WarehouseName string `gorm:"type:varchar(64);not null;comment:仓库名称"`
	Abbreviation  string `gorm:"type:varchar(64);not null;comment:仓库简称"`
}

// DictType 字典类型表
type DictType struct {
	Code       string    `gorm:"type:varchar(50);primaryKey;comment:字典类型编码"`
	Name       string    `gorm:"type:varchar(20);not null;comment:字典类型名称"`
	CreateTime time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
}

// Dict 字典表
type Dict struct {
	TypeCode   string    `gorm:"type:varchar(20);not null;primaryKey;comment:字典类型编码"`
	Code       string    `gorm:"type:varchar(50);not null;primaryKey;comment:字典编码"`
	Name       string    `gorm:"type:varchar(20);not null;comment:字典名称"`
	Value      string    `gorm:"type:varchar(20);not null;comment:字典值"`
	IsEdit     bool      `gorm:"type:tinyint;not null;default:0;comment:是否可编辑:0:否,1:是"`
	CreateTime time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
}

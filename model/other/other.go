package other

import (
	"pick_v2/model"
	"time"
)

//店铺 先同步后勾选批量设置线路
type Shop struct {
	model.Base
	ShopName string `gorm:"type:varchar(64);not null;comment:店铺名称"`
	ShopType string `gorm:"type:varchar(64);not null;comment:店铺类型"`
	ShopCode string `gorm:"type:varchar(255);not null;comment:店铺编号"`
	LineName string `gorm:"type:varchar(64);default:'';comment:线路名称"`
}

//分类
type Classification struct {
	model.Base
	GoodsClass     string `gorm:"type:varchar(64);not null;comment:商品分类"`
	WarehouseClass string `gorm:"type:varchar(64);default:'';comment:仓库分类"`
}

//仓库
type Warehouse struct {
	model.Base
	WarehouseName string `gorm:"type:varchar(64);not null;comment:仓库名称"`
	Abbreviation  string `gorm:"type:varchar(64);not null;comment:仓库简称"`
}

//DictType 字典类型表
type DictType struct {
	Code       string    `gorm:"type:varchar(50);primaryKey;comment:字典类型编码"`
	Name       string    `gorm:"type:varchar(20);not null;comment:字典类型名称"`
	CreateTime time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
}

//Dict 字典表
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

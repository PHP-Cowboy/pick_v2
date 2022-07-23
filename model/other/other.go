package other

import (
	"pick_v2/model"
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
	Abbreviation string `gorm:"type:varchar(64);not null;comment:仓库简称"`
}



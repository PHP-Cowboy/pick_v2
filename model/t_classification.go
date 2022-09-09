package model

// 分类
type Classification struct {
	Base
	GoodsClass     string `gorm:"type:varchar(64);not null;comment:商品分类" json:"types"`
	WarehouseClass string `gorm:"type:varchar(64);default:'';comment:仓库分类" json:"warehouse_class"`
}

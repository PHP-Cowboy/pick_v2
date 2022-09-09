package model

// 预拣货备注明细
type PrePickRemark struct {
	Base
	WarehouseId  int    `gorm:"type:int(11);comment:仓库"`
	BatchId      int    `gorm:"type:int(11) unsigned;comment:批次表id"`
	OrderGoodsId int    `gorm:"type:int(11) unsigned;comment:订单商品表ID"`
	ShopId       int    `gorm:"type:int(11);comment:店铺id"`
	PrePickId    int    `gorm:"type:int(11) unsigned;index;comment:预拣货表id"`
	Number       string `gorm:"type:varchar(64);comment:订单编号"`
	OrderRemark  string `gorm:"type:varchar(512);comment:订单备注"`
	GoodsRemark  string `gorm:"type:varchar(255);comment:商品备注"`
	ShopName     string `gorm:"type:varchar(64);not null;comment:店铺名称"`
	Line         string `gorm:"type:varchar(255);not null;comment:线路"`
	Status       int    `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:已进入拣货池"`
}

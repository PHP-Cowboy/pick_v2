package model

// 拣货备注明细
type PickRemark struct {
	Base
	WarehouseId     int    `gorm:"type:int(11);comment:仓库"`
	BatchId         int    `gorm:"type:int(11) unsigned;comment:批次表id"`
	PickId          int    `gorm:"type:int(11) unsigned;comment:拣货表id"`
	PrePickRemarkId int    `gorm:"type:int(11);comment:预拣货备注表id"`
	OrderGoodsId    int    `gorm:"type:int(11) unsigned;comment:订单商品表ID"`
	Number          string `gorm:"type:varchar(64);comment:订单编号"`
	OrderRemark     string `gorm:"type:varchar(512);comment:订单备注"`
	GoodsRemark     string `gorm:"type:varchar(255);comment:商品备注"`
	ShopName        string `gorm:"type:varchar(64);not null;comment:店铺名称"`
	Line            string `gorm:"type:varchar(255);not null;comment:线路"`
}

package model

import "gorm.io/gorm"

// 预拣货商品明细
type PrePickGoods struct {
	Base
	WarehouseId      int    `gorm:"type:int(11);comment:仓库"`
	BatchId          int    `gorm:"type:int(11) unsigned;comment:批次表id"`
	OrderGoodsId     int    `gorm:"type:int(11) unsigned;comment:订单商品表ID"` //t_pick_order_goods 表 id
	Number           string `gorm:"type:varchar(32);comment:订单编号"`
	ShopId           int    `gorm:"type:int(11);comment:店铺id"`
	PrePickId        int    `gorm:"type:int(11) unsigned;index;comment:预拣货表id"` //index
	DistributionType int    `gorm:"type:tinyint unsigned;comment:配送方式:1:公司配送,2:用户自提,3:三方物流,4:快递配送,5:首批物料|设备单"`
	Sku              string `gorm:"type:varchar(64);comment:sku"`
	GoodsName        string `gorm:"type:varchar(64);comment:商品名称"`
	GoodsType        string `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe         string `gorm:"type:varchar(128);comment:商品规格"`
	Unit             string `gorm:"type:varchar(64);comment:单位"`
	Shelves          string `gorm:"type:varchar(64);comment:货架"`
	DiscountPrice    int    `gorm:"comment:折扣价"`
	NeedNum          int    `gorm:"type:int;not null;comment:需拣数量"`
	CloseNum         int    `gorm:"type:int;not null;comment:关闭数量"`
	OutCount         int    `gorm:"type:int;comment:出库数量"`
	NeedOutNum       int    `gorm:"type:int;comment:需出库数量"`
	Status           int    `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:已进入拣货池"`
}

func PrePickGoodsBatchSave(db *gorm.DB, list []PrePickGoods) (err error) {
	result := db.Model(&PrePickGoods{}).Save(&list)

	return result.Error
}

package model

import "gorm.io/gorm"

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
	Typ          int    `gorm:"type:tinyint;default:1;comment:批次类型:1:常规批次,2:快递批次"`
}

const (
	PrePickRemarkStatusUnhandled  = iota //未处理
	PrePickRemarkStatusProcessing        //处理中(已进入拣货池)
)

func PrePickRemarkBatchSave(db *gorm.DB, list *[]PrePickRemark) (err error) {
	err = db.Model(&PrePickRemark{}).Save(list).Error

	return
}

func UpdatePrePickRemarkByPrePickIds(db *gorm.DB, prePickIds []int, mp map[string]interface{}) (err error) {
	err = db.Model(&PrePickRemark{}).
		Where("pre_pick_id in (?)", prePickIds).
		Updates(mp).
		Error

	return
}

func UpdatePrePickRemarkByIds(db *gorm.DB, ids []int, mp map[string]interface{}) (err error) {
	err = db.Model(&PrePickRemark{}).
		Where("id in (?)", ids).
		Updates(mp).
		Error

	return
}

func GetPrePickRemarkByOrderGoodsIds(db *gorm.DB, orderGoodsIds []int) (err error, prePickRemarks []PrePickRemark) {
	err = db.Model(&PrePickRemark{}).Where("order_goods_id in (?)", orderGoodsIds).Find(&prePickRemarks).Error
	return
}

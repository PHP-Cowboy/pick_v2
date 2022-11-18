package model

import "gorm.io/gorm"

// 预拣货列表
type PrePick struct {
	Base
	WarehouseId int    `gorm:"type:int(11);comment:仓库"`
	BatchId     int    `gorm:"type:int(11) unsigned;index;comment:批次表id"`
	ShopId      int    `gorm:"type:int(11);comment:店铺id"`
	ShopCode    string `gorm:"type:varchar(255);not null;comment:店铺编号"`
	ShopName    string `gorm:"type:varchar(64);not null;comment:店铺名称"`
	Line        string `gorm:"type:varchar(255);not null;comment:线路"`
	Status      int    `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:已进入拣货池,2:关闭"`
}

const (
	PrePickStatusUnhandled  = iota //未处理
	PrePickStatusProcessing        //处理中(已进入拣货池)
	PrePickStatusClose             //关闭
)

func PrePickBatchSave(db *gorm.DB, list []PrePick) (err error, res []PrePick) {
	result := db.Model(&PrePick{}).Save(&list)

	return result.Error, list
}

func UpdatePrePickStatusByIds(db *gorm.DB, ids []int, status int) error {
	result := db.Model(&PrePick{}).
		Where("id in (?)", ids).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

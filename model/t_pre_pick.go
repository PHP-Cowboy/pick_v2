package model

import (
	"errors"
	"gorm.io/gorm"
	"pick_v2/utils/ecode"
)

// 预拣货列表
type PrePick struct {
	Base
	WarehouseId int    `gorm:"type:int(11);comment:仓库"`
	TaskId      int    `gorm:"type:int(11) unsigned;comment:出库任务id"`
	BatchId     int    `gorm:"type:int(11) unsigned;index;comment:批次表id"`
	ShopId      int    `gorm:"type:int(11);comment:店铺id"`
	ShopCode    string `gorm:"type:varchar(255);not null;comment:店铺编号"`
	ShopName    string `gorm:"type:varchar(64);not null;comment:店铺名称"`
	Line        string `gorm:"type:varchar(255);not null;comment:线路"`
	Status      int    `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:已进入拣货池,2:关闭"`
	Typ         int    `gorm:"type:tinyint;default:1;comment:批次类型:1:常规批次,2:快递批次"`
}

const (
	PrePickStatusUnhandled  = iota //未处理
	PrePickStatusProcessing        //处理中(已进入拣货池)
	PrePickStatusClose             //关闭
)

func PrePickBatchSave(db *gorm.DB, list *[]PrePick) (err error) {
	err = db.Model(&PrePick{}).CreateInBatches(list, BatchSize).Error

	return
}

func UpdatePrePickStatusByIds(db *gorm.DB, ids []int, status int) (err error) {
	err = db.Model(&PrePick{}).
		Where("id in (?)", ids).
		Update("status", status).
		Error

	return
}

func UpdatePrePickByIds(db *gorm.DB, ids []int, mp map[string]interface{}) (err error) {
	err = db.Model(&PrePick{}).
		Where("id in (?)", ids).
		Updates(mp).
		Error

	return
}

// 根据id和status状态获取拣货池数据
func GetPrePickByPk(db *gorm.DB, id int) (err error, prePick PrePick) {

	err = db.Model(&PrePick{}).First(&prePick, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = ecode.DataNotExist
		}
	}

	return
}

// 根据id和status状态获取拣货池数据
func GetPrePickByIdsAndStatus(db *gorm.DB, ids []int, status int) (err error, prePick []PrePick) {
	//status 0:未处理,1:已进入拣货池
	err = db.Model(&PrePick{}).Where("id in (?) and status = ?", ids, status).Find(&prePick).Error

	return
}

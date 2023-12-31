package model

import (
	"errors"

	"gorm.io/gorm"
	"pick_v2/utils/ecode"
)

// 批次
type Batch struct {
	Base
	TaskId            int     `gorm:"type:int(11);index;comment:出库任务ID"`
	WarehouseId       int     `gorm:"type:int(11);comment:仓库"`
	BatchName         string  `gorm:"type:varchar(64);comment:批次名称"`
	DeliveryStartTime *MyTime `gorm:"type:datetime;default:null;comment:发货起始时间"`
	DeliveryEndTime   *MyTime `gorm:"type:datetime;comment:发货截止时间"`
	UserName          string  `gorm:"type:varchar(32);comment:用户名称"`
	Line              string  `gorm:"type:varchar(255);comment:所属路线"`
	DeliveryMethod    int     `gorm:"type:tinyint;not null;comment:配送方式"`
	EndTime           *MyTime `gorm:"type:datetime;comment:结束时间"`
	Status            int     `gorm:"type:tinyint;comment:状态:0:进行中,1:已结束,2:暂停"`
	Sort              int     `gorm:"type:int(11) unsigned;comment:排序"`
	PayEndTime        *MyTime `gorm:"type:datetime;comment:支付截止时间"`
	Version           int     `gorm:"type:int(11);default:0;comment:版本"`
	Typ               int     `gorm:"type:tinyint;default:1;comment:批次类型:1:常规批次,2:快递批次"`
}

const (
	BatchOngoingStatus = iota //进行中
	BatchClosedStatus         //已结束
	BatchSuspendStatus        //暂停
)

const (
	_                       = iota
	RegularBatchTyp         //常规批次
	ExpressDeliveryBatchTyp //快递批次
)

func BatchSave(db *gorm.DB, batch *Batch) (err error) {

	err = db.Model(&Batch{}).Save(batch).Error

	return
}

// 根据主键批量更新批次数据
func UpdateBatchByIds(db *gorm.DB, ids []int, mp map[string]interface{}) (err error) {
	err = db.Model(&Batch{}).Where("id in (?)", ids).Updates(mp).Error
	return
}

// 根据主键更新批次数据
func UpdateBatchByPk(db *gorm.DB, pk int, mp map[string]interface{}) (err error) {
	err = db.Model(&Batch{}).Where("id = ?", pk).Updates(mp).Error
	return
}

// 通过主键查询数据
func GetBatchByPk(db *gorm.DB, pk int) (err error, batch Batch) {
	err = db.Model(&Batch{}).First(&batch, pk).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = ecode.DataNotExist
			return
		}
	}

	return
}

// 根据条件查找第一条数据
func GetBatchFirstByStatus(db *gorm.DB, status int) (err error, batch *Batch) {
	err = db.Model(&Batch{}).Where("status = ?", status).First(&batch).Error
	return
}

// 根据出库任务获取批次列表
func GetBatchListByTaskId(db *gorm.DB, taskId int) (err error, list []Batch) {

	err = db.Model(&Batch{}).
		Where(&Batch{TaskId: taskId}).
		Find(&list).
		Error

	return
}

// 快递批次进行中或暂停的单数量
func GetBatchListByTyp(db *gorm.DB, typ int) (err error, list []Batch) {
	err = db.Model(&Batch{}).
		Where("typ = ? and ( status = ? or status = ? )", typ, BatchOngoingStatus, BatchSuspendStatus).
		Find(&list).
		Error

	return
}

func GetBatchListByIds(db *gorm.DB, ids []int) (err error, list []Batch) {
	err = db.Model(&Batch{}).
		Where("id in (?)", ids).
		Find(&list).
		Error

	return
}

func GetBatchList(db *gorm.DB, cond *Batch) (err error, list []Batch) {
	err = db.Model(&Batch{}).
		Where(cond).
		Find(&list).
		Error

	return
}

func GetBatchListByStatus(db *gorm.DB, status int) (err error, list []Batch) {
	err = db.Model(&Batch{}).
		Where("status = ?", status).
		Find(&list).
		Error

	return
}

// 获取进行中的批次
func GetBatchListByStatusAndTyp(db *gorm.DB, status, typ int) (err error, list []Batch) {
	err = db.Model(&Batch{}).
		Where("status = ? and typ = ?", status, typ).
		Find(&list).
		Error

	return
}

// 根据批次ID获取根据sort排序的批次数据
func GetBatchListByBatchIdsAndSort(db *gorm.DB, ids []int, sort string) (err error, batch Batch) {
	err = db.Model(&Batch{}).
		Where("id in (?)", ids).
		Order(sort).
		First(&batch).
		Error

	return
}

// 获取进行中的批次 或 用户已接单但未完成的批次
func GetBatchListByIdsOrPending(db *gorm.DB, ids []int, typ int) (err error, list []Batch) {
	err = db.Model(&Batch{}).
		Where("id in (?) or status = ? and typ = ?", ids, BatchOngoingStatus, typ).
		Find(&list).
		Error

	return
}

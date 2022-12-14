package model

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type OutboundOrder struct {
	TaskId            int       `gorm:"primaryKey;type:int(11);not null;comment:t_outbound_task表ID"`
	Number            string    `gorm:"primaryKey;type:varchar(64);index;comment:订单编号"`
	CreateTime        time.Time `gorm:"autoCreateTime;type:datetime;comment:创建时间"`
	UpdateTime        time.Time `gorm:"autoCreateTime;type:datetime;comment:更新时间"`
	DeleteTime        time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
	OrderId           int       `gorm:"type:int(11);comment:订单id"`
	PayAt             *MyTime   `gorm:"type:datetime;comment:支付时间"`
	ShopId            int       `gorm:"type:int(11);not null;comment:店铺id"`
	ShopName          string    `gorm:"type:varchar(64);not null;comment:店铺名称"`
	ShopType          string    `gorm:"type:varchar(64);not null;comment:店铺类型"`
	ShopCode          string    `gorm:"type:varchar(255);not null;comment:店铺编号"`
	HouseCode         string    `gorm:"type:varchar(64);not null;comment:仓库编码"`
	DistributionType  int       `gorm:"type:tinyint;comment:配送方式"`
	Line              string    `gorm:"type:varchar(255);not null;comment:线路"`
	Province          string    `gorm:"type:varchar(64);comment:省"`
	City              string    `gorm:"type:varchar(64);comment:市"`
	District          string    `gorm:"type:varchar(64);comment:区"`
	Address           string    `gorm:"type:varchar(255);comment:地址"`
	ConsigneeName     string    `gorm:"type:varchar(64);comment:收货人名称"`
	ConsigneeTel      string    `gorm:"type:varchar(64);comment:收货人电话"`
	OrderType         int       `gorm:"type:tinyint;default:1;comment:订单类型:1:新订单,2:拣货中,3:已完成,4:已关闭"`
	LatestPickingTime *MyTime   `gorm:"type:datetime;default:null;comment:最近拣货时间"`
	HasRemark         int       `gorm:"type:tinyint;default:1;comment:是否有备注:1:否,2:是"`
	OrderRemark       string    `gorm:"type:varchar(512);comment:订单备注"`
}

type OutboundOrderTypeCount struct {
	OrderType int `json:"order_type"`
	Count     int `json:"count"`
}

const (
	OutboundOrderType         = iota
	OutboundOrderTypeNew      //新订单
	OutboundOrderTypePicking  //拣货中
	OutboundOrderTypeComplete //已完成
	OutboundOrderTypeClose    //已关闭
)

func OutboundOrderBatchSave(db *gorm.DB, list []OutboundOrder) (err error) {

	err = db.Model(&OutboundOrder{}).CreateInBatches(&list, BatchSize).Error

	return
}

func OutboundOrderReplaceSave(db *gorm.DB, list []OutboundOrder, values []string) (err error) {
	err = db.Model(&OutboundOrder{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "task_id,number"}},
			DoUpdates: clause.AssignmentColumns(values),
		}).
		CreateInBatches(&list, BatchSize).Error

	return
}

func OutboundOrderBatchUpdate(db *gorm.DB, where OutboundOrder, mp map[string]interface{}) (err error) {
	err = db.Model(&OutboundOrder{}).Where(&where).Updates(mp).Error

	return
}

func UpdateOutboundOrderByTaskIdAndNumbers(db *gorm.DB, taskId int, numbers []string, mp map[string]interface{}) (err error) {
	err = db.Model(&OutboundOrder{}).Where("task_id = ? and number in (?)", taskId, numbers).Updates(mp).Error

	return err
}

func GetOutboundNumber(taskId int, number string) string {
	return fmt.Sprintf("%v%s", taskId, number)
}

func GetOutboundOrderByPk(db *gorm.DB, taskId int, number string) (err error, outboundOrder OutboundOrder) {
	err = db.Model(&OutboundOrder{}).
		Where("task_id = ? and number = ?", taskId, number).
		Find(&outboundOrder).
		Error
	return
}

func GetOutboundOrderByNumberFirstSortByTaskId(db *gorm.DB, number string) (err error, outboundOrder OutboundOrder) {
	err = db.Model(&OutboundOrder{}).
		Where("number = ?", number).
		Order("task_id desc").
		First(&outboundOrder).
		Error
	return
}

func GetOutboundOrderByTaskIdAndNumbers(db *gorm.DB, taskId int, numbers []string) (err error, list []OutboundOrder) {
	err = db.Model(&OutboundOrder{}).Where("task_id = ? and number in (?)", taskId, numbers).Find(&list).Error
	return
}

func GetOutboundOrderInMultiColumn(db *gorm.DB, multiColumn [][]interface{}) (err error, list []OutboundOrder) {
	err = db.Model(&OutboundOrder{}).Where("(task_id , number ) IN ?", multiColumn).Find(&list).Error
	return
}

func GetOutboundOrderByNumbers(db *gorm.DB, numbers []string) (err error, list []OutboundOrder) {
	err = db.Model(&OutboundOrder{}).Where("number in (?)", numbers).Find(&list).Error
	return
}

// 根据status分组统计任务条数
func OutboundOrderOrderTypeCount(db *gorm.DB, taskId int) (err error, count []OutboundOrderTypeCount) {

	err = db.Model(&OutboundOrder{}).
		Select("count(1) as count, order_type").
		Where("task_id = ? and order_type != ?", taskId, OutboundOrderTypeClose).
		Group("order_type").
		Find(&count).
		Error

	return
}

// 根据任务ID查询出库任务订单新订单类型数据
func GetOutboundOrderByTaskId(db *gorm.DB, taskId int) (err error, list []OutboundOrder) {
	err = db.Model(&OutboundOrder{}).
		Where("task_id = ? and order_type = ?", taskId, OutboundOrderTypeNew).
		Find(&list).
		Error

	return
}

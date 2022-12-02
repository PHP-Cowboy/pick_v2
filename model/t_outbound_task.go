package model

import "gorm.io/gorm"

type OutboundTask struct {
	Base
	TaskName          string  `gorm:"type:varchar(64);comment:任务名称"`
	DeliveryStartTime *MyTime `gorm:"type:datetime;default:null;comment:发货起始时间"`
	DeliveryEndTime   *MyTime `gorm:"type:datetime;comment:发货截止时间"`
	Line              string  `gorm:"type:varchar(255);comment:所属路线"`
	DistributionType  int     `gorm:"type:tinyint;comment:配送方式"`
	PayEndTime        *MyTime `gorm:"type:datetime;comment:结束时间"`
	Status            int     `gorm:"type:tinyint;default:1;comment:状态:1:进行中,2:已结束"`
	IsPush            int     `gorm:"type:tinyint;default:1;comment:推送状态:1:未推送,2:已推送"`
	Sku               string  `gorm:"type:varchar(255);comment:筛选sku"`
	GoodsName         string  `gorm:"type:varchar(255);comment:筛选商品名称"`
	Creator
}

type OutboundTaskCountStatus struct {
	Status int `json:"status"`
	Count  int `json:"count"`
}

const (
	OutboundTaskStatus        = iota
	OutboundTaskStatusOngoing //进行中
	OutboundTaskStatusClosed  //已结束
)

func OutboundTaskCreate(db *gorm.DB, task *OutboundTask) error {
	result := db.Model(&OutboundTask{}).Create(task)
	return result.Error
}

func UpdateOutboundTaskStatusById(db *gorm.DB, taskId int) error {
	result := db.Model(&OutboundTask{}).
		Where("id = ?", taskId).
		Update("status", OutboundTaskStatusClosed)

	return result.Error
}

// 根据status分组统计任务条数
func OutboundTaskCountGroupStatus(db *gorm.DB) (err error, count []OutboundTaskCountStatus) {

	result := db.Model(&OutboundTask{}).
		Select("count(1) as count, status").
		Group("status").
		Find(&count)

	if result.Error != nil {
		return result.Error, nil
	}

	return nil, count
}

func GetOutboundTaskStatusOngoingList(db *gorm.DB) (err error, list []OutboundTask) {
	result := db.Model(&OutboundTask{}).
		Where("status = ?", OutboundTaskStatusOngoing).
		Find(&list)

	return result.Error, list
}

// 根据ID查找出库任务数据
func GetOutboundTaskById(db *gorm.DB, id int) (err error, task OutboundTask) {
	result := db.Model(&OutboundTask{}).First(&task, id)

	return result.Error, task
}

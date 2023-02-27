package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// u8 推送日志
type StockLog struct {
	Base
	Number      string `gorm:"type:varchar(64);default:'';comment:订单编号" json:"number"`
	BatchId     int    `gorm:"type:int(11) unsigned;comment:批次表id" json:"batch_id"`
	PickId      int    `gorm:"type:int(11) unsigned;comment:拣货批次id" json:"pick_id"`
	Status      int    `gorm:"type:tinyint;default:0;comment:状态:0:已创建,1:推送成功,2:推送失败,3:手工补单" json:"status"`
	RequestXml  string `gorm:"type:text;comment:请求xml" json:"request_xml"`
	ResponseXml string `gorm:"type:text;comment:响应xml" json:"response_xml"`
	ResponseNo  string `gorm:"type:varchar(64);default:'';comment:U8返回单号" json:"response_no"`
	Msg         string `gorm:"type:text;comment:信息" json:"msg"`
	ShopName    string `gorm:"type:varchar(64);default:'';not null;comment:店铺名称" json:"shop_name"`
}

const (
	StockLogCreatedStatus             = iota //已创建
	StockLogPushSucceededStatus              //推送成功
	StockLogPushFailedStatus                 //推送失败
	StockLogManualReplenishmentStatus        //手工补单
)

func BatchSaveStockLog(db *gorm.DB, list *[]StockLog) (err error) {
	err = db.Model(&StockLog{}).CreateInBatches(list, BatchSize).Error
	return
}

func StockLogReplaceSave(db *gorm.DB, list *StockLog, values []string) (err error) {
	err = db.Model(&StockLog{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns(values),
		}).
		CreateInBatches(list, BatchSize).
		Error
	return
}

func GerStockLogListByIds(db *gorm.DB, ids []int, fields string) (err error, list []StockLog) {
	err = db.Model(&StockLog{}).Select(fields).Where("id in (?)", ids).Find(&list).Error
	return
}

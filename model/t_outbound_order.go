package model

import (
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
	PayAt             *MyTime   `gorm:"type:datetime;comment:支付时间"`
	ShopId            int       `gorm:"type:int(11);not null;comment:店铺id"`
	ShopName          string    `gorm:"type:varchar(64);not null;comment:店铺名称"`
	ShopType          string    `gorm:"type:varchar(64);not null;comment:店铺类型"`
	ShopCode          string    `gorm:"type:varchar(255);not null;comment:店铺编号"`
	HouseCode         string    `gorm:"type:varchar(64);not null;comment:仓库编码"`
	DistributionType  int       `gorm:"type:tinyint;comment:配送方式"`
	GoodsNum          int       `gorm:"type:int;default:0;comment:下单商品总数"`
	LimitNum          int       `gorm:"type:int;default:0;comment:限发数量"`
	CloseNum          int       `gorm:"type:int;default:0;comment:关闭数量"`
	Line              string    `gorm:"type:varchar(255);not null;comment:线路"`
	Province          string    `gorm:"type:varchar(64);comment:省"`
	City              string    `gorm:"type:varchar(64);comment:市"`
	District          string    `gorm:"type:varchar(64);comment:区"`
	Address           string    `gorm:"type:varchar(255);comment:地址"`
	ConsigneeName     string    `gorm:"type:varchar(64);comment:收货人名称"`
	ConsigneeTel      string    `gorm:"type:varchar(64);comment:收货人电话"`
	OrderType         int       `gorm:"type:tinyint;default:1;comment:订单类型:1:新订单,2:拣货中,3:已完成,4:已关闭"`
	LatestPickingTime *MyTime   `gorm:"type:datetime;default:null;comment:最近拣货时间"`
	HasRemark         int       `gorm:"type:tinyint;default:0;comment:是否有备注:0:否,1:是"`
	OrderRemark       string    `gorm:"type:varchar(512);comment:订单备注"`
}

const (
	OutboundOrderType         = iota
	OutboundOrderTypeNew      //新订单
	OutboundOrderTypePicking  //拣货中
	OutboundOrderTypeComplete //已完成
	OutboundOrderTypeClose    //已关闭
)

func OutboundOrderBatchSave(db *gorm.DB, data []OutboundOrder) error {

	result := db.Save(&data)

	return result.Error
}

func OutboundOrderBatchUpdate(db *gorm.DB, where OutboundOrder, mp map[string]interface{}) error {
	result := db.Model(&OutboundOrder{}).Where(&where).Updates(mp)

	return result.Error
}

func OutboundOrderReplaceSave(db *gorm.DB, list []OutboundOrder, values []string) error {
	result := db.Model(&OutboundOrder{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "task_id,number,sku"}},
			DoUpdates: clause.AssignmentColumns(values),
		}).
		Save(&list)

	return result.Error
}

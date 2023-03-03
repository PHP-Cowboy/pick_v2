package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// 订单表
type OrderNew struct {
	Id                int       `gorm:"primaryKey;type:int(11) unsigned;comment:id"`
	CreateTime        time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime        time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime        time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
	ShopId            int       `gorm:"type:int(11);not null;comment:店铺id"`
	ShopName          string    `gorm:"type:varchar(64);not null;comment:店铺名称"`
	ShopType          string    `gorm:"type:varchar(64);not null;comment:店铺类型"`
	ShopCode          string    `gorm:"type:varchar(255);not null;comment:店铺编号"`
	Number            string    `gorm:"type:varchar(64);unique;comment:订单编号"`
	HouseCode         string    `gorm:"type:varchar(64);not null;comment:仓库编码"`
	Line              string    `gorm:"type:varchar(255);not null;comment:线路"`
	DistributionType  int       `gorm:"type:tinyint;comment:配送方式"`
	OrderRemark       string    `gorm:"type:varchar(512);comment:订单备注"`
	PayAt             *MyTime   `gorm:"type:datetime;comment:支付时间"`
	DeliveryAt        MyTime    `gorm:"type:date;comment:配送时间"`
	Province          string    `gorm:"type:varchar(64);comment:省"`
	City              string    `gorm:"type:varchar(64);comment:市"`
	District          string    `gorm:"type:varchar(64);comment:区"`
	Address           string    `gorm:"type:varchar(255);comment:地址"`
	ConsigneeName     string    `gorm:"type:varchar(64);comment:收货人名称"`
	ConsigneeTel      string    `gorm:"type:varchar(64);comment:收货人电话"`
	OrderType         int       `gorm:"type:tinyint;default:1;comment:订单类型:1:新订单,2:拣货中,3:欠货单,4:已关闭"`
	HasRemark         int       `gorm:"type:tinyint;default:1;comment:是否备注:1:否,2:是"`
	LatestPickingTime *MyTime   `gorm:"type:datetime;default:null;comment:最近拣货时间"`
}

func OrderNewBatchSave(db *gorm.DB, list []OrderNew) (err error) {
	err = db.Model(&OrderNew{}).CreateInBatches(&list, BatchSize).Error
	return
}

func OrderNewReplaceSave(db *gorm.DB, list *[]OrderNew, values []string) (err error) {
	//[]string{"shop_id", "shop_name", "shop_type", "shop_code", "house_code", "line"}

	if len(*list) == 0 {
		return
	}

	err = db.Model(&OrderNew{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns(values),
		}).
		CreateInBatches(list, BatchSize).
		Error

	return
}

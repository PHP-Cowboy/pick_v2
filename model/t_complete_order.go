package model

import (
	"errors"
	"gorm.io/gorm"
)

// 完成订单表
type CompleteOrder struct {
	Base
	Number         string  `gorm:"type:varchar(64);unique;comment:订单编号"`
	OrderRemark    string  `gorm:"type:varchar(512);comment:订单备注"`
	ShopId         int     `gorm:"type:int(11);not null;comment:店铺id"`
	ShopName       string  `gorm:"type:varchar(64);not null;comment:店铺名称"`
	ShopType       string  `gorm:"type:varchar(64);not null;comment:店铺类型"`
	ShopCode       string  `gorm:"type:varchar(255);not null;comment:店铺编号"`
	Line           string  `gorm:"type:varchar(255);not null;comment:线路"`
	DeliveryMethod int     `gorm:"type:tinyint;not null;comment:配送方式"`
	HouseCode      string  `gorm:"type:varchar(64);not null;comment:仓库编码"`
	Province       string  `gorm:"type:varchar(64);comment:省"`
	City           string  `gorm:"type:varchar(64);comment:市"`
	District       string  `gorm:"type:varchar(64);comment:区"`
	PickTime       *MyTime `gorm:"type:datetime;default: null;comment:最近拣货时间"`
	PayAt          *MyTime `gorm:"type:datetime;comment:支付时间"`
}

func CompleteOrderSave(db *gorm.DB, list *CompleteOrder) (err error) {
	err = db.Model(&CompleteOrder{}).Save(list).Error

	return
}

func CompleteOrderBatchSave(db *gorm.DB, list *[]CompleteOrder) (err error) {
	if len(*list) == 0 {
		return
	}

	err = db.Model(&CompleteOrder{}).CreateInBatches(list, BatchSize).Error

	return
}

func GetCompleteOrderList(db *gorm.DB, cond *CompleteOrder) (err error, list []CompleteOrder) {
	err = db.Model(&CompleteOrder{}).Where(cond).Find(&list).Error
	return
}

func FindCompleteOrderExist(db *gorm.DB, number string) (err error, exist bool) {

	var completeOrder CompleteOrder

	err = db.Model(&CompleteOrder{}).Where("number = ?", number).First(&completeOrder).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false
		}
		return
	}
	//查询到了数据，即为存在
	exist = true

	return
}

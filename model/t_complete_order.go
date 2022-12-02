package model

import "gorm.io/gorm"

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
	PayCount       int     `gorm:"comment:下单数量"`
	CloseCount     int     `gorm:"type:int;comment:关闭数量"`
	OutCount       int     `gorm:"type:int;comment:出库数量"`
	Province       string  `gorm:"type:varchar(64);comment:省"`
	City           string  `gorm:"type:varchar(64);comment:市"`
	District       string  `gorm:"type:varchar(64);comment:区"`
	PickTime       *MyTime `gorm:"type:datetime;default: null;comment:最近拣货时间"`
	PayAt          MyTime  `gorm:"type:datetime;comment:支付时间"`
}

func CompleteOrderSave(db *gorm.DB, list *CompleteOrder) error {
	result := db.Model(&CompleteOrder{}).Save(list)

	return result.Error
}

func CompleteOrderBatchSave(db *gorm.DB, list *[]CompleteOrder) error {
	result := db.Model(&CompleteOrder{}).Save(list)

	return result.Error
}

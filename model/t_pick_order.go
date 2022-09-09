package model

import "time"

// 拣货单
type PickOrder struct {
	Base
	OrderId           int        `gorm:"type:int(11);index;not null;comment:订单id"`
	ShopId            int        `gorm:"type:int(11);not null;comment:店铺id"`
	ShopName          string     `gorm:"type:varchar(64);not null;comment:店铺名称"`
	ShopType          string     `gorm:"type:varchar(64);not null;comment:店铺类型"`
	ShopCode          string     `gorm:"type:varchar(255);not null;comment:店铺编号"`
	Number            string     `gorm:"type:varchar(64);unique;comment:订单编号"`
	PickNumber        string     `gorm:"type:varchar(64);unique;comment:拣货单编号"`
	HouseCode         string     `gorm:"type:varchar(64);not null;comment:仓库编码"`
	Line              string     `gorm:"type:varchar(255);not null;comment:线路"`
	DistributionType  int        `gorm:"type:tinyint;comment:配送方式"`
	OrderRemark       string     `gorm:"type:varchar(512);comment:订单备注"`
	PayAt             string     `gorm:"type:datetime;comment:支付时间"`
	ShipmentsNum      int        `gorm:"type:int;default:0;comment:发货总数"`
	LimitNum          int        `gorm:"type:int;default:0;comment:限发数量"`
	CloseNum          int        `gorm:"type:int;default:0;comment:关闭数量"`
	DeliveryAt        string     `gorm:"type:date;comment:配送时间"`
	Province          string     `gorm:"type:varchar(64);comment:省"`
	City              string     `gorm:"type:varchar(64);comment:市"`
	District          string     `gorm:"type:varchar(64);comment:区"`
	Address           string     `gorm:"type:varchar(255);comment:地址"`
	ConsigneeName     string     `gorm:"type:varchar(64);comment:收货人名称"`
	ConsigneeTel      string     `gorm:"type:varchar(64);comment:收货人电话"`
	OrderType         int        `gorm:"type:tinyint;default:1;comment:订单类型:1:新订单,2:拣货中,3:已关闭,4:已完成"`
	HasRemark         int        `gorm:"type:tinyint;default:0;comment:是否备注:0:否,1:是"`
	LatestPickingTime *time.Time `gorm:"type:datetime;default:null;comment:最近拣货时间"`
}

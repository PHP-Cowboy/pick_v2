package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type OutboundGoods struct {
	TaskId          int       `gorm:"primaryKey;type:int(11);comment:t_outbound_task表ID"`
	Number          string    `gorm:"primaryKey;type:varchar(64);comment:订单编号"`
	Sku             string    `gorm:"primaryKey;type:varchar(64);comment:sku"`
	CreateTime      time.Time `gorm:"autoCreateTime;type:datetime;comment:创建时间"`
	UpdateTime      time.Time `gorm:"autoCreateTime;type:datetime;comment:更新时间"`
	DeleteTime      time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
	OrderGoodsId    int       `gorm:"type:int(11) unsigned;index;comment:订单商品表id"`
	BatchId         int       `gorm:"type:int(11);index;comment:批次id"`
	GoodsName       string    `gorm:"type:varchar(64);comment:商品名称"`
	GoodsType       string    `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe        string    `gorm:"type:varchar(128);comment:商品规格"`
	Shelves         string    `gorm:"type:varchar(64);comment:货架"`
	DiscountPrice   int       `gorm:"comment:折扣价"`
	GoodsUnit       string    `gorm:"type:varchar(64);comment:商品单位"`
	SaleUnit        string    `gorm:"type:varchar(64);comment:销售单位"`
	SaleCode        string    `gorm:"comment:销售编码"`
	PayCount        int       `gorm:"comment:下单数量"`
	CloseCount      int       `gorm:"type:int;default:0;comment:关闭数量"`
	LackCount       int       `gorm:"type:int;comment:欠货数量"`
	OutCount        int       `gorm:"type:int;comment:出库数量"`
	LimitNum        int       `gorm:"type:int;default:0;comment:限发数量"`
	GoodsRemark     string    `gorm:"type:varchar(255);comment:商品备注"`
	Status          int       `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:拣货中,2:已出库"`
	DeliveryOrderNo GormList  `gorm:"type:varchar(255);comment:出库单号"`
}

type OutboundGoodsJoinOrder struct {
	TaskId            int      `gorm:"primaryKey;type:int(11);not null;comment:t_outbound_task表ID"`
	Number            string   `gorm:"primaryKey;type:varchar(64);index;comment:订单编号"`
	PayAt             *MyTime  `gorm:"type:datetime;comment:支付时间"`
	ShopId            int      `gorm:"type:int(11);not null;comment:店铺id"`
	ShopName          string   `gorm:"type:varchar(64);not null;comment:店铺名称"`
	ShopType          string   `gorm:"type:varchar(64);not null;comment:店铺类型"`
	ShopCode          string   `gorm:"type:varchar(255);not null;comment:店铺编号"`
	HouseCode         string   `gorm:"type:varchar(64);not null;comment:仓库编码"`
	DistributionType  int      `gorm:"type:tinyint;comment:配送方式"`
	GoodsNum          int      `gorm:"type:int;default:0;comment:下单商品总数"`
	CloseNum          int      `gorm:"type:int;default:0;comment:关闭数量"`
	Line              string   `gorm:"type:varchar(255);not null;comment:线路"`
	Province          string   `gorm:"type:varchar(64);comment:省"`
	City              string   `gorm:"type:varchar(64);comment:市"`
	District          string   `gorm:"type:varchar(64);comment:区"`
	Address           string   `gorm:"type:varchar(255);comment:地址"`
	ConsigneeName     string   `gorm:"type:varchar(64);comment:收货人名称"`
	ConsigneeTel      string   `gorm:"type:varchar(64);comment:收货人电话"`
	OrderType         int      `gorm:"type:tinyint;default:1;comment:订单类型:1:新订单,2:拣货中,3:已完成,4:已关闭"`
	LatestPickingTime *MyTime  `gorm:"type:datetime;default:null;comment:最近拣货时间"`
	HasRemark         int      `gorm:"type:tinyint;default:0;comment:是否有备注:0:否,1:是"`
	OrderRemark       string   `gorm:"type:varchar(512);comment:订单备注"`
	Sku               string   `gorm:"primaryKey;type:varchar(64);comment:sku"` //t_outbound_goods
	OrderGoodsId      int      `gorm:"type:int(11) unsigned;index;comment:订单商品表id"`
	BatchId           int      `gorm:"type:int(11);index;comment:批次id"`
	GoodsName         string   `gorm:"type:varchar(64);comment:商品名称"`
	GoodsType         string   `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe          string   `gorm:"type:varchar(128);comment:商品规格"`
	Shelves           string   `gorm:"type:varchar(64);comment:货架"`
	DiscountPrice     int      `gorm:"comment:折扣价"`
	GoodsUnit         string   `gorm:"type:varchar(64);comment:商品单位"`
	SaleUnit          string   `gorm:"type:varchar(64);comment:销售单位"`
	SaleCode          string   `gorm:"comment:销售编码"`
	PayCount          int      `gorm:"comment:下单数量"`
	CloseCount        int      `gorm:"type:int;default:0;comment:关闭数量"`
	LackCount         int      `gorm:"type:int;comment:欠货数量"`
	OutCount          int      `gorm:"type:int;comment:出库数量"`
	LimitNum          int      `gorm:"type:int;default:0;comment:限发数量"`
	GoodsRemark       string   `gorm:"type:varchar(255);comment:商品备注"`
	Status            int      `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:拣货中,2:已出库"`
	DeliveryOrderNo   GormList `gorm:"type:varchar(255);comment:出库单号"`
}

const (
	OutboundGoodsStatusUnhandled        = iota //未处理
	OutboundGoodsStatusPicking                 //拣货中
	OutboundGoodsStatusOutboundDelivery        //已出库
)

func OutboundGoodsBatchSave(db *gorm.DB, list []OutboundGoods) error {

	result := db.Model(&OutboundGoods{}).Save(list)

	return result.Error
}

func GetOutboundGoodsJoinOrderList(db *gorm.DB, taskId int, number string) (err error, list []OutboundGoodsJoinOrder) {

	result := db.Table("t_outbound_goods og").
		Select("*").
		Joins("left join t_outbound_order oo on og.task_id = oo.task_id and og.number = oo.number").
		Where("oo.task_id = ? and oo.number = ?", taskId, number).
		Find(&list)

	if result.Error != nil {
		return result.Error, nil
	}

	return nil, list
}

func ReplaceSave(db *gorm.DB, list []OutboundGoods, values []string) error {
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "task_id,number,sku"}},
		DoUpdates: clause.AssignmentColumns(values),
	}).Save(&list)

	return result.Error
}

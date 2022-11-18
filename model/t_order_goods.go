package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// 订单商品表
type OrderGoods struct {
	Id              int       `gorm:"primaryKey;type:int(11) unsigned;comment:id"`
	CreateTime      time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime      time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime      time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
	Number          string    `gorm:"type:varchar(64);index:number_sku_idx;comment:订单编号"`
	GoodsName       string    `gorm:"type:varchar(64);comment:商品名称"`
	Sku             string    `gorm:"type:varchar(64);index:number_sku_idx;comment:sku"`
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
	GoodsRemark     string    `gorm:"type:varchar(255);comment:商品备注"`
	BatchId         int       `gorm:"type:int(11);index;comment:批次id"`
	Status          int       `gorm:"type:tinyint;default:0;comment:状态:1:未处理,2:处理中"`
	DeliveryOrderNo GormList  `gorm:"type:varchar(255);comment:出库单号"`
}

type OrderJoinGoods struct {
	ShopId            int      `json:"shop_id"` //order表
	ShopName          string   `json:"shop_name"`
	ShopType          string   `json:"shop_type"`
	ShopCode          string   `json:"shop_code"`
	Number            string   `json:"number"`
	HouseCode         string   `json:"house_code"`
	Line              string   `json:"line"`
	DistributionType  int      `json:"distribution_type"`
	OrderRemark       string   `json:"order_remark"`
	PayAt             MyTime   `json:"pay_at"`
	PayTotal          int      `json:"pay_total"`
	Picked            int      `json:"picked"`
	UnPicked          int      `json:"un_picked"`
	CloseNum          int      `json:"close_num"`
	DeliveryAt        MyTime   `json:"delivery_at"`
	Province          string   `json:"province"`
	City              string   `json:"city"`
	District          string   `json:"district"`
	Address           string   `json:"address"`
	ConsigneeName     string   `json:"consignee_name"`
	ConsigneeTel      string   `json:"consignee_tel"`
	OrderType         int      `json:"order_type"`
	HasRemark         int      `json:"has_remark"`
	LatestPickingTime *MyTime  `json:"latest_picking_time"`
	Id                int      `json:"id"` //order_goods 表 这里的查询需要注意一下，别查到了order表id
	GoodsName         string   `gorm:"type:varchar(64);comment:商品名称"`
	Sku               string   `gorm:"type:varchar(64);index:number_sku_idx;comment:sku"`
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
	GoodsRemark       string   `gorm:"type:varchar(255);comment:商品备注"`
	BatchId           int      `gorm:"type:int(11);index;comment:批次id"`
	Status            int      `gorm:"type:tinyint;default:0;comment:状态:1:未处理,2:处理中"` //订单中有欠货状态
	DeliveryOrderNo   GormList `gorm:"type:varchar(255);comment:出库单号"`
}

// 状态
const (
	OrderGoodsDefaultStatus    = iota
	OrderGoodsUnhandledStatus  //未处理
	OrderGoodsProcessingStatus //处理中
)

// 批量更新订单商品数据
func UpdateOrderGoodsByIds(db *gorm.DB, ids []int, mp map[string]interface{}) error {
	result := db.Model(&OrderGoods{}).Where("id in (?)", ids).Updates(mp)

	return result.Error
}

// 订单&&商品信息
func GetOrderJoinGoodsList(db *gorm.DB, number []string) (err error, list []OrderJoinGoods) {
	result := db.Table("t_order_goods og").
		Omit("o.id").
		Select("*").
		Joins("left join t_order o on og.number = o.number").
		Where("og.number in (?)", number).
		Find(&list)

	if result.Error != nil {
		return result.Error, nil
	}

	return nil, list
}

func UpdateOrderGoodsStatus(db *gorm.DB, list []OrderGoods, values []string) error {
	result := db.Model(&OrderGoods{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(values),
	}).Save(&list)

	return result.Error
}

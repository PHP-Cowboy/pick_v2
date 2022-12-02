package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// 订单表
type Order struct {
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
	PayAt             MyTime    `gorm:"type:datetime;comment:支付时间"`
	DeliveryAt        MyTime    `gorm:"type:date;comment:配送时间"`
	Province          string    `gorm:"type:varchar(64);comment:省"`
	City              string    `gorm:"type:varchar(64);comment:市"`
	District          string    `gorm:"type:varchar(64);comment:区"`
	Address           string    `gorm:"type:varchar(255);comment:地址"`
	ConsigneeName     string    `gorm:"type:varchar(64);comment:收货人名称"`
	ConsigneeTel      string    `gorm:"type:varchar(64);comment:收货人电话"`
	OrderType         int       `gorm:"type:tinyint;default:1;comment:订单类型:1:新订单,2:拣货中,3:欠货单,4:已关闭"`
	HasRemark         int       `gorm:"type:tinyint;default:0;comment:是否备注:0:否,1:是"`
	LatestPickingTime *MyTime   `gorm:"type:datetime;default:null;comment:最近拣货时间"`

	//PayTotal          int       `gorm:"type:int;default:0;comment:支付商品总数"`
	//Picked            int       `gorm:"type:int;default:0;comment:已拣数量"`
	//UnPicked          int       `gorm:"type:int;default:0;comment:未拣数量"`
	//CloseNum          int       `gorm:"type:int;default:0;comment:关闭数量"`
}

const (
	_                = iota
	NewOrderType     //新订单
	PickingOrderType //拣货中
	LackOrderType    //欠货单
	CloseOrderType   //已关闭
)

const (
	DistributionType           = iota //配送方式
	DistributionTypeCompany           // 1:公司配送
	DistributionTypeUser              // 2:用户自提
	DistributionTypeThreeParty        // 3:三方物流
	DistributionTypeCourier           // 4:快递配送
	DistributionTypeFirst             // 5:首批物料|设备单
)

func OrderSave(db *gorm.DB, order *Order) error {
	result := db.Model(&Order{}).Save(order)
	return result.Error
}

func OrderBatchSave(db *gorm.DB, list []Order) error {
	result := db.Model(&Order{}).Save(&list)
	return result.Error
}

func OrderReplaceSave(db *gorm.DB, list []Order, values []string) error {
	//[]string{"shop_id", "shop_name", "shop_type", "shop_code", "house_code", "line"}

	result := db.Model(&Order{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns(values),
		}).Save(&list)

	return result.Error
}

func UpdateOrderByIds(db *gorm.DB, ids []int, mp map[string]interface{}) error {
	result := db.Model(&Order{}).Where("id in (?)", ids).Updates(mp)

	return result.Error
}

func UpdateOrderByNumbers(db *gorm.DB, numbers []string, mp map[string]interface{}) error {
	result := db.Model(&Order{}).Where("number in (?)", numbers).Updates(mp)

	return result.Error
}

func DeleteOrderByNumbers(db *gorm.DB, numbers []string) error {
	result := db.Delete(&Order{}, "number in (?)", numbers)

	return result.Error
}

func OrderOrCompleteOrderExist(db *gorm.DB, ids []int, number string) (err error, exist bool) {
	var (
		order         []Order
		completeOrder []CompleteOrder
	)

	// 查询是否已存在 存在的过滤掉
	result := db.Where("id in (?)", ids).Find(&order)

	if result.Error != nil {
		return result.Error, false
	}

	if len(order) > 0 {
		return nil, true
	}

	//查看完成订单里有没有
	result = db.Where("number = ?", number).Find(&completeOrder)

	if result.Error != nil {
		return result.Error, false
	}

	if len(completeOrder) > 0 {
		return nil, true
	}

	return nil, false
}

func GetOrderListByNumbers(db *gorm.DB, numbers []string) (err error, list []Order) {
	result := db.Model(&Order{}).Where("number in (?)", numbers).Find(&list)

	return result.Error, list
}

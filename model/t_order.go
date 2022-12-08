package model

import (
	"errors"
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
	PayAt             *MyTime   `gorm:"type:datetime;comment:支付时间"`
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

func OrderCreate(db *gorm.DB, order *Order) (err error) {
	err = db.Model(&Order{}).CreateInBatches(order, BatchSize).Error
	return
}

func OrderSave(db *gorm.DB, order *Order) (err error) {
	err = db.Model(&Order{}).Save(order).Error
	return
}

func OrderBatchSave(db *gorm.DB, list []Order) (err error) {
	err = db.Model(&Order{}).CreateInBatches(&list, BatchSize).Error
	return
}

func OrderReplaceSave(db *gorm.DB, list *[]Order, values []string) (err error) {
	//[]string{"shop_id", "shop_name", "shop_type", "shop_code", "house_code", "line"}

	if len(*list) == 0 {
		return
	}

	err = db.Model(&Order{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns(values),
		}).
		CreateInBatches(list, BatchSize).
		Error

	return
}

func UpdateOrderByIds(db *gorm.DB, ids []int, mp map[string]interface{}) (err error) {
	err = db.Model(&Order{}).Where("id in (?)", ids).Updates(mp).Error

	return
}

func UpdateOrderByNumbers(db *gorm.DB, numbers []string, mp map[string]interface{}) (err error) {
	err = db.Model(&Order{}).Where("number in (?)", numbers).Updates(mp).Error

	return
}

func DeleteOrderByNumbers(db *gorm.DB, numbers []string) (err error) {
	err = db.Delete(&Order{}, "number in (?)", numbers).Error

	return
}

func DeleteOrderByIds(db *gorm.DB, ids []int) (err error) {
	if len(ids) == 0 {
		return
	}

	err = db.Delete(&Order{}, "id in (?)", ids).Error

	return
}

// 查询订单表中存在的ids的条数
func FindOrderExistByIds(db *gorm.DB, ids []int) (err error, exist bool) {
	var order Order
	err = db.Model(&Order{}).Where("id in (?)", ids).First(&order).Error

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

// 根据订单编号查询订单数据
func GetOrderListByNumbers(db *gorm.DB, numbers []string) (err error, list []Order) {
	err = db.Model(&Order{}).Where("number in (?)", numbers).Find(&list).Error

	return
}

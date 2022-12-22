package model

import "gorm.io/gorm"

// 关闭订单表
type CloseOrder struct {
	Base
	Number           string  `gorm:"type:varchar(64);index;comment:订单编号"`
	ShopName         string  `gorm:"type:varchar(64);not null;comment:店铺名称"`
	PayAt            *MyTime `gorm:"type:datetime;comment:支付时间"`
	PayTotal         int     `gorm:"type:int;default:0;comment:下单总数"`
	NeedCloseTotal   int     `gorm:"type:int;default:0;comment:需关闭总数"`
	ShopType         string  `gorm:"type:varchar(64);not null;comment:店铺类型"`
	DistributionType int     `gorm:"type:tinyint;comment:配送方式"`
	Province         string  `gorm:"type:varchar(64);comment:省"`
	City             string  `gorm:"type:varchar(64);comment:市"`
	District         string  `gorm:"type:varchar(64);comment:区"`
	OrderRemark      string  `gorm:"type:varchar(512);comment:订单备注"`
	Typ              int     `gorm:"type:tinyint;default:1;comment:状态:1:部分关闭,2:全单关闭"`
	Status           int     `gorm:"type:tinyint;default:1;comment:状态:1:处理中,2:已完成,3:异常"`
	Applicant        string  `gorm:"type:varchar(64);comment:申请人"`
	ApplyTime        *MyTime `gorm:"type:datetime;comment:申请时间"`
	CloseUser        string  `gorm:"type:varchar(32);comment:关单处理人昵称"`
}

const (
	CloseOrderStatus          = iota
	CloseOrderStatusPending   //处理中
	CloseOrderStatusComplete  //已完成
	CloseOrderStatusException //异常
)

const (
	_                 = iota
	CloseOrderTypPart //部分关闭
	CloseOrderTypAll  //全单关闭
)

func SaveCloseOrder(db *gorm.DB, data *CloseOrder) (err error) {
	err = db.Model(&CloseOrder{}).Save(data).Error
	return
}

func UpdateCloseOrderByPk(db *gorm.DB, id int, mp map[string]interface{}) (err error) {
	err = db.Model(&CloseOrder{}).Where("id = ?", id).Updates(mp).Error
	return
}

func UpdateCloseOrderByNumbers(db *gorm.DB, numbers []string, mp map[string]interface{}) (err error) {
	err = db.Model(&CloseOrder{}).Where("number in (?)", numbers).Updates(mp).Error
	return
}

func GetCloseOrderByPk(db *gorm.DB, id int) (err error, closeOrder CloseOrder) {
	err = db.Model(&CloseOrder{}).First(&closeOrder, id).Error
	return
}

func GetCloseOrderList(db *gorm.DB, cond CloseOrder) (err error, list []CloseOrder) {
	err = db.Model(&CloseOrder{}).Where(&cond).Find(&list).Error
	return
}

func GetCloseOrderPageList(db *gorm.DB, cond CloseOrder, page, size int) (err error, list []CloseOrder) {
	err = db.Model(&CloseOrder{}).Where(&cond).Scopes(Paginate(page, size)).Find(&list).Error
	return
}

func GetCloseOrderByNumbers(db *gorm.DB, numbers []string) (err error, list []CloseOrder) {
	err = db.Model(&CloseOrder{}).Where("number in (?)", numbers).Find(&list).Error
	return
}

func CountCloseOrderByCond(db *gorm.DB, cond CloseOrder) (err error, count int64) {
	err = db.Model(&CloseOrder{}).Where(&cond).Count(&count).Error
	return
}

type CountCloseOrder struct {
	Status int `json:"status"`
	Count  int `json:"count"`
}

func CountCloseOrderStatus(db *gorm.DB) (err error, countCloseOrder []CountCloseOrder) {
	err = db.Model(&CloseOrder{}).Select("status,count(1) as count").Group("status").Find(&countCloseOrder).Error
	return
}

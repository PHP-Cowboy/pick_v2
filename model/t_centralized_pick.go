package model

import "gorm.io/gorm"

// 集中拣货
type CentralizedPick struct {
	Base
	BatchId        int     `gorm:"type:int(11) unsigned;index;comment:批次表id"`
	Sku            string  `gorm:"type:varchar(64);comment:sku"`
	GoodsName      string  `gorm:"type:varchar(64);comment:商品名称"`
	GoodsType      string  `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe       string  `gorm:"type:varchar(128);comment:商品规格"`
	NeedNum        int     `gorm:"type:int;default:0;comment:需拣数量"`
	PickNum        int     `gorm:"type:int;default:0;comment:拣货数量"`
	PickUser       string  `gorm:"type:varchar(32);default:'';comment:拣货人"`
	TakeOrdersTime *MyTime `gorm:"type:datetime;default:null;comment:接单时间"`
	HasRemark      int     `gorm:"type:tinyint;default:0;comment:是否备注:0:否,1:是"`
}

func CentralizedPickSave(db *gorm.DB, list *[]CentralizedPick) error {
	result := db.Model(&CentralizedPick{}).Save(list)

	return result.Error
}

func GetCentralizedPickList(db *gorm.DB) (err error, total int64, list []CentralizedPick) {
	return
}

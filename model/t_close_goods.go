package model

import (
	"gorm.io/gorm"
)

// 关闭订单商品表
type CloseGoods struct {
	Base
	CloseOrderId   int    `gorm:"type:int;index;default:0;comment:关闭订单表ID"`
	OrderGoodsId   int    `gorm:"type:int(11) unsigned;comment:订单商品ID"`
	GoodsName      string `gorm:"type:varchar(64);comment:商品名称"`
	Sku            string `gorm:"type:varchar(64);index:number_sku_idx;comment:sku"`
	GoodsSpe       string `gorm:"type:varchar(128);comment:商品规格"`
	PayCount       int    `gorm:"comment:下单数量"`
	CloseCount     int    `gorm:"type:int;default:0;comment:已关闭数量"`
	NeedCloseCount int    `gorm:"type:int;default:0;comment:需关闭数量"`
	GoodsRemark    string `gorm:"type:varchar(255);comment:商品备注"`
}

func BatchSaveCloseGoods(db *gorm.DB, list *[]CloseGoods) (err error) {
	err = db.Model(&CloseGoods{}).CreateInBatches(list, BatchSize).Error
	return
}

func GetCloseGoodsListByCond(db *gorm.DB, cond CloseGoods) (err error, list []CloseGoods) {
	err = db.Model(&CloseGoods{}).Where(&cond).Find(&list).Error
	return
}

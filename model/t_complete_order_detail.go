package model

import "gorm.io/gorm"

// 完成订单明细表
type CompleteOrderDetail struct {
	Base
	Number          string   `gorm:"type:varchar(64);index;comment:订单编号"`
	GoodsName       string   `gorm:"type:varchar(64);comment:商品名称"`
	Sku             string   `gorm:"type:varchar(64);comment:sku"`
	GoodsSpe        string   `gorm:"type:varchar(128);comment:商品规格"`
	GoodsType       string   `gorm:"type:varchar(64);comment:商品类型"`
	Shelves         string   `gorm:"type:varchar(64);comment:货架"`
	PayCount        int      `gorm:"comment:下单数量"`
	CloseCount      int      `gorm:"type:int;comment:关闭数量"`
	ReviewCount     int      `gorm:"type:int;comment:出库数量"`
	GoodsRemark     string   `gorm:"type:varchar(255);comment:商品备注"`
	DeliveryOrderNo GormList `gorm:"type:varchar(255);comment:出库单号"`
}

func CompleteOrderDetailBatchSave(db *gorm.DB, list *[]CompleteOrderDetail) error {
	result := db.Model(&CompleteOrderDetail{}).Save(list)

	return result.Error
}

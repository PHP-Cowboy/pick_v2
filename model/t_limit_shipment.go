package model

import "gorm.io/gorm"

type LimitShipment struct {
	TaskId    int    `gorm:"primaryKey;type:int(11);not null;comment:t_outbound_task表ID" json:"task_id"`
	Number    string `gorm:"primaryKey;type:varchar(64);comment:订单编号" json:"number"`
	Sku       string `gorm:"primaryKey;type:varchar(64);comment:sku" json:"sku"`
	ShopName  string `gorm:"type:varchar(64);comment:门店名称" json:"shop_name"`
	GoodsName string `gorm:"type:varchar(64);comment:商品名称" json:"goods_name"`
	GoodsSpe  string `gorm:"type:varchar(128);comment:商品规格" json:"goods_spe"`
	LimitNum  int    `gorm:"type:int;default:0;comment:限发数量" json:"limit_num"`
	Status    int    `gorm:"type:tinyint;default:1;comment:状态:0:撤销,1:正常" json:"status"`
}

func LimitShipmentSave(db *gorm.DB, list []LimitShipment) error {
	result := db.Model(&LimitShipment{}).Save(&list)
	return result.Error
}

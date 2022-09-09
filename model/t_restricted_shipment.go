package model

// 限制发货
type RestrictedShipment struct {
	PickOrderGoodsId int    `gorm:"primaryKey;type:int(11) unsigned;comment:拣货单商品表id" json:"pick_order_goods_id"`
	PickNumber       string `gorm:"type:varchar(64);index;comment:拣货编号" json:"pick_number"`
	ShopName         string `gorm:"type:varchar(64);comment:门店名称" json:"shop_name"`
	GoodsName        string `gorm:"type:varchar(64);comment:商品名称" json:"goods_name"`
	GoodsSpe         string `gorm:"type:varchar(128);comment:商品规格" json:"goods_spe"`
	LimitNum         int    `gorm:"type:int;default:0;comment:限发数量" json:"limit_num"`
	Status           int    `gorm:"type:tinyint;default:1;comment:状态:0:撤销,1:正常" json:"status"`
}

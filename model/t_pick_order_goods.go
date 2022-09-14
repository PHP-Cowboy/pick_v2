package model

// 拣货单商品
type PickOrderGoods struct {
	Base
	PickOrderId     int      `gorm:"type:int(11) unsigned;index;comment:拣货单表id"`
	OrderGoodsId    int      `gorm:"type:int(11) unsigned;index;comment:订单商品表id"`
	Number          string   `gorm:"type:varchar(64);index;comment:订单编号"`
	GoodsName       string   `gorm:"type:varchar(64);comment:商品名称"`
	Sku             string   `gorm:"type:varchar(64);index;comment:sku"`
	GoodsType       string   `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe        string   `gorm:"type:varchar(128);comment:商品规格"`
	Shelves         string   `gorm:"type:varchar(64);comment:货架"`
	DiscountPrice   int      `gorm:"comment:折扣价"`
	GoodsUnit       string   `gorm:"type:varchar(64);comment:商品单位"`
	SaleUnit        string   `gorm:"type:varchar(64);comment:销售单位"`
	SaleCode        string   `gorm:"comment:销售编码"`
	PayCount        int      `gorm:"comment:下单数量"`
	CloseCount      int      `gorm:"type:int;default:0;comment:关闭数量"`
	LackCount       int      `gorm:"type:int;comment:欠货数量"`
	OutCount        int      `gorm:"type:int;comment:出库数量"`
	LimitNum        int      `gorm:"type:int;default:0;comment:限发数量"`
	GoodsRemark     string   `gorm:"type:varchar(255);comment:商品备注"`
	Status          int      `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:拣货中,2:已出库"`
	BatchId         int      `gorm:"type:int(11);index;comment:批次id"`
	DeliveryOrderNo GormList `gorm:"type:varchar(255);comment:出库单号"`
}

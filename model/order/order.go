package order

import (
	"pick_v2/model"
	"time"
)

type Order struct {
	model.Base
	Number         string    `gorm:"type:varchar(32);comment:订单编号"`
	PayTime        time.Time `gorm:"type:datetime;not null;comment:支付时间"`
	ShopNo         string    `gorm:"type:varchar(32);not null;comment:店铺编号"`
	ShopName       string    `gorm:"type:varchar(32);not null;comment:店铺名称"`
	ShopType       string    `gorm:"type:varchar(32);not null;comment:店铺类型"`
	DeliveryMethod int       `gorm:"type:tinyint;not null;comment:配送方式"`
	GoodsNum       int       `gorm:"type:int;not null;comment:商品数量"`
	GoodsUnit      int       `gorm:"type:int;not null;comment:商品单位"`
	Line           string    `gorm:"type:varchar(32);not null;comment:线路"`
	Region         string    `gorm:"type:varchar(255);not null;comment:地区"`
	Remark         string    `gorm:"type:varchar(255);not null;comment:地区"`
}

type OrderInfo struct {
	model.Base
	BatchId          int    `gorm:"type:int(11);not null;index;comment:批次id"`
	ShopId           int    `gorm:"type:int(11);not null;comment:店铺id"`
	ShopName         string `gorm:"type:varchar(64);not null;comment:店铺名称"`
	ShopType         string `gorm:"type:varchar(64);not null;comment:店铺类型"`
	ShopCode         string `gorm:"type:varchar(255);not null;comment:店铺编号"`
	HouseCode        string `gorm:"type:varchar(64);not null;comment:仓库编码"`
	Line             string `gorm:"type:varchar(255);not null;comment:线路"`
	Number           string `gorm:"type:varchar(64);comment:订单编号"`
	Status           int    `gorm:"type:tinyint;comment:订单状态"`
	DeliveryAt       string `gorm:"type:date;comment:配送时间"`
	DistributionType int    `gorm:"type:tinyint;comment:配送方式"`
	OrderRemark      string `gorm:"type:varchar(512);comment:订单备注"`
	Province         string `gorm:"type:varchar(64);comment:省"`
	City             string `gorm:"type:varchar(64);comment:市"`
	District         string `gorm:"type:varchar(64);comment:区"`
	Address          string `gorm:"type:varchar(255);comment:地址"`
	ConsigneeName    string `gorm:"type:varchar(64);comment:收货人名称"`
	ConsigneeTel     string `gorm:"type:varchar(64);comment:收货人电话"`
	Name             string `gorm:"type:varchar(64);comment:商品名称"`
	Sku              string `gorm:"type:varchar(64);comment:sku"`
	GoodsSpe         string `gorm:"type:varchar(128);comment:商品规格"`
	GoodsType        string `gorm:"type:varchar(64);comment:商品类型"`
	Shelves          string `gorm:"type:varchar(64);comment:货架"`
	OriginalPrice    int    `gorm:"comment:原价"`
	DiscountPrice    int    `gorm:"comment:折扣价"`
	GoodsUnit        string `gorm:"type:varchar(64);comment:商品单位"`
	SaleUnit         string `gorm:"type:varchar(64);comment:销售单位"`
	SaleCode         string `gorm:"comment:销售编码"`
	PayCount         int    `gorm:"comment:下单数量"`
	CloseCount       int    `gorm:"type:int;comment:关闭数量"`
	OutCount         int    `gorm:"type:int;comment:出库数量"`
	GoodsRemark      string `gorm:"type:varchar(255);comment:商品备注"`
	PickStatus       int    `gorm:"type:tinyint;comment:拣货状态"`
	PayAt            string `gorm:"type:datetime;comment:支付时间"`
	LackCount        int    `gorm:"type:int;comment:欠货数量"`
}

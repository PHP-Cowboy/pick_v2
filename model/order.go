package model

import (
	"time"
)

// 订单表
type Order struct {
	Id                int        `gorm:"primaryKey;type:int(11) unsigned;comment:id"`
	CreateTime        time.Time  `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime        time.Time  `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime        time.Time  `gorm:"type:datetime;default:null;comment:删除时间"`
	ShopId            int        `gorm:"type:int(11);not null;comment:店铺id"`
	ShopName          string     `gorm:"type:varchar(64);not null;comment:店铺名称"`
	ShopType          string     `gorm:"type:varchar(64);not null;comment:店铺类型"`
	ShopCode          string     `gorm:"type:varchar(255);not null;comment:店铺编号"`
	Number            string     `gorm:"type:varchar(64);unique;comment:订单编号"`
	HouseCode         string     `gorm:"type:varchar(64);not null;comment:仓库编码"`
	Line              string     `gorm:"type:varchar(255);not null;comment:线路"`
	DistributionType  int        `gorm:"type:tinyint;comment:配送方式"`
	OrderRemark       string     `gorm:"type:varchar(512);comment:订单备注"`
	PayAt             string     `gorm:"type:datetime;comment:支付时间"`
	PayTotal          int        `gorm:"type:int;default:0;comment:支付商品总数"`
	Picked            int        `gorm:"type:int;default:0;comment:已拣数量"`
	UnPicked          int        `gorm:"type:int;default:0;comment:未拣数量"`
	CloseNum          int        `gorm:"type:int;default:0;comment:关闭数量"`
	DeliveryAt        string     `gorm:"type:date;comment:配送时间"`
	Province          string     `gorm:"type:varchar(64);comment:省"`
	City              string     `gorm:"type:varchar(64);comment:市"`
	District          string     `gorm:"type:varchar(64);comment:区"`
	Address           string     `gorm:"type:varchar(255);comment:地址"`
	ConsigneeName     string     `gorm:"type:varchar(64);comment:收货人名称"`
	ConsigneeTel      string     `gorm:"type:varchar(64);comment:收货人电话"`
	OrderType         int        `gorm:"type:tinyint;default:1;comment:订单类型:1:新订单,2:拣货中,3:欠货单,4:已关闭"`
	HasRemark         int        `gorm:"type:tinyint;default:0;comment:是否备注:0:否,1:是"`
	LatestPickingTime *time.Time `gorm:"type:datetime;default:null;comment:最近拣货时间"`
}

// 订单商品表
type OrderGoods struct {
	Id              int       `gorm:"primaryKey;type:int(11) unsigned;comment:id"`
	CreateTime      time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime      time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime      time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
	Number          string    `gorm:"type:varchar(64);index:number_sku_idx;comment:订单编号"`
	GoodsName       string    `gorm:"type:varchar(64);comment:商品名称"`
	Sku             string    `gorm:"type:varchar(64);index:number_sku_idx;comment:sku"`
	GoodsType       string    `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe        string    `gorm:"type:varchar(128);comment:商品规格"`
	Shelves         string    `gorm:"type:varchar(64);comment:货架"`
	DiscountPrice   int       `gorm:"comment:折扣价"`
	GoodsUnit       string    `gorm:"type:varchar(64);comment:商品单位"`
	SaleUnit        string    `gorm:"type:varchar(64);comment:销售单位"`
	SaleCode        string    `gorm:"comment:销售编码"`
	PayCount        int       `gorm:"comment:下单数量"`
	CloseCount      int       `gorm:"type:int;default:0;comment:关闭数量"`
	LackCount       int       `gorm:"type:int;comment:欠货数量"`
	OutCount        int       `gorm:"type:int;comment:出库数量"`
	GoodsRemark     string    `gorm:"type:varchar(255);comment:商品备注"`
	BatchId         int       `gorm:"type:int(11);index;comment:批次id"`
	DeliveryOrderNo []*string `gorm:"type:varchar(255);comment:出库单号"`
}

// 拣货单
type PickOrder struct {
	Base
	OrderId           int        `gorm:"type:int(11);index;not null;comment:订单id"`
	ShopId            int        `gorm:"type:int(11);not null;comment:店铺id"`
	ShopName          string     `gorm:"type:varchar(64);not null;comment:店铺名称"`
	ShopType          string     `gorm:"type:varchar(64);not null;comment:店铺类型"`
	ShopCode          string     `gorm:"type:varchar(255);not null;comment:店铺编号"`
	Number            string     `gorm:"type:varchar(64);unique;comment:订单编号"`
	PickNumber        string     `gorm:"type:varchar(64);unique;comment:拣货单编号"`
	HouseCode         string     `gorm:"type:varchar(64);not null;comment:仓库编码"`
	Line              string     `gorm:"type:varchar(255);not null;comment:线路"`
	DistributionType  int        `gorm:"type:tinyint;comment:配送方式"`
	OrderRemark       string     `gorm:"type:varchar(512);comment:订单备注"`
	PayAt             string     `gorm:"type:datetime;comment:支付时间"`
	ShipmentsNum      int        `gorm:"type:int;default:0;comment:发货总数"`
	LimitNum          int        `gorm:"type:int;default:0;comment:限发数量"`
	CloseNum          int        `gorm:"type:int;default:0;comment:关闭数量"`
	DeliveryAt        string     `gorm:"type:date;comment:配送时间"`
	Province          string     `gorm:"type:varchar(64);comment:省"`
	City              string     `gorm:"type:varchar(64);comment:市"`
	District          string     `gorm:"type:varchar(64);comment:区"`
	Address           string     `gorm:"type:varchar(255);comment:地址"`
	ConsigneeName     string     `gorm:"type:varchar(64);comment:收货人名称"`
	ConsigneeTel      string     `gorm:"type:varchar(64);comment:收货人电话"`
	OrderType         int        `gorm:"type:tinyint;default:1;comment:订单类型:1:新订单,2:拣货中,3:已关闭,4:已完成"`
	HasRemark         int        `gorm:"type:tinyint;default:0;comment:是否备注:0:否,1:是"`
	LatestPickingTime *time.Time `gorm:"type:datetime;default:null;comment:最近拣货时间"`
}

// 拣货单商品
type PickOrderGoods struct {
	Base
	OrderGoodsId    int       `gorm:"type:int(11) unsigned;index;comment:订单商品表id"`
	Number          string    `gorm:"type:varchar(64);index;comment:订单编号"`
	GoodsName       string    `gorm:"type:varchar(64);comment:商品名称"`
	Sku             string    `gorm:"type:varchar(64);index;comment:sku"`
	GoodsType       string    `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe        string    `gorm:"type:varchar(128);comment:商品规格"`
	Shelves         string    `gorm:"type:varchar(64);comment:货架"`
	DiscountPrice   int       `gorm:"comment:折扣价"`
	GoodsUnit       string    `gorm:"type:varchar(64);comment:商品单位"`
	SaleUnit        string    `gorm:"type:varchar(64);comment:销售单位"`
	SaleCode        string    `gorm:"comment:销售编码"`
	PayCount        int       `gorm:"comment:下单数量"`
	CloseCount      int       `gorm:"type:int;default:0;comment:关闭数量"`
	LackCount       int       `gorm:"type:int;comment:欠货数量"`
	OutCount        int       `gorm:"type:int;comment:出库数量"`
	LimitNum        int       `gorm:"type:int;default:0;comment:限发数量"`
	GoodsRemark     string    `gorm:"type:varchar(255);comment:商品备注"`
	Status          int       `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:拣货中,2:已出库"`
	BatchId         int       `gorm:"type:int(11);index;comment:批次id"`
	DeliveryOrderNo []*string `gorm:"type:varchar(255);comment:出库单号"`
}

// 完成订单表
type CompleteOrder struct {
	Base
	Number         string     `gorm:"type:varchar(64);unique;comment:订单编号"`
	OrderRemark    string     `gorm:"type:varchar(512);comment:订单备注"`
	ShopId         int        `gorm:"type:int(11);not null;comment:店铺id"`
	ShopName       string     `gorm:"type:varchar(64);not null;comment:店铺名称"`
	ShopType       string     `gorm:"type:varchar(64);not null;comment:店铺类型"`
	ShopCode       string     `gorm:"type:varchar(255);not null;comment:店铺编号"`
	Line           string     `gorm:"type:varchar(255);not null;comment:线路"`
	DeliveryMethod int        `gorm:"type:tinyint;not null;comment:配送方式"`
	PayCount       int        `gorm:"comment:下单数量"`
	CloseCount     int        `gorm:"type:int;comment:关闭数量"`
	OutCount       int        `gorm:"type:int;comment:出库数量"`
	Province       string     `gorm:"type:varchar(64);comment:省"`
	City           string     `gorm:"type:varchar(64);comment:市"`
	District       string     `gorm:"type:varchar(64);comment:区"`
	PickTime       *time.Time `gorm:"type:datetime;not null;comment:最近拣货时间"`
	PayAt          string     `gorm:"type:datetime;comment:支付时间"`
}

// 完成订单明细表
type CompleteOrderDetail struct {
	Base
	Number          string    `gorm:"type:varchar(64);index;comment:订单编号"`
	GoodsName       string    `gorm:"type:varchar(64);comment:商品名称"`
	Sku             string    `gorm:"type:varchar(64);comment:sku"`
	GoodsSpe        string    `gorm:"type:varchar(128);comment:商品规格"`
	GoodsType       string    `gorm:"type:varchar(64);comment:商品类型"`
	Shelves         string    `gorm:"type:varchar(64);comment:货架"`
	PayCount        int       `gorm:"comment:下单数量"`
	CloseCount      int       `gorm:"type:int;comment:关闭数量"`
	ReviewCount     int       `gorm:"type:int;comment:出库数量"`
	GoodsRemark     string    `gorm:"type:varchar(255);comment:商品备注"`
	DeliveryOrderNo []*string `gorm:"type:varchar(16);comment:出库单号"`
}

type RestrictedShipment struct {
	PickOrderGoodsId int    `gorm:"primaryKey;type:int(11) unsigned;comment:拣货单商品表id"`
	Number           string `gorm:"type:varchar(64);index;comment:订单编号"`
	ShopName         string `gorm:"type:varchar(64);comment:门店名称"`
	GoodsSpe         string `gorm:"type:varchar(128);comment:商品规格"`
	LimitNum         int    `gorm:"type:int;default:0;comment:限发数量"`
	Status           int    `gorm:"type:tinyint;default:1;comment:状态:0:撤销,1:正常"`
}

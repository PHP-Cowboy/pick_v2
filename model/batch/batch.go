package batch

import (
	"pick_v2/model"
	"time"
)

//
type BatchCondition struct {
	model.Base
	BatchId           int        `gorm:"type:int(11);index;comment:批次id"`
	WarehouseId       int        `gorm:"type:int(11);comment:仓库"`
	PayEndTime        *time.Time `gorm:"type:datetime;comment:支付截止时间"`
	DeliveryStartTime *time.Time `gorm:"type:datetime;default:null;comment:发货起始时间"`
	DeliveryEndTime   *time.Time `gorm:"type:datetime;comment:发货截止时间"`
	Line              string     `gorm:"type:varchar(255);comment:线路;default:null"`
	DeliveryMethod    int        `gorm:"type:tinyint;not null;comment:配送方式"`
	Sku               string     `gorm:"type:varchar(255);default:null;comment:商品sku"`
	Goods             string     `gorm:"type:varchar(255);comment:商品"`
}

//批次
type Batch struct {
	model.Base
	WarehouseId       int        `gorm:"type:int(11);comment:仓库"`
	BatchName         string     `gorm:"type:varchar(64);comment:批次名称"`
	DeliveryStartTime *time.Time `gorm:"type:datetime;default:null;comment:发货起始时间"`
	DeliveryEndTime   *time.Time `gorm:"type:datetime;comment:发货截止时间"`
	ShopNum           int        `gorm:"type:int(11);comment:门店数量"`
	OrderNum          int        `gorm:"type:int(11);comment:订单数量"`
	GoodsNum          int        `gorm:"type:int(11);comment:商品数量"`
	UserName          string     `gorm:"type:varchar(32);comment:用户名称"`
	Line              string     `gorm:"type:varchar(255);comment:所属路线"`
	DeliveryMethod    int        `gorm:"type:tinyint;not null;comment:配送方式"`
	EndTime           *time.Time `gorm:"type:datetime;comment:结束时间"`
	Status            int        `gorm:"type:tinyint;comment:状态"`
	PickNum           int        `gorm:"type:tinyint;comment:拣货单"`
	RecheckSheetNum   int        `gorm:"type:tinyint;comment:复核单数量"`
	Sort              int        `gorm:"type:int(11) unsigned;comment:排序"`
}

//预拣货列表
type PrePick struct {
	model.Base
	WarehouseId int    `gorm:"type:int(11);comment:仓库"`
	BatchId     int    `gorm:"type:int(11) unsigned;comment:批次表id"`
	ShopId      int    `gorm:"type:int(11);comment:店铺id"`
	ShopCode    string `gorm:"type:varchar(255);not null;comment:店铺编号"`
	ShopName    string `gorm:"type:varchar(64);not null;comment:店铺名称"`
	Line        string `gorm:"type:varchar(255);not null;comment:线路"`
	OrderNum    int    `gorm:"type:int(11);comment:订单数量"`
	Status      int    `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:已进入拣货池"`
}

//预拣货商品明细
type PrePickGoods struct {
	model.Base
	WarehouseId int    `gorm:"type:int(11);comment:仓库"`
	BatchId     int    `gorm:"type:int(11) unsigned;comment:批次表id"`
	OrderInfoId int    `gorm:"type:int(11) unsigned;comment:订单商品接口返回ID"`
	Number      string `gorm:"type:varchar(32);comment:订单编号"`
	ShopId      int    `gorm:"type:int(11);comment:店铺id"`
	PrePickId   int    `gorm:"type:int(11) unsigned;index;comment:预拣货表id"`
	GoodsName   string `gorm:"type:varchar(64);comment:商品名称"`
	GoodsType   string `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe    string `gorm:"type:varchar(128);comment:商品规格"`
	Shelves     string `gorm:"type:varchar(64);comment:货架"`
	NeedNum     int    `gorm:"type:int;not null;comment:需拣数量"`
	CloseNum    int    `gorm:"type:int;not null;comment:关闭数量"`
	OutCount    int    `gorm:"type:int;comment:出库数量"`
	NeedOutNum  int    `gorm:"type:int;comment:需出库数量"`
	Status      int    `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:已进入拣货池"`
}

//预拣货备注明细
type PrePickRemark struct {
	model.Base
	WarehouseId int    `gorm:"type:int(11);comment:仓库"`
	BatchId     int    `gorm:"type:int(11) unsigned;comment:批次表id"`
	OrderInfoId int    `gorm:"type:int(11) unsigned;comment:订单商品接口返回ID"`
	ShopId      int    `gorm:"type:int(11);comment:店铺id"`
	PrePickId   int    `gorm:"type:int(11) unsigned;index;comment:预拣货表id"`
	Number      string `gorm:"type:varchar(64);comment:订单编号"`
	OrderRemark string `gorm:"type:varchar(512);comment:订单备注"`
	GoodsRemark string `gorm:"type:varchar(255);comment:商品备注"`
	ShopName    string `gorm:"type:varchar(64);not null;comment:店铺名称"`
	Line        string `gorm:"type:varchar(255);not null;comment:线路"`
	Status      int    `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:已进入拣货池"`
}

//拣货列表
type Pick struct {
	model.Base
	WarehouseId    int        `gorm:"type:int(11);comment:仓库"`
	BatchId        int        `gorm:"type:int(11) unsigned;comment:批次表id"`
	PrePickIds     string     `gorm:"type:varchar(255);comment:预拣货ids"`
	TaskName       string     `gorm:"type:varchar(64);comment:任务名称"`
	ShopCode       string     `gorm:"type:varchar(255);not null;comment:店铺编号"`
	ShopName       string     `gorm:"type:varchar(64);not null;comment:店铺名称"`
	Line           string     `gorm:"type:varchar(255);not null;comment:线路"`
	ShopNum        int        `gorm:"type:int;not null;comment:门店数"`
	OrderNum       int        `gorm:"type:int;not null;comment:订单数"`
	NeedNum        int        `gorm:"type:int;not null;comment:需拣数量"`
	PickUser       string     `gorm:"type:varchar(32);default:'';comment:拣货人"`
	ReviewUser     string     `gorm:"type:varchar(32);default:'';comment:复核人"`
	TakeOrdersTime *time.Time `gorm:"type:datetime;default:null;comment:接单时间"`
	Sort           int        `gorm:"type:int(11) unsigned;comment:排序"`
	Version        int        `gorm:"type:tinyint(1);default:0;comment:版本"`
	Status         int        `gorm:"type:tinyint;comment:状态"`
}

//拣货商品明细
type PickGoods struct {
	model.Base
	WarehouseId    int    `gorm:"type:int(11);comment:仓库"`
	BatchId        int    `gorm:"type:int(11) unsigned;comment:批次表id"`
	PickId         int    `gorm:"type:int(11) unsigned;comment:拣货表id"`
	PrePickGoodsId int    `gorm:"type:int(11);comment:预拣货商品表id"`
	Number         string `gorm:"type:varchar(64);comment:订单编号"`
	ShopId         int    `gorm:"type:int(11);comment:店铺id"`
	GoodsName      string `gorm:"type:varchar(64);comment:商品名称"`
	GoodsType      string `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe       string `gorm:"type:varchar(128);comment:商品规格"`
	Shelves        string `gorm:"type:varchar(64);comment:货架"`
	NeedNum        int    `gorm:"type:int;not null;comment:需拣数量"`
}

//拣货备注明细
type PickRemark struct {
	model.Base
	WarehouseId     int    `gorm:"type:int(11);comment:仓库"`
	BatchId         int    `gorm:"type:int(11) unsigned;comment:批次表id"`
	PickId          int    `gorm:"type:int(11) unsigned;comment:拣货表id"`
	PrePickRemarkId int    `gorm:"type:int(11);comment:预拣货备注表id"`
	Number          string `gorm:"type:varchar(64);comment:订单编号"`
	OrderRemark     string `gorm:"type:varchar(512);comment:订单备注"`
	GoodsRemark     string `gorm:"type:varchar(255);comment:商品备注"`
	ShopName        string `gorm:"type:varchar(64);not null;comment:店铺名称"`
	Line            string `gorm:"type:varchar(255);not null;comment:线路"`
}

//完成订单 打标结束 改状态

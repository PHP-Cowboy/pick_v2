package model

import "time"

// 批次
type Batch struct {
	Base
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
	Status            int        `gorm:"type:tinyint;comment:状态:0:进行中,1:已结束,2:暂停"`
	PrePickNum        int        `gorm:"type:tinyint;comment:预拣单数"`
	PickNum           int        `gorm:"type:tinyint;comment:拣货单数"`
	RecheckSheetNum   int        `gorm:"type:tinyint;comment:复核单数"`
	Sort              int        `gorm:"type:int(11) unsigned;comment:排序"`
	PayEndTime        *time.Time `gorm:"type:datetime;comment:支付截止时间"`
	Version           int        `gorm:"type:int;default:0;comment:版本"`
}

const (
	BatchOngoingStatus = iota //进行中
	BatchClosedStatus
	BatchSuspendStatus
)

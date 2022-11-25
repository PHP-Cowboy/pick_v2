package model

import "gorm.io/gorm"

// 拣货列表
type Pick struct {
	Base
	DeliveryOrderNo GormList `gorm:"type:varchar(255);comment:出库单号"`
	WarehouseId     int      `gorm:"type:int(11);comment:仓库"`
	BatchId         int      `gorm:"type:int(11) unsigned;index;comment:批次表id"`
	PrePickIds      string   `gorm:"type:varchar(255);comment:预拣货ids"`
	TaskName        string   `gorm:"type:varchar(64);comment:任务名称"`
	ShopCode        string   `gorm:"type:varchar(255);not null;comment:店铺编号"`
	ShopName        string   `gorm:"type:varchar(64);not null;comment:店铺名称"`
	Line            string   `gorm:"type:varchar(255);not null;comment:线路"`
	ShopNum         int      `gorm:"type:int;default:0;comment:门店数"`
	OrderNum        int      `gorm:"type:int;default:0;comment:订单数"`
	NeedNum         int      `gorm:"type:int;default:0;comment:需拣数量"`
	PickNum         int      `gorm:"type:int;default:0;comment:拣货数量"`
	ReviewNum       int      `gorm:"type:int;default:0;comment:复核数量"`
	Num             int      `gorm:"type:int;default:0;comment:件数"`
	PrintNum        int      `gorm:"type:int;default:0;comment:打印次数"`
	PickUser        string   `gorm:"type:varchar(32);default:'';comment:拣货人"`
	TakeOrdersTime  *MyTime  `gorm:"type:datetime;default:null;comment:接单时间"`
	ReviewUser      string   `gorm:"type:varchar(32);default:'';comment:复核人"`
	ReviewTime      *MyTime  `gorm:"type:datetime;default:null;comment:复核时间"`
	Sort            int      `gorm:"type:int(11) unsigned;comment:排序"`
	Version         int      `gorm:"type:int;default:0;comment:版本"`
	Status          int      `gorm:"type:tinyint;comment:状态:0:待拣货,1:待复核,2:复核完成,3:停止拣货,4:终止拣货,5:返回预拣池"`
	DeliveryNo      string   `gorm:"type:varchar(255);index:delivery_no_idx;comment:出库单号"`
}

const (
	ToBePickedStatus         = iota //待拣货
	ToBeReviewedStatus              //待复核
	ReviewCompletedStatus           //复核完成
	StopPickingStatus               //停止拣货
	TerminationPickingStatus        //终止拣货
	ReturnPrePickStatus             //返回预拣池
)

type PickAndGoods struct {
	PickId   int    `json:"pick_id"`
	Status   int    `json:"status"`
	Number   string `json:"number"`
	PickUser string `json:"pick_user"`
}

func GetPickListByIds(db *gorm.DB, ids []int) (err error, list []Pick) {

	result := db.Model(&Pick{}).Where("id in (?)", ids).Find(&list)

	return result.Error, list
}

func UpdatePickByIds(db *gorm.DB, ids []int, mp map[string]interface{}) error {
	result := db.Model(&Pick{}).Where("id in (?)", ids).Updates(mp)

	return result.Error
}

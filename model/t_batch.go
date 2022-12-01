package model

import "gorm.io/gorm"

// 批次
type Batch struct {
	Base
	TaskId            int     `gorm:"type:int(11);index;comment:出库任务ID"`
	WarehouseId       int     `gorm:"type:int(11);comment:仓库"`
	BatchName         string  `gorm:"type:varchar(64);comment:批次名称"`
	DeliveryStartTime *MyTime `gorm:"type:datetime;default:null;comment:发货起始时间"`
	DeliveryEndTime   *MyTime `gorm:"type:datetime;comment:发货截止时间"`
	UserName          string  `gorm:"type:varchar(32);comment:用户名称"`
	Line              string  `gorm:"type:varchar(255);comment:所属路线"`
	DeliveryMethod    int     `gorm:"type:tinyint;not null;comment:配送方式"`
	EndTime           *MyTime `gorm:"type:datetime;comment:结束时间"`
	Status            int     `gorm:"type:tinyint;comment:状态:0:进行中,1:已结束,2:暂停"`
	Sort              int     `gorm:"type:int(11) unsigned;comment:排序"`
	PayEndTime        *MyTime `gorm:"type:datetime;comment:支付截止时间"`
	Version           int     `gorm:"type:int;default:0;comment:版本"`
	Typ               int     `gorm:"type:tinyint;default:1;comment:批次类型:1:常规批次,2:快递批次"`

	//ShopNum         int `gorm:"type:int(11);comment:门店数量"`
	//OrderNum        int `gorm:"type:int(11);comment:订单数量"`
	//GoodsNum        int `gorm:"type:int(11);comment:商品数量"`
	//PrePickNum      int `gorm:"type:tinyint;comment:预拣单数"`
	//PickNum         int `gorm:"type:tinyint;comment:拣货单数"`
	//RecheckSheetNum int `gorm:"type:tinyint;comment:复核单数"`
}

const (
	BatchOngoingStatus = iota //进行中
	BatchClosedStatus         //已结束
	BatchSuspendStatus        //暂停
)

const (
	_                       = iota
	RegularBatchTyp         //常规批次
	ExpressDeliveryBatchTyp //快递批次
)

// 获取配送方式
func GetDeliveryMethod(method int) string {
	var methodMp = map[int]string{1: "公司配送", 2: "用户自提", 3: "三方物流", 4: "快递配送", 5: "首批物料|设备单"}

	s, ok := methodMp[method]

	if !ok {
		return ""
	}

	return s
}

func BatchSave(db *gorm.DB, list Batch) (err error, b Batch) {

	result := db.Model(&Batch{}).Save(&list)

	return result.Error, list
}

// 通过主键查询数据
func GetBatchByPk(db *gorm.DB, pk int) (err error, batch Batch) {
	result := db.Model(&Batch{}).First(&batch, pk)

	return result.Error, batch
}

// 根据出库任务获取批次列表
func GetBatchListByTaskId(db *gorm.DB, taskId int) (err error, list []Batch) {

	result := db.Model(&Batch{}).Where(&Batch{TaskId: taskId}).Find(&list)

	return result.Error, list
}

// 快递批次进行中或暂停的单数量
func GetBatchListByTyp(db *gorm.DB, typ int) (err error, list []Batch) {
	result := db.Model(&Batch{}).Where("typ = ? and ( status = 0 or status = 2 )", typ).Find(&list)

	return result.Error, list
}

// 获取进行中或已被当前用户接取但为拣货完成的批次
func GetPendingOrUserReceivedNotPickCompleteBatchList(db *gorm.DB, userName string) {
	//db.Table("t_batch b").
	//	Select("").
	//	Joins("").Where("").Find()
}

func GetBatchList(db *gorm.DB, cond Batch) (err error, list []Batch) {
	result := db.Model(&Batch{}).Where(&cond).Find(&list)

	return result.Error, list
}

func GetBatchListByIdsOrPending(db *gorm.DB, ids []int) (err error, list []Batch) {
	result := db.Model(&Batch{}).Where("id in (?) or status = ? and typ = ?", ids, BatchOngoingStatus, ExpressDeliveryBatchTyp).Find(&list)

	return result.Error, list
}

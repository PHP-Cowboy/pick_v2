package model

import (
	"errors"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
	"pick_v2/utils/ecode"
)

// 拣货列表
type Pick struct {
	Base
	WarehouseId     int      `gorm:"type:int(11);comment:仓库"`
	TaskId          int      `gorm:"type:int(11) unsigned;comment:出库任务id"`
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
	OutboundType    int      `gorm:"type:tinyint;default:1;comment:状态:1.正常,2.无需出库"`
	DeliveryNo      string   `gorm:"type:varchar(255);index:delivery_no_idx;comment:出库单号"`
	Typ             int      `gorm:"type:tinyint;default:1;comment:批次类型:1:常规批次,2:快递批次"`
	DeliveryOrderNo GormList `gorm:"type:varchar(255);comment:出库单号"`
}

const (
	_                         = iota
	OutboundTypeNormal        //正常
	OutboundTypeNoNeedToIssue //无需出库
)

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

func PickBatchSave(db *gorm.DB, picks *[]Pick) (err error) {
	err = db.Model(&Pick{}).CreateInBatches(picks, BatchSize).Error
	return
}

func PickSave(db *gorm.DB, picks *Pick) (err error) {
	err = db.Model(&Pick{}).Save(picks).Error
	return
}

func PickReplaceSave(db *gorm.DB, list *[]Pick, values []string) (err error) {
	err = db.Model(&Pick{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns(values),
		}).
		CreateInBatches(list, BatchSize).
		Error

	return
}

func UpdatePickByPk(db *gorm.DB, pk int, mp map[string]interface{}) (err error) {
	err = db.Model(&Pick{}).Where("id = ?", pk).Updates(mp).Error

	return
}

// 根据主键和版本号更新，乐观锁
func UpdatePickByPkAndVersion(db *gorm.DB, pk, version int, mp map[string]interface{}) (err error) {
	err = db.Model(&Pick{}).Where("id = ? and version = ?", pk, version).Updates(mp).Error

	return
}

func UpdatePickByIds(db *gorm.DB, ids []int, mp map[string]interface{}) (err error) {
	err = db.Model(&Pick{}).Where("id in (?)", ids).Updates(mp).Error

	return
}

func GetPickByPk(db *gorm.DB, id int) (err error, pick Pick) {
	err = db.Model(&Pick{}).First(&pick, id).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = ecode.DataNotExist
			return
		}
	}

	return
}

func GetPickListByIds(db *gorm.DB, ids []int) (err error, list []Pick) {

	err = db.Model(&Pick{}).Where("id in (?)", ids).Find(&list).Error

	return
}

// 根据条件获取拣货列表
func GetPickList(db *gorm.DB, cond *Pick) (err error, list []Pick) {

	err = db.Model(&Pick{}).Where(cond).Find(&list).Error

	return
}

// 查询当前拣货员被分配的任务或已经接单 且未完成拣货 的数据
func GetPickListByPickUserAndStatusAndTyp(db *gorm.DB, pickUser string, status int, typ int) (err error, list []Pick) {
	err = db.Model(&Pick{}).
		Where("pick_user = ? and status = ? and typ = ?", pickUser, status, typ).
		Find(&list).
		Error

	return
}

// 查询未被接单的拣货池数据
func GetPickListNoOrderReceived(db *gorm.DB, batchIds []int, typ int) (err error, list []Pick) {
	err = db.Model(&Pick{}).
		Where("batch_id in (?) and pick_user = '' and status = ? and typ = ?", batchIds, ToBePickedStatus, typ).
		Find(&list).
		Error

	return
}

// 获取当前用户已接单但未复核完成的批次
func GetPickListByPickUserAndNotReviewCompleted(db *gorm.DB, userName string, typ int) (err error, list []Pick) {
	err = db.Model(&Pick{}).
		Where("pick_user = ? and status in (?) and typ = ?", userName, []int{ToBePickedStatus, ToBeReviewedStatus}, typ).
		Find(&list).
		Error

	return
}

// 根据出库单号查询拣货数据
func GetPickListByDeliveryNo(db *gorm.DB, deliveryNo []string) (err error, list []Pick) {
	err = db.Model(&Pick{}).Where("delivery_no in (?)", deliveryNo).Find(&list).Error
	return
}

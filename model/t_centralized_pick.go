package model

import "gorm.io/gorm"

// 集中拣货
type CentralizedPick struct {
	Base
	BatchId        int     `gorm:"type:int(11) unsigned;uniqueIndex:batchSku;comment:批次表id"`
	Sku            string  `gorm:"type:varchar(64);uniqueIndex:batchSku;comment:sku"`
	TaskId         int     `gorm:"type:int(11) unsigned;comment:任务表id"`
	GoodsName      string  `gorm:"type:varchar(64);comment:商品名称"`
	GoodsType      string  `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe       string  `gorm:"type:varchar(128);comment:商品规格"`
	Shelves        string  `gorm:"type:varchar(64);comment:货架"`
	NeedNum        int     `gorm:"type:int;default:0;comment:需拣数量"`
	CompleteNum    int     `gorm:"type:int;default:0;comment:拣货数量"`
	PickUser       string  `gorm:"type:varchar(32);default:'';comment:拣货人"`
	TakeOrdersTime *MyTime `gorm:"type:datetime;default:null;comment:接单时间"`
	GoodsRemark    string  `gorm:"type:varchar(255);comment:商品备注"`
	GoodsUnit      string  `gorm:"type:varchar(64);comment:商品单位"`
	Status         int     `gorm:"type:tinyint;comment:状态:0:待拣货,1:已完成"`
	Version        int     `gorm:"type:int;default:0;comment:版本"`
	Sort           int     `gorm:"type:int(11) unsigned;comment:排序"`
}

const (
	CentralizedPickStatusPending   = iota //待拣货
	CentralizedPickStatusCompleted        //已完成
)

func CentralizedPickSave(db *gorm.DB, list *[]CentralizedPick) (err error) {
	err = db.Model(&CentralizedPick{}).Save(list).Error

	return
}

func CentralizedPickReplaceSave(db *gorm.DB, cp CentralizedPick) (err error) {
	err = db.Model(&CentralizedPick{}).Save(&cp).Error

	return
}

func UpdateCentralizedPickById(db *gorm.DB, id int, mp map[string]interface{}) (err error) {
	err = db.Model(&CentralizedPick{}).Where("id = ?", id).Updates(mp).Error

	return
}

// 集中拣货分页列表
func GetCentralizedPickPageList(db *gorm.DB, batchId, page, size, isRemark int, goodsName, goodsType string) (err error, total int64, list []CentralizedPick) {

	local := db.Model(&CentralizedPick{}).
		Where(&CentralizedPick{BatchId: batchId, GoodsName: goodsName, GoodsType: goodsType})

	if isRemark == 0 {
		local.Where("goods_remark = ''")
	} else if isRemark == 1 {
		local.Where("goods_remark != ''")
	}

	err = local.Count(&total).Error

	if err != nil {
		return
	}

	//条数为0，直接返回
	if total == 0 {
		return
	}

	err = local.Scopes(Paginate(page, size)).Find(&list).Error

	return
}

// 集中拣货列表
func GetCentralizedPickList(db *gorm.DB, cond CentralizedPick) (err error, list []CentralizedPick) {

	err = db.Model(&CentralizedPick{}).
		Where(&cond).
		Find(&list).
		Error

	return
}

func GetCentralizedPickByBatchSku(db *gorm.DB, batch int, sku string) (err error, list []CentralizedPick) {
	err = db.Model(&CentralizedPick{}).
		Where("batch = ? and sku = ?", batch, sku).
		Find(&list).
		Error

	return
}

// 查询是否有当前拣货员被分配的任务或已经接单且未完成拣货的数据,如果被分配多条，第一按批次优先级，第二按拣货池优先级 优先拣货
func GetCentralizedPickByPickUser(db *gorm.DB, userName string) (err error, list []CentralizedPick) {
	err = db.Model(&CentralizedPick{}).
		Where("pick_user = ? and status = ?", userName, CentralizedPickStatusPending).
		Find(&list).
		Error

	return
}

func CountCentralizedPickByBatchAndUser(db *gorm.DB, batchIds []int, userName string) (err error, count int64) {
	err = db.Model(&CentralizedPick{}).
		Where("batch_id in (?) and (pick_user = '' or pick_user = ?)", batchIds, userName).
		Count(&count).
		Error

	return
}

type CountPickNums struct {
	BatchId          int `json:"batch_id"`
	SumNeedNum       int `json:"sum_need_num"`
	SumPickNum       int `json:"sum_pick_num"`
	SumCompleteNum   int `json:"sum_complete_num"`
	CountNeedNum     int `json:"count_need_num"`
	CountPickNum     int `json:"count_pick_num"`
	CountCompleteNum int `json:"count_complete_num"`
}

// 统计集中拣货相关数量
func CountCentralizedPickByBatch(db *gorm.DB, batchIds []int) (err error, countCentralizedPick []CountPickNums) {

	err = db.Model(&CentralizedPick{}).
		Select("batch_id,sum(need_num) as sum_need_num,sum(complete_num) as sum_pick_num,count(need_num) as count_need_num,count(complete_num) as count_pick_num").
		Where("batch_id in (?)", batchIds).
		Group("batch_id").
		Find(&countCentralizedPick).
		Error

	return
}

func GetCentralizedPickById(db *gorm.DB, id int) (err error, first CentralizedPick) {
	err = db.Model(&CentralizedPick{}).First(&first, id).Error

	return
}

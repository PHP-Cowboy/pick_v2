package model

import "gorm.io/gorm"

// 集中拣货
type CentralizedPick struct {
	Base
	BatchId        int     `gorm:"type:int(11) unsigned;uniqueIndex:batchSku;comment:批次表id"`
	Sku            string  `gorm:"type:varchar(64);uniqueIndex:batchSku;comment:sku"`
	GoodsName      string  `gorm:"type:varchar(64);comment:商品名称"`
	GoodsType      string  `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe       string  `gorm:"type:varchar(128);comment:商品规格"`
	Shelves        string  `gorm:"type:varchar(64);comment:货架"`
	NeedNum        int     `gorm:"type:int;default:0;comment:需拣数量"`
	PickNum        int     `gorm:"type:int;default:0;comment:拣货数量"`
	PickUser       string  `gorm:"type:varchar(32);default:'';comment:拣货人"`
	TakeOrdersTime *MyTime `gorm:"type:datetime;default:null;comment:接单时间"`
	GoodsRemark    string  `gorm:"type:varchar(255);comment:商品备注"`
	GoodsUnit      string  `gorm:"type:varchar(64);comment:商品单位"`
	Status         int     `gorm:"type:tinyint;comment:状态:0:待拣货,1:已完成"`
	Version        int     `gorm:"type:int;default:0;comment:版本"`
	Sort           int     `gorm:"type:int(11) unsigned;comment:排序"`
	PickType       int     `gorm:"type:tinyint;default:0;comment:状态:0:正常拣货,1:无需拣货"`
}

const (
	CentralizedPickStatusPending   = iota //待拣货
	CentralizedPickStatusCompleted        //已完成
)

func CentralizedPickSave(db *gorm.DB, list *[]CentralizedPick) error {
	result := db.Model(&CentralizedPick{}).Save(list)

	return result.Error
}

func CentralizedPickReplaceSave(db *gorm.DB, cp CentralizedPick) error {
	result := db.Model(&CentralizedPick{}).Save(&cp)

	return result.Error
}

func UpdateCentralizedPickById(db *gorm.DB, id int, mp map[string]interface{}) error {
	result := db.Model(&CentralizedPick{}).Where("id = ?", id).Updates(mp)

	return result.Error
}

// 集中拣货分页列表
func GetCentralizedPickPageList(db *gorm.DB, page, size, isRemark int, goodsName, goodsType string) (err error, total int64, list []CentralizedPick) {

	local := db.Model(&CentralizedPick{}).Where(&CentralizedPick{GoodsName: goodsName, GoodsType: goodsType})

	if isRemark == 0 {
		local.Where("goods_remark = ''")
	} else if isRemark == 1 {
		local.Where("goods_remark != ''")
	}

	result := local.Find(&list)

	if result.Error != nil {
		return result.Error, 0, nil
	}

	total = result.RowsAffected

	result = local.Scopes(Paginate(page, size)).Find(&list)

	if result.Error != nil {
		return result.Error, 0, nil
	}

	return nil, total, list
}

// 集中拣货列表
func GetCentralizedPickList(db *gorm.DB, cond CentralizedPick) (err error, list []CentralizedPick) {

	result := db.Model(&CentralizedPick{}).Where(&cond).Find(&list)

	return result.Error, list
}

func GetCentralizedPickByBatchSku(db *gorm.DB, batch int, sku string) (err error, list []CentralizedPick) {
	result := db.Model(&CentralizedPick{}).Where("batch = ? and sku = ?", batch, sku).Find(&list)

	return result.Error, list
}

// 查询是否有当前拣货员被分配的任务或已经接单且未完成拣货的数据,如果被分配多条，第一按批次优先级，第二按拣货池优先级 优先拣货
func GetCentralizedPickByPickUser(db *gorm.DB, userName string) (err error, list []CentralizedPick) {
	result := db.Model(&CentralizedPick{}).
		Where("pick_user = ? and status = ?", userName, CentralizedPickStatusPending).
		Find(&list)

	if result.Error != nil {
		return result.Error, nil
	}

	return nil, list
}

func CountCentralizedPickByBatchAndUser(db *gorm.DB, batchIds []int, userName string) (err error, count int64) {
	result := db.Model(&CentralizedPick{}).
		Where("batch_id in (?) and (pick_user = '' or pick_user = ?)", batchIds, userName).
		Count(&count)

	return result.Error, count
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

	result := db.Model(&CentralizedPick{}).
		Select("batch_id,sum(need_num) as sum_need_num,sum(pick_num) as sum_pick_num,count(need_num) as count_need_num,count(pick_num) as count_pick_num").
		Where("batch_id in (?)", batchIds).
		Group("batch_id").
		Find(&countCentralizedPick)

	return result.Error, countCentralizedPick
}

func GetCentralizedPickById(db *gorm.DB, id int) (err error, first CentralizedPick) {
	result := db.Model(&CentralizedPick{}).First(&first, id)

	return result.Error, first
}

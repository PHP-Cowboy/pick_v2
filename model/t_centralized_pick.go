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
	NeedNum        int     `gorm:"type:int;default:0;comment:需拣数量"`
	PickNum        int     `gorm:"type:int;default:0;comment:拣货数量"`
	PickUser       string  `gorm:"type:varchar(32);default:'';comment:拣货人"`
	TakeOrdersTime *MyTime `gorm:"type:datetime;default:null;comment:接单时间"`
	GoodsRemark    string  `gorm:"type:varchar(255);comment:商品备注"`
	GoodsUnit      string  `gorm:"type:varchar(64);comment:商品单位"`
	Status         int     `gorm:"type:tinyint;comment:状态:0:待拣货,1:已完成"`
}

const (
	CentralizedPickStatusPending   = iota //待拣货
	CentralizedPickStatusCompleted        //已完成
)

func CentralizedPickSave(db *gorm.DB, list *[]CentralizedPick) error {
	result := db.Model(&CentralizedPick{}).Save(list)

	return result.Error
}

func GetCentralizedPickList(db *gorm.DB, page, size int, query interface{}, args ...interface{}) (err error, total int64, list []CentralizedPick) {
	result := db.Model(&CentralizedPick{}).Where(query, args).Find(&list)

	if result.Error != nil {
		return result.Error, 0, nil
	}

	total = result.RowsAffected

	result = db.Model(&CentralizedPick{}).Scopes(Paginate(page, size)).Where(query, args).Find(&list)

	if result.Error != nil {
		return result.Error, 0, nil
	}

	return nil, total, list
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

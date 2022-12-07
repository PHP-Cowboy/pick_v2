package model

import "gorm.io/gorm"

type LimitShipment struct {
	TaskId    int    `gorm:"primaryKey;type:int(11);not null;comment:t_outbound_task表ID" json:"task_id"`
	Number    string `gorm:"primaryKey;type:varchar(64);comment:订单编号" json:"number"`
	Sku       string `gorm:"primaryKey;type:varchar(64);comment:sku" json:"sku"`
	ShopName  string `gorm:"type:varchar(64);comment:门店名称" json:"shop_name"`
	GoodsName string `gorm:"type:varchar(64);comment:商品名称" json:"goods_name"`
	GoodsSpe  string `gorm:"type:varchar(128);comment:商品规格" json:"goods_spe"`
	LimitNum  int    `gorm:"type:int;default:0;comment:限发数量" json:"limit_num"`
	Status    int    `gorm:"type:tinyint;default:1;comment:状态:0:撤销,1:正常" json:"status"`
	Typ       int    `gorm:"type:tinyint;default:1;comment:状态:1:订单限发,2:任务限发" json:"typ"`
}

const (
	LimitShipmentStatusRevoke = iota //撤销
	LimitShipmentStatusNormal        //正常
)

const (
	LimitShipmentTyp      = iota
	LimitShipmentTypOrder //订单限发
	LimitShipmentTypTask  //任务限发
)

func LimitShipmentSave(db *gorm.DB, list []LimitShipment) (err error) {
	err = db.Model(&LimitShipment{}).CreateInBatches(&list, BatchSize).Error
	return
}

// 根据任务ID和订单编号查询限发列表
func GetLimitShipmentListByTaskIdAndNumbers(db *gorm.DB, taskId int, number []string) (err error, list []LimitShipment) {

	err = db.Model(&LimitShipment{}).Where("task_id = ? and number in (?)", taskId, number).Find(&list).Error

	return
}

// 查询 订单
func GetLimitShipmentPageListByTaskIdAndNumber(db *gorm.DB, taskId int, number string, page, size int) (err error, total int64, list []LimitShipment) {
	local := db.Model(&LimitShipment{}).
		Where(&LimitShipment{
			TaskId: taskId,
			Number: number,
			Status: LimitShipmentStatusNormal,
		})

	err = local.Count(&total).Error

	if err != nil {
		return
	}

	//没查到直接返回
	if total == 0 {
		return
	}

	err = local.Scopes(Paginate(page, size)).Find(&list).Error

	return
}

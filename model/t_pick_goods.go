package model

import (
	"errors"
	"gorm.io/gorm"
)

// 拣货商品明细
type PickGoods struct {
	Base
	WarehouseId      int    `gorm:"type:int(11);comment:仓库"`
	PickId           int    `gorm:"type:int(11) unsigned;index:pick_and_batch_idx;comment:拣货表id"`
	BatchId          int    `gorm:"type:int(11) unsigned;index:pick_and_batch_idx;comment:批次表id"`
	PrePickGoodsId   int    `gorm:"type:int(11);comment:预拣货商品表id"`
	OrderGoodsId     int    `gorm:"type:int(11) unsigned;comment:订单商品表ID"` //t_pick_order_goods 表 id
	Number           string `gorm:"type:varchar(64);comment:订单编号"`
	ShopId           int    `gorm:"type:int(11);comment:店铺id"`
	DistributionType int    `gorm:"type:tinyint unsigned;comment:配送方式:1:公司配送,2:用户自提,3:三方物流,4:快递配送,5:首批物料|设备单"`
	Sku              string `gorm:"type:varchar(64);comment:sku"`
	GoodsName        string `gorm:"type:varchar(64);comment:商品名称"`
	GoodsType        string `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe         string `gorm:"type:varchar(128);comment:商品规格"`
	Shelves          string `gorm:"type:varchar(64);comment:货架"`
	DiscountPrice    int    `gorm:"comment:折扣价"`
	NeedNum          int    `gorm:"type:int;not null;comment:需拣数量"`
	CompleteNum      int    `gorm:"type:int;default:null;comment:已拣数量"` //默认为null，无需拣货或者拣货数量为0时更新为0
	ReviewNum        int    `gorm:"type:int;default:0;comment:复核数量"`
	Unit             string `gorm:"type:varchar(64);comment:单位"`
}

func GetPickGoodsByNumber(db *gorm.DB, numbers []string) (err error, list []PickGoods) {
	result := db.Model(&PickGoods{}).Where("number in (?)", numbers).Find(&list)

	return result.Error, list
}

// 根据id查拣货池商品数据
func GetPickGoodsByPickIds(db *gorm.DB, pickIds []int) (err error, list []PickGoods) {
	result := db.Model(&PickGoods{}).Where("pick_id in (?)", pickIds).Find(&list)

	return result.Error, list
}

// 查询拣货池商品订单是否有已拣的
func GetFirstPickGoodsByNumbers(db *gorm.DB, numbers []string) (err error, exist bool) {

	var pickGoods PickGoods

	//
	result := db.Model(&PickGoods{}).Where("number in (?) and complete_num >= 0", numbers).First(&pickGoods)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, false
		}

		return result.Error, true
	}

	return nil, false
}

func UpdatePickGoodsByPickIds(db *gorm.DB, pickIds []int, mp map[string]interface{}) error {
	result := db.Model(&Pick{}).Where("pick_id in (?)", pickIds).Updates(mp)

	return result.Error
}

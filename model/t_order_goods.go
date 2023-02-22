package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// 订单商品表
type OrderGoods struct {
	Id              int       `gorm:"primaryKey;type:int(11) unsigned;comment:id"`
	CreateTime      time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime      time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime      time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
	Number          string    `gorm:"type:varchar(64);index:number_sku_idx;comment:订单编号"`
	GoodsName       string    `gorm:"type:varchar(64);comment:商品名称"`
	Sku             string    `gorm:"type:varchar(64);index:number_sku_idx;comment:sku"`
	GoodsType       string    `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe        string    `gorm:"type:varchar(128);comment:商品规格"`
	Shelves         string    `gorm:"type:varchar(64);comment:货架"`
	DiscountPrice   int       `gorm:"type:int(11);comment:折扣价"`
	GoodsUnit       string    `gorm:"type:varchar(64);comment:商品单位"`
	SaleUnit        string    `gorm:"type:varchar(64);comment:销售单位"`
	SaleCode        string    `gorm:"type:varchar(255);comment:销售编码"`
	PayCount        int       `gorm:"type:int(11);comment:下单数量"`
	CloseCount      int       `gorm:"type:int(11);default:0;comment:关闭数量"`
	LackCount       int       `gorm:"type:int(11);comment:欠货数量"`
	OutCount        int       `gorm:"type:int(11);comment:出库数量"`
	GoodsRemark     string    `gorm:"type:varchar(255);comment:商品备注"`
	BatchId         int       `gorm:"type:int(11);index;comment:批次id"`
	Status          int       `gorm:"type:tinyint;default:1;comment:状态:1:未处理,2:处理中"`
	DeliveryOrderNo GormList  `gorm:"type:varchar(255);comment:出库单号"`
}

// 商品相关数量统计
type OrderGoodsNumsStatistical struct {
	Number     string `json:"number"`
	PayCount   int    `json:"pay_count"`
	CloseCount int    `json:"close_count"`
	OutCount   int    `json:"out_count"`
	LackCount  int    `json:"lack_count"`
}

// 历史订单商品
type HistoryOrderGoods struct {
	DeliveryNumber GormList `json:"delivery_number"` //出库单号
	OutCount       int      `json:"out_count"`       //出库数量
}

// 状态
const (
	_                          = iota
	OrderGoodsUnhandledStatus  //未处理(商品中不区分是欠货还是新订单，订单表中能区分，未完成的统一定义为未处理)
	OrderGoodsProcessingStatus //处理中
)

// 批量保存订单商品
func OrderGoodsBatchSave(db *gorm.DB, list *[]OrderGoods) (err error) {
	err = db.Model(&OrderGoods{}).CreateInBatches(list, BatchSize).Error

	return
}

func OrderGoodsReplaceSave(db *gorm.DB, list *[]OrderGoods, values []string) (err error) {

	if len(*list) == 0 {
		return
	}

	err = db.Model(&OrderGoods{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns(values),
		}).
		CreateInBatches(list, BatchSize).
		Error

	return
}

// 批量更新订单商品数据
func UpdateOrderGoodsByIds(db *gorm.DB, ids []int, mp map[string]interface{}) (err error) {
	err = db.Model(&OrderGoods{}).Where("id in (?)", ids).Updates(mp).Error

	return
}

// 批量更新订单商品数据
func UpdateOrderGoodsByNumbers(db *gorm.DB, numbers []string, mp map[string]interface{}) (err error) {
	err = db.Model(&OrderGoods{}).Where("number in (?)", numbers).Updates(mp).Error

	return
}

// 通过ids变更订单类型&&商品状态
func UpdateOrderAndGoodsByIds(db *gorm.DB, orderIds []int, orderGoodsIds []int, orderType, status int) (err error) {

	err = UpdateOrderByIds(db, orderIds, map[string]interface{}{"order_type": orderType})

	if err != nil {
		return err
	}

	err = UpdateOrderGoodsByIds(db, orderGoodsIds, map[string]interface{}{"status": status})

	if err != nil {
		return err
	}

	return nil
}

// 通过numbers变更订单类型&&商品状态
func UpdateOrderAndGoodsByNumbers(db *gorm.DB, numbers []string, orderType, status int) (err error) {

	err = UpdateOrderByNumbers(db, numbers, map[string]interface{}{"order_type": orderType})

	if err != nil {
		return err
	}

	err = UpdateOrderGoodsByNumbers(db, numbers, map[string]interface{}{"status": status})

	if err != nil {
		return err
	}

	return nil
}

// 变更订单商品状态
func UpdateOrderGoodsStatus(db *gorm.DB, list []OrderGoods, values []string) (err error) {
	err = db.Model(&OrderGoods{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(values),
	}).
		Save(&list).
		Error

	return
}

func DeleteOrderGoodsByNumbers(db *gorm.DB, numbers []string) (err error) {
	err = db.Delete(&OrderGoods{}, "number in (?)", numbers).Error
	return
}

func DeleteOrderGoodsByIds(db *gorm.DB, ids []int) (err error) {
	if len(ids) == 0 {
		return
	}

	err = db.Delete(&OrderGoods{}, "id in (?)", ids).Error
	return
}

// 订单&&商品信息
func GetOrderJoinGoodsListByNumbers(db *gorm.DB, number []string) (err error, list []GoodsJoinOrder) {
	err = db.Table("t_order_goods og").
		Select("o.*,o.id as order_id,og.*,og.id as order_goods_id").
		Joins("left join t_order o on og.number = o.number").
		Where("og.number in (?)", number).
		Find(&list).
		Error

	return
}

// 根据orderGoodsId查订单商品数据并根据支付时间排序
func GetOrderGoodsJoinOrderByIds(db *gorm.DB, ids []int) (err error, list []GoodsJoinOrder) {
	err = db.Table("t_order_goods og").
		Select("o.*,o.id as order_id,og.*,og.id as order_goods_id").
		Joins("left join t_order o on og.number = o.number").
		Where("og.id in (?)", ids).
		Order("pay_at ASC").
		Find(&list).
		Error

	return
}

// 根据orderGoodsId查订单商品数据并根据支付时间排序
func GetOrderGoodsJoinOrderByIdsNoSort(db *gorm.DB, ids []int) (err error, list []GoodsJoinOrder) {
	err = db.Table("t_order_goods og").
		Select("o.*,o.id as order_id,og.*,og.id as order_goods_id").
		Joins("left join t_order o on og.number = o.number").
		Where("og.id in (?)", ids).
		Find(&list).
		Error

	return
}

// 根据批次id查询订单&&订单商品数据
func GetOrderGoodsJoinOrderByBatchId(db *gorm.DB, batchId int) (err error, list []GoodsJoinOrder) {
	err = db.Table("t_order_goods og").
		Select("o.*,o.id as order_id,og.*,og.id as order_goods_id").
		Joins("left join t_order o on og.number = o.number").
		Where("og.batch_id = ?", batchId).
		Find(&list).
		Error

	return
}

// 根据订单商品id查询数据
func GetOrderGoodsListByIds(db *gorm.DB, ids []int) (err error, list []OrderGoods) {
	err = db.Model(&OrderGoods{}).Where("id in (?)", ids).Find(&list).Error
	return
}

func GetOrderGoodsListByNumbers(db *gorm.DB, numbers []string) (err error, list []OrderGoods) {
	err = db.Model(&OrderGoods{}).Where("number in (?)", numbers).Find(&list).Error
	return
}

// 获取出库任务订单 商品相关数量
func OrderGoodsNumsStatisticalByNumbers(db *gorm.DB, query string, number []string) (err error, mp map[string]OrderGoodsNumsStatistical) {
	var nums []OrderGoodsNumsStatistical

	err = db.Model(&OrderGoods{}).
		Select(query).
		Where("number in (?)", number).
		Group("number").
		Find(&nums).
		Error

	if err != nil {
		return
	}

	mp = make(map[string]OrderGoodsNumsStatistical, 0)

	for _, n := range nums {
		mp[n.Number] = n
	}

	for _, s := range number {
		num, ok := mp[s]
		//订单统计没有数据时赋值为0
		if !ok {
			num.LackCount = 0
			num.OutCount = 0
			num.PayCount = 0
			num.CloseCount = 0

			mp[s] = num
		}
	}

	return err, mp
}

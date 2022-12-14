package model

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 拣货商品明细
type PickGoods struct {
	Base
	WarehouseId      int    `gorm:"type:int(11);comment:仓库"`
	PickId           int    `gorm:"type:int(11) unsigned;index:pick_and_batch_idx;comment:拣货表id"`
	BatchId          int    `gorm:"type:int(11) unsigned;index:pick_and_batch_idx;comment:批次表id"`
	PrePickGoodsId   int    `gorm:"type:int(11);comment:预拣货商品表id"`
	OrderGoodsId     int    `gorm:"type:int(11) unsigned;index;comment:订单商品表ID"` //t_pick_order_goods 表 id
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
	CloseNum         int    `gorm:"type:int;not null;comment:关闭数量"`
	Status           int    `gorm:"type:tinyint;default:1;comment:状态:1:正常,2:关闭"`
	Unit             string `gorm:"type:varchar(64);comment:单位"`
}

type PickGoodsJoinOrder struct {
	ShopCode  string `json:"shop_code"`
	GoodsName string `json:"goods_name"`
	GoodsSpe  string `json:"goods_spe"`
	Unit      string `json:"unit"`
	Sku       string `json:"sku"`
	NeedNum   int    `json:"need_num"`
	ReviewNum int    `json:"review_num"`
}

const (
	_                     = iota
	PickGoodsStatusNormal //正常
	PickGoodsStatusClosed //关闭
)

func PickGoodsSave(db *gorm.DB, list *[]PickGoods) (err error) {
	err = db.Model(&PickGoods{}).CreateInBatches(list, BatchSize).Error
	return
}

func PickGoodsReplaceSave(db *gorm.DB, list *[]PickGoods, values []string) (err error) {

	err = db.Model(&PickGoods{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns(values),
		}).CreateInBatches(list, BatchSize).
		Error

	return
}

func UpdatePickGoodsByIds(db *gorm.DB, ids []int, mp map[string]interface{}) (err error) {
	err = db.Model(&PickGoods{}).Where("id in (?)", ids).Updates(mp).Error
	return
}

// 根据pickId更新拣货池商品数据
func UpdatePickGoodsByPickId(db *gorm.DB, pickId int, mp map[string]interface{}) (err error) {
	err = db.Model(&PickGoods{}).Where("pick_id = ?", pickId).Updates(mp).Error
	return
}

func UpdatePickGoodsByPickIds(db *gorm.DB, pickIds []int, mp map[string]interface{}) (err error) {
	err = db.Model(&PickGoods{}).Where("pick_id in (?)", pickIds).Updates(mp).Error

	return
}

func GetPickGoodsList(db *gorm.DB, cond *PickGoods) (err error, list []PickGoods) {
	err = db.Model(&PickGoods{}).Where(cond).Find(&list).Error
	return
}

func GetPickGoodsByIds(db *gorm.DB, ids []int) (err error, list []PickGoods) {
	err = db.Model(&PickGoods{}).Where("id in (?)", ids).Find(&list).Error

	return
}

func GetPickGoodsByNumber(db *gorm.DB, numbers []string) (err error, list []PickGoods) {
	err = db.Model(&PickGoods{}).Where("number in (?)", numbers).Find(&list).Error

	return
}

// 根据拣货id查拣货池商品数据
func GetPickGoodsByPickIds(db *gorm.DB, pickIds []int) (err error, list []PickGoods) {
	err = db.Model(&PickGoods{}).Where("pick_id in (?)", pickIds).Find(&list).Error

	return
}

// 根据拣货id查拣货池商品数据
func GetPickGoodsByOrderGoodsIds(db *gorm.DB, orderGoodsId []int) (err error, list []PickGoods) {
	err = db.Model(&PickGoods{}).Where("order_goods_id in (?)", orderGoodsId).Find(&list).Error

	return
}

// 根据拣货id查拣货池商品数据
func GetPickGoodsByPickIdAndSku(db *gorm.DB, pickId int, skus []string) (err error, list []PickGoods) {
	err = db.Model(&PickGoods{}).Where("pick_id = ? and sku in (?)", pickId, skus).Find(&list).Error

	return
}

// 根据批次ID查询拣货商品数据
func GetPickGoodsByBatchId(db *gorm.DB, batchId int) (err error, list []PickGoods) {
	err = db.Model(&PickGoods{}).Where("batch_id = ?", batchId).Find(&list).Error

	return
}

// 根据订单商品表订单编号查询拣货表数据
func GetPickGoodsJoinPickByNumbers(db *gorm.DB, numbers []string) (err error, list []PickAndGoods) {
	err = db.Table("t_pick_goods pg").
		Select("p.id as pick_id,p.status,pg.number,p.pick_user").
		Joins("left join t_pick p on pg.pick_id = p.id").
		Where("number in (?)", numbers).
		Find(&list).
		Error

	return
}

// 根据拣货ID查询拣货池商品数据并获取相关订单数据
func GetPickGoodsJoinOrderByPickId(db *gorm.DB, pickId int) (err error, list []PickGoodsJoinOrder) {
	err = db.Table("t_pick_goods pg").
		Select("shop_code,goods_name,goods_spe,unit,sku,pg.need_num,pg.review_num").
		Where("pg.pick_id = ?", pickId).
		Joins("left join t_pick_order po on po.number = pg.number").
		Scan(&list).
		Error

	return
}

// 查询拣货池商品订单是否有已拣的
func GetFirstPickGoodsByNumbers(db *gorm.DB, numbers []string) (err error, exist bool) {

	var pickGoods PickGoods

	//
	err = db.Model(&PickGoods{}).Where("number in (?) and complete_num >= 0", numbers).First(&pickGoods).Error

	if err != nil {
		//未查到
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false
		}

		//有错误
		return
	}

	//查询到已有

	return nil, true
}

// 统计集中拣货相关数量
func CountPickByBatch(db *gorm.DB, batchIds []int) (err error, countCentralizedPick []CountPickNums) {

	err = db.Model(&PickGoods{}).
		Select("batch_id,sum(need_num) as sum_need_num,sum(complete_num) as sum_complete_num,count(need_num) as count_need_num,count(complete_num) as count_complete_num").
		Where("batch_id in (?)", batchIds).
		Group("batch_id").
		Find(&countCentralizedPick).Error

	return
}

type CountBatchNums struct {
	BatchId    int `json:"batch_id"`     //批次ID
	ShopNum    int `json:"shop_num"`     //门店数
	OrderNum   int `json:"order_num"`    //订单数
	GoodsNum   int `json:"goods_num"`    //商品数
	PrePickNum int `json:"pre_pick_num"` //预拣单数
	PickNum    int `json:"pick_num"`     //拣货单数
	ReviewNum  int `json:"review_num"`   //复核单数
	Count      int `json:"count"`
	Status     int `json:"status"`
}

// 统计批次 门店数量、预拣单数等
func CountBatchNumsByBatchIds(db *gorm.DB, batchIds []int) (err error, mp map[int]CountBatchNums) {
	var (
		count        []CountBatchNums
		pickCount    []CountBatchNums
		prePickCount []CountBatchNums
	)

	mp = make(map[int]CountBatchNums, 0)

	//预拣池商品表，统计各个批次 店铺数，订单数，商品数
	err = db.Model(&PrePickGoods{}).
		Select("batch_id,count(distinct(shop_id)) as shop_num,count(distinct(number)) as order_num,sum(need_num) as goods_num").
		Where("batch_id in (?)", batchIds).
		Group("batch_id").
		Find(&count).
		Error

	if err != nil {
		return
	}

	for _, n := range count {
		mp[n.BatchId] = n
	}

	//统计预拣池任务数量
	err = db.Model(&PrePick{}).
		Select("batch_id,count(1) as pre_pick_num").
		Where("batch_id in (?) and status = ?", batchIds, PrePickStatusUnhandled).
		Group("batch_id").
		Find(&pickCount).Error

	if err != nil {
		return
	}

	for _, pc := range pickCount {
		nums, ok := mp[pc.BatchId]

		if !ok {
			nums = pc
			continue
		}
		//替换 预拣单数
		nums.PrePickNum = pc.PrePickNum

		mp[pc.BatchId] = nums
	}

	//统计拣货池 待拣货任务数 待复核任务数
	err = db.Model(&Pick{}).
		Select("batch_id,status,count(1) as count").
		Where("batch_id in (?) and status in (?)", batchIds, []int{ToBePickedStatus, ToBeReviewedStatus}).
		Group("batch_id,status").
		Find(&prePickCount).
		Error

	if err != nil {
		return
	}

	for _, ppc := range prePickCount {
		nums, ok := mp[ppc.BatchId]

		if !ok {
			nums = ppc
			continue
		}

		//替换 拣货单数 复核单数
		switch ppc.Status {
		case ToBePickedStatus: //拣货单数
			nums.PickNum = ppc.Count
			break
		case ToBeReviewedStatus: //复核单数
			nums.ReviewNum = ppc.Count
			break
		}

		mp[ppc.BatchId] = nums
	}

	for _, id := range batchIds {
		//新批次时 拣货池 复核池为空 对应的批次没有数据 给默认值
		_, ok := mp[id]

		if !ok {
			mp[id] = CountBatchNums{}
		}
	}

	return
}

type CountPickPoolNums struct {
	PickId      int `json:"pick_id"`      //拣货id
	ShopNum     int `json:"shop_num"`     //门店数
	OrderNum    int `json:"order_num"`    //订单数
	NeedNum     int `json:"need_num"`     //商品数
	CompleteNum int `json:"complete_num"` //已拣数量
	ReviewNum   int `json:"review_num"`   //复核数量
}

// 统计拣货 门店、订单、需拣数量
func CountPickPoolNumsByPickIds(db *gorm.DB, pickIds []int, query string) (err error, mp map[int]CountPickPoolNums) {
	var count []CountPickPoolNums

	mp = make(map[int]CountPickPoolNums, 0)

	err = db.Model(&PickGoods{}).
		Select(query).
		Where("pick_id in (?)", pickIds).
		Group("pick_id").
		Find(&count).
		Error

	if err != nil {
		return
	}

	for _, n := range count {
		mp[n.PickId] = n
	}

	return
}

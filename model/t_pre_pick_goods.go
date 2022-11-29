package model

import "gorm.io/gorm"

// 预拣货商品明细
type PrePickGoods struct {
	Base
	WarehouseId      int    `gorm:"type:int(11);comment:仓库"`
	BatchId          int    `gorm:"type:int(11) unsigned;index;comment:批次表id"`
	OrderGoodsId     int    `gorm:"type:int(11) unsigned;comment:订单商品表ID"` //t_pick_order_goods 表 id
	Number           string `gorm:"type:varchar(32);comment:订单编号"`
	ShopId           int    `gorm:"type:int(11);comment:店铺id"`
	PrePickId        int    `gorm:"type:int(11) unsigned;index;comment:预拣货表id"` //index
	DistributionType int    `gorm:"type:tinyint unsigned;comment:配送方式:1:公司配送,2:用户自提,3:三方物流,4:快递配送,5:首批物料|设备单"`
	Sku              string `gorm:"type:varchar(64);comment:sku"`
	GoodsName        string `gorm:"type:varchar(64);comment:商品名称"`
	GoodsType        string `gorm:"type:varchar(64);comment:商品类型"`
	GoodsSpe         string `gorm:"type:varchar(128);comment:商品规格"`
	Unit             string `gorm:"type:varchar(64);comment:单位"`
	Shelves          string `gorm:"type:varchar(64);comment:货架"`
	DiscountPrice    int    `gorm:"comment:折扣价"`
	NeedNum          int    `gorm:"type:int;not null;comment:需拣数量"`
	CloseNum         int    `gorm:"type:int;not null;comment:关闭数量"`
	OutCount         int    `gorm:"type:int;comment:出库数量"`
	NeedOutNum       int    `gorm:"type:int;comment:需出库数量"`
	Status           int    `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:已进入拣货池,2:关闭"`
	Typ              int    `gorm:"type:tinyint;default:1;comment:批次类型:1:常规批次,2:快递批次"`
}

type PrePickGoodsJoinPrePick struct {
	PrePickGoodsId   int    `json:"pre_pick_goods_id"`
	WarehouseId      int    `json:"warehouse_id"`
	BatchId          int    `json:"batch_id"`
	OrderGoodsId     int    `json:"order_goods_id"` //t_pick_order_goods 表 id
	Number           string `json:"number"`
	ShopId           int    `json:"shop_id"`
	PrePickId        int    `json:"pre_pick_id"`       //预拣货表id
	DistributionType int    `json:"distribution_type"` //配送方式:1:公司配送,2:用户自提,3:三方物流,4:快递配送,5:首批物料|设备单
	Sku              string `json:"sku"`
	GoodsName        string `json:"goods_name"`
	GoodsType        string `json:"goods_type"`
	GoodsSpe         string `json:"goods_spe"`
	Unit             string `json:"unit"`
	Shelves          string `json:"shelves"`
	DiscountPrice    int    `json:"discount_price"`
	NeedNum          int    `json:"need_num"`
	CloseNum         int    `json:"close_num"`
	OutCount         int    `json:"out_count"`
	NeedOutNum       int    `json:"need_out_num"`
	Status           int    `json:"status"`
	ShopCode         string `gorm:"type:varchar(255);not null;comment:店铺编号"` //pre_pick
	ShopName         string `gorm:"type:varchar(64);not null;comment:店铺名称"`
	Line             string `gorm:"type:varchar(255);not null;comment:线路"`
	PrePickStatus    int    `gorm:"type:tinyint;default:0;comment:状态:0:未处理,1:已进入拣货池,2:关闭"`
}

type PrePickGoodsJoinPrePickRemark struct {
	Number      string `gorm:"type:varchar(32);comment:订单编号"`
	Sku         string `gorm:"type:varchar(64);comment:sku"`
	NeedNum     int    `gorm:"type:int;not null;comment:需拣数量"`
	GoodsRemark string `gorm:"type:varchar(255);comment:商品备注"`
	Unit        string `json:"unit"`
}

const (
	PrePickGoodsStatusUnhandled  = iota //未处理
	PrePickGoodsStatusProcessing        //处理中(已进入拣货池)
	PrePickGoodsStatusClose             //关闭
)

func PrePickGoodsBatchSave(db *gorm.DB, list []PrePickGoods) (err error) {
	result := db.Model(&PrePickGoods{}).Save(&list)

	return result.Error
}

func UpdatePrePickGoodsByIds(db *gorm.DB, ids []int, mp map[string]interface{}) error {
	result := db.Model(&PrePickGoods{}).
		Where("id in (?)", ids).
		Updates(mp)

	return result.Error
}

func UpdatePrePickGoodsStatusByIds(db *gorm.DB, ids []int, status int) error {
	result := db.Model(&PrePickGoods{}).
		Where("id in (?)", ids).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdatePrePickGoodsByPrePickIds(db *gorm.DB, prePickIds []int, mp map[string]interface{}) error {
	result := db.Model(&PrePickGoods{}).
		Where("pre_pick_id in (?)", prePickIds).
		Updates(mp)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetPrePickGoodsJoinPrePickListByNumber(db *gorm.DB, numbers []string) (err error, list []PrePickGoodsJoinPrePick) {
	result := db.Table("t_pre_pick_goods pg").
		Select("pg.id as pre_pick_goods_id,pg.*,pp.*").
		Joins("left join t_pre_pick pp on pg.pre_pick_id = pp.id").
		Where("pg.number in (?)", numbers).
		Find(&list)

	if result.Error != nil {
		return result.Error, list
	}

	return result.Error, list
}

func GetPrePickGoodsJoinPrePickListByBatchId(db *gorm.DB, batchId int) (err error, list []PrePickGoodsJoinPrePick) {
	result := db.Table("t_pre_pick_goods pg").
		Select("pg.id as pre_pick_goods_id,pg.*,pp.*").
		Joins("left join t_pre_pick pp on pg.pre_pick_id = pp.id").
		Where("pg.batch_id = ?", batchId).
		Find(&list)

	if result.Error != nil {
		return result.Error, list
	}

	return result.Error, list
}

// 按分类或商品获取未进入拣货池的商品数据
func GetPrePickGoodsByTypeParam(db *gorm.DB, ids []int, formType int, typeParam []string) (err error, prePickGoods []PrePickGoods) {

	local := db.Model(&PrePickGoods{}).Where("pre_pick_id in (?) and `status` = 0", ids)

	//默认全单
	if formType == 2 { //按分类
		local.Where("goods_type in (?)", typeParam)
	} else if formType == 3 { //按商品
		local.Where("sku in (?)", typeParam)
	}

	result := local.Find(&prePickGoods)

	return result.Error, prePickGoods
}

func GetPrePickGoodsList(db *gorm.DB, cond *PrePickGoods) (err error, list []PrePickGoods) {
	result := db.Model(&PrePickGoods{}).Where(cond).Find(&list)

	if result.Error != nil {
		return result.Error, nil
	}

	return nil, list
}

func GetPrePickGoodsAndRemark(db *gorm.DB, batchId int, sku string) (err error, list []PrePickGoodsJoinPrePickRemark) {
	result := db.Table("t_pre_pick_goods pg").
		Select("pg.number,pg.sku,pg.need_num,pg.unit,pr.goods_remark").
		Joins("left join t_pre_pick_remark pr on pg.pre_pick_id = pr.pre_pick_id and pg.order_goods_id = pr.order_goods_id").
		Where("pg.batch_id = ? and pg.sku = ?", batchId, sku).
		Find(&list)

	return result.Error, list
}

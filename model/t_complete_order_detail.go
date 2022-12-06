package model

import "gorm.io/gorm"

// 完成订单明细表
type CompleteOrderDetail struct {
	Base
	Number          string   `gorm:"type:varchar(64);index;comment:订单编号"`
	GoodsName       string   `gorm:"type:varchar(64);comment:商品名称"`
	Sku             string   `gorm:"type:varchar(64);comment:sku"`
	GoodsSpe        string   `gorm:"type:varchar(128);comment:商品规格"`
	GoodsType       string   `gorm:"type:varchar(64);comment:商品类型"`
	Shelves         string   `gorm:"type:varchar(64);comment:货架"`
	PayCount        int      `gorm:"comment:下单数量"`
	CloseCount      int      `gorm:"type:int;comment:关闭数量"`
	ReviewCount     int      `gorm:"type:int;comment:出库数量"`
	GoodsRemark     string   `gorm:"type:varchar(255);comment:商品备注"`
	DeliveryOrderNo GormList `gorm:"type:varchar(255);comment:出库单号"`
}

func CompleteOrderDetailBatchSave(db *gorm.DB, list *[]CompleteOrderDetail) error {
	result := db.Model(&CompleteOrderDetail{}).Save(list)

	return result.Error
}

func GetCompleteOrderDetailBySku(db *gorm.DB, sku string) (err error, list []CompleteOrderDetail) {
	result := db.Model(&CompleteOrderDetail{}).Where("sku = ?", sku).Find(&list)

	return result.Error, list
}

type CompleteOrderNums struct {
	Number         string `json:"number"`
	SumPayCount    int    `json:"sum_pay_count"`
	SumCloseCount  int    `json:"sum_close_count"`
	SumReviewCount int    `json:"sum_review_count"`
}

// 根据订单号统计订单 下单数量 关闭数量 出库数量
func CountCompleteOrderNumsByNumber(db *gorm.DB, numbers []string) (err error, mp map[string]CompleteOrderNums) {
	mp = make(map[string]CompleteOrderNums, 0)

	var completeOrderNums []CompleteOrderNums

	err = db.Model(&CompleteOrderDetail{}).
		Select("number,sum(pay_count) as sum_pay_count,sum(close_count) as sum_close_count,sum(review_count) as sum_review_count").
		Where("number in (?)", numbers).
		Find(&completeOrderNums).
		Error

	if err != nil {
		return
	}

	//填充mp数据
	for _, nums := range completeOrderNums {
		mp[nums.Number] = nums
	}

	//传递过来的订单号如果没有查到，填充0
	for _, ns := range numbers {
		nums, mpOk := mp[ns]
		if !mpOk {
			nums = CompleteOrderNums{
				Number:         ns,
				SumPayCount:    0,
				SumCloseCount:  0,
				SumReviewCount: 0,
			}

			mp[ns] = nums
		}
	}

	return
}

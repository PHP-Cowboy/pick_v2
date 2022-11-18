package dao

import (
	"gorm.io/gorm"
	"pick_v2/model"
)

// 变更订单类型&&商品状态
func UpdateOrderAndGoods(db *gorm.DB, order []model.Order, orderGoods []model.OrderGoods) error {

	err := model.UpdateOrder(db, order)

	if err != nil {
		return err
	}

	err = model.UpdateOrderGoodsStatus(db, orderGoods, []string{"status"})

	if err != nil {
		return err
	}

	return nil
}

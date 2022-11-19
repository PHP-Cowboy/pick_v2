package dao

import (
	"gorm.io/gorm"
	"pick_v2/model"
)

// 变更订单类型&&商品状态
func UpdateOrderAndGoods(db *gorm.DB, orderIds []int, orderGoodsIds []int) error {

	err := model.UpdateOrderByIds(db, orderIds, map[string]interface{}{"order_type": model.PickingOrderType})

	if err != nil {
		return err
	}

	err = model.UpdateOrderGoodsByIds(db, orderGoodsIds, map[string]interface{}{"status": model.OrderGoodsProcessingStatus})

	if err != nil {
		return err
	}

	return nil
}

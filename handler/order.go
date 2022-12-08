package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
)

// 关闭订单
func CloseOrder(c *gin.Context) {
	var form req.CloseOrderForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		pickOrder []model.PickOrder
	)

	db := global.DB

	//根据ID倒序，查最新的记录，避免欠货单更新错误
	result := db.Model(&model.PickOrder{}).Where("order_id = ?", form.OrderId).Find(&pickOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, order := range pickOrder {
		if order.OrderType != 1 { //只要有 != 1 的 即 不允许更新
			xsq_net.ErrorJSON(c, errors.New("当前订单不允许更新"))
			return
		}
	}

	tx := db.Begin()

	// 拣货单查到数据 且是新订单
	if len(pickOrder) > 0 {
		result = tx.Model(&model.PickOrder{}).Where("order_id = ?", form.OrderId).Update("order_type", 3)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	}

	var order model.Order

	//查询订单表
	result = db.Model(&model.Order{}).First(&order, form.OrderId)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//更新订单表
	result = tx.Model(&model.Order{}).Where("id = ?", form.OrderId).Update("order_type", 4)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
	return
}

// 关闭商品
func CloseOrderGoods(c *gin.Context) {
	var form req.CloseOrderGoodsForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		order      model.Order
		orderGoods model.OrderGoods
	)

	db := global.DB

	err, outboundGoods := model.GetOutboundGoodsFirstByOrderGoodsIdSortByTaskId(db, form.GoodsId)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	tx := db.Begin()
	//查到数据
	if outboundGoods.Status != model.OutboundGoodsStatusUnhandled {
		xsq_net.ErrorJSON(c, errors.New("当前商品不允许关闭"))
		return
	}

	if outboundGoods.LackCount < form.CloseNum {
		xsq_net.ErrorJSON(c, errors.New("关闭数量大于欠货数量"))
		return
	}

	// 增加关闭数量 && 扣减发货数量
	result := tx.Model(&model.OutboundGoods{}).
		Where("task_id = ? and number = ? and sku = ?", outboundGoods.TaskId, outboundGoods.Number, outboundGoods.Sku).
		Updates(map[string]interface{}{
			"close_count": gorm.Expr("close_count + ?", form.CloseNum),
			"lack_count":  gorm.Expr("lack_count - ?", form.CloseNum),
		})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//result = tx.Model(&model.OutboundOrder{}).
	//	Where("task_id = ? and number = ?", outboundGoods.TaskId, outboundGoods.Number).
	//	Updates(map[string]interface{}{
	//		"close_count":   gorm.Expr("close_count + ?", form.CloseNum),
	//		"shipments_num": gorm.Expr("shipments_num - ?", form.CloseNum),
	//	})
	//
	//if result.Error != nil {
	//	tx.Rollback()
	//	xsq_net.ErrorJSON(c, result.Error)
	//	return
	//}

	result = db.Model(&model.OrderGoods{}).First(&orderGoods, form.GoodsId)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			//mq中订单没到拣货系统中，这个其实是不合理的状态
			//没查到，先什么都不做，让订货关商品通过，后续再拉单
			xsq_net.Success(c)
			return
		}

		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.Order{}).Where("number = ?", orderGoods.Number).First(&order)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			//mq中订单没到拣货系统中，这个其实是不合理的状态
			//没查到，先什么都不做，让订货关商品通过，后续再拉单
			xsq_net.Success(c)
			return
		}
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//2拣货中 4已关闭
	if order.OrderType == 2 || order.OrderType == 4 {
		xsq_net.ErrorJSON(c, errors.New("当前商品不允许关闭"))
		return
	}

	// 增加关闭数量 && 扣减发货数量
	result = tx.Model(&model.OrderGoods{}).
		Where("id = ?", orderGoods.Id).
		Updates(map[string]interface{}{
			"close_count": gorm.Expr("close_count + ?", form.CloseNum),
			"lack_count":  gorm.Expr("lack_count - ?", form.CloseNum),
		})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Model(&model.Order{}).
		Where("number = ?", orderGoods.Number).
		Updates(map[string]interface{}{
			"close_num": gorm.Expr("close_num + ?", form.CloseNum),
			"un_picked": gorm.Expr("un_picked - ?", form.CloseNum),
		})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

func Test(c *gin.Context) {
	sign := middlewares.Generate()
	xsq_net.SucJson(c, gin.H{"sign": sign})
}

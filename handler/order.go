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
func CloseOrderDiscard(c *gin.Context) {
	var form req.CloseOrderForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		err           error
		order         model.Order
		outboundOrder model.OutboundOrder
	)

	db := global.DB

	err, order = model.GetOrderByPk(db, form.OrderId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//没查到订单，可能是还在消息中，未同步到拣货系统中
			xsq_net.Success(c)
			return
		}
		xsq_net.ErrorJSON(c, err)
		return
	}

	//新订单 直接关闭
	if order.OrderType != model.NewOrderType {
		//查询出库任务
		err, outboundOrder = model.GetOutboundOrderByNumberFirstSortByTaskId(db, order.Number)

		if err != nil {
			xsq_net.ErrorJSON(c, err)
			return
		}

		//出库任务中不是新订单
		if outboundOrder.OrderType != model.OutboundOrderTypeNew {
			xsq_net.ErrorJSON(c, errors.New("当前订单不允许关闭"))
			return
		}
	}

	tx := db.Begin()

	err = CloseOrderLogic(tx, form, outboundOrder)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
	return
}

// 订单关闭
func CloseOrderLogic(tx *gorm.DB, form req.CloseOrderForm, outboundOrder model.OutboundOrder) (err error) {
	//更新订单表
	err = model.UpdateOrderByIds(tx, []int{form.OrderId}, map[string]interface{}{"order_type": model.CloseOrderType})

	if err != nil {
		return
	}

	if outboundOrder.TaskId > 0 {
		err = model.UpdateOutboundOrderByTaskIdAndNumbers(tx, outboundOrder.TaskId, []string{outboundOrder.Number}, map[string]interface{}{"order_type": model.OutboundOrderTypeClose})
		return
	}

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
		err           error
		order         model.Order
		outboundGoods model.OutboundGoods
		outboundOrder model.OutboundOrder
	)

	db := global.DB

	//是否进入批次为准

	//订单数据
	err, order = model.GetOrderByPk(db, form.OrderId)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	//订单已关闭
	if order.OrderType == model.CloseOrderType {
		xsq_net.ErrorJSON(c, errors.New("订单已被关闭"))
		return
	}

	//订单是否在拣货中
	//是 在哪个任务中[要查最新的任务]，任务是否已进入批次 如果已进入批次则不允许关闭

	//订单在拣货中
	if order.OrderType == 2 {
		//获取出库任务商品列表
		err, outboundGoods = model.GetOutboundGoodsFirstByOrderGoodsIdSortByTaskId(db, form.GoodsId)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				//订单拣货中 但是任务中没查到，说明订单进入任务的没有这个品 更新订单欠货数量&&关闭数量
				goto CloseLogic
			}

			xsq_net.ErrorJSON(c, err)
			return
		}

		if outboundGoods.LackCount < form.CloseNum {
			xsq_net.ErrorJSON(c, errors.New("关闭数量大于欠货数量"))
			return
		}

		//获取出库任务订单数据
		err, outboundOrder = model.GetOutboundOrderByPk(db, outboundGoods.TaskId, outboundGoods.Number)

		if err != nil {
			xsq_net.ErrorJSON(c, err)
			return
		}

		//出库任务订单状态不是新订单，则已进入批次
		if outboundOrder.OrderType != model.OutboundOrderTypeNew {
			xsq_net.ErrorJSON(c, errors.New("订单已进入批次中，不允许关闭"))
			return
		}
	}

CloseLogic:
	OrderGoodsCloseLogic(db, form, outboundGoods)

	xsq_net.Success(c)
}

// 订单商品关闭
func OrderGoodsCloseLogic(db *gorm.DB, form req.CloseOrderGoodsForm, outboundGoods model.OutboundGoods) (err error) {
	var orderGoods model.OrderGoods

	tx := db.Begin()

	if outboundGoods.TaskId > 0 {
		// 增加关闭数量 && 扣减发货数量
		err = tx.Model(&model.OutboundGoods{}).
			Where("task_id = ? and number = ? and sku = ?", outboundGoods.TaskId, outboundGoods.Number, outboundGoods.Sku).
			Updates(map[string]interface{}{
				"close_count": gorm.Expr("close_count + ?", form.CloseNum),
				"lack_count":  gorm.Expr("lack_count - ?", form.CloseNum),
			}).Error

		if err != nil {
			tx.Rollback()
			return
		}
	}

	err = db.Model(&model.OrderGoods{}).First(&orderGoods, form.GoodsId).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			//mq中订单没到拣货系统中，这个其实是不合理的状态
			//没查到，先什么都不做，让订货关商品通过，后续再拉单
			err = nil
			return
		}

		return
	}

	// 增加关闭数量 && 扣减发货数量
	err = tx.Model(&model.OrderGoods{}).
		Where("id = ?", orderGoods.Id).
		Updates(map[string]interface{}{
			"close_count": gorm.Expr("close_count + ?", form.CloseNum),
			"lack_count":  gorm.Expr("lack_count - ?", form.CloseNum),
		}).Error

	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Model(&model.Order{}).
		Where("number = ?", orderGoods.Number).
		Updates(map[string]interface{}{
			"close_num": gorm.Expr("close_num + ?", form.CloseNum),
			"un_picked": gorm.Expr("un_picked - ?", form.CloseNum),
		}).
		Error

	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return
}

func Test(c *gin.Context) {
	sign := middlewares.Generate()
	xsq_net.SucJson(c, gin.H{"sign": sign})
}

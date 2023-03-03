package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"pick_v2/dao"
	"pick_v2/forms/req"
	"pick_v2/global"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
)

// 生成出库任务
func CreateOutboundTask(c *gin.Context) {
	var form req.CreateOutboundForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	rdsKey := c.Request.URL.Path

	err := cache.AntiRepeatedClick(rdsKey, 30)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	//执行完成后删除锁定时间
	defer func(key string) {
		_, _ = cache.Del(key)
	}(rdsKey)

	db := global.DB

	tx := db.Begin()

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, ecode.GetContextUserInfoFailed)
		return
	}

	//出库任务相关保存
	taskId, err := dao.OutboundTaskSave(tx, form, userInfo)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	//出库订单相关保存
	err = dao.OutboundOrderBatchSave(tx, form, taskId)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
	return
}

// 出库单任务列表
func OutboundTaskList(c *gin.Context) {
	var form req.OutboundTaskListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, res := dao.OutboundTaskList(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, res)
}

// 简化版任务列表
func OutboundTaskListSimple(c *gin.Context) {

	err, res := dao.OutboundTaskListSimple(global.DB)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, res)
}

// 出库任务状态数量
func OutboundTaskCount(c *gin.Context) {

	err, res := dao.OutboundTaskCount(global.DB)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, res)
}

// 出库单订单列表
func OutboundOrderList(c *gin.Context) {
	var form req.OutboundOrderListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, rsp := dao.OutboundOrderList(global.DB, form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, rsp)
}

// 出库订单数量
func OutboundOrderCount(c *gin.Context) {
	var form req.OutboundOrderCountForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, res := dao.OutboundOrderCount(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, res)
}

// 获取任务某个商品的发货数量
func GetTaskSkuNum(c *gin.Context) {

	var form req.GetTaskSkuNumForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	//获取任务某个商品的发货数量
	err, num := dao.GetTaskSkuNum(global.DB, form)
	if err != nil {
		return
	}

	xsq_net.SucJson(c, gin.H{"num": num})
}

// 出库单订单详情
func OutboundOrderDetail(c *gin.Context) {
	var form req.OutboundOrderDetailForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, rsp := dao.OutboundOrderDetail(global.DB, form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, rsp)
}

// 出库任务商品列表
func OutboundOrderGoodsList(c *gin.Context) {
	var form req.OutboundOrderGoodsListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, rsp := dao.OutboundOrderGoodsList(global.DB, form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, rsp)
}

// 结束任务
func EndOutboundTask(c *gin.Context) {
	var form req.EndOutboundTaskForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := dao.EndOutboundTask(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

// 关闭订单
func OutboundTaskCloseOrder(c *gin.Context) {
	var form req.OutboundTaskCloseOrderForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	tx := global.DB.Begin()

	err := dao.OutboundTaskCloseOrder(tx, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

// 加单
func OutboundTaskAddOrder(c *gin.Context) {
	var form req.OutboundTaskAddOrderForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	tx := global.DB.Begin()

	err := dao.OutboundTaskAddOrder(tx, form)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

// 订单出库记录
func OrderOutboundRecord(c *gin.Context) {
	var form req.OrderOutboundRecordForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, list := dao.OrderOutboundRecord(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, list)
}

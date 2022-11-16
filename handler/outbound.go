package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"pick_v2/dao"
	"pick_v2/forms/req"
	"pick_v2/global"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
)

// 生成出库任务
func CreateOutboundTask(c *gin.Context) {
	var form req.CreateOutboundForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	tx := db.Begin()

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, errors.New("获取上下文用户数据失败"))
		return
	}

	taskId, err := dao.OutboundTaskSave(tx, form, userInfo)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	err = dao.OutboundOrderBatchSave(tx, form, taskId)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
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

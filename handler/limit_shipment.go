package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"pick_v2/dao"
	"pick_v2/forms/req"
	"pick_v2/global"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
)

// 订单批量限发
func OrderLimit(c *gin.Context) {
	var form req.OrderLimitForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if err := dao.OrderLimit(global.DB, form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}
	xsq_net.Success(c)
}

// 任务批量限发
func TaskLimit(c *gin.Context) {
	var form req.TaskLimitForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	if err := dao.TaskLimit(global.DB, form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}
	xsq_net.Success(c)
}

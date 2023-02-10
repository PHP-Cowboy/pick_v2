package handler

import (
	"github.com/gin-gonic/gin"
	"pick_v2/dao"
	"pick_v2/forms/req"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
)

// 后台拣货列表
func AdminPickList(c *gin.Context) {
	var form req.AdminPickListReq

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, res := dao.AdminPickList(form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, res)
}

// 后台拣货数据详情
func AdminPickDetail(c *gin.Context) {
	var form req.AdminPickDetailReq

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, res := dao.AdminPickDetail(form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, res)
}

// 后台拣货数据详情
func BatchShopGoodsList(c *gin.Context) {
	var form req.BatchShopGoodsListReq

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, res := dao.BatchShopGoodsList(form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, res)
}

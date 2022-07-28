package handler

import (
	"github.com/gin-gonic/gin"
	"pick_v2/forms/req"
	"pick_v2/global"
	"pick_v2/model/batch"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"time"
)

//生成拣货批次
func CreateBatch(c *gin.Context) {
	var form req.CreateBatchForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		condition batch.BatchCondition
	)

	payEndTime, errPayEndTime := time.ParseInLocation(timeutil.TimeFormat, form.PayEndTime, time.Local)

	deliveryEndTime, errDeliveryEndTime := time.ParseInLocation(timeutil.TimeFormat, form.DeliveryEndTime, time.Local)

	deliveryStartTime, errDeliveryStartTime := time.ParseInLocation(timeutil.TimeFormat, form.DeliveryStartTime, time.Local)

	if errPayEndTime != nil || errDeliveryEndTime != nil || errDeliveryStartTime != nil {
		xsq_net.ErrorJSON(c, ecode.DataTransformationError)
		return
	}

	condition.WarehouseId = form.WarehouseId
	condition.PayEndTime = payEndTime
	condition.DeliveryEndTime = deliveryEndTime
	condition.DeliveryStartTime = deliveryStartTime
	condition.DeliveryMethod = form.DeliveryMethod
	condition.Line = form.Line
	condition.Goods = form.Goods

	//todo 根据条件筛选 如果查到 调用锁单接口 同时存储相关数据

	db := global.DB

	result := db.Save(&condition)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
	}
}

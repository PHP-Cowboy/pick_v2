package handler

import (
	"github.com/gin-gonic/gin"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
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

//获取批次列表
func GetBatchList(c *gin.Context) {
	var (
		form req.GetBatchListForm
		res  rsp.GetBatchListRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var batches []batch.Batch

	db := global.DB

	result := db.Find(&batches)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	db.Scopes(model.Paginate(form.Page, form.Size)).Find(&batches)

	for _, b := range batches {
		res.List = append(res.List, &rsp.Batch{
			BatchName:         b.BatchName,
			DeliveryStartTime: b.DeliveryStartTime.Format(timeutil.TimeFormat),
			DeliveryEndTime:   b.DeliveryEndTime.Format(timeutil.TimeFormat),
			ShopNum:           b.ShopNum,
			OrderNum:          b.OrderNum,
			UserName:          b.UserName,
			Line:              b.Line,
			DeliveryMethod:    b.DeliveryMethod,
			EndTime:           b.EndTime.Format(timeutil.TimeFormat),
			Status:            b.Status,
			PickNum:           b.PickNum,
			RecheckSheetNum:   b.RecheckSheetNum,
		})
	}

	xsq_net.SucJson(c, res)
}

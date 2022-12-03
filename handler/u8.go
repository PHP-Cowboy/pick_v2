package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"pick_v2/common/constant"
	"pick_v2/dao"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
)

// u8推送日志列表
func LogList(c *gin.Context) {
	var form req.LogListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var (
		total    int64
		stockLog []model.StockLog
		res      rsp.LogListRsp
	)

	db := global.DB

	localDb := db.Model(&model.StockLog{}).Where(model.StockLog{Status: form.Status})

	if form.StartTime != "" {
		localDb = localDb.Where("create_time >= ?", form.StartTime)
	}

	if form.EndTime != "" {
		localDb = localDb.Where("create_time <= ?", form.EndTime)
	}

	result := localDb.Count(&total)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = total

	localDb.Scopes(model.Paginate(form.Page, form.Size)).Find(&stockLog)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.LogList, 0, len(stockLog))

	for _, log := range stockLog {
		list = append(list, rsp.LogList{
			Id:          log.Id,
			CreateTime:  log.CreateTime.Format(timeutil.TimeFormat),
			UpdateTime:  log.UpdateTime.Format(timeutil.TimeFormat),
			Number:      log.Number,
			BatchId:     log.BatchId,
			PickId:      log.PickId,
			Status:      log.Status,
			RequestXml:  log.RequestXml,
			ResponseXml: log.ResponseXml,
			ResponseNo:  log.ResponseNo,
			Msg:         log.Msg,
			ShopName:    log.ShopName,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// u8 批量补单
func BatchSupplement(c *gin.Context) {
	var form req.BatchSupplementForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	second, err := cache.TTL(constant.BATCH_SUPPLEMENT)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	if second > 0 {
		xsq_net.ErrorJSON(c, errors.New(fmt.Sprintf("推送u8处理中，请%v秒后再试", second.Seconds())))
		return
	}

	expire := dao.BaseNum * (len(form.Ids) + 10) //channel 没读取到数据时 等待了 BaseNum * 10 秒

	_, err = cache.Set(constant.BATCH_SUPPLEMENT, "1", expire)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	for _, id := range form.Ids {
		dao.YongYouProducer(id)
	}

	xsq_net.Success(c)
}

// 推送u8拣货详情
func LogDetail(c *gin.Context) {
	var form req.LogDetailForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		res       rsp.LogDetailRsp
		pickGoods []model.PickGoods
		order     model.Order
		pick      model.Pick
	)

	db := global.DB

	result := db.Model(&model.Order{}).Where("number = ?", form.Number).First(&order)

	if result.Error != nil {

		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.First(&pick, form.PickId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.ShopName = order.ShopName
	res.PayAt = order.PayAt

	res.PickUser = pick.PickUser
	res.TakeOrdersTime = pick.TakeOrdersTime
	res.ReviewUser = pick.ReviewUser
	res.ReviewTime = pick.ReviewTime

	result = db.Model(&model.PickGoods{}).Where(model.PickGoods{PickId: form.PickId, Number: form.Number}).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.LogDetail, 0, len(pickGoods))

	for _, pg := range pickGoods {
		list = append(list, rsp.LogDetail{
			Id:               pg.Id,
			UpdateTime:       pg.UpdateTime.Format(timeutil.TimeFormat),
			PickId:           pg.PickId,
			BatchId:          pg.BatchId,
			PrePickGoodsId:   pg.PrePickGoodsId,
			OrderGoodsId:     pg.OrderGoodsId,
			Number:           pg.Number,
			ShopId:           pg.ShopId,
			DistributionType: pg.DistributionType,
			Sku:              pg.Sku,
			GoodsName:        pg.GoodsName,
			GoodsType:        pg.GoodsType,
			GoodsSpe:         pg.GoodsSpe,
			Shelves:          pg.Shelves,
			DiscountPrice:    float64(pg.DiscountPrice) / 100,
			NeedNum:          pg.NeedNum,
			CompleteNum:      pg.CompleteNum,
			ReviewNum:        pg.ReviewNum,
			Unit:             pg.Unit,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

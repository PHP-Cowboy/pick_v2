package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/model/batch"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"time"
)

// 接单拣货
func ReceivingOrders(c *gin.Context) {
	var (
		res     rsp.ReceivingOrdersRsp
		pick    batch.Pick
		batches batch.Batch
	)

	db := global.DB

	claims, ok := c.Get("claims")

	if !ok {
		xsq_net.ErrorJSON(c, errors.New("claims 获取失败"))
		return
	}

	userInfo := claims.(*middlewares.CustomClaims)

	result := db.Model(&batch.Pick{}).Where("pick_user = ? and status = 0", userInfo.Name).First(&pick)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if result.RowsAffected > 0 {
		//todo 如果是首批物料则需要到后台拣货
		res.Id = pick.Id
		xsq_net.SucJson(c, res)
		return
	}

	result = db.Model(&batch.Pick{}).Where("pick_user = '' and status = 0").First(&pick)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if result.RowsAffected > 0 {
		result = db.First(&batches, pick.Id)
		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		if batches.Status != 0 {
			xsq_net.ErrorJSON(c, errors.New("批次不在进行中"))
			return
		}

		res.Id = pick.Id

		now := time.Now()
		//更新拣货员为当前用户
		result = db.Model(batch.Pick{}).Where("id = ?", pick.Id).Updates(map[string]interface{}{"pick_user": userInfo.Name, "take_orders_time": &now})

		if result.Error != nil {
			xsq_net.ErrorJSON(c, ecode.DataSaveError)
			return
		}

		xsq_net.SucJson(c, res)
		return
	} else {
		xsq_net.ErrorJSON(c, errors.New("暂无拣货单"))
		return
	}
}

// 完成拣货
func CompletePick(c *gin.Context) {
	var form req.CompletePickForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	// todo 这里需要做并发处理
	var (
		pick      batch.Pick
		pickGoods []batch.PickGoods
	)

	db := global.DB

	result := db.First(&pick, form.PickId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			xsq_net.ErrorJSON(c, ecode.DataNotExist)
			return
		}
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if pick.Status == 1 {
		xsq_net.ErrorJSON(c, ecode.OrderPickingCompleted)
		return
	}

	claims, ok := c.Get("claims")

	if !ok {
		xsq_net.ErrorJSON(c, errors.New("claims 获取失败"))
		return
	}

	userInfo := claims.(*middlewares.CustomClaims)

	if pick.PickUser != userInfo.Name {
		xsq_net.ErrorJSON(c, ecode.DataNotExist)
		return
	}

	pickGoodsMap := make(map[int]int, len(form.CompletePick))

	for _, cp := range form.CompletePick {
		pickGoodsMap[cp.Id] = cp.CompleteNum
	}

	result = db.Where("pick_id = ?", form.PickId).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for k, pg := range pickGoods {
		num, ok := pickGoodsMap[pg.Id]

		if !ok {
			continue
		}

		pickGoods[k].CompleteNum = num
	}

	tx := db.Begin()

	status := 1

	if form.Type == 2 { //无需拣货
		status = 2
	} else {
		//正常拣货的更新拣货数量，无需拣货不更新
		result = tx.Save(&pickGoods)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	}

	//更新主表
	result = tx.Model(&batch.Pick{}).Where("id = ?", pick.Id).Update("status", status)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

// 剩余数量 放拣货池那边
func RemainingQuantity(c *gin.Context) {
	var form req.ReceivingOrdersForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}
}

// 拣货记录
func PickingRecord(c *gin.Context) {
	var (
		form      req.PickingRecordForm
		res       rsp.PickingRecordRsp
		pick      []batch.Pick
		pickGoods []batch.PickGoods
		pickIds   []int
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	claims, ok := c.Get("claims")

	if !ok {
		xsq_net.ErrorJSON(c, errors.New("claims 获取失败"))
		return
	}

	userInfo := claims.(*middlewares.CustomClaims)

	//两天前日期
	towDaysAgo := timeutil.GetTimeAroundByDays(-2)

	db := global.DB

	local := db.Where("pick_user = ? and take_orders_time >= ?", userInfo.Name, towDaysAgo)

	if form.Status == 0 {
		local.Where("status = ?", form.Status)
	} else {
		local.Where("status != 0")
	}

	result := local.Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	result = local.Scopes(model.Paginate(form.Page, form.Size)).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, p := range pick {
		pickIds = append(pickIds, p.Id)
	}

	result = db.Where("pick_id in (?) ", pickIds).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	type Goods struct {
		CompleteNum      int
		DistributionType int
	}

	pickGoodsMp := make(map[int]Goods, 0)

	for _, pg := range pickGoods {
		_, ok := pickGoodsMp[pg.PickId]

		g := Goods{
			CompleteNum:      pg.CompleteNum,
			DistributionType: pg.DistributionType,
		}

		if !ok {
			pickGoodsMp[pg.PickId] = g
		} else {
			g.CompleteNum += pickGoodsMp[pg.PickId].CompleteNum
			pickGoodsMp[pg.PickId] = g
		}
	}

	list := make([]rsp.PickingRecord, 0)

	for _, p := range pick {
		pgMp, ok := pickGoodsMp[p.Id]

		outNum := 0
		distributionType := 0

		if ok {
			outNum = pgMp.CompleteNum
			distributionType = pgMp.DistributionType
		}

		reviewStatus := "未复核"
		if p.ReviewTime != nil {
			reviewStatus = "已复核"
		}

		list = append(list, rsp.PickingRecord{
			Id:               p.Id,
			TaskName:         p.TaskName,
			ShopCode:         p.ShopCode,
			ShopNum:          p.ShopNum,
			OrderNum:         p.OrderNum,
			NeedNum:          p.NeedNum,
			TakeOrdersTime:   p.TakeOrdersTime.Format(timeutil.TimeFormat),
			ReviewUser:       p.ReviewUser,
			OutNum:           outNum,
			ReviewStatus:     reviewStatus,
			DistributionType: distributionType,
			IsRemark:         false,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// 拣货记录明细
func PickingRecordDetail(c *gin.Context) {
	var (
		form req.PickingRecordDetailForm
		res  rsp.PickingRecordDetailRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		pick       batch.Pick
		pickGoods  []batch.PickGoods
		pickRemark []batch.PickRemark
	)

	db := global.DB

	result := db.Where("id = ?", form.PickId).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.TaskName = pick.TaskName
	res.OutTotal = 0
	res.UnselectedTotal = 0
	res.PickUser = pick.PickUser
	res.TakeOrdersTime = pick.TakeOrdersTime.Format(timeutil.TimeFormat)
	res.ReviewUser = pick.ReviewUser

	var reviewTime string

	if pick.ReviewTime != nil {
		reviewTime = pick.ReviewTime.Format(timeutil.TimeFormat)
	}
	res.ReviewTime = reviewTime

	result = db.Where("pick_id = ?", form.PickId).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	goodsMap := make(map[string][]rsp.PickGoods, 0)

	needTotal := 0
	completeTotal := 0
	for _, goods := range pickGoods {
		completeTotal += goods.CompleteNum
		needTotal += goods.NeedNum
		goodsMap[goods.GoodsType] = append(goodsMap[goods.GoodsType], rsp.PickGoods{
			Id:          goods.Id,
			GoodsName:   goods.GoodsName,
			GoodsSpe:    goods.GoodsSpe,
			Shelves:     goods.Shelves,
			NeedNum:     goods.NeedNum,
			CompleteNum: goods.CompleteNum,
			Unit:        goods.Unit,
		})
	}

	res.OutTotal = completeTotal
	res.UnselectedTotal = needTotal - completeTotal

	res.Goods = goodsMap

	result = db.Where("pick_id = ?", form.PickId).Find(&pickRemark)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := []rsp.PickRemark{}
	for _, remark := range pickRemark {
		list = append(list, rsp.PickRemark{
			Number:      remark.Number,
			OrderRemark: remark.OrderRemark,
			GoodsRemark: remark.GoodsRemark,
		})
	}

	res.RemarkList = list

	xsq_net.SucJson(c, res)
}

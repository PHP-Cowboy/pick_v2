package handler

import (
	"errors"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"pick_v2/dao"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
)

// 接单拣货
func ReceivingOrders(c *gin.Context) {

	var (
		form req.ReceivingOrdersForm
	)

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, ecode.GetContextUserInfoFailed)
		return
	}

	form.UserName = userInfo.Name

	err, res := dao.ReceivingOrders(global.DB, form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, res)

}

// 集中拣货接单
func ConcentratedPickReceivingOrders(c *gin.Context) {
	var form req.ConcentratedPickReceivingOrdersForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, ecode.GetContextUserInfoFailed)
		return
	}

	err, res := dao.ConcentratedPickReceivingOrders(global.DB, form, userInfo.Name)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, res)
}

// 完成拣货
func CompletePick(c *gin.Context) {
	var form req.CompletePickForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, ecode.GetContextUserInfoFailed)
		return
	}

	form.UserName = userInfo.Name

	err := dao.CompletePick(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

// 完成集中拣货
func CompleteConcentratedPick(c *gin.Context) {
	var form req.CompleteConcentratedPickForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := dao.CompleteConcentratedPick(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

// 剩余数量 放拣货池那边
func RemainingQuantity(c *gin.Context) {

	var (
		form     req.RemainingQuantityForm
		count    int64
		batches  []model.Batch
		batchIds []int
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	db := global.DB

	//批次进行中或暂停的单数量
	result := db.Where("typ = ? and ( status = 0 or status = 2 )", form.Typ).Find(&batches)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, ecode.GetContextUserInfoFailed)
		return
	}

	for _, b := range batches {
		batchIds = append(batchIds, b.Id)
	}

	if len(batchIds) > 0 {
		result = db.Model(&model.Pick{}).
			Where(
				"batch_id in (?) and status = ? and typ = ? and (pick_user = '' or pick_user = ?)",
				batchIds,
				model.ToBePickedStatus,
				form.Typ,
				userInfo.Name,
			).
			Count(&count)

		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	}

	xsq_net.SucJson(c, gin.H{"count": count})
}

// 集中拣货剩余数量
func CentralizedPickRemainingQuantity(c *gin.Context) {

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, ecode.GetContextUserInfoFailed)
		return
	}

	//集中拣货剩余数量统计
	err, count := dao.CentralizedPickRemainingQuantity(global.DB, userInfo.Name)

	if err != nil {
		return
	}

	xsq_net.SucJson(c, gin.H{"count": count})
}

// 拣货记录
func PickingRecord(c *gin.Context) {
	var (
		form      req.PickingRecordForm
		res       rsp.PickingRecordRsp
		pick      []model.Pick
		pickGoods []model.PickGoods
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

	if form.Status != nil && *form.Status == 0 {
		local.Where("status = ?", *form.Status)
	} else {
		local.Where("status != 0")
	}

	result := local.Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	result = local.Scopes(model.Paginate(form.Page, form.Size)).Order("take_orders_time desc").Find(&pick)

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
		_, pgMpOk := pickGoodsMp[pg.PickId]

		g := Goods{
			CompleteNum:      pg.CompleteNum,
			DistributionType: pg.DistributionType,
		}

		if !pgMpOk {
			pickGoodsMp[pg.PickId] = g
		} else {
			g.CompleteNum += pickGoodsMp[pg.PickId].CompleteNum
			pickGoodsMp[pg.PickId] = g
		}
	}

	list := make([]rsp.PickingRecord, 0)

	for _, p := range pick {
		pgMp, pgMpOk := pickGoodsMp[p.Id]

		outNum := 0
		distributionType := 0

		if pgMpOk {
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
			TakeOrdersTime:   p.TakeOrdersTime,
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
		pick       model.Pick
		pickGoods  []model.PickGoods
		pickRemark []model.PickRemark
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

	res.TakeOrdersTime = pick.TakeOrdersTime
	res.ReviewUser = pick.ReviewUser
	res.ReviewTime = pick.ReviewTime

	result = db.Where("pick_id = ?", form.PickId).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	pickGoodsSkuMp := make(map[string]rsp.MergePickGoods, 0)
	//相同sku合并处理
	for _, goods := range pickGoods {
		val, ok := pickGoodsSkuMp[goods.Sku]

		paramsId := rsp.ParamsId{
			PickGoodsId:  goods.Id,
			OrderGoodsId: goods.OrderGoodsId,
		}

		if !ok {

			pickGoodsSkuMp[goods.Sku] = rsp.MergePickGoods{
				Id:          goods.Id,
				Sku:         goods.Sku,
				GoodsName:   goods.GoodsName,
				GoodsType:   goods.GoodsType,
				GoodsSpe:    goods.GoodsSpe,
				Shelves:     goods.Shelves,
				NeedNum:     goods.NeedNum,
				CompleteNum: goods.CompleteNum,
				ReviewNum:   goods.ReviewNum,
				Unit:        goods.Unit,
				ParamsId:    []rsp.ParamsId{paramsId},
			}
		} else {
			val.NeedNum += goods.NeedNum
			val.CompleteNum += goods.CompleteNum
			val.ParamsId = append(val.ParamsId, paramsId)
			pickGoodsSkuMp[goods.Sku] = val
		}
	}

	goodsMap := make(map[string][]rsp.MergePickGoods, 0)

	needTotal := 0
	completeTotal := 0
	for _, goods := range pickGoodsSkuMp {
		completeTotal += goods.CompleteNum
		needTotal += goods.NeedNum

		goodsMap[goods.GoodsType] = append(goodsMap[goods.GoodsType], rsp.MergePickGoods{
			Id:          goods.Id,
			Sku:         goods.Sku,
			GoodsName:   goods.GoodsName,
			GoodsType:   goods.GoodsType,
			GoodsSpe:    goods.GoodsSpe,
			Shelves:     goods.Shelves,
			NeedNum:     goods.NeedNum,
			CompleteNum: goods.CompleteNum,
			ReviewNum:   goods.ReviewNum,
			Unit:        goods.Unit,
			ParamsId:    goods.ParamsId,
		})
	}

	res.ShopCode = pick.ShopCode
	res.OutTotal = completeTotal
	res.UnselectedTotal = needTotal - completeTotal

	//按货架号排序
	for s, goods := range goodsMap {

		ret := rsp.MyMergePickGoods(goods)

		sort.Sort(ret)

		goodsMap[s] = ret
	}

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

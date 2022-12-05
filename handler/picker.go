package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"pick_v2/dao"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"sort"
	"time"
)

func getPick(pick []model.Pick) (res rsp.ReceivingOrdersRsp, err error) {

	db := global.DB

	if len(pick) == 1 { //只查到一条
		res.Id = pick[0].Id
		res.BatchId = pick[0].BatchId
		res.Version = pick[0].Version
		res.TakeOrdersTime = pick[0].TakeOrdersTime
	} else { //查到多条
		//排序
		var (
			batchIds []int
			batchMp  = make(map[int]struct{}, 0)
			pickMp   = make(map[int][]model.Pick, 0)
		)

		//去重，构造批次id切片
		for _, b := range pick {
			//构造批次下的拣货池数据map
			//批次排序后，直接获取某个批次的全部拣货池数据。
			//然后对这部分数据排序
			pickMp[b.BatchId] = append(pickMp[b.BatchId], b)
			//已经存入了批次map的，跳过
			_, bMpOk := batchMp[b.BatchId]
			if bMpOk {
				continue
			}
			//写入批次mp
			batchMp[b.BatchId] = struct{}{}
			//存入批次id切片
			batchIds = append(batchIds, b.BatchId)
		}

		var (
			bat    model.Batch
			result *gorm.DB
		)

		if len(batchIds) == 0 { //只有一个批次
			bat.Id = batchIds[0]
		} else {
			//多个批次
			result = db.Select("id").Where("id in (?)", batchIds).Order("sort desc").First(&bat)

			if result.Error != nil {
				return rsp.ReceivingOrdersRsp{}, result.Error
			}
		}

		maxSort := 0

		res.BatchId = bat.Id

		//循环排序最大的批次下的拣货数据，并取出sort最大的那个的id
		for _, pm := range pickMp[bat.Id] {
			if pm.Sort >= maxSort {
				res.Id = pm.Id
				res.Version = pm.Version
				res.TakeOrdersTime = pm.TakeOrdersTime
			}
		}
	}

	return res, nil
}

// 接单拣货
func ReceivingOrders(c *gin.Context) {

	var (
		form    req.ReceivingOrdersForm
		res     rsp.ReceivingOrdersRsp
		pick    []model.Pick
		err     error
		batches []model.Batch
	)

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err = c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, ecode.GetContextUserInfoFailed)
		return
	}

	// 先查询是否有当前拣货员被分配的任务或已经接单且未完成拣货的数据,如果被分配多条，第一按批次优先级，第二按拣货池优先级 优先拣货
	result := db.Model(&model.Pick{}).Where("pick_user = ? and status = 0 and typ = ?", userInfo.Name, form.Typ).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	now := time.Now()

	//有分配的拣货任务
	if result.RowsAffected > 0 {
		res, err = getPick(pick)
		if err != nil {
			xsq_net.ErrorJSON(c, err)
			return
		}
		//后台分配的单没有接单时间,更新接单时间
		if res.TakeOrdersTime == nil {
			result = db.Model(&model.Pick{}).Where("id = ?", res.Id).Update("take_orders_time", &now)

			if result.Error != nil {
				xsq_net.ErrorJSON(c, result.Error)
				return
			}
		}
		xsq_net.SucJson(c, res)
		return
	}

	//进行中的批次
	result = db.Model(&model.Batch{}).Where("status = 0 and typ = ?", form.Typ).Find(&batches)

	batchIds := make([]int, 0)

	for _, b := range batches {
		batchIds = append(batchIds, b.Id)
	}

	if len(batchIds) == 0 {
		xsq_net.ErrorJSON(c, errors.New("没有进行中的批次,无法接单"))
		return
	}

	//查询未被接单的拣货池数据
	result = db.Model(&model.Pick{}).Where("batch_id in (?) and pick_user = '' and status = 0 and typ = ?", batchIds, form.Typ).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//拣货池有未接单的数据
	if result.RowsAffected > 0 {

		res, err = getPick(pick)
		if err != nil {
			xsq_net.ErrorJSON(c, err)
		}

		tx := db.Begin()

		//更新拣货池 + version 防并发
		result = tx.Model(&model.Pick{}).
			Where("id = ? and version = ?", res.Id, res.Version).
			Updates(map[string]interface{}{
				"pick_user":        userInfo.Name,
				"take_orders_time": &now,
				"version":          gorm.Expr("version + ?", 1),
			})

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, ecode.DataSaveError)
			return
		}

		tx.Commit()

		xsq_net.SucJson(c, res)
		return
	} else {
		xsq_net.ErrorJSON(c, errors.New("暂无拣货单"))
		return
	}
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

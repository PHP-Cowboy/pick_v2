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
		res     rsp.ReceivingOrdersRsp
		pick    []model.Pick
		err     error
		batches []model.Batch
	)

	db := global.DB

	claims, ok := c.Get("claims")

	if !ok {
		xsq_net.ErrorJSON(c, errors.New("claims 获取失败"))
		return
	}

	userInfo := claims.(*middlewares.CustomClaims)

	// 先查询是否有当前拣货员被分配的任务或已经接单且未完成拣货的数据,如果被分配多条，第一按批次优先级，第二按拣货池优先级 优先拣货
	result := db.Model(&model.Pick{}).Where("pick_user = ? and status = 0", userInfo.Name).Find(&pick)

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
	result = db.Where("status = 0").Find(&batches)

	batchIds := make([]int, 0)

	for _, b := range batches {
		batchIds = append(batchIds, b.Id)
	}

	if len(batchIds) == 0 {
		xsq_net.ErrorJSON(c, errors.New("没有进行中的批次,无法接单"))
		return
	}

	//查询未被接单的拣货池数据
	result = db.Model(&model.Pick{}).Where("batch_id in (?) and pick_user = '' and status = 0", batchIds).Find(&pick)

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

// 完成拣货
func CompletePick(c *gin.Context) {
	var form req.CompletePickForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	// 这里是否需要做并发处理
	var (
		pick       model.Pick
		pickGoods  []model.PickGoods
		orderGoods []model.OrderGoods
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
		xsq_net.ErrorJSON(c, errors.New("请确认拣货单是否被分配给其他拣货员"))
		return
	}

	tx := db.Begin()

	//****************************** 无需拣货 ******************************//
	if form.Type == 2 {
		//更新主表 无需拣货直接更新为复核完成
		result = tx.Model(&model.Pick{}).Where("id = ?", pick.Id).Updates(map[string]interface{}{"status": 2})
		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		err := UpdateBatchPickNums(tx, pick.BatchId)

		if err != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		tx.Commit()

		xsq_net.Success(c)
		return
	}
	//****************************** 无需拣货逻辑完成 ******************************//

	//****************************** 正常拣货逻辑 ******************************//
	//step:处理前端传递的拣货数据，构造[订单表id切片,订单表id和拣货商品表id map,sku完成数量 map]
	//step: 根据 订单表id切片 查出订单数据 根据支付时间升序
	//step: 构造 拣货商品表 id, 完成数量 并扣减 sku 完成数量
	//step: 更新拣货商品表

	var (
		orderGoodsIds      []int
		orderPickGoodsIdMp = make(map[int]int, 0)
		skuCompleteNumMp   = make(map[string]int, 0)
		totalNum           int //更新拣货池拣货数量
	)

	//step:处理前端传递的拣货数据，构造[订单表id切片,订单表id和拣货商品表id映射,sku完成数量映射]
	for _, cp := range form.CompletePick {
		//全部订单数据id
		for _, ids := range cp.ParamsId {
			orderGoodsIds = append(orderGoodsIds, ids.OrderGoodsId)
			//map[订单表id]拣货商品表id
			orderPickGoodsIdMp[ids.OrderGoodsId] = ids.PickGoodsId
		}
		//sku完成数量
		skuCompleteNumMp[cp.Sku] = cp.CompleteNum
		totalNum += cp.CompleteNum //总拣货数量
	}

	//step: 根据 订单表id切片 查出订单数据 根据支付时间升序
	result = db.Table("t_pick_order_goods og").
		Select("og.*").
		Joins("left join t_pick_order o on og.pick_order_id = o.id").
		Where("og.id in (?)", orderGoodsIds).
		Order("pay_at ASC").
		Find(&orderGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//拣货表 id 和 拣货数量
	mp := make(map[int]int, 0)

	var pickGoodsIds []int

	//step: 构造 拣货商品表 id, 完成数量 并扣减 sku 完成数量
	for _, info := range orderGoods {
		//完成数量
		completeNum, completeOk := skuCompleteNumMp[info.Sku]

		if !completeOk {
			continue
		}

		pickGoodsId, mpOk := orderPickGoodsIdMp[info.Id]

		if !mpOk {
			continue
		}

		pickCompleteNum := 0

		if completeNum >= info.LackCount { //完成数量大于等于需拣数量
			pickCompleteNum = info.LackCount
			skuCompleteNumMp[info.Sku] = completeNum - info.LackCount //减
		} else {
			//按下单时间拣货少于需拣时
			pickCompleteNum = completeNum
			skuCompleteNumMp[info.Sku] = 0
		}
		pickGoodsIds = append(pickGoodsIds, pickGoodsId)
		mp[pickGoodsId] = pickCompleteNum

	}

	//查出拣货商品数据
	result = tx.Where("id in (?)", pickGoodsIds).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//更新拣货数量数据
	for i, good := range pickGoods {
		completeNum, mpOk := mp[good.Id]

		if !mpOk {
			continue
		}

		pickGoods[i].CompleteNum = completeNum
	}

	//正常拣货 更新拣货数量
	result = tx.Save(&pickGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//更新主表
	result = tx.Model(&model.Pick{}).Where("id = ?", pick.Id).Updates(map[string]interface{}{"status": 1, "pick_num": totalNum})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	err := UpdateBatchPickNums(tx, pick.BatchId)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

// 剩余数量 放拣货池那边
func RemainingQuantity(c *gin.Context) {

	var (
		count    int64
		batches  []model.Batch
		batchIds []int
	)

	db := global.DB

	//批次进行中或暂停的单数量
	result := db.Where("status = 0 or status = 2").Find(&batches)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, b := range batches {
		batchIds = append(batchIds, b.Id)
	}

	if len(batchIds) > 0 {
		result = db.Model(&model.Pick{}).Where("batch_id in (?) and status = ? ", batchIds, model.ToBePickedStatus).Count(&count)

		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
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

	takeOrdersTime := ""
	if pick.TakeOrdersTime != nil {
		takeOrdersTime = pick.TakeOrdersTime.Format(timeutil.TimeFormat)
	}
	res.TakeOrdersTime = takeOrdersTime
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

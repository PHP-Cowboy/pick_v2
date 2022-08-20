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

func getPick(pick []batch.Pick, pickUser string) (res rsp.ReceivingOrdersRsp, err error) {

	if len(pick) == 1 { //只查到一条
		res.Id = pick[0].Id
	} else { //查到多条
		//排序
		var (
			batchIds []int
			batchMp  = make(map[int]struct{}, 0)
			pickMp   = make(map[int][]batch.Pick, 0)
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
			bat    batch.Batch
			result *gorm.DB
		)

		now := time.Now()

		db := global.DB

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
			}
		}

		//更新拣货池
		result = db.Model(&batch.Pick{}).Where("id = ?", res.Id).Updates(map[string]interface{}{"pick_user": pickUser, "take_orders_time": &now})

		if result.Error != nil {
			return rsp.ReceivingOrdersRsp{}, result.Error
		}
	}

	return res, nil
}

// 接单拣货
func ReceivingOrders(c *gin.Context) {
	var (
		res     rsp.ReceivingOrdersRsp
		pick    []batch.Pick
		err     error
		batches []batch.Batch
	)

	db := global.DB

	claims, ok := c.Get("claims")

	if !ok {
		xsq_net.ErrorJSON(c, errors.New("claims 获取失败"))
		return
	}

	userInfo := claims.(*middlewares.CustomClaims)

	// 先查询是否有当前拣货员被分配的任务或已经接单且未完成拣货的数据,如果被分配多条，第一按批次优先级，第二按拣货池优先级 优先拣货
	result := db.Model(&batch.Pick{}).Where("pick_user = ? and status = 0", userInfo.Name).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//有分配的拣货任务
	if result.RowsAffected > 0 {
		res, err = getPick(pick, userInfo.Name)
		if err != nil {
			xsq_net.ErrorJSON(c, err)
			return
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
		xsq_net.ErrorJSON(c, errors.New("已停止拣货,无法接单"))
		return
	}

	//查询未被接单的拣货池数据
	result = db.Model(&batch.Pick{}).Where("batch_id in (?) and pick_user = '' and status = 0", batchIds).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if result.RowsAffected > 0 {

		res, err = getPick(pick, userInfo.Name)
		if err != nil {
			xsq_net.ErrorJSON(c, err)
		}

		//更新拣货单数量
		result = db.Model(batch.Batch{}).Where("id = ?", res.BatchId).Update("pick_num", gorm.Expr("pick_num + ?", 1))

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

	// todo 这里是否需要做并发处理
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
		xsq_net.ErrorJSON(c, errors.New("请确认改单是否被分配给其他拣货员"))
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

	totalNum := 0 //更新拣货池拣货数量

	for k, pg := range pickGoods {
		num, pgMpOk := pickGoodsMap[pg.Id]

		if !pgMpOk {
			continue
		}

		pickGoods[k].CompleteNum = num

		totalNum += num
	}

	tx := db.Begin()

	status := 1

	updates := make(map[string]interface{}, 0)

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

		updates["pick_num"] = totalNum
	}

	updates["status"] = status

	//更新主表
	result = tx.Model(&batch.Pick{}).Where("id = ?", pick.Id).Updates(updates)

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

	var (
		count    int64
		batches  []batch.Batch
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

	result = db.Model(&batch.Pick{}).Where("batch_id in (?) and status = 0 and pick_user = ''", batchIds).Count(&count)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.SucJson(c, gin.H{"count": count})
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

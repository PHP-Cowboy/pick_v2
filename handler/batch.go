package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"pick_v2/common/constant"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/model/batch"
	"pick_v2/model/order"
	"pick_v2/utils/ecode"
	"pick_v2/utils/helper"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"strconv"
	"strings"
	"time"
)

//生成拣货批次
func CreateBatch(c *gin.Context) {
	var form req.CreateBatchForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	tx := global.DB.Begin()

	claims, ok := c.Get("claims")

	if !ok {
		xsq_net.ErrorJSON(c, ecode.DataNotExist)
		return
	}

	if len(form.Goods) > 0 {
		var sku, goodsName string
		for _, good := range form.Goods {
			sku = good.Sku + ","
			goodsName = good.Name + ","
		}
		form.Sku = strings.TrimRight(sku, ",")
		form.GoodsName = strings.TrimRight(goodsName, ",")
	}

	userInfo := claims.(*middlewares.CustomClaims)

	deliveryEndTime, errDeliveryEndTime := time.ParseInLocation(timeutil.TimeFormat, form.DeliveryEndTime, time.Local)
	deliveryStartTime, errDeliveryStartTime := time.ParseInLocation(timeutil.TimeFormat, form.DeliveryStartTime, time.Local)
	payEndTime, errPayEndTime := time.ParseInLocation(timeutil.TimeFormat, form.PayEndTime, time.Local)

	if errDeliveryEndTime != nil || errDeliveryStartTime != nil || errPayEndTime != nil {
		xsq_net.ErrorJSON(c, ecode.DataTransformationError)
		return
	}

	//批次数据
	batches := batch.Batch{
		WarehouseId:     form.WarehouseId,
		BatchName:       form.Lines + helper.GetDeliveryMethod(form.DeType),
		DeliveryEndTime: &deliveryEndTime,
		ShopNum:         0,
		OrderNum:        0,
		GoodsNum:        0,
		UserName:        userInfo.Name,
		Line:            form.Lines,
		DeliveryMethod:  form.DeType,
		EndTime:         &payEndTime,
		Status:          0,
		PickNum:         0,
		RecheckSheetNum: 0,
		Sort:            0,
	}

	result := tx.Save(&batches)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	condition := batch.BatchCondition{
		BatchId:           batches.Id,
		WarehouseId:       form.WarehouseId,
		PayEndTime:        &payEndTime,
		DeliveryStartTime: &deliveryStartTime,
		DeliveryEndTime:   &deliveryEndTime,
		Line:              form.Lines,
		DeliveryMethod:    form.DeType,
		Sku:               form.Sku,
		Goods:             form.GoodsName,
	}

	//筛选条件保存
	result = tx.Save(&condition)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	form.BatchNumber = strconv.Itoa(batches.Id)

	goodsRes, err := RequestGoodsList(form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	if len(goodsRes.Data.List) <= 0 {
		xsq_net.ErrorJSON(c, ecode.NoOrderFound)
		return
	}

	var (
		orders        []order.OrderInfo
		prePicks      []batch.PrePick
		prePickGoods  []*batch.PrePickGoods
		prePickRemark []*batch.PrePickRemark
		shopMap       = make(map[int]int, 0)
	)

	//订单相关数据
	for _, goods := range goodsRes.Data.List {
		orders = append(orders, order.OrderInfo{
			BatchId:          batches.Id,
			ShopId:           goods.ShopId,
			ShopName:         goods.ShopName,
			ShopType:         goods.ShopType,
			ShopCode:         goods.ShopCode,
			HouseCode:        goods.HouseCode,
			Line:             goods.Line,
			Number:           goods.Number,
			Status:           goods.Status,
			DeliveryAt:       goods.DeliveryAt,
			DistributionType: goods.DistributionType,
			OrderRemark:      goods.OrderRemark,
			Province:         goods.Province,
			City:             goods.City,
			District:         goods.District,
			Address:          goods.Address,
			ConsigneeName:    goods.ConsigneeName,
			ConsigneeTel:     goods.ConsigneeTel,
			Name:             goods.Name,
			Sku:              goods.Sku,
			GoodsSpe:         goods.GoodsSpe,
			GoodsType:        goods.GoodsType,
			Shelves:          goods.Shelves,
			OriginalPrice:    goods.OriginalPrice,
			DiscountPrice:    int(goods.DiscountPrice * 100),
			GoodsUnit:        goods.GoodsUnit,
			SaleUnit:         goods.SaleUnit,
			SaleCode:         goods.SaleCode,
			PayCount:         goods.PayCount,
			CloseCount:       goods.CloseCount,
			OutCount:         goods.OutCount,
			GoodsRemark:      goods.GoodsRemark,
			PickStatus:       goods.PickStatus,
			PayAt:            goods.PayAt,
			LackCount:        goods.LackCount,
		})

		prePickGoods = append(prePickGoods, &batch.PrePickGoods{
			WarehouseId: form.WarehouseId,
			BatchId:     batches.Id,
			OrderInfoId: goods.Id,
			Number:      goods.Number,
			PrePickId:   0,
			ShopId:      goods.ShopId,
			GoodsName:   goods.Name,
			GoodsType:   goods.GoodsType,
			GoodsSpe:    goods.GoodsSpe,
			Shelves:     goods.Shelves,
			NeedNum:     0,
			CloseNum:    goods.CloseCount,
			OutCount:    goods.OutCount,
			NeedOutNum:  0,
		})

		if goods.GoodsRemark != "" || goods.OrderRemark != "" {
			prePickRemark = append(prePickRemark, &batch.PrePickRemark{
				WarehouseId: form.WarehouseId,
				BatchId:     batches.Id,
				OrderInfoId: goods.Id,
				ShopId:      goods.ShopId,
				Number:      goods.Number,
				OrderRemark: goods.OrderRemark,
				GoodsRemark: goods.GoodsRemark,
				ShopName:    goods.ShopName,
				Line:        goods.Line,
				PrePickId:   0,
			})
		}

		_, ok := shopMap[goods.ShopId]
		if ok {
			continue
		}
		shopMap[goods.ShopId] = 0
		prePicks = append(prePicks, batch.PrePick{
			WarehouseId: form.WarehouseId,
			BatchId:     batches.Id,
			ShopId:      goods.ShopId,
			ShopCode:    goods.ShopCode,
			ShopName:    goods.ShopName,
			Line:        goods.Line,
			Status:      0,
		})

	}

	result = tx.Save(&orders)
	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Save(&prePicks)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, pick := range prePicks {
		shopMap[pick.ShopId] = pick.Id
	}

	for k, good := range prePickGoods {
		val, ok := shopMap[good.ShopId]
		if !ok {
			xsq_net.ErrorJSON(c, ecode.MapKeyNotExist)
			return
		}
		prePickGoods[k].PrePickId = val
	}

	result = tx.Save(&prePickGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if len(prePickRemark) > 0 {
		for k, remark := range prePickRemark {
			val, ok := shopMap[remark.ShopId]
			if !ok {
				xsq_net.ErrorJSON(c, ecode.MapKeyNotExist)
				return
			}
			prePickRemark[k].PrePickId = val
		}

		result = tx.Save(&prePickRemark)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	}

	tx.Commit()

	xsq_net.Success(c)
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

	var (
		batches      []batch.Batch
		batchIdSlice []int
	)

	db := global.DB

	//子表数据
	if form.GoodsName != "" || form.Number != "" || form.ShopId > 0 {
		var batchIds []struct {
			BatchId int
		}

		preGoodsRes := global.DB.Model(batch.PrePickGoods{}).
			Where(batch.PrePickGoods{
				GoodsName: form.GoodsName,
				Number:    form.Number,
				ShopId:    form.ShopId,
			}).
			Select("batch_id").
			Find(&batchIds)

		if preGoodsRes.Error != nil {
			xsq_net.ErrorJSON(c, preGoodsRes.Error)
			return
		}

		//利用map键唯一，去重
		uMap := make(map[int]struct{}, 0)
		for _, b := range batchIds {
			_, ok := uMap[b.BatchId]
			if ok {
				continue
			}
			uMap[b.BatchId] = struct{}{}
			batchIdSlice = append(batchIdSlice, b.BatchId)
		}
	}

	if form.Line != "" {
		db = db.Where("line like ?", form.Line+"%")
	}

	if len(batchIdSlice) > 0 {
		db = db.Where("id in (?)", batchIdSlice)
	}

	if form.CreateTime != "" {
		db = db.Where("create_time <= ?", form.CreateTime)
	}

	if form.EndTime != "" {
		db = db.Where("end_time <= ?", form.EndTime)
	}

	result := db.Where(map[string]interface{}{"status": form.Status}).Find(&batches)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	db.Scopes(model.Paginate(form.Page, form.Size)).Order("sort desc").Find(&batches)

	list := make([]*rsp.Batch, 0, len(batches))
	for _, b := range batches {

		deliveryStartTime := ""
		if b.DeliveryStartTime != nil {
			deliveryStartTime = b.DeliveryStartTime.Format(timeutil.TimeFormat)
		}

		list = append(list, &rsp.Batch{
			Id:                b.Id,
			BatchName:         b.BatchName,
			DeliveryStartTime: deliveryStartTime,
			DeliveryEndTime:   b.DeliveryEndTime.Format(timeutil.TimeFormat),
			ShopNum:           b.ShopNum,
			OrderNum:          b.OrderNum,
			GoodsNum:          b.GoodsNum,
			UserName:          b.UserName,
			Line:              b.Line,
			DeliveryMethod:    b.DeliveryMethod,
			EndTime:           b.EndTime.Format(timeutil.TimeFormat),
			Status:            b.Status,
			PickNum:           b.PickNum,
			RecheckSheetNum:   b.RecheckSheetNum,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

type Ret struct {
	OutC      int
	NeedC     int
	ShopId    int
	GoodsType string
}

//预拣池基础信息
func GetBase(c *gin.Context) {

	var (
		form      req.GetBaseForm
		batchCond batch.BatchCondition
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result := global.DB.Where("batch_id = ?", form.BatchId).First(&batchCond)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	deliveryStartTime := ""
	if batchCond.DeliveryStartTime != nil {
		deliveryStartTime = batchCond.DeliveryStartTime.Format(timeutil.TimeFormat)
	}

	ret := rsp.GetBaseRsp{
		CreateTime:        batchCond.CreateTime.Format(timeutil.TimeFormat),
		PayEndTime:        batchCond.PayEndTime.Format(timeutil.TimeFormat),
		DeliveryStartTime: deliveryStartTime,
		DeliveryEndTime:   batchCond.DeliveryEndTime.Format(timeutil.TimeFormat),
		DeliveryMethod:    batchCond.DeliveryMethod,
		Line:              batchCond.Line,
		Goods:             batchCond.Goods,
	}

	xsq_net.SucJson(c, ret)
}

//预拣池列表
func GetPrePickList(c *gin.Context) {
	var (
		form req.GetPrePickListForm
		res  rsp.GetPrePickListRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		prePicks []batch.PrePick
		//prePickGoods []batch.PrePickGoods
		prePickIds []int
	)

	db := global.DB

	result := db.Find(&prePicks)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	db.Scopes(model.Paginate(form.Page, form.Size)).Find(&prePicks)

	for _, pick := range prePicks {
		prePickIds = append(prePickIds, pick.Id)
	}

	retCount := []Ret{}

	result = db.Model(&batch.PrePickGoods{}).
		Select("SUM(out_count) as outC, SUM(need_num) AS needC, shop_id, goods_type").
		Where("pre_pick_id in (?)", prePickIds).
		Group("shop_id, goods_type").
		Find(&retCount)

	typeMap := make(map[int]map[string]rsp.PickCount, 0)

	for _, r := range retCount {
		_, ok := typeMap[r.ShopId]
		if !ok {
			countMap := make(map[string]rsp.PickCount, 0)
			typeMap[r.ShopId] = countMap
			countMap[r.GoodsType] = rsp.PickCount{
				WaitingPick: r.NeedC,
				PickedCount: r.OutC,
			}
			typeMap[r.ShopId][r.GoodsType] = countMap[r.GoodsType]
		}
	}

	for _, pick := range prePicks {
		res.List = append(res.List, &rsp.PrePick{
			Id:           pick.Id,
			ShopCode:     pick.ShopCode,
			ShopName:     pick.ShopName,
			Line:         pick.Line,
			CategoryInfo: typeMap[pick.ShopId],
		})
	}

	xsq_net.SucJson(c, res)

}

//预拣货明细
func GetPrePickDetail(c *gin.Context) {
	var (
		form req.GetPrePickDetailForm
		res  rsp.GetPrePickDetailRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		prePickGoods  []batch.PrePickGoods
		prePickRemark []batch.PrePickRemark
	)

	db := global.DB

	result := db.Where("pre_pick_id = ?", form.PrePickId).Find(&prePickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	goodsMap := make(map[string][]rsp.PrePickGoods, 0)

	for _, goods := range prePickGoods {
		goodsMap[goods.GoodsType] = append(goodsMap[goods.GoodsType], rsp.PrePickGoods{
			GoodsName:  goods.GoodsName,
			GoodsSpe:   goods.GoodsSpe,
			Shelves:    goods.Shelves,
			NeedNum:    goods.NeedNum,
			CloseNum:   goods.CloseNum,
			OutCount:   goods.OutCount,
			NeedOutNum: goods.NeedOutNum,
		})
	}

	res.Goods = goodsMap

	result = db.Where("pre_pick_id = ?", form.PrePickId).Find(&prePickRemark)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := []rsp.Remark{}
	for _, remark := range prePickRemark {
		list = append(list, rsp.Remark{
			Number:      remark.Number,
			OrderRemark: remark.OrderRemark,
			GoodsRemark: remark.GoodsRemark,
		})
	}

	res.RemarkList = list

	xsq_net.SucJson(c, res)
}

//置顶
func Topping(c *gin.Context) {
	var form req.ToppingForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	//redis
	redis := global.Redis

	redisKey := constant.BATCH_TOPPING

	val, err := redis.Do(context.Background(), "incr", redisKey).Result()
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	sort := int(val.(int64))

	result := global.DB.Model(batch.Batch{}).Where("id = ?", form.Id).Update("sort", sort)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

//批次池内单数量
func GetPoolNum(c *gin.Context) {
	var res rsp.GetPoolNumRsp

	res = rsp.GetPoolNumRsp{
		PrePickNum:  100,
		PickNum:     100,
		ToReviewNum: 100,
		CompleteNum: 100,
	}

	xsq_net.SucJson(c, res)
}

//批量拣货
func BatchPick(c *gin.Context) {
	var form req.BatchPickForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	warehouseId := c.GetInt("warehouseId")

	db := global.DB
	var (
		prePick        []batch.PrePick
		prePickGoods   []batch.PrePickGoods
		prePickRemarks []batch.PrePickRemark
	)

	result := db.Where("id in (?)", form.Ids).Find(&prePick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Where("pre_pick_id in (?)", form.Ids).Find(&prePickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	prePickGoodsMap := make(map[int][]batch.PrePickGoods, 0)

	for _, goods := range prePickGoods {
		prePickGoodsMap[goods.PrePickId] = append(prePickGoodsMap[goods.PrePickId], goods)
	}

	result = db.Where("pre_pick_id in (?)", form.Ids).Find(&prePickRemarks)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	prePickRemarksMap := make(map[int][]batch.PrePickRemark, 0)

	for _, remark := range prePickRemarks {
		prePickRemarksMap[remark.PrePickId] = append(prePickRemarksMap[remark.PrePickId], remark)
	}

	var (
		pickGoods  []batch.PickGoods
		pickRemark []batch.PickRemark
	)

	tx := db.Begin()

	for _, pre := range prePick {
		pick := batch.Pick{
			WarehouseId:    warehouseId,
			BatchId:        pre.BatchId,
			PrePickIds:     strconv.Itoa(pre.Id),
			TaskName:       pre.ShopName,
			ShopCode:       pre.ShopCode,
			ShopName:       pre.ShopName,
			Line:           pre.Line,
			ShopNum:        0,
			OrderNum:       pre.OrderNum,
			NeedNum:        0,
			PickUser:       "",
			ReviewUser:     "",
			TakeOrdersTime: nil,
			Sort:           0,
			Version:        0,
		}

		result = tx.Save(&pick)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, goods := range prePickGoodsMap[pre.Id] {
			pickGoods = append(pickGoods, batch.PickGoods{
				WarehouseId:    warehouseId,
				BatchId:        pre.BatchId,
				PickId:         pick.Id,
				PrePickGoodsId: goods.Id,
				GoodsName:      goods.GoodsName,
				GoodsSpe:       goods.GoodsSpe,
				Shelves:        goods.Shelves,
				NeedNum:        goods.NeedNum,
			})
		}

		result = tx.Save(&pickGoods)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, remark := range prePickRemarksMap[pre.Id] {
			pickRemark = append(pickRemark, batch.PickRemark{
				WarehouseId:     warehouseId,
				BatchId:         pre.BatchId,
				PickId:          pick.Id,
				PrePickRemarkId: remark.Id,
				Number:          remark.Number,
				OrderRemark:     remark.OrderRemark,
				GoodsRemark:     remark.GoodsRemark,
				ShopName:        remark.ShopName,
				Line:            remark.Line,
			})
		}

		if len(pickRemark) > 0 {
			result = tx.Save(&pickRemark)

			if result.Error != nil {
				tx.Rollback()
				xsq_net.ErrorJSON(c, result.Error)
				return
			}
		}
	}

	//批量更新
	result = tx.Model(batch.PrePick{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"status": 1})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Model(batch.PrePickGoods{}).Where("pre_pick_id in (?)", form.Ids).Updates(map[string]interface{}{"status": 1})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Model(batch.PrePickRemark{}).Where("pre_pick_id in (?)", form.Ids).Updates(map[string]interface{}{"status": 1})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

//合并拣货
func MergePick(c *gin.Context) {
	var form req.MergePickForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	form.WarehouseId = c.GetInt("warehouseId")

	var err error

	switch form.Type {
	case 1:
		err = ByAllOrder(form)
		break
	case 2:
		err = ByClassification(form)
		break
	case 3:
		err = ByGoods(form)
		break
	default:
		err = errors.New("类型不合法")
	}

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

//全单拣货
func ByAllOrder(form req.MergePickForm) error {
	var (
		prePickGoods   []batch.PrePickGoods
		prePickRemarks []batch.PrePickRemark
		pickGoods      []batch.PickGoods
		pickRemarks    []batch.PickRemark
	)

	db := global.DB

	prePickIds := ""

	for _, id := range form.Ids {
		prePickIds = strconv.Itoa(id) + ","
	}

	prePickIds = strings.TrimRight(prePickIds, ",")

	result := db.Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickGoods)

	if result.Error != nil {
		return result.Error
	}

	result = db.Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickRemarks)

	if result.Error != nil {
		return result.Error
	}

	tx := db.Begin()

	pick := batch.Pick{
		WarehouseId:    form.WarehouseId,
		BatchId:        form.BatchId,
		PrePickIds:     prePickIds,
		TaskName:       form.TaskName,
		ShopCode:       "",
		ShopName:       form.TaskName,
		Line:           "",
		ShopNum:        0,
		OrderNum:       0,
		NeedNum:        0,
		PickUser:       "",
		ReviewUser:     "",
		TakeOrdersTime: nil,
		Sort:           0,
		Version:        0,
	}

	result = tx.Save(&pick)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	var (
		prePickGoodsIds, prePickRemarksIds []int
	)

	for _, goods := range prePickGoods {
		//更新预拣货池商品使用
		prePickGoodsIds = append(prePickGoodsIds, goods.Id)

		pickGoods = append(pickGoods, batch.PickGoods{
			WarehouseId:    form.WarehouseId,
			BatchId:        form.BatchId,
			PickId:         pick.Id,
			PrePickGoodsId: goods.Id,
			GoodsName:      goods.GoodsName,
			GoodsSpe:       goods.GoodsSpe,
			Shelves:        goods.Shelves,
			NeedNum:        goods.NeedNum,
		})
	}

	result = tx.Save(&pickGoods)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	for _, remark := range prePickRemarks {
		//更新预拣货池商品备注使用
		prePickRemarksIds = append(prePickRemarksIds, remark.Id)

		pickRemarks = append(pickRemarks, batch.PickRemark{
			WarehouseId:     form.WarehouseId,
			BatchId:         form.BatchId,
			PickId:          pick.Id,
			PrePickRemarkId: remark.Id,
			Number:          remark.Number,
			OrderRemark:     remark.OrderRemark,
			GoodsRemark:     remark.GoodsRemark,
			ShopName:        remark.ShopName,
			Line:            remark.Line,
		})
	}

	result = tx.Save(&pickRemarks)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	//全单拣货，全部状态更新
	result = tx.Model(batch.PrePick{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"status": 1})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	//可能没查到相关数据
	if len(prePickGoodsIds) > 0 {
		result = tx.Model(batch.PrePickGoods{}).Where("id in (?)", prePickGoodsIds).Updates(map[string]interface{}{"status": 1})

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	if len(prePickRemarksIds) > 0 {
		result = tx.Model(batch.PrePickRemark{}).Where("id in (?)", prePickRemarksIds).Updates(map[string]interface{}{"status": 1})

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	tx.Commit()

	return nil
}

//按分类拣货
func ByClassification(form req.MergePickForm) error {
	var (
		prePickGoods   []batch.PrePickGoods
		prePickRemarks []batch.PrePickRemark

		pickGoods   []batch.PickGoods
		pickRemarks []batch.PickRemark
	)

	db := global.DB

	var (
		prePickIds string
		prePickGoodsIds,
		orderInfoIds,
		prePickRemarksIds []int
	)

	//查全部数据，程序过滤出 goods_type = form.TypeParam 的数据 并计算是否更新 pre_pick 表状态
	result := db.Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickGoods)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("仓库分类:" + form.TypeParam + "数据未找到")
	}

	tx := db.Begin()

	pick := batch.Pick{
		WarehouseId:    form.WarehouseId,
		BatchId:        form.BatchId,
		PrePickIds:     prePickIds,
		TaskName:       form.TaskName,
		ShopCode:       "",
		ShopName:       form.TaskName,
		Line:           "",
		ShopNum:        0,
		OrderNum:       0,
		NeedNum:        0,
		PickUser:       "",
		ReviewUser:     "",
		TakeOrdersTime: nil,
		Sort:           0,
		Version:        0,
	}

	result = tx.Save(&pick)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	for _, goods := range prePickGoods {

		//只保留相关仓库类型数据
		if goods.GoodsType != form.TypeParam {
			continue
		}

		prePickGoodsIds = append(prePickGoodsIds, goods.Id)

		orderInfoIds = append(orderInfoIds, goods.Id)

		pickGoods = append(pickGoods, batch.PickGoods{
			WarehouseId:    form.WarehouseId,
			BatchId:        form.BatchId,
			PickId:         pick.Id,
			PrePickGoodsId: goods.Id,
			GoodsName:      goods.GoodsName,
			GoodsSpe:       goods.GoodsSpe,
			Shelves:        goods.Shelves,
			NeedNum:        goods.NeedNum,
		})
	}

	// 待拣货商品表数据 全部是 form.TypeParam 类型的时，更新 预拣货表状态
	if len(prePickGoods) == len(prePickGoodsIds) {
		//待优化 如果两个待拣货单，其中一个的仓库类型完成另一个未完成，这时这里就无法更新已完成的那个的状态
		//更新预拣货池预拣货表状态
		result = tx.Model(batch.PrePick{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"status": 1})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	result = tx.Save(&pickGoods)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	//更新预拣货池商品相关数据状态
	result = tx.Model(batch.PrePickGoods{}).Where("id in (?)", prePickGoodsIds).Updates(map[string]interface{}{"status": 1})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = db.Where("order_info_id in (?)", orderInfoIds).Find(&prePickRemarks)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if len(prePickRemarks) > 0 {
		for _, remark := range prePickRemarks {

			prePickRemarksIds = append(prePickRemarksIds, remark.Id)

			pickRemarks = append(pickRemarks, batch.PickRemark{
				WarehouseId:     form.WarehouseId,
				BatchId:         form.BatchId,
				PickId:          pick.Id,
				PrePickRemarkId: remark.Id,
				Number:          remark.Number,
				OrderRemark:     remark.OrderRemark,
				GoodsRemark:     remark.GoodsRemark,
				ShopName:        remark.ShopName,
				Line:            remark.Line,
			})
		}

		result = tx.Save(&pickRemarks)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		//更新预拣货池备注相关数据状态
		result = tx.Model(batch.PrePickRemark{}).Where("id in (?)", prePickRemarksIds).Updates(map[string]interface{}{"status": 1})

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	return nil
}

//按单品拣货
func ByGoods(form req.MergePickForm) error {
	var (
		prePickGoods   []batch.PrePickGoods
		prePickRemarks []batch.PrePickRemark

		pickGoods   []batch.PickGoods
		pickRemarks []batch.PickRemark
	)

	db := global.DB

	var (
		prePickIds string
		prePickGoodsIds,
		orderInfoIds,
		prePickRemarksIds []int
	)

	//查全部数据，程序过滤出 goods_name = form.TypeParam 的数据 并计算是否更新 pre_pick 表状态
	result := db.Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickGoods)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("仓库分类:" + form.TypeParam + "数据未找到")
	}

	tx := db.Begin()

	pick := batch.Pick{
		WarehouseId:    form.WarehouseId,
		BatchId:        form.BatchId,
		PrePickIds:     prePickIds,
		TaskName:       form.TaskName,
		ShopCode:       "",
		ShopName:       form.TaskName,
		Line:           "",
		ShopNum:        0,
		OrderNum:       0,
		NeedNum:        0,
		PickUser:       "",
		ReviewUser:     "",
		TakeOrdersTime: nil,
		Sort:           0,
		Version:        0,
	}

	result = tx.Save(&pick)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	for _, goods := range prePickGoods {

		//只保留相关商品数据
		if goods.GoodsType != form.TypeParam {
			continue
		}

		prePickGoodsIds = append(prePickGoodsIds, goods.Id)

		orderInfoIds = append(orderInfoIds, goods.Id)

		pickGoods = append(pickGoods, batch.PickGoods{
			WarehouseId:    form.WarehouseId,
			BatchId:        form.BatchId,
			PickId:         pick.Id,
			PrePickGoodsId: goods.Id,
			GoodsName:      goods.GoodsName,
			GoodsSpe:       goods.GoodsSpe,
			Shelves:        goods.Shelves,
			NeedNum:        goods.NeedNum,
		})
	}

	// 待拣货商品表数据 全部是 form.TypeParam 类型的时，更新 预拣货表状态
	if len(prePickGoods) == len(prePickGoodsIds) {
		//待优化 如果两个待拣货单，其中一个的仓库类型完成另一个未完成，这时这里就无法更新已完成的那个的状态
		//更新预拣货池预拣货表状态
		result = tx.Model(batch.PrePick{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"status": 1})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	result = tx.Save(&pickGoods)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	//更新预拣货池商品相关数据状态
	result = tx.Model(batch.PrePickGoods{}).Where("id in (?)", prePickGoodsIds).Updates(map[string]interface{}{"status": 1})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = db.Where("order_info_id in (?)", orderInfoIds).Find(&prePickRemarks)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if len(prePickRemarks) > 0 {
		for _, remark := range prePickRemarks {

			prePickRemarksIds = append(prePickRemarksIds, remark.Id)

			pickRemarks = append(pickRemarks, batch.PickRemark{
				WarehouseId:     form.WarehouseId,
				BatchId:         form.BatchId,
				PickId:          pick.Id,
				PrePickRemarkId: remark.Id,
				Number:          remark.Number,
				OrderRemark:     remark.OrderRemark,
				GoodsRemark:     remark.GoodsRemark,
				ShopName:        remark.ShopName,
				Line:            remark.Line,
			})
		}

		result = tx.Save(&pickRemarks)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		//更新预拣货池备注相关数据状态
		result = tx.Model(batch.PrePickRemark{}).Where("id in (?)", prePickRemarksIds).Updates(map[string]interface{}{"status": 1})

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	return nil
}

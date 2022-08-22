package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pick_v2/common/constant"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/model/batch"
	"pick_v2/model/order"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/helper"
	"pick_v2/utils/request"
	"pick_v2/utils/slice"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"strconv"
	"strings"
	"time"
)

// 生成拣货批次
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

	var (
		deliveryStartTime    time.Time
		errDeliveryStartTime error
	)

	deliveryEndTime, errDeliveryEndTime := time.ParseInLocation(timeutil.TimeFormat, form.DeliveryEndTime, time.Local)

	payEndTime, errPayEndTime := time.ParseInLocation(timeutil.TimeFormat, form.PayEndTime, time.Local)

	if errDeliveryEndTime != nil || errPayEndTime != nil {
		xsq_net.ErrorJSON(c, ecode.DataTransformationError)
		return
	}

	//批次数据
	batches := batch.Batch{
		WarehouseId:     form.WarehouseId,
		BatchName:       form.Lines + helper.GetDeliveryMethod(form.DeType),
		DeliveryEndTime: &deliveryEndTime,
		ShopNum:         0, //在后续的逻辑中更新处理，调用接口时需传批次id，只能更新，其他方式可能导致订货系统锁住数据，而拣货系统获取不到被锁住的数据
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
		BatchId:         batches.Id,
		WarehouseId:     form.WarehouseId,
		PayEndTime:      &payEndTime,
		DeliveryEndTime: &deliveryEndTime,
		Line:            form.Lines,
		DeliveryMethod:  form.DeType,
		Sku:             form.Sku,
		Goods:           form.GoodsName,
	}

	if form.DeliveryStartTime != "" {
		deliveryStartTime, errDeliveryStartTime = time.ParseInLocation(timeutil.TimeFormat, form.DeliveryStartTime, time.Local)
		if errDeliveryStartTime != nil {
			xsq_net.ErrorJSON(c, ecode.DataTransformationError)
			return
		}
		condition.DeliveryStartTime = &deliveryStartTime
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
	)

	mp, err := cache.GetClassification()

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	//订单相关数据 -店铺数 订单数 商品数
	goodsNum := 0                              //商品数
	shopNumMp := make(map[int]struct{}, 0)     //店铺
	orderNumMp := make(map[string]struct{}, 0) //订单

	//缓存中的线路数据
	lineCacheMp, errCache := cache.GetShopLine()

	if errCache != nil {
		xsq_net.ErrorJSON(c, errors.New("线路缓存获取失败"))
		return
	}

	for _, goods := range goodsRes.Data.List {
		goodsNum += goods.LackCount

		orderNumMp[goods.Number] = struct{}{}

		goodsType, mpOk := mp[goods.SecondType]

		if !mpOk {
			xsq_net.ErrorJSON(c, errors.New("商品类型:"+goods.SecondType+"数据未同步"))
			return
		}

		cacheMpLine, cacheMpOk := lineCacheMp[goods.ShopId]

		if !cacheMpOk {
			xsq_net.ErrorJSON(c, errors.New("店铺:"+goods.ShopName+"线路未同步，请先同步"))
			return
		}

		orders = append(orders, order.OrderInfo{
			Id:               goods.Id,
			BatchId:          batches.Id,
			ShopId:           goods.ShopId,
			ShopName:         goods.ShopName,
			ShopType:         goods.ShopType,
			ShopCode:         goods.ShopCode,
			HouseCode:        goods.HouseCode,
			Line:             cacheMpLine,
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
			GoodsType:        goodsType,
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
			PickTime:         nil,
		})

		prePickGoods = append(prePickGoods, &batch.PrePickGoods{
			WarehouseId:      form.WarehouseId,
			BatchId:          batches.Id,
			OrderInfoId:      goods.Id,
			Number:           goods.Number,
			PrePickId:        0, //后续逻辑变更
			ShopId:           goods.ShopId,
			DistributionType: goods.DistributionType,
			Sku:              goods.Sku,
			GoodsName:        goods.Name,
			GoodsType:        goodsType,
			GoodsSpe:         goods.GoodsSpe,
			Shelves:          goods.Shelves,
			NeedNum:          goods.LackCount,
			CloseNum:         goods.CloseCount,
			OutCount:         goods.OutCount,
			NeedOutNum:       goods.LackCount,
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
				Line:        cacheMpLine,
				PrePickId:   0,
			})
		}

		_, shopMpOk := shopNumMp[goods.ShopId]

		if shopMpOk {
			continue
		}

		shopNumMp[goods.ShopId] = struct{}{}

		prePicks = append(prePicks, batch.PrePick{
			WarehouseId: form.WarehouseId,
			BatchId:     batches.Id,
			ShopId:      goods.ShopId,
			ShopCode:    goods.ShopCode,
			ShopName:    goods.ShopName,
			Line:        cacheMpLine,
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

	shopMap := make(map[int]int, 0)

	for _, pick := range prePicks {
		shopMap[pick.ShopId] = pick.Id
	}

	for k, good := range prePickGoods {
		val, shopMapOk := shopMap[good.ShopId]
		if !shopMapOk {
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
			val, shopMapOk := shopMap[remark.ShopId]
			if !shopMapOk {
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

	shopNum := len(shopNumMp)
	orderNum := len(orderNumMp)

	result = tx.Model(&batch.Batch{}).Where("id = ?", batches.Id).Updates(map[string]interface{}{"goods_num": goodsNum, "shop_num": shopNum, "order_num": orderNum})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

// 根据订单生成批次
func CreateByOrder(c *gin.Context) {
	var form req.CreateByOrderReq

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	claims, ok := c.Get("claims")

	if !ok {
		xsq_net.ErrorJSON(c, errors.New("获取上下文用户数据失败"))
		return
	}

	userInfo := claims.(*middlewares.CustomClaims)

	tx := global.DB.Begin()

	//创建批次
	batches := batch.Batch{
		WarehouseId:     userInfo.WarehouseId,
		BatchName:       form.Number,
		DeliveryEndTime: nil,
		ShopNum:         0, //在后续的逻辑中更新处理，调用接口时需传批次id，只能更新，其他方式可能导致订货系统锁住数据，而拣货系统获取不到被锁住的数据
		OrderNum:        0,
		GoodsNum:        0,
		UserName:        userInfo.Name,
		Line:            "",
		DeliveryMethod:  0,
		EndTime:         nil,
		Status:          0,
		PickNum:         0,
		RecheckSheetNum: 0,
		Sort:            0,
	}

	result := tx.Save(&batches)

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
	)

	mp, err := cache.GetClassification()

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	//订单相关数据 -店铺数 订单数 商品数
	goodsNum := 0                              //商品数
	shopNumMp := make(map[int]struct{}, 0)     //店铺
	orderNumMp := make(map[string]struct{}, 0) //订单

	params := make(map[string]req.BatchParams, 0)

	//缓存中的线路数据
	lineCacheMp, errCache := cache.GetShopLine()

	if errCache != nil {
		xsq_net.ErrorJSON(c, errors.New("线路缓存获取失败"))
		return
	}

	for _, goods := range goodsRes.Data.List {

		//更新批次和批次条件的相关数据使用
		_, pOk := params[goods.Number]
		if !pOk {
			params[goods.Number] = req.BatchParams{
				DeliveryEndTime: goods.DeliveryAt,
				PayEndTime:      goods.PayAt,
				DeliveryMethod:  goods.DistributionType,
			}
		}

		goodsNum += goods.LackCount

		orderNumMp[goods.Number] = struct{}{}

		goodsType, mpOk := mp[goods.SecondType]

		if !mpOk {
			xsq_net.ErrorJSON(c, errors.New("商品类型:"+goods.SecondType+"数据未同步"))
			return
		}

		cacheMpLine, cacheMpOk := lineCacheMp[goods.ShopId]

		if !cacheMpOk {
			xsq_net.ErrorJSON(c, errors.New("店铺:"+goods.ShopName+"线路未同步，请先同步"))
			return
		}

		orders = append(orders, order.OrderInfo{
			Id:               goods.Id,
			BatchId:          batches.Id,
			ShopId:           goods.ShopId,
			ShopName:         goods.ShopName,
			ShopType:         goods.ShopType,
			ShopCode:         goods.ShopCode,
			HouseCode:        goods.HouseCode,
			Line:             cacheMpLine,
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
			GoodsType:        goodsType,
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
			PickTime:         nil,
		})

		prePickGoods = append(prePickGoods, &batch.PrePickGoods{
			WarehouseId:      userInfo.WarehouseId,
			BatchId:          batches.Id,
			OrderInfoId:      goods.Id,
			Number:           goods.Number,
			PrePickId:        0, //后续逻辑变更
			ShopId:           goods.ShopId,
			DistributionType: goods.DistributionType,
			Sku:              goods.Sku,
			GoodsName:        goods.Name,
			GoodsType:        goodsType,
			GoodsSpe:         goods.GoodsSpe,
			Shelves:          goods.Shelves,
			NeedNum:          goods.LackCount,
			CloseNum:         goods.CloseCount,
			OutCount:         goods.OutCount,
			NeedOutNum:       goods.LackCount,
		})

		if goods.GoodsRemark != "" || goods.OrderRemark != "" {
			prePickRemark = append(prePickRemark, &batch.PrePickRemark{
				WarehouseId: userInfo.WarehouseId,
				BatchId:     batches.Id,
				OrderInfoId: goods.Id,
				ShopId:      goods.ShopId,
				Number:      goods.Number,
				OrderRemark: goods.OrderRemark,
				GoodsRemark: goods.GoodsRemark,
				ShopName:    goods.ShopName,
				Line:        cacheMpLine,
				PrePickId:   0,
			})
		}

		_, shopMpOk := shopNumMp[goods.ShopId]

		if shopMpOk {
			continue
		}

		shopNumMp[goods.ShopId] = struct{}{}

		prePicks = append(prePicks, batch.PrePick{
			WarehouseId: userInfo.WarehouseId,
			BatchId:     batches.Id,
			ShopId:      goods.ShopId,
			ShopCode:    goods.ShopCode,
			ShopName:    goods.ShopName,
			Line:        cacheMpLine,
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

	shopMap := make(map[int]int, 0)

	for _, pick := range prePicks {
		shopMap[pick.ShopId] = pick.Id
	}

	for k, good := range prePickGoods {
		val, shopMapOk := shopMap[good.ShopId]
		if !shopMapOk {
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
			val, shopMapOk := shopMap[remark.ShopId]
			if !shopMapOk {
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

	shopNum := len(shopNumMp)
	orderNum := len(orderNumMp)

	updates := gin.H{}

	var (
		deliveryEndTime    time.Time
		payEndTime         time.Time
		errDeliveryEndTime error
		errPayEndTime      error
		deliveryMethod     int
	)

	if p, pOk := params[form.Number]; pOk {

		fmt.Println(p.DeliveryEndTime)
		fmt.Println(p.PayEndTime)

		deliveryEndTime, errDeliveryEndTime = time.ParseInLocation(timeutil.TimeFormat, p.DeliveryEndTime, time.Local)

		payEndTime, errPayEndTime = time.ParseInLocation(timeutil.TimeFormat, p.PayEndTime, time.Local)

		if errDeliveryEndTime != nil || errPayEndTime != nil {
			xsq_net.ErrorJSON(c, ecode.DataTransformationError)
			return
		}

		updates["delivery_end_time"] = &deliveryEndTime
		updates["pay_end_time"] = &payEndTime
		deliveryMethod = p.DeliveryMethod
		updates["delivery_method"] = deliveryMethod
	}

	updates["goods_num"] = goodsNum
	updates["shop_num"] = shopNum
	updates["order_num"] = orderNum

	result = tx.Model(&batch.Batch{}).Where("id = ?", batches.Id).Updates(updates)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//批次创建条件
	condition := batch.BatchCondition{
		BatchId:         batches.Id,
		WarehouseId:     userInfo.WarehouseId,
		PayEndTime:      &payEndTime,
		DeliveryEndTime: &deliveryEndTime,
		Line:            "",
		DeliveryMethod:  deliveryMethod,
		Sku:             "",
		Goods:           "",
	}

	tx.Save(condition)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

// 结束拣货批次
func EndBatch(c *gin.Context) {
	var form req.EndBatchForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		batches   batch.Batch
		pickGoods []batch.PickGoods
		pick      []batch.Pick
		orderInfo []order.OrderInfo
		outGoods  req.OutGoods
	)

	db := global.DB

	result := db.First(&batches, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if batches.Status != 2 {
		xsq_net.ErrorJSON(c, errors.New("请先停止拣货"))
		return
	}

	//修改批次状态为已结束
	result = db.Model(&batch.Batch{}).Where("id = ?", batches.Id).Updates(map[string]interface{}{"status": 1})

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&order.OrderInfo{}).Where("batch_id = ?", form.Id).Find(&orderInfo)
	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//查询批次下全部订单
	result = db.Model(&batch.PickGoods{}).Where("batch_id = ?", form.Id).Find(&pickGoods)
	if result.Error != nil {
		global.SugarLogger.Error("批次结束成功，但推送u8拣货数据查询失败:" + result.Error.Error())
		xsq_net.ErrorJSON(c, errors.New("批次结束成功，但推送u8拣货数据查询失败"))
		return
	}

	result = db.Model(&batch.Pick{}).Where("batch_id = ?", form.Id).Find(&pick)
	if result.Error != nil {
		global.SugarLogger.Error("批次结束成功，但推送u8拣货数据查询失败:" + result.Error.Error())
		xsq_net.ErrorJSON(c, errors.New("批次结束成功，但推送u8拣货主表数据查询失败"))
		return
	}

	//拣货表数据map
	mpPick := make(map[int]batch.Pick, 0)

	for _, p := range pick {
		mpPick[p.Id] = p
	}

	//拣货商品map
	mpGoods := make(map[int]batch.PickGoods, 0)

	for _, good := range pickGoods {
		mpGoods[good.OrderInfoId] = good
	}

	outGoods.BatchNumber = strconv.Itoa(batches.Id)

	//批次未完成订单map
	notCompleteOrderMp := make(map[string]struct{}, 0)
	//批次全部订单map
	allOrderMp := make(map[string][]order.OrderInfo, 0)
	for _, info := range orderInfo {

		//构造未完成订单,有出库单号认为已完成（确认出库时，需拣和复核数一致的会写入出库单号）
		if info.DeliveryOrderNo == "" {
			notCompleteOrderMp[info.Number] = struct{}{}
		}

		_, compOk := allOrderMp[info.Number]

		if !compOk {
			allOrderMp[info.Number] = make([]order.OrderInfo, 0)
		}

		//全部订单 把未完成的过滤掉
		allOrderMp[info.Number] = append(allOrderMp[info.Number], info)

		goods, ok := mpGoods[info.Id]

		if !ok {
			//
			continue
		}

		var outNumber, exWarehouse string

		pMp, pickOk := mpPick[goods.PickId]

		if pickOk {
			outNumber = strconv.Itoa(pMp.Id) //拣货单号
			exWarehouse = pMp.DeliveryOrderNo
		}

		outGoods.List = append(outGoods.List, req.OutGoodsList{
			GoodsLogId:   info.Id,
			Number:       info.Number,
			OutNumber:    outNumber, //拣货单号
			CkNumber:     exWarehouse,
			Sku:          info.Sku,
			Name:         info.Name,
			OutCount:     goods.CompleteNum,
			Price:        info.DiscountPrice,
			SumPrice:     goods.CompleteNum * info.DiscountPrice,
			OutAt:        goods.UpdateTime.Format(timeutil.TimeFormat),
			PayAt:        info.PayAt,
			GoodsSpe:     info.GoodsSpe,
			GoodsUnit:    info.GoodsUnit,
			DeliveryType: info.DistributionType,
		})
	}

	//请求接口 释放锁单
	err := OutGoods(outGoods)
	if err != nil {
		xsq_net.ErrorJSON(c, errors.New("归还订货系统欠货信息失败"))
		return
	}

	//将要被删除订单表number
	deleteNumbers := make([]string, 0)

	tx := db.Begin()

	completeOrder := make([]order.CompleteOrder, 0)
	completeOrderDetail := make([]order.CompleteOrderDetail, 0)

	//已完成订单转存入完成订单表 同时删除 订单商品表数据
	for k, orderSlice := range allOrderMp {
		//订单号在未完成订单map中，过滤掉
		_, ok := notCompleteOrderMp[k]
		if ok {
			continue
		}
		//不在未完成订单map中，
		//step1:存入 将要被删除订单主键
		deleteNumbers = append(deleteNumbers, k)
		//step2:构造存入完成订单表相关数据
		for i, o := range orderSlice {
			if i == 0 { //第一条中取出完成订单表数据
				payAt, payAtErr := time.ParseInLocation(timeutil.TimeZoneFormat, o.PayAt, time.Local)

				if payAtErr != nil {
					xsq_net.ErrorJSON(c, ecode.DataTransformationError)
					return
				}

				completeOrder = append(completeOrder, order.CompleteOrder{
					Number:         o.Number,
					OrderRemark:    o.OrderRemark,
					ShopId:         o.ShopId,
					ShopName:       o.ShopName,
					ShopType:       o.ShopType,
					ShopCode:       o.ShopCode,
					Line:           o.Line,
					DeliveryMethod: o.DistributionType,
					PayCount:       o.PayCount,
					CloseCount:     o.CloseCount,
					OutCount:       o.OutCount,
					Province:       o.Province,
					City:           o.City,
					District:       o.District,
					PickTime:       o.PickTime,
					PayAt:          payAt.Format(timeutil.TimeFormat),
				})
			}
			completeOrderDetail = append(completeOrderDetail, order.CompleteOrderDetail{
				Number:          o.Number,
				Name:            o.Name,
				Sku:             o.Sku,
				GoodsSpe:        o.GoodsSpe,
				GoodsType:       o.GoodsType,
				Shelves:         o.Shelves,
				PayCount:        o.PayCount,
				CloseCount:      o.CloseCount,
				ReviewCount:     o.LackCount, //这里可能有问题
				GoodsRemark:     o.GoodsRemark,
				DeliveryOrderNo: o.DeliveryOrderNo,
			})
		}
	}

	//添加完成订单主表数据
	result = tx.Create(&completeOrder)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//添加完成订单明细表数据
	result = tx.Create(&completeOrderDetail)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//删除订单商品表数据
	result = tx.Delete(&order.OrderInfo{}, "number in (?)", deleteNumbers)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	err = PushU8(pickGoods)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

func OutGoods(responseData interface{}) error {
	var result rsp.OutGoodsRsp

	path := "api/v1/remote/sync/out/goods"
	body, err := request.Post(path, responseData)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return err
	}

	if result.Code != 200 {
		global.SugarLogger.Errorf("path:%s,errMsg:%s", path, result.Msg)
		return errors.New(result.Msg)
	}

	return nil
}

func PushU8(pickGoods []batch.PickGoods) error {

	var orderInfo []order.OrderInfo

	numbers := []string{}

	mpGoods := make(map[string]batch.PickGoods, 0)

	for _, good := range pickGoods {
		numbers = append(numbers, good.Number)
		//订单数据
		//step1 map[number+sku]{...}
		mpGoods[good.Number+good.Sku] = good
		//step2 合并到 map[number][]PickGoods{...}
	}

	global.SugarLogger.Info(mpGoods)

	//去重
	numbers = slice.UniqueStringSlice(numbers)

	result := global.DB.Model(&order.OrderInfo{}).Where("number in (?)", numbers).Find(&orderInfo)
	if result.Error != nil {
		global.SugarLogger.Error("批次结束成功，但推送u8订单数据查询失败:" + result.Error.Error())
		return errors.New("批次结束成功，但推送u8订单数据查询失败")
	}

	mpPgv := make(map[string]PickGoodsView, 0)

	for _, info := range orderInfo {
		mp, goodsOk := mpGoods[info.Number+info.Sku]

		if !goodsOk { //订单商品未拣货
			continue
		}

		pgv, ok := mpPgv[info.Number]

		if !ok {
			pgv = PickGoodsView{}
		}
		pgv.SaleNumber = info.Number
		pgv.ShopId = int64(info.ShopId)
		pgv.ShopName = info.ShopName
		pgv.Date = info.PayAt
		pgv.Remark = info.OrderRemark
		pgv.DeliveryType = info.DistributionType //配送方式
		pgv.Line = info.Line
		pgv.List = append(pgv.List, PickGoods{
			GoodsName:    mp.GoodsName,
			Sku:          mp.Sku,
			Price:        int64(info.OriginalPrice),
			GoodsSpe:     mp.GoodsSpe,
			Shelves:      mp.Shelves,
			RealOutCount: mp.ReviewNum,
			SlaveCode:    info.SaleCode,
			GoodsUnit:    info.GoodsUnit,
			SlaveUnit:    info.SaleUnit,
		})
	}

	global.SugarLogger.Info(mpPgv)

	for _, view := range mpPgv {
		//推送u8
		xml := GenU8Xml(view, view.ShopId, view.ShopName, view.HouseCode) //店铺属性中获 HouseCode
		SendShopXml(xml)
	}

	return nil
}

// 当前批次是否有接单
func IsPick(c *gin.Context) {
	var (
		form   req.EndBatchForm
		pick   batch.Pick
		status bool
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result := global.DB.Where("batch_id = ? and status = 1", form.Id).First(&pick)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if pick.Id > 0 {
		status = true
	}

	xsq_net.SucJson(c, gin.H{"status": status})
}

// 变更批次状态
func ChangeBatch(c *gin.Context) {
	//todo 把状态为0的更新为停止拣货，其他的正常操作
	var form req.StopPickForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var batches batch.Batch

	db := global.DB

	result := db.First(&batches, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//默认为更新为进行中
	updateStatus := 0

	if *form.Status == 0 {
		//如果传递过来的是进行中，则更新为暂停
		updateStatus = 2
	}

	//查询条件是传递过来的值
	result = db.Model(&batch.Batch{}).Where("id = ? and status = ?", form.Id, form.Status).Update("status", updateStatus)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 获取批次列表
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
	if form.Sku != "" || form.Number != "" || form.ShopId > 0 {

		var prePickGoods []batch.PrePickGoods

		preGoodsRes := global.DB.Model(&batch.PrePickGoods{}).
			Where(batch.PrePickGoods{
				Sku:    form.Sku,
				Number: form.Number,
				ShopId: form.ShopId,
			}).
			Select("batch_id").
			Find(&prePickGoods)

		if preGoodsRes.Error != nil {
			xsq_net.ErrorJSON(c, preGoodsRes.Error)
			return
		}

		//未找到，直接返回
		if preGoodsRes.RowsAffected == 0 {
			xsq_net.SucJson(c, res)
			return
		}

		//利用map键唯一，去重
		uMap := make(map[int]struct{}, 0)
		for _, b := range prePickGoods {
			_, ok := uMap[b.BatchId]
			if ok {
				continue
			}
			uMap[b.BatchId] = struct{}{}
			batchIdSlice = append(batchIdSlice, b.BatchId)
		}

		db = db.Where("id in (?)", batchIdSlice)
	}

	if form.Line != "" {
		db = db.Where("line like ?", form.Line+"%")
	}

	if form.CreateTime != "" {
		db = db.Where("create_time <= ?", form.CreateTime)
	}

	if form.EndTime != "" {
		db = db.Where("end_time <= ?", form.EndTime)
	}

	if *form.Status == 0 {
		db = db.Where("status in (0,2)")
	} else {
		db = db.Where("status = 1")
	}

	db.Where(&batch.Batch{DeliveryMethod: form.DeliveryMethod})

	result := db.Find(&batches)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	result = db.Scopes(model.Paginate(form.Page, form.Size)).Order("sort desc").Find(&batches)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]*rsp.Batch, 0, len(batches))
	for _, b := range batches {

		deliveryStartTime := ""
		if b.DeliveryStartTime != nil {
			deliveryStartTime = b.DeliveryStartTime.Format(timeutil.TimeFormat)
		}

		list = append(list, &rsp.Batch{
			Id:                b.Id,
			UpdateTime:        b.UpdateTime.Format(timeutil.TimeFormat),
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

// 批次池数量
func GetBatchPoolNum(c *gin.Context) {
	var (
		batchPool []rsp.BatchPoolNum
		res       rsp.GetBatchPoolNumRsp
		ongoing   int
		suspend   int
		finished  int
	)

	result := global.DB.Model(&batch.Batch{}).
		Select("count(id) as count, status").
		Group("status").
		Find(&batchPool)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, bp := range batchPool {
		switch bp.Status {
		case 0: //进行中
			ongoing = bp.Count
			break
		case 1: //已结束
			finished = bp.Count
			break
		case 2: //暂停 也属于进行中
			suspend = bp.Count
		}
	}

	res.Ongoing = ongoing + suspend
	res.Finished = finished

	xsq_net.SucJson(c, res)
}

// 预拣池基础信息
func GetBase(c *gin.Context) {

	var (
		form      req.GetBaseForm
		batchCond batch.BatchCondition
		batches   batch.Batch
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	result := db.Where("batch_id = ?", form.BatchId).First(&batchCond)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.First(&batches, form.BatchId)

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
		Status:            batches.Status,
	}

	xsq_net.SucJson(c, ret)
}

// 预拣池列表
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

	result := db.Where(batch.PrePick{ShopId: form.ShopId, Line: form.Line}).Where("status = 0").Find(&prePicks)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	db.Where(batch.PrePick{ShopId: form.ShopId, Line: form.Line}).Where("status = 0").Scopes(model.Paginate(form.Page, form.Size)).Find(&prePicks)

	for _, pick := range prePicks {
		prePickIds = append(prePickIds, pick.Id)
	}

	retCount := []rsp.Ret{}

	result = db.Model(&batch.PrePickGoods{}).
		Select("SUM(out_count) as out_c, SUM(need_num) AS need_c, shop_id, goods_type").
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

	list := make([]*rsp.PrePick, 0, len(prePicks))

	for _, pick := range prePicks {
		list = append(list, &rsp.PrePick{
			Id:           pick.Id,
			ShopCode:     pick.ShopCode,
			ShopName:     pick.ShopName,
			Line:         pick.Line,
			Status:       pick.Status,
			CategoryInfo: typeMap[pick.ShopId],
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)

}

// 预拣货明细
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
		prePick       batch.PrePick
		prePickGoods  []batch.PrePickGoods
		prePickRemark []batch.PrePickRemark
	)

	db := global.DB

	result := db.First(&prePick, form.PrePickId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.TaskName = prePick.ShopName
	res.OrderNum = prePick.OrderNum
	res.Line = prePick.Line

	result = db.Where("pre_pick_id = ? and status = 0", form.PrePickId).Find(&prePickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	goodsMap := make(map[string][]rsp.PrePickGoods, 0)

	//商品数
	goodsNum := 0
	//订单数map
	orderMp := make(map[string]struct{}, 0)

	for _, goods := range prePickGoods {
		orderMp[goods.Number] = struct{}{}
		goodsNum += goods.NeedNum
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

	//商品数
	res.GoodsNum = goodsNum
	//订单数
	res.OrderNum = len(orderMp)

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

// 置顶
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

// 批次池内单数量
func GetPoolNum(c *gin.Context) {
	var (
		form         req.GetPoolNumReq
		res          rsp.GetPoolNumRsp
		count        int64
		poolNumCount []rsp.PoolNumCount
		pickNum,
		toReviewNum,
		completeNum int
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	db := global.DB

	result := db.Model(&batch.PrePick{}).Select("id").Where("batch_id = ? and status = 0", form.BatchId).Count(&count)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&batch.Pick{}).
		Select("count(id) as count, status").
		Where("batch_id = ?", form.BatchId).
		Group("status").
		Find(&poolNumCount)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, pc := range poolNumCount {
		switch pc.Status {
		case 0: //待拣货
			pickNum = pc.Count
			break
		case 1: //待复核
			toReviewNum = pc.Count
			break
		case 2: //已完成
			completeNum = pc.Count
		}
	}

	res = rsp.GetPoolNumRsp{
		PrePickNum:  count,
		PickNum:     pickNum,
		ToReviewNum: toReviewNum,
		CompleteNum: completeNum,
	}

	xsq_net.SucJson(c, res)
}

// 批量拣货
func BatchPick(c *gin.Context) {
	var (
		form req.BatchPickForm
		err  error
	)

	if err = c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var batches batch.Batch

	result := global.DB.First(&batches, form.BatchId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if batches.Status == 1 { //状态:0:进行中,1:已结束,2:暂停
		xsq_net.ErrorJSON(c, errors.New("请先开启拣货"))
		return
	}

	form.WarehouseId = c.GetInt("warehouseId")

	switch form.Type {
	case 1:
		err = BatchPickByParams(form)
		break
	case 2:
		err = BatchPickByParams(form)
		break
	case 3:
		err = BatchPickByParams(form)
		break
	}

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

func BatchPickAll(form req.BatchPickForm) error {
	db := global.DB
	var (
		prePick        []batch.PrePick
		prePickGoods   []batch.PrePickGoods
		prePickRemarks []batch.PrePickRemark
	)

	//0:未处理,1:已进入拣货池
	result := db.Where("id in (?) and status = 0", form.Ids).Find(&prePick)

	if result.Error != nil {
		return result.Error
	}

	//0:未处理,1:已进入拣货池
	result = db.Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickGoods)

	if result.Error != nil {
		return result.Error
	}

	prePickGoodsMap := make(map[int][]batch.PrePickGoods, 0)

	//查到的 prePickGoods prePickRemarks 表id数据，更新状态使用
	var (
		prePickGoodsIds   []int
		prePickRemarksIds []int
	)

	for _, goods := range prePickGoods {
		prePickGoodsIds = append(prePickGoodsIds, goods.Id)
		prePickGoodsMap[goods.PrePickId] = append(prePickGoodsMap[goods.PrePickId], goods)
	}

	result = db.Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickRemarks)

	if result.Error != nil {
		return result.Error
	}

	prePickRemarksMap := make(map[int][]batch.PrePickRemark, 0)

	for _, remark := range prePickRemarks {
		prePickRemarksIds = append(prePickRemarksIds, remark.Id)
		prePickRemarksMap[remark.PrePickId] = append(prePickRemarksMap[remark.PrePickId], remark)
	}

	var (
		pickGoods  []batch.PickGoods
		pickRemark []batch.PickRemark
		pickNums   []rsp.PickNums
	)

	//拣货池 订单、门店、需拣 数量
	result = db.Model(&batch.PrePickGoods{}).
		Select("pre_pick_id,count(DISTINCT(number)) as order_num,count(DISTINCT(shop_id)) as shop_num,sum(need_num) as need_num").
		Where("pre_pick_id in (?) and status = 0", form.Ids).
		Group("pre_pick_id").
		Find(&pickNums)

	if result.Error != nil {
		return result.Error
	}

	//拣货池 订单、门店、需拣 数量 mp
	pickNumsMp := make(map[int]rsp.PickNums, 0)

	for _, nums := range pickNums {
		pickNumsMp[nums.PrePickId] = nums
	}

	tx := db.Begin()

	for _, pre := range prePick {

		//预拣池商品中未找到相关数据
		_, pgMpOk := prePickGoodsMap[pre.Id]

		if !pgMpOk {
			continue
		}

		var (
			shopNum  = 0
			orderNum = 0
			needNum  = 0
		)

		val, ok := pickNumsMp[pre.Id]

		if ok {
			shopNum = val.ShopNum
			orderNum = val.OrderNum
			needNum = val.NeedNum
		}

		pick := batch.Pick{
			WarehouseId:    form.WarehouseId,
			BatchId:        pre.BatchId,
			PrePickIds:     strconv.Itoa(pre.Id),
			TaskName:       pre.ShopName,
			ShopCode:       pre.ShopCode,
			ShopName:       pre.ShopName,
			Line:           pre.Line,
			ShopNum:        shopNum,
			OrderNum:       orderNum,
			NeedNum:        needNum,
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

		for _, goods := range prePickGoodsMap[pre.Id] {
			pickGoods = append(pickGoods, batch.PickGoods{
				WarehouseId:    form.WarehouseId,
				BatchId:        pre.BatchId,
				PickId:         pick.Id,
				PrePickGoodsId: goods.Id,
				OrderInfoId:    goods.OrderInfoId,
				GoodsName:      goods.GoodsName,
				GoodsSpe:       goods.GoodsSpe,
				Shelves:        goods.Shelves,
				NeedNum:        goods.NeedNum,
				Number:         goods.Number,
				ShopId:         goods.ShopId,
				GoodsType:      goods.GoodsType,
			})
		}

		result = tx.Save(&pickGoods)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		for _, remark := range prePickRemarksMap[pre.Id] {
			pickRemark = append(pickRemark, batch.PickRemark{
				WarehouseId:     form.WarehouseId,
				BatchId:         pre.BatchId,
				PickId:          pick.Id,
				PrePickRemarkId: remark.Id,
				OrderInfoId:     remark.OrderInfoId,
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
				return result.Error
			}
		}
	}

	//批量更新
	result = tx.Model(batch.PrePick{}).Where("id in (?) and status = 0", form.Ids).Updates(map[string]interface{}{"status": 1})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	//通过ID更新
	result = tx.Model(batch.PrePickGoods{}).Where("id in (?)", prePickGoodsIds).Updates(map[string]interface{}{"status": 1})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	//通过ID更新
	result = tx.Model(batch.PrePickRemark{}).Where("id in (?)", prePickRemarksIds).Updates(map[string]interface{}{"status": 1})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()

	return nil
}

// 批量拣货 - 根据参数类型
func BatchPickByParams(form req.BatchPickForm) error {

	db := global.DB
	var (
		prePick        []batch.PrePick
		prePickGoods   []batch.PrePickGoods
		prePickRemarks []batch.PrePickRemark
		pickNums       []rsp.PickNums
	)

	//0:未处理,1:已进入拣货池
	result := db.Where("id in (?) and status = 0", form.Ids).Find(&prePick)

	if result.Error != nil {
		return result.Error
	}

	local := db.Where("pre_pick_id in (?) and status = 0", form.Ids)

	//计算拣货池 订单、门店、需拣 数量 sql 拼接
	numCountLocal := db.Model(&batch.PrePickGoods{}).
		Select("pre_pick_id,count(DISTINCT(number)) as order_num,count(DISTINCT(shop_id)) as shop_num,sum(need_num) as need_num").
		Where("pre_pick_id in (?) and status = 0", form.Ids)

	if form.Type == 2 { //按分类
		local.Where("goods_type in (?)", form.TypeParam)
		numCountLocal.Where("goods_type in (?)", form.TypeParam)
	} else if form.Type == 3 { //按商品
		local.Where("sku in (?)", form.TypeParam)
		numCountLocal.Where("sku in (?)", form.TypeParam)
	}

	result = local.Find(&prePickGoods)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("对应的拣货池商品不存在")
	}

	prePickGoodsMap := make(map[int][]batch.PrePickGoods, 0)

	for _, goods := range prePickGoods {
		prePickGoodsMap[goods.PrePickId] = append(prePickGoodsMap[goods.PrePickId], goods)
	}

	//拣货池 订单、门店、需拣 数量
	result = numCountLocal.Group("pre_pick_id").Find(&pickNums)

	if result.Error != nil {
		return result.Error
	}

	//拣货池 订单、门店、需拣 数量 mp
	pickNumsMp := make(map[int]rsp.PickNums, 0)

	for _, nums := range pickNums {
		pickNumsMp[nums.PrePickId] = nums
	}

	tx := db.Begin()

	var (
		prePickGoodsIds   []int
		prePickRemarksIds []int
		pickGoods         []batch.PickGoods
		pickRemark        []batch.PickRemark
	)

	for _, pre := range prePick {

		//预拣池商品中未找到相关数据
		_, pgMpOk := prePickGoodsMap[pre.Id]

		if !pgMpOk {
			continue
		}

		var (
			shopNum  = 0
			orderNum = 0
			needNum  = 0
		)

		val, ok := pickNumsMp[pre.Id]

		if ok {
			shopNum = val.ShopNum
			orderNum = val.OrderNum
			needNum = val.NeedNum
		}

		pick := batch.Pick{
			WarehouseId:    form.WarehouseId,
			BatchId:        pre.BatchId,
			PrePickIds:     strconv.Itoa(pre.Id),
			TaskName:       pre.ShopName,
			ShopCode:       pre.ShopCode,
			ShopName:       pre.ShopName,
			Line:           pre.Line,
			ShopNum:        shopNum,
			OrderNum:       orderNum,
			NeedNum:        needNum,
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

		var orderInfoIds []int

		for _, goods := range prePickGoodsMap[pre.Id] {

			orderInfoIds = append(orderInfoIds, goods.OrderInfoId)

			//更新 prePickGoods 使用
			prePickGoodsIds = append(prePickGoodsIds, goods.Id)

			pickGoods = append(pickGoods, batch.PickGoods{
				WarehouseId:    form.WarehouseId,
				BatchId:        pre.BatchId,
				PickId:         pick.Id,
				PrePickGoodsId: goods.Id,
				OrderInfoId:    goods.OrderInfoId,
				GoodsName:      goods.GoodsName,
				GoodsSpe:       goods.GoodsSpe,
				Shelves:        goods.Shelves,
				NeedNum:        goods.NeedNum,
				Number:         goods.Number,
				ShopId:         goods.ShopId,
				GoodsType:      goods.GoodsType,
			})
		}

		result = db.Where("order_info_id in (?)", orderInfoIds).Find(&prePickRemarks)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		for _, remark := range prePickRemarks {
			//更新 prePickRemarks 使用
			prePickRemarksIds = append(prePickRemarksIds, remark.Id)

			pickRemark = append(pickRemark, batch.PickRemark{
				WarehouseId:     form.WarehouseId,
				BatchId:         pre.BatchId,
				PickId:          pick.Id,
				PrePickRemarkId: remark.Id,
				OrderInfoId:     remark.OrderInfoId,
				Number:          remark.Number,
				OrderRemark:     remark.OrderRemark,
				GoodsRemark:     remark.GoodsRemark,
				ShopName:        remark.ShopName,
				Line:            remark.Line,
			})
		}
	}

	//商品数据保存
	result = tx.Save(&pickGoods)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	//订单备注数据
	if len(pickRemark) > 0 {
		result = tx.Save(&pickRemark)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	//更新预拣池商品表的商品数据状态
	if len(prePickGoodsIds) > 0 {
		result = tx.Model(batch.PrePickGoods{}).Where("id in (?)", prePickGoodsIds).Updates(map[string]interface{}{"status": 1})

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	//预拣池内商品全部进入拣货池时 更新 对应的 预拣池状态
	if form.Type == 1 { //全单拣货
		result = tx.Model(batch.PrePick{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"status": 1})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	} else {
		//0:未处理,1:已进入拣货池
		result = tx.Model(&batch.PrePickGoods{}).Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickGoods)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		//将传过来的id转换成map
		idsMp := make(map[int]struct{}, 0)

		for _, id := range form.Ids {
			idsMp[id] = struct{}{}
		}

		//去除未处理的预拣池id
		for _, good := range prePickGoods {
			delete(idsMp, good.PrePickId)
		}

		//将map转回切片
		prePickIds := []int{}
		for id, _ := range idsMp {
			prePickIds = append(prePickIds, id)
		}

		if len(prePickIds) > 0 {
			result = tx.Model(batch.PrePick{}).Where("id in (?)", prePickIds).Updates(map[string]interface{}{"status": 1})
			if result.Error != nil {
				tx.Rollback()
				return result.Error
			}
		}
	}

	//更新预拣池商品备注表的数据状态
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

// 批量拣货-按照分类
func BatchPickByClassification(form req.BatchPickForm) error {
	db := global.DB
	var (
		prePick        []batch.PrePick
		prePickGoods   []batch.PrePickGoods
		prePickRemarks []batch.PrePickRemark
	)

	result := db.Where("id in (?) and status = 0", form.Ids).Find(&prePick)

	if result.Error != nil {
		return result.Error
	}

	result = db.Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickGoods)

	if result.Error != nil {
		return result.Error
	}

	prePickGoodsMap := make(map[int][]batch.PrePickGoods, 0)

	for _, goods := range prePickGoods {
		prePickGoodsMap[goods.PrePickId] = append(prePickGoodsMap[goods.PrePickId], goods)
	}

	var (
		pickNums []rsp.PickNums
	)

	mp := make(map[string]struct{}, 0)

	for _, tp := range form.TypeParam {
		mp[tp] = struct{}{}
	}

	//拣货池 订单、门店、需拣 数量
	result = db.Model(&batch.PrePickGoods{}).
		Select("pre_pick_id,count(DISTINCT(number)) as order_num,count(DISTINCT(shop_id)) as shop_num,sum(need_num) as need_num").
		Where("pre_pick_id in (?) and status = 0", form.Ids).
		Where("goods_type in (?)", form.TypeParam).
		Group("pre_pick_id").
		Find(&pickNums)

	if result.Error != nil {
		return result.Error
	}

	//拣货池 订单、门店、需拣 数量 mp
	pickNumsMp := make(map[int]rsp.PickNums, 0)

	for _, nums := range pickNums {
		pickNumsMp[nums.PrePickId] = nums
	}

	tx := db.Begin()

	var prePickIds,
		prePickGoodsIds,
		prePickRemarksIds []int

	for _, pre := range prePick {

		var (
			shopNum    = 0
			orderNum   = 0
			needNum    = 0
			pickGoods  []batch.PickGoods
			pickRemark []batch.PickRemark
		)

		val, ok := pickNumsMp[pre.Id]

		if ok {
			shopNum = val.ShopNum
			orderNum = val.OrderNum
			needNum = val.NeedNum
		}

		pick := batch.Pick{
			WarehouseId:    form.WarehouseId,
			BatchId:        pre.BatchId,
			PrePickIds:     strconv.Itoa(pre.Id),
			TaskName:       pre.ShopName,
			ShopCode:       pre.ShopCode,
			ShopName:       pre.ShopName,
			Line:           pre.Line,
			ShopNum:        shopNum,
			OrderNum:       orderNum,
			NeedNum:        needNum,
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

		var orderInfoIds []int

		for _, goods := range prePickGoodsMap[pre.Id] {
			//类型不在mp中
			_, ok := mp[goods.GoodsType]
			if !ok {
				continue
			}

			orderInfoIds = append(orderInfoIds, goods.OrderInfoId)

			//更新 prePickGoods 使用
			prePickGoodsIds = append(prePickGoodsIds, goods.Id)

			pickGoods = append(pickGoods, batch.PickGoods{
				WarehouseId:    form.WarehouseId,
				BatchId:        pre.BatchId,
				PickId:         pick.Id,
				PrePickGoodsId: goods.Id,
				GoodsName:      goods.GoodsName,
				GoodsSpe:       goods.GoodsSpe,
				Shelves:        goods.Shelves,
				NeedNum:        goods.NeedNum,
				Number:         goods.Number,
				ShopId:         goods.ShopId,
				GoodsType:      goods.GoodsType,
			})
		}

		if len(pickGoods) == 0 { //对应的类型商品数据不存在
			continue
		}

		result = tx.Save(&pickGoods)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		result = db.Where("order_info_id in (?)", orderInfoIds).Find(&prePickRemarks)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		for _, remark := range prePickRemarks {
			//更新 prePickRemarks 使用
			prePickRemarksIds = append(prePickRemarksIds, remark.Id)
			pickRemark = append(pickRemark, batch.PickRemark{
				WarehouseId:     form.WarehouseId,
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
				return result.Error
			}
		}

		if len(prePickGoodsMap[pre.Id]) == len(pickGoods) {
			//当前类型的商品数据 总条数 和 当前类型存入拣货池的商品数据条数相等 批量更新
			prePickIds = append(prePickIds, pre.Id)
		}
	}

	if len(prePickIds) > 0 {
		result = tx.Model(batch.PrePick{}).Where("id in (?)", prePickIds).Updates(map[string]interface{}{"status": 1})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

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

func BatchPickByGoods(form req.BatchPickForm) error {
	db := global.DB
	var (
		prePick        []batch.PrePick
		prePickGoods   []batch.PrePickGoods
		prePickRemarks []batch.PrePickRemark
	)

	result := db.Where("id in (?) and status = 0", form.Ids).Find(&prePick)

	if result.Error != nil {
		return result.Error
	}

	result = db.Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickGoods)

	if result.Error != nil {
		return result.Error
	}

	prePickGoodsMap := make(map[int][]batch.PrePickGoods, 0)

	for _, goods := range prePickGoods {
		prePickGoodsMap[goods.PrePickId] = append(prePickGoodsMap[goods.PrePickId], goods)
	}

	var (
		pickNums []rsp.PickNums
	)

	mp := make(map[string]struct{}, 0)

	for _, tp := range form.TypeParam {
		mp[tp] = struct{}{}
	}

	//拣货池 订单、门店、需拣 数量
	result = db.Model(&batch.PrePickGoods{}).
		Select("pre_pick_id,count(DISTINCT(number)) as order_num,count(DISTINCT(shop_id)) as shop_num,sum(need_num) as need_num").
		Where("pre_pick_id in (?) and status = 0", form.Ids).
		Where("goods_name in (?)", form.TypeParam).
		Group("pre_pick_id").
		Find(&pickNums)

	if result.Error != nil {
		return result.Error
	}

	//拣货池 订单、门店、需拣 数量 mp
	pickNumsMp := make(map[int]rsp.PickNums, 0)

	for _, nums := range pickNums {
		pickNumsMp[nums.PrePickId] = nums
	}

	tx := db.Begin()

	var prePickIds,
		prePickGoodsIds,
		prePickRemarksIds []int

	for _, pre := range prePick {

		var (
			shopNum    = 0
			orderNum   = 0
			needNum    = 0
			pickGoods  []batch.PickGoods
			pickRemark []batch.PickRemark
		)

		val, ok := pickNumsMp[pre.Id]

		if ok {
			shopNum = val.ShopNum
			orderNum = val.OrderNum
			needNum = val.NeedNum
		}

		pick := batch.Pick{
			WarehouseId:    form.WarehouseId,
			BatchId:        pre.BatchId,
			PrePickIds:     strconv.Itoa(pre.Id),
			TaskName:       pre.ShopName,
			ShopCode:       pre.ShopCode,
			ShopName:       pre.ShopName,
			Line:           pre.Line,
			ShopNum:        shopNum,
			OrderNum:       orderNum,
			NeedNum:        needNum,
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

		var orderInfoIds []int

		for _, goods := range prePickGoodsMap[pre.Id] {
			//商品不在mp中
			_, ok := mp[goods.GoodsName]
			if !ok {
				continue
			}

			orderInfoIds = append(orderInfoIds, goods.OrderInfoId)

			//更新prePickGoods表状态使用
			prePickGoodsIds = append(prePickGoodsIds, goods.Id)

			pickGoods = append(pickGoods, batch.PickGoods{
				WarehouseId:    form.WarehouseId,
				BatchId:        pre.BatchId,
				PickId:         pick.Id,
				PrePickGoodsId: goods.Id,
				GoodsName:      goods.GoodsName,
				GoodsSpe:       goods.GoodsSpe,
				Shelves:        goods.Shelves,
				NeedNum:        goods.NeedNum,
				Number:         goods.Number,
				ShopId:         goods.ShopId,
				GoodsType:      goods.GoodsType,
			})
		}

		if len(pickGoods) == 0 { //对应的类型商品数据不存在
			continue
		}

		result = tx.Save(&pickGoods)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		result = db.Where("order_info_id in (?)", orderInfoIds).Find(&prePickRemarks)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		for _, remark := range prePickRemarks {
			prePickRemarksIds = append(prePickRemarksIds, remark.Id)
			pickRemark = append(pickRemark, batch.PickRemark{
				WarehouseId:     form.WarehouseId,
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
				return result.Error
			}
		}

		if len(prePickGoodsMap[pre.Id]) == len(pickGoods) {
			//当前类型的商品数据 总条数 和 当前类型存入拣货池的商品数据条数相等 批量更新
			prePickIds = append(prePickIds, pre.Id)
		}
	}

	if len(prePickIds) > 0 {
		result = tx.Model(batch.PrePick{}).Where("id in (?)", prePickIds).Updates(map[string]interface{}{"status": 1})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

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

// 合并拣货
func MergePick(c *gin.Context) {
	var form req.MergePickForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var batches batch.Batch

	result := global.DB.First(&batches, form.BatchId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if batches.Status == 1 {
		xsq_net.ErrorJSON(c, errors.New("请先开启拣货"))
		return
	}

	form.WarehouseId = c.GetInt("warehouseId")

	var err error

	switch form.Type {
	case 1:
		err = MergePickByParams(form)
		break
	case 2:
		err = MergePickByParams(form)
		break
	case 3:
		err = MergePickByParams(form)
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

func MergePickByParams(form req.MergePickForm) error {
	var (
		prePickGoods   []batch.PrePickGoods
		prePickRemarks []batch.PrePickRemark
		pickGoods      []batch.PickGoods
		pickRemarks    []batch.PickRemark
		pickNums       rsp.MergePickNums
	)

	db := global.DB

	var (
		prePickIds string
		prePickGoodsIds,
		orderInfoIds,
		prePickRemarksIds []int
	)

	local := db.Where("pre_pick_id in (?) and status = 0", form.Ids)

	//计算拣货池 订单、门店、需拣 数量 sql 拼接
	numCountLocal := db.Model(&batch.PrePickGoods{}).
		Select("pre_pick_id,count(DISTINCT(number)) as order_num,count(DISTINCT(shop_id)) as shop_num,sum(need_num) as need_num").
		Where("pre_pick_id in (?) and status = 0", form.Ids)

	if form.Type == 2 { //按分类
		local.Where("goods_type in (?)", form.TypeParam)
		numCountLocal.Where("goods_type in (?)", form.TypeParam)
	} else if form.Type == 3 { //按商品
		local.Where("sku in (?)", form.TypeParam)
		numCountLocal.Where("sku in (?)", form.TypeParam)
	}

	result := local.Find(&prePickGoods)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("商品数据未找到")
	}

	//拣货池 订单、门店、需拣 数量
	result = numCountLocal.Find(&pickNums)

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
		ShopNum:        pickNums.ShopNum,
		OrderNum:       pickNums.OrderNum,
		NeedNum:        pickNums.NeedNum,
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

		prePickGoodsIds = append(prePickGoodsIds, goods.Id)

		orderInfoIds = append(orderInfoIds, goods.Id)

		pickGoods = append(pickGoods, batch.PickGoods{
			WarehouseId:    form.WarehouseId,
			BatchId:        form.BatchId,
			PickId:         pick.Id,
			PrePickGoodsId: goods.Id,
			OrderInfoId:    goods.OrderInfoId,
			GoodsName:      goods.GoodsName,
			GoodsType:      goods.GoodsType,
			GoodsSpe:       goods.GoodsSpe,
			Shelves:        goods.Shelves,
			NeedNum:        goods.NeedNum,
			Number:         goods.Number,
			ShopId:         goods.ShopId,
		})
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

	//预拣池内商品全部进入拣货池时 更新 对应的 预拣池状态
	if form.Type == 1 { //全单拣货
		result = tx.Model(batch.PrePick{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"status": 1})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	} else {
		//0:未处理,1:已进入拣货池
		result = tx.Model(batch.PrePickGoods{}).Where("pre_pick_id in (?) and status = 0", form.Ids).Find(&prePickGoods)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		//将传过来的id转换成map
		idsMp := make(map[int]struct{}, 0)

		for _, id := range form.Ids {
			idsMp[id] = struct{}{}
		}

		//去除未处理的预拣池id
		for _, good := range prePickGoods {
			delete(idsMp, good.PrePickId)
		}

		//将map转回切片
		prePickIdSlice := []int{}
		for id, _ := range idsMp {
			prePickIdSlice = append(prePickIdSlice, id)
		}

		if len(prePickIdSlice) > 0 {
			result = tx.Model(batch.PrePick{}).Where("id in (?)", prePickIdSlice).Updates(map[string]interface{}{"status": 1})
			if result.Error != nil {
				tx.Rollback()
				return result.Error
			}
		}
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
				OrderInfoId:     remark.OrderInfoId,
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

	tx.Commit()

	return nil
}

// 全单拣货
func ByAllOrder(form req.MergePickForm) error {
	var (
		prePickGoods   []batch.PrePickGoods
		prePickRemarks []batch.PrePickRemark
		pickGoods      []batch.PickGoods
		pickRemarks    []batch.PickRemark
		pickNums       rsp.PickNums
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

	//拣货池 订单、门店、需拣 数量
	result = db.Model(&batch.PrePickGoods{}).
		Select("pre_pick_id,count(DISTINCT(number)) as order_num,count(DISTINCT(shop_id)) as shop_num,sum(need_num) as need_num").
		Where("pre_pick_id in (?)", form.Ids).
		Find(&pickNums)

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
		ShopNum:        pickNums.ShopNum,
		OrderNum:       pickNums.OrderNum,
		NeedNum:        pickNums.NeedNum,
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
		prePickGoodsIds,
		prePickRemarksIds []int
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
			GoodsType:      goods.GoodsType,
			GoodsSpe:       goods.GoodsSpe,
			Shelves:        goods.Shelves,
			NeedNum:        goods.NeedNum,
			Number:         goods.Number,
			ShopId:         goods.ShopId,
		})
	}

	//验证
	if len(pickGoods) == 0 {
		tx.Rollback()
		return errors.New("没有未拣的商品")
	}

	result = tx.Save(&pickGoods)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if len(prePickRemarks) > 0 {
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

// 按分类拣货
func ByClassification(form req.MergePickForm) error {
	var (
		prePickGoods   []batch.PrePickGoods
		prePickRemarks []batch.PrePickRemark
		pickGoods      []batch.PickGoods
		pickRemarks    []batch.PickRemark
		pickNums       rsp.MergePickNums
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
		return errors.New("仓库分类:" + strings.Join(form.TypeParam, ",") + "数据未找到")
	}

	//拣货池 订单、门店、需拣 数量
	result = db.Model(&batch.PrePickGoods{}).
		Select("pre_pick_id,count(DISTINCT(number)) as order_num,count(DISTINCT(shop_id)) as shop_num,sum(need_num) as need_num").
		Where("pre_pick_id in (?)", form.Ids).
		Where("goods_type in (?)", form.TypeParam).
		Find(&pickNums)

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
		ShopNum:        pickNums.ShopNum,
		OrderNum:       pickNums.OrderNum,
		NeedNum:        pickNums.NeedNum,
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

	mp := make(map[string]struct{}, 0)

	for _, tp := range form.TypeParam {
		mp[tp] = struct{}{}
	}

	for _, goods := range prePickGoods {

		_, ok := mp[goods.GoodsType]
		//只保留相关仓库类型数据
		if !ok {
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
			GoodsType:      goods.GoodsType,
			GoodsSpe:       goods.GoodsSpe,
			Shelves:        goods.Shelves,
			NeedNum:        goods.NeedNum,
			Number:         goods.Number,
			ShopId:         goods.ShopId,
		})
	}

	//验证
	if len(pickGoods) == 0 {
		tx.Rollback()
		return errors.New("没有未拣的商品")
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

// 按单品拣货
func ByGoods(form req.MergePickForm) error {
	var (
		prePickGoods   []batch.PrePickGoods
		prePickRemarks []batch.PrePickRemark
		pickGoods      []batch.PickGoods
		pickRemarks    []batch.PickRemark
		pickNums       rsp.MergePickNums
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
		return errors.New("仓库分类:" + strings.Join(form.TypeParam, ",") + "数据未找到")
	}

	//拣货池 订单、门店、需拣 数量
	result = db.Model(&batch.PrePickGoods{}).
		Select("pre_pick_id,count(DISTINCT(number)) as order_num,count(DISTINCT(shop_id)) as shop_num,sum(need_num) as need_num").
		Where("pre_pick_id in (?)", form.Ids).
		Where("goods_name in (?)", form.TypeParam).
		Find(&pickNums)

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
		ShopNum:        pickNums.ShopNum,
		OrderNum:       pickNums.OrderNum,
		NeedNum:        pickNums.NeedNum,
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

	mp := make(map[string]struct{}, 0)

	for _, tp := range form.TypeParam {
		mp[tp] = struct{}{}
	}

	for _, goods := range prePickGoods {

		_, ok := mp[goods.GoodsName]

		//只保留相关商品数据
		if !ok {
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
			GoodsType:      goods.GoodsType,
			GoodsSpe:       goods.GoodsSpe,
			Shelves:        goods.Shelves,
			NeedNum:        goods.NeedNum,
			Number:         goods.Number,
			ShopId:         goods.ShopId,
		})
	}

	//验证
	if len(pickGoods) == 0 {
		tx.Rollback()
		return errors.New("没有未拣的商品")
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

// 打印
func PrintCallGet(c *gin.Context) {
	var (
		form req.PrintCallGetReq
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	printCh := GetPrintJobMap(form.HouseCode)

	//通道中没有任务
	if printCh == nil {
		xsq_net.SucJson(c, nil)
		return
	}

	global.SugarLogger.Infof("%+v", printCh)

	var (
		pick       batch.Pick
		pickGoods  []batch.PickGoods
		orderInfos []order.OrderInfo
	)

	db := global.DB

	result := db.Model(&batch.Pick{}).Where("delivery_order_no = ?", printCh.DeliveryOrderNo).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&batch.PickGoods{}).Where("pick_id = ? and shop_id = ?", pick.Id, printCh.ShopId).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	length := len(pickGoods) //有多少条pickGoods就有多少条OrderInfo数据，map数也是

	orderInfoIds := make([]int, 0, length)

	goodsMp := make(map[int]batch.PickGoods, length)

	for _, good := range pickGoods {
		orderInfoIds = append(orderInfoIds, good.OrderInfoId)

		goodsMp[good.OrderInfoId] = good
	}

	result = db.Model(&order.OrderInfo{}).Where("id in (?)", orderInfoIds).Find(&orderInfos)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if len(orderInfos) <= 0 {
		xsq_net.ErrorJSON(c, ecode.OrderDataNotFound)
		return
	}

	item := rsp.PrintCallGetRsp{
		ShopName:    pick.ShopName,
		JHNumber:    strconv.Itoa(pick.Id),
		PickName:    pick.PickUser, //拣货人
		ShopType:    orderInfos[0].ShopType,
		CheckName:   pick.ReviewUser,                                          //复核员
		HouseName:   TransferHouse(orderInfos[0].HouseCode),                   //TransferHouse(info.HouseCode)
		Delivery:    TransferDistributionType(orderInfos[0].DistributionType), //TransferDistributionType(info.DistributionType),
		OrderRemark: orderInfos[0].OrderRemark,
		Consignee:   orderInfos[0].ConsigneeName, //info.ConsigneeName
		Shop_code:   pick.ShopCode,
		Packages:    0,
		Phone:       orderInfos[0].ConsigneeTel, //info.ConsigneeTel,
		PriType:     1,
	}

	if orderInfos[0].ShopCode != "" {
		item.ShopName = orderInfos[0].ShopCode + "--" + orderInfos[0].ShopName
	}

	item2 := rsp.CallGetGoodsView{
		SaleNumber:  orderInfos[0].Number,
		Date:        orderInfos[0].PayAt,
		OrderRemark: orderInfos[0].OrderRemark,
	}

	for _, info := range orderInfos {

		pgs, ok := goodsMp[info.Id]

		if !ok {
			continue
		}

		item3 := rsp.CallGetGoods{
			GoodsName:    info.Name,
			GoodsSpe:     info.GoodsSpe,
			GoodsCount:   info.PayCount,
			RealOutCount: pgs.ReviewNum,
			GoodsUnit:    info.GoodsUnit,
			Price:        int64(info.DiscountPrice) * int64(pgs.ReviewNum),
			LackCount:    info.PayCount - pgs.ReviewNum,
		}
		item2.List = append(item2.List, item3)
	}

	item.GoodsList = append(item.GoodsList, item2)

	ret := make([]rsp.PrintCallGetRsp, 0, 1)

	ret = append(ret, item)

	xsq_net.SucJson(c, ret)
}

func TransferHouse(s string) string {
	switch s {
	case constant.JH_HUOSE_CODE:
		return constant.JH_HUOSE_NAME
	default:
		return constant.OT_HUOSE_NAME
	}
}

func TransferDistributionType(t int) (method string) {
	switch t {
	case 1:
		method = "公司配送"
		break
	case 2:
		method = "用户自提"
		break
	case 3:
		method = "三方物流"
		break
	case 4:
		method = "快递配送"
		break
	case 5:
		method = "首批物料|设备单"
		break
	default:
		method = "其他"
		break
	}

	return method
}

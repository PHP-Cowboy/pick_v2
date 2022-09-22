package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"pick_v2/common/constant"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/slice"
	"pick_v2/utils/timeutil"
	"time"

	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/utils/ecode"
	"pick_v2/utils/request"
	"pick_v2/utils/xsq_net"
)

// 订单列表
func GetGoodsList(c *gin.Context) {
	var form req.GoodsListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		orders     []model.Order
		orderGoods []model.OrderGoods
		numbers    []string
	)

	db := global.DB

	if form.Sku != "" {
		result := db.Where("sku = ?", form.Sku).Find(&orderGoods)

		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, good := range orderGoods {
			numbers = append(numbers, good.Number)
		}
	}

	localDb := db.Model(&model.Order{})

	if len(numbers) > 0 {
		localDb = localDb.Where("number in (?)", numbers)
	}

	localDb.Where(&model.Order{
		OrderType:        form.OrderType,
		ShopId:           form.ShopId,
		Number:           form.Number,
		Line:             form.Lines,
		DistributionType: form.DistributionType,
		ShopType:         form.ShopType,
		Province:         form.Province,
		City:             form.City,
		District:         form.District,
		HasRemark:        form.HasRemark,
	})

	if form.PayEndTime != "" {
		localDb = localDb.Where("pay_at >= ?", form.PayEndTime)
	}

	localDb.Where("order_type != 2") //不要拣货中的

	var (
		res   rsp.GoodsListRsp
		dbRes []OrderNum
		total int64
	)

	result := localDb.Count(&total)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	localDb.Order(fmt.Sprintf("pay_at %s", form.PayAtSort))

	if form.ShopCodeSort != "" {
		localDb.Order(fmt.Sprintf("shop_code  %s", form.ShopCodeSort))
	}

	result = localDb.Scopes(model.Paginate(form.Page, form.Size)).Find(&orders)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.Order, 0, form.Size)

	for _, o := range orders {
		latestPickingTime := ""

		if o.LatestPickingTime != nil {
			latestPickingTime = o.LatestPickingTime.Format(timeutil.TimeFormat)
		}

		payAt := ""

		if o.PayAt != "" {

			at, payAtErr := time.ParseInLocation(timeutil.TimeZoneFormat, o.PayAt, time.Local)

			if payAtErr != nil {
				xsq_net.ErrorJSON(c, ecode.DataTransformationError)
				return
			}
			payAt = at.Format(timeutil.TimeFormat)
		}

		list = append(list, rsp.Order{
			Number:            o.Number,
			PayAt:             payAt,
			ShopCode:          o.ShopCode,
			ShopName:          o.ShopName,
			ShopType:          o.ShopType,
			DistributionType:  o.DistributionType,
			PayCount:          o.PayTotal,
			Picked:            o.Picked,
			UnPicked:          o.UnPicked,
			CloseNum:          o.CloseNum,
			Line:              o.Line,
			Region:            o.Province + o.City + o.District,
			OrderRemark:       o.OrderRemark,
			OrderType:         o.OrderType,
			LatestPickingTime: latestPickingTime,
		})
	}

	result = localDb.Model(&model.Order{}).Select("count(id) as count, order_type").Group("order_type").Find(&dbRes)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// 订单明细
func GetOrderDetail(c *gin.Context) {

	var form req.GetOrderDetailForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		orders     model.Order
		orderGoods []model.OrderGoods
	)

	db := global.DB

	result := db.Model(&model.Order{}).Where("number = ?", form.Number).First(&orders)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	payAt, payAtErr := time.ParseInLocation(timeutil.TimeZoneFormat, orders.PayAt, time.Local)

	if payAtErr != nil {
		xsq_net.ErrorJSON(c, ecode.DataTransformationError)
		return
	}

	result = db.Model(&model.OrderGoods{}).Where("number = ?", form.Number).Find(&orderGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	mp, err := cache.GetClassification()

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var res rsp.OrderDetail

	detailMap := make(map[string]*rsp.Detail, 0)

	deliveryOrderNoArr := make(model.GormList, 0)

	for _, og := range orderGoods {

		if form.IsLack > 0 && og.LackCount == 0 {
			//是欠货单时只显示欠货的
			continue
		}

		deliveryOrderNoArr = append(deliveryOrderNoArr, og.DeliveryOrderNo...)

		goodsType, ok := mp[og.GoodsType]

		if !ok {
			xsq_net.ErrorJSON(c, errors.New("商品类型:"+og.GoodsType+"数据未同步"))
			return
		}

		if _, detailOk := detailMap[goodsType]; !detailOk {
			detailMap[goodsType] = &rsp.Detail{
				Total: 0,
				List:  make([]*rsp.GoodsDetail, 0),
			}
		}

		detailMap[goodsType].Total += og.PayCount

		detailMap[goodsType].List = append(detailMap[goodsType].List, &rsp.GoodsDetail{
			Name:        og.GoodsName,
			GoodsSpe:    og.GoodsSpe,
			Shelves:     og.Shelves,
			PayCount:    og.PayCount,
			CloseCount:  og.CloseCount,
			LackCount:   og.LackCount,
			OutCount:    og.OutCount,
			GoodsRemark: og.GoodsRemark,
		})
	}

	res.Number = orders.Number
	res.PayAt = payAt.Format(timeutil.TimeFormat)
	res.ShopCode = orders.ShopCode
	res.ShopName = orders.ShopName
	res.Line = orders.Line
	res.Region = orders.Province + orders.City + orders.District
	res.ShopType = orders.ShopType
	res.OrderRemark = orders.OrderRemark

	res.Detail = detailMap

	deliveryOrderNoArr = slice.UniqueStringSlice(deliveryOrderNoArr)
	//历史出库单号
	res.DeliveryOrderNo = deliveryOrderNoArr

	xsq_net.SucJson(c, res)
}

// 商品列表
func CommodityList(c *gin.Context) {
	var result rsp.CommodityListRsp

	body, err := request.Post("api/v1/remote/shop/sku", nil)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	if result.Code != 200 {
		xsq_net.ErrorJSON(c, errors.New(result.Msg))
		return
	}

	xsq_net.SucJson(c, result.Data)
}

// 订单出货记录
func OrderShippingRecord(c *gin.Context) {
	var (
		form req.OrderShippingRecordReq
		res  rsp.OrderShippingRecordRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var pick []model.Pick

	result := global.DB.Where("delivery_no in (?)", form.DeliveryOrderNo).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.OrderShippingRecord, 0, result.RowsAffected)

	for _, p := range pick {
		list = append(list, rsp.OrderShippingRecord{
			Id:              p.Id,
			TakeOrdersTime:  p.TakeOrdersTime.Format(timeutil.TimeFormat),
			PickUser:        p.PickUser,
			ReviewUser:      p.ReviewUser,
			ReviewTime:      p.ReviewTime.Format(timeutil.TimeFormat),
			ReviewNum:       p.ReviewNum,
			DeliveryOrderNo: p.DeliveryOrderNo,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// 订单出货记录明细
func ShippingRecordDetail(c *gin.Context) {
	var (
		form req.ShippingRecordDetailReq
		//res  rsp.ShippingRecordDetailRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var pickGoods []model.PickGoods

	result := global.DB.Where("pick_id in (?) and review_num > 0", form.Id).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	type Goods struct {
		GoodsName string `json:"goods_name"`
		GoodsSpe  string `json:"goods_spe"`
		ReviewNum int    `json:"review_num"`
	}

	mp := make(map[string][]Goods, 0)

	for _, good := range pickGoods {

		mp[good.GoodsType] = append(mp[good.GoodsType], Goods{
			GoodsName: good.GoodsName,
			GoodsSpe:  good.GoodsSpe,
			ReviewNum: good.ReviewNum,
		})
	}

	xsq_net.SucJson(c, mp)
}

// 完成订单
func CompleteOrder(c *gin.Context) {
	var (
		form req.CompleteOrderReq
		res  rsp.CompleteOrderRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var completeOrder []model.CompleteOrder

	db := global.DB

	numbers := []string{}

	if form.Sku != "" {
		var completeOrderDetail []model.CompleteOrderDetail
		result := db.Model(&model.CompleteOrderDetail{}).Where("sku = ?", form.Sku).Find(&completeOrderDetail)
		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, detail := range completeOrderDetail {
			numbers = append(numbers, detail.Number)
		}

		numbers = slice.UniqueStringSlice(numbers)
	}

	//商品
	local := db.
		Model(&model.CompleteOrder{}).
		Where(&model.CompleteOrder{
			ShopId:         form.ShopId,
			Number:         form.Number,
			Line:           form.Line,
			DeliveryMethod: form.DeliveryMethod,
			ShopType:       form.ShopType,
			Province:       form.Province,
			City:           form.City,
			District:       form.District,
		})

	if len(numbers) > 0 {
		local.Where("number in (?)", numbers)
	}

	if form.IsRemark == 1 { //没有备注
		local.Where("order_remark == ''")
	} else if form.IsRemark == 2 { //有备注
		local.Where("order_remark != ''")
	}

	if form.PayAt != "" {
		t, err := time.ParseInLocation(timeutil.DateFormat, form.PayAt, time.Local)

		if err != nil {
			xsq_net.ErrorJSON(c, ecode.DataTransformationError)
			return
		}

		local.Where("pay_at <= ", timeutil.GetLastTime(t))
	}

	result := local.Find(&completeOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	result = local.Scopes(model.Paginate(form.Page, form.Size)).Find(&completeOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.CompleteOrder, 0, form.Size)

	for _, o := range completeOrder {

		pickTime := ""

		if o.PickTime != nil {
			pickTime = o.PickTime.Format(timeutil.TimeFormat)
		}

		list = append(list, rsp.CompleteOrder{
			Number:         o.Number,
			PayAt:          o.PayAt,
			ShopCode:       o.ShopCode,
			ShopName:       o.ShopName,
			ShopType:       o.ShopType,
			PayCount:       o.PayCount,
			OutCount:       o.OutCount,
			CloseCount:     o.CloseCount,
			Line:           o.Line,
			DeliveryMethod: o.DeliveryMethod,
			Region:         fmt.Sprintf("%s-%s-%s", o.Province, o.City, o.District),
			PickTime:       pickTime,
			OrderRemark:    o.OrderRemark,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// 完成订单详情
func CompleteOrderDetail(c *gin.Context) {
	var (
		form req.CompleteOrderDetailReq
		res  rsp.CompleteOrderDetailRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		completeOrder       model.CompleteOrder
		completeOrderDetail []model.CompleteOrderDetail
	)

	db := global.DB

	result := db.Model(&model.CompleteOrder{}).Where("number = ?", form.Number).First(&completeOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Where("number = ?", form.Number).Find(&completeOrderDetail)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	goodsMap := make(map[string][]rsp.PrePickGoods, 0)

	deliveryOrderNoArr := make(model.GormList, 0)

	for _, goods := range completeOrderDetail {
		goodsMap[goods.GoodsType] = append(goodsMap[goods.GoodsType], rsp.PrePickGoods{
			GoodsName:   goods.GoodsName,
			GoodsSpe:    goods.GoodsSpe,
			Shelves:     goods.Shelves,
			NeedNum:     goods.PayCount,
			CloseNum:    goods.CloseCount,
			OutCount:    goods.ReviewCount,
			NeedOutNum:  goods.PayCount,
			GoodsRemark: goods.GoodsRemark,
		})
		deliveryOrderNoArr = append(deliveryOrderNoArr, goods.DeliveryOrderNo...)
	}

	deliveryOrderNoArr = slice.UniqueStringSlice(deliveryOrderNoArr)

	res.Goods = goodsMap

	res.ShopName = completeOrder.ShopName
	res.ShopCode = completeOrder.ShopCode
	res.Line = completeOrder.Line
	res.Region = fmt.Sprintf("%s-%s-%s", completeOrder.Province, completeOrder.City, completeOrder.District)
	res.ShopType = completeOrder.ShopType
	res.Number = completeOrder.Number
	res.OrderRemark = completeOrder.OrderRemark
	res.DeliveryOrderNo = deliveryOrderNoArr

	xsq_net.SucJson(c, res)
}

type OrderNum struct {
	Count     int `json:"count"`
	OrderType int `json:"order_type"`
}

func Count(c *gin.Context) {

	var form req.CountFrom

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		orderGoods []model.OrderGoods
		numbers    []string
	)

	db := global.DB

	if form.Sku != "" {
		result := db.Where("sku = ?", form.Sku).Find(&orderGoods)

		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, good := range orderGoods {
			numbers = append(numbers, good.Number)
		}
	}

	localDb := db.Model(&model.Order{})

	if len(numbers) > 0 {
		localDb = localDb.Where("number in (?)", numbers)
	}

	localDb.Where(&model.Order{
		ShopId:           form.ShopId,
		Number:           form.Number,
		Line:             form.Lines,
		DistributionType: form.DistributionType,
		ShopType:         form.ShopType,
		Province:         form.Province,
		City:             form.City,
		District:         form.District,
		HasRemark:        form.HasRemark,
	})

	if form.PayEndTime != "" {
		localDb = localDb.Where("pay_at >= ?", form.PayEndTime)
	}

	localDb.Where("order_type != 2") //不要拣货中的

	type OrderNum struct {
		Count     int `json:"count"`
		OrderType int `json:"order_type"`
	}

	var (
		dbRes []OrderNum
		res   rsp.CountRes
	)

	result := localDb.Model(&model.Order{}).Select("count(id) as count, order_type").Group("order_type").Find(&dbRes)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, r := range dbRes {
		switch r.OrderType {
		case 1: //1:新订单
			res.NewCount = r.Count
			break
		case 2: //2:拣货中
			continue //不统计拣货中的
			//res.PickCount = r.Count
			//break
		case 3: //3:欠货单
			res.OldCount = r.Count
			break
		case 4: //已关闭
			res.CloseCount = r.Count
			break
		}
		res.AllCount += r.Count
	}

	xsq_net.SucJson(c, res)
}

// 创建拣货单
func CreatePickOrder(c *gin.Context) {
	var form req.CreatePickOrderForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var (
		order      []model.Order
		orderGoods []model.OrderGoods
	)

	db := global.DB

	result := db.Model(&model.Order{}).Where("number in (?)", form.Numbers).Find(&order)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.OrderGoods{}).Where("number in (?)", form.Numbers).Find(&orderGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	var (
		pickOrder      = make([]model.PickOrder, 0)
		pickOrderGoods = make([]model.PickOrderGoods, 0)
	)

	for _, o := range order {
		payAt, err := time.ParseInLocation(timeutil.TimeZoneFormat, o.PayAt, time.Local)

		if err != nil {
			xsq_net.ErrorJSON(c, ecode.DataTransformationError)
			return
		}

		deliveryAt, err := time.ParseInLocation(timeutil.TimeZoneFormat, o.DeliveryAt, time.Local)

		if err != nil {
			xsq_net.ErrorJSON(c, ecode.DataTransformationError)
			return
		}

		pickNumber, err := cache.GetIncrNumberByKey(constant.PICK_ORDER_NO, 4)

		if err != nil {
			xsq_net.ErrorJSON(c, errors.New("拣货单号生成失败"))
			return
		}

		total := o.PayTotal - o.Picked

		pickOrder = append(pickOrder, model.PickOrder{
			OrderId:           o.Id,
			ShopId:            o.ShopId,
			ShopName:          o.ShopName,
			ShopType:          o.ShopType,
			ShopCode:          o.ShopCode,
			Number:            o.Number,
			PickNumber:        "JHD" + pickNumber,
			HouseCode:         o.HouseCode,
			Line:              o.Line,
			DistributionType:  o.DistributionType,
			OrderRemark:       o.OrderRemark,
			PayAt:             payAt.Format(timeutil.TimeFormat),
			ShipmentsNum:      total,
			LimitNum:          total, //限发总数 默认等于应发总数
			CloseNum:          o.CloseNum,
			DeliveryAt:        deliveryAt.Format(timeutil.TimeFormat),
			Province:          o.Province,
			City:              o.City,
			District:          o.District,
			Address:           o.Address,
			ConsigneeName:     o.ConsigneeName,
			ConsigneeTel:      o.ConsigneeTel,
			OrderType:         1, //重新进入时，改为新订单
			HasRemark:         o.HasRemark,
			LatestPickingTime: o.LatestPickingTime,
		})
	}

	for _, og := range orderGoods {

		if og.LackCount <= 0 {
			continue
		}

		pickOrderGoods = append(pickOrderGoods, model.PickOrderGoods{
			OrderGoodsId:    og.Id,
			Number:          og.Number,
			GoodsName:       og.GoodsName,
			Sku:             og.Sku,
			GoodsType:       og.GoodsType,
			GoodsSpe:        og.GoodsSpe,
			Shelves:         og.Shelves,
			DiscountPrice:   og.DiscountPrice,
			GoodsUnit:       og.GoodsUnit,
			SaleUnit:        og.SaleUnit,
			SaleCode:        og.SaleCode,
			PayCount:        og.PayCount,
			CloseCount:      og.CloseCount,
			LackCount:       og.LackCount,
			OutCount:        og.OutCount,
			LimitNum:        og.LackCount,
			GoodsRemark:     og.GoodsRemark,
			BatchId:         0,
			DeliveryOrderNo: og.DeliveryOrderNo,
		})
	}

	tx := db.Begin()

	result = tx.Save(&pickOrder)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//取出pick_order 表中 number 对应的 id 存入 到 pick_order_goods 表中
	var pickOrderNumberIdMp = make(map[string]int, len(pickOrder))

	for _, po := range pickOrder {
		pickOrderNumberIdMp[po.Number] = po.Id
	}

	for i, goods := range pickOrderGoods {
		pickOrderId, ok := pickOrderNumberIdMp[goods.Number]

		if !ok {
			tx.Rollback()
			xsq_net.ErrorJSON(c, errors.New("number:"+goods.Number+"拣货单表id数据不存在，请联系管理员"))
			return
		}

		pickOrderGoods[i].PickOrderId = pickOrderId
	}

	result = tx.Save(&pickOrderGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Model(&model.Order{}).Where("number in (?)", form.Numbers).Updates(map[string]interface{}{"order_type": 2})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

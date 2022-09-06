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
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/slice"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
)

// 拣货单列表
func PickOrderList(c *gin.Context) {
	var form req.PickOrderListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var (
		pickOrder      []model.PickOrder
		pickOrderGoods []model.PickOrderGoods
		numbers        []string
	)

	db := global.DB

	if form.Sku != "" {
		result := db.Where("sku = ?", form.Sku).Find(&pickOrderGoods)

		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, good := range pickOrderGoods {
			numbers = append(numbers, good.Number)
		}
	}

	localDb := db.Model(&model.PickOrder{})

	if len(numbers) > 0 {
		numbers = slice.UniqueStringSlice(numbers)
		localDb = localDb.Where("number in (?)", numbers)
	}

	localDb.Where(&model.PickOrder{
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

	var (
		total int64
		res   rsp.PickOrderGoodsListRsp
	)

	result := localDb.Count(&total)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = localDb.Scopes(model.Paginate(form.Page, form.Size)).Order("id desc").Find(&pickOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.PickOrder, 0, form.Size)

	for _, o := range pickOrder {
		latestPickingTime := ""

		if o.LatestPickingTime != nil {
			latestPickingTime = o.LatestPickingTime.Format(timeutil.TimeFormat)
		}

		list = append(list, rsp.PickOrder{
			Number:            o.Number,
			PickNumber:        o.PickNumber,
			PayAt:             o.PayAt,
			ShopCode:          o.ShopCode,
			ShopName:          o.ShopName,
			ShopType:          o.ShopType,
			DistributionType:  o.DistributionType,
			ShipmentsNum:      o.ShipmentsNum,
			LimitNum:          o.LimitNum,
			CloseNum:          o.CloseNum,
			Line:              o.Line,
			Region:            o.Province + o.City + o.District,
			OrderRemark:       o.OrderRemark,
			OrderType:         o.OrderType,
			LatestPickingTime: latestPickingTime,
		})
	}

	res.Total = total
	res.List = list

	xsq_net.SucJson(c, res)
}

// 拣货单明细
func GetPickOrderDetail(c *gin.Context) {

	var form req.GetOrderDetailForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		orders     model.PickOrder
		orderGoods []model.PickOrderGoods
	)

	db := global.DB

	result := db.Model(&model.PickOrder{}).Where("number = ?", form.Number).First(&orders)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.PickOrderGoods{}).Where("number = ?", form.Number).Find(&orderGoods)

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

	deliveryOrderNoArr := make([]*string, 0)

	for _, og := range orderGoods {

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
			LackCount:   og.LimitNum, //需拣数 以限发数为准
			GoodsRemark: og.GoodsRemark,
		})
	}

	res.Number = orders.Number
	res.PayAt = orders.PayAt
	res.ShopCode = orders.ShopCode
	res.ShopName = orders.ShopName
	res.Line = orders.Line
	res.Region = orders.Province + orders.City + orders.District
	res.ShopType = orders.ShopType
	res.OrderRemark = orders.OrderRemark

	res.Detail = detailMap

	deliveryOrderNoArr = slice.UniqueStringSlicePtr(deliveryOrderNoArr)
	//历史出库单号
	res.DeliveryOrderNo = deliveryOrderNoArr

	xsq_net.SucJson(c, res)
}

// 配送方式明细
func DeliveryMethodInfo(c *gin.Context) {
	var form req.DeliveryMethodInfoForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		pickOrder model.PickOrder
		res       rsp.DeliveryMethodInfoRsp
	)

	result := global.DB.Model(&model.PickOrder{}).Where("id = ?", form.Id).First(&pickOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.UserName = pickOrder.ConsigneeName
	res.Tel = pickOrder.ConsigneeTel
	res.Province = pickOrder.Province
	res.City = pickOrder.City
	res.District = pickOrder.District
	res.Address = pickOrder.Address

	xsq_net.SucJson(c, res)
}

// 修改配送方式
func ChangeDeliveryMethod(c *gin.Context) {
	var form req.ChangeDeliveryMethodForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var dict []int

	db := global.DB

	result := db.Model(&model.Dict{}).Select("value").Where("type_code = 'delivery_method'").Find(&dict)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//校验是否在字典中
	if ok, _ := slice.InArray(form.DeliveryMethod, dict); !ok {
		xsq_net.ErrorJSON(c, ecode.DataCannotBeModified)
		return
	}

	result = db.Model(&model.PickOrder{}).Where("id = ?", form.Id).Update("distribution_type", form.DeliveryMethod)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 订单商品列表
func OrderGoodsList(c *gin.Context) {
	var form req.OrderGoodsListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var pickOrderGood []model.PickOrderGoods

	result := global.DB.Model(&model.PickOrderGoods{}).Where("number = ?", form.Number).Find(&pickOrderGood)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.PickOrderGoods, 0, len(pickOrderGood))

	for _, goods := range pickOrderGood {
		list = append(list, rsp.PickOrderGoods{
			Id:        goods.Id,
			GoodsName: goods.GoodsName,
			LackCount: goods.LackCount,
			LimitNum:  goods.LimitNum,
		})
	}

	xsq_net.SucJson(c, list)
}

// 限发
func RestrictedShipment(c *gin.Context) {
	var form req.RestrictedShipmentForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	var (
		pickOrder          model.PickOrder
		goodsIds           []int                                 //拣货单商品表id
		mp                 = make(map[int]int, len(form.Params)) // 拣货单商品表id 和 限发数量 map
		restrictedShipment = make([]model.RestrictedShipment, 0, len(form.Params))
	)

	result := db.Model(&model.PickOrder{}).Where("number = ?", form.Number).First(&pickOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if pickOrder.OrderType != 1 {
		xsq_net.ErrorJSON(c, errors.New("当前订单不允许设置限发"))
		return
	}

	//限发总数
	limitTotal := 0
	for _, p := range form.Params {
		goodsIds = append(goodsIds, p.Id)

		limitTotal += p.LimitNum
		mp[p.Id] = p.LimitNum
	}

	var pickOrderGoods []model.PickOrderGoods

	result = db.Model(&model.PickOrderGoods{}).Where("id in (?)", goodsIds).Find(&pickOrderGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for i, good := range pickOrderGoods {
		num, ok := mp[good.Id]

		if !ok {
			continue
		}

		pickOrderGoods[i].LimitNum = num

		restrictedShipment = append(restrictedShipment, model.RestrictedShipment{
			PickOrderGoodsId: good.Id,
			Number:           good.Number,
			ShopName:         good.GoodsName,
			GoodsSpe:         good.GoodsSpe,
			LimitNum:         num,
		})
	}

	tx := db.Begin()

	result = tx.Save(&pickOrderGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Save(&restrictedShipment)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Model(&model.PickOrder{}).Where("id = ?", pickOrder.Id).Update("limit_num", limitTotal)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

// 批量设置限发
func BatchRestrictedShipment(c *gin.Context) {
	var form req.BatchRestrictedShipmentForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		pickOrderGoods     []model.PickOrderGoods
		restrictedShipment = make([]model.RestrictedShipment, 0)
	)

	db := global.DB

	result := db.Model(&model.PickOrderGoods{}).Where("sku = ? and status = 0", form.Sku).Find(&pickOrderGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, pg := range pickOrderGoods {
		restrictedShipment = append(restrictedShipment, model.RestrictedShipment{
			PickOrderGoodsId: pg.Id,
			Number:           pg.Number,
			ShopName:         pg.GoodsName,
			GoodsSpe:         pg.GoodsSpe,
			LimitNum:         form.LimitNum,
		})
	}

	if len(restrictedShipment) > 0 {
		result = db.Save(&restrictedShipment)

		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	} else {
		xsq_net.ErrorJSON(c, errors.New("没有待拣货的sku"))
		return
	}
	xsq_net.Success(c)
}

// 批量设置限发商品数量
func GoodsNum(c *gin.Context) {
	var form req.GoodsNumForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var sum int

	result := global.DB.Model(&model.PickOrderGoods{}).Select("sum(lack_count) as sum").Where("sku = ? and status = 0", form.Sku).Scan(&sum)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.SucJson(c, gin.H{"total": sum})
}

// 限发列表
func RestrictedShipmentList(c *gin.Context) {
	var form req.RestrictedShipmentListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		total              int64
		restrictedShipment []model.RestrictedShipment
	)

	db := global.DB

	result := db.Model(&model.RestrictedShipment{}).Where("status = 1").Count(&total)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.RestrictedShipment{}).Where("status = 1").Scopes(model.Paginate(form.Page, form.Size)).Find(&restrictedShipment)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.SucJson(c, restrictedShipment)
}

// 撤销限发
func RevokeRestrictedShipment(c *gin.Context) {
	var form req.RevokeRestrictedShipmentForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var pickOrderGoods model.PickOrderGoods

	db := global.DB

	result := db.Model(&model.PickOrderGoods{}).Where("id = ?", form.PickOrderGoodsId).First(&pickOrderGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if pickOrderGoods.Status != 0 {
		xsq_net.ErrorJSON(c, errors.New("当前订单不允许撤销限发"))
		return
	}

	tx := db.Begin()

	result = tx.Model(&model.RestrictedShipment{}).Where("pick_order_goods_id = ?", form.PickOrderGoodsId).Update("status", 0)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//更新限发数量未原始值(欠货数)
	pickOrderGoods.LimitNum = pickOrderGoods.LackCount

	result = tx.Model(&model.PickOrderGoods{}).Where("id = ?", pickOrderGoods.Id).Update("limit_num", pickOrderGoods.LackCount)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

// 关闭订单
func CloseOrder(c *gin.Context) {
	var form req.CloseOrderForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		pickOrder []model.PickOrder
	)

	db := global.DB

	//根据ID倒叙，查最新的记录，避免欠货单更新错误
	result := db.Model(&model.PickOrder{}).Where("order_id = ?", form.OrderId).Find(&pickOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, order := range pickOrder {
		if order.OrderType != 1 { //只要有 != 1 的 即 不允许更新
			xsq_net.ErrorJSON(c, errors.New("当前订单不允许更新"))
			return
		}
	}

	tx := db.Begin()

	// 拣货单查到数据 且是新订单
	if len(pickOrder) > 0 {
		result = tx.Model(&model.PickOrder{}).Where("order_id = ?", form.OrderId).Update("order_type", 3)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	}

	var order model.Order

	//查询订单表
	result = db.Model(&model.Order{}).First(&order, form.OrderId)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//更新订单表
	result = tx.Model(&model.Order{}).Where("id = ?", form.OrderId).Update("order_type", 4)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
	return
}

// 关闭商品
func CloseOrderGoods(c *gin.Context) {
	var form req.CloseOrderGoodsForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		order          model.Order
		orderGoods     model.OrderGoods
		pickOrderGoods model.PickOrderGoods
	)

	db := global.DB

	result := db.Model(&model.PickOrderGoods{}).Where("order_goods_id = ?", form.GoodsId).Order("id desc").First(&pickOrderGoods)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx := db.Begin()
	//查到数据
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {

		if pickOrderGoods.Status != 0 {
			xsq_net.ErrorJSON(c, errors.New("当前商品不允许关闭"))
			return
		}

		if pickOrderGoods.LackCount < form.CloseNum {
			xsq_net.ErrorJSON(c, errors.New("关闭数量大于欠货数量"))
			return
		}

		// 增加关闭数量 && 扣减发货数量
		result = tx.Model(&model.PickOrderGoods{}).
			Where("id = ?", pickOrderGoods.Id).
			Updates(map[string]interface{}{
				"close_count": gorm.Expr("close_count + ?", form.CloseNum),
				"lack_count":  gorm.Expr("lack_count - ?", form.CloseNum),
			})

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		result = tx.Model(&model.PickOrder{}).
			Where("number = ?", pickOrderGoods.Number).
			Updates(map[string]interface{}{
				"close_count":   gorm.Expr("close_count + ?", form.CloseNum),
				"shipments_num": gorm.Expr("shipments_num - ?", form.CloseNum),
			})

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	}

	result = db.Model(&model.OrderGoods{}).First(&orderGoods, form.GoodsId)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			//mq中订单没到拣货系统中，这个其实是不合理的状态
			//没查到，先什么都不做，让订货关商品通过，后续再拉单
			xsq_net.Success(c)
			return
		}

		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.Order{}).Where("number = ?", orderGoods.Number).First(&order)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			//mq中订单没到拣货系统中，这个其实是不合理的状态
			//没查到，先什么都不做，让订货关商品通过，后续再拉单
			xsq_net.Success(c)
			return
		}
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//2拣货中 4已关闭
	if order.OrderType == 2 || order.OrderType == 4 {
		xsq_net.ErrorJSON(c, errors.New("当前商品不允许关闭"))
		return
	}

	// 增加关闭数量 && 扣减发货数量
	result = tx.Model(&model.OrderGoods{}).
		Where("id = ?", orderGoods.Id).
		Updates(map[string]interface{}{
			"close_count": gorm.Expr("close_count + ?", form.CloseNum),
			"lack_count":  gorm.Expr("lack_count - ?", form.CloseNum),
		})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Model(&model.Order{}).
		Where("number = ?", orderGoods.Number).
		Updates(map[string]interface{}{
			"close_num": gorm.Expr("close_num + ?", form.CloseNum),
			"un_picked": gorm.Expr("un_picked - ?", form.CloseNum),
		})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

// 拣货单统计
func PickOrderCount(c *gin.Context) {

	type OrderNum struct {
		Count     int `json:"count"`
		OrderType int `json:"order_type"`
	}

	var (
		dbRes []OrderNum
		res   rsp.PickOrderCount
	)

	result := global.DB.Model(&model.PickOrder{}).Select("count(id) as count, order_type").Group("order_type").Find(&dbRes)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//1:新订单,2:拣货中,3:已关闭,4:已完成
	for _, r := range dbRes {
		switch r.OrderType {
		case 1: //1:新订单
			res.NewCount = r.Count
			break
		case 2: //2:拣货中
			res.PickCount = r.Count
			break
		case 3: //3:已关闭
			res.CloseCount = r.Count
			break
		case 4: //已关闭
			res.CompleteCount = r.Count
			break
		}
		res.AllCount += r.Count
	}

	xsq_net.SucJson(c, res)
}

func Test(c *gin.Context) {
	salt, sign := middlewares.Generate()
	xsq_net.SucJson(c, gin.H{"salt": salt, "sign": sign})
}

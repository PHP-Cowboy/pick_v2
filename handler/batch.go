package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"pick_v2/common/constant"
	"pick_v2/dao"
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
	"strconv"
	"strings"
	"time"
)

// 全量拣货 -按任务创建批次
func CreateBatchByTask(c *gin.Context) {
	var form req.CreateBatchByTaskForm

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

	form.Typ = 1 // 常规批次

	err := dao.CreateBatchByTask(global.DB, form, userInfo)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

// 创建批次
func NewBatch(c *gin.Context) {
	var form req.NewCreateBatchForm

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

	form.Typ = 1 // 常规批次

	err := dao.CreateBatch(global.DB, form, userInfo)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

// 创建快递批次
func CourierBatch(c *gin.Context) {
	var form req.NewCreateBatchForm

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

	form.Typ = 2 //快递批次

	err := dao.CourierBatch(global.DB, form, userInfo)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)

}

// 集中拣货列表
func CentralizedPickList(c *gin.Context) {

	var form req.CentralizedPickListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, list := dao.CentralizedPickList(global.DB, form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, list)
}

// 集中拣货详情
func CentralizedPickDetail(c *gin.Context) {
	var form req.CentralizedPickDetailForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	err, list := dao.CentralizedPickDetail(global.DB, form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, gin.H{"list": list})
}

// 生成拣货批次
func CreateBatch(c *gin.Context) {
	var form req.CreateBatchForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	tx := global.DB.Begin()

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, ecode.GetContextUserInfoFailed)
		return
	}

	var sku, goodsName string

	if len(form.Sku) > 0 {
		sku = strings.Join(form.Sku, ",")
		goodsName = strings.Join(form.GoodsName, ",")
	}

	var (
		deliveryStartTime *time.Time
	)

	deliveryEndTime, errDeliveryEndTime := time.ParseInLocation(timeutil.TimeFormat, form.DeliveryEndTime, time.Local)

	payTime, errPayTime := time.ParseInLocation(timeutil.TimeFormat, form.PayTime, time.Local)

	if errDeliveryEndTime != nil || errPayTime != nil {
		xsq_net.ErrorJSON(c, ecode.DataTransformationError)
		return
	}

	batchName := form.BatchName

	lines := "全部线路"

	if len(form.Lines) > 0 {
		lines = strings.Join(form.Lines, ",")
	}

	if form.BatchName == "" {
		batchName = lines
	}

	if form.DeliveryStartTime != "" {
		deliveryStart, errDeliveryStartTime := time.ParseInLocation(timeutil.TimeFormat, form.DeliveryStartTime, time.Local)
		if errDeliveryStartTime != nil {
			xsq_net.ErrorJSON(c, ecode.DataTransformationError)
			return
		}
		deliveryStartTime = &deliveryStart
	}

	batchId, err := SaveBatch(tx, userInfo, batchName, lines, sku, goodsName, form.DistributionType, (*model.MyTime)(&payTime), (*model.MyTime)(deliveryStartTime), (*model.MyTime)(&deliveryEndTime))
	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, errors.New("批次数据保存失败:"+err.Error()))
		return
	}

	numbers, pickOrderGoodsId, orderGoodsIds, err := SavePrePickPool(tx, userInfo, batchId, "", &form)
	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	err = UpdateOrder(tx, numbers, pickOrderGoodsId, orderGoodsIds, batchId)
	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, errors.New("订单数据更新失败:"+err.Error()))
		return
	}

	err = UpdateBatchPickNums(tx, batchId)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, errors.New("批次拣货池数量更新失败:"+err.Error()))
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

// 根据订单生成批次
func CreateByOrder(c *gin.Context) {
	var form req.CreateByOrderReq

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

	var (
		pickOrder model.PickOrder
	)

	db := global.DB

	result := db.Where("pick_number = ?", form.PickNumber).First(&pickOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if pickOrder.OrderType == 2 {
		xsq_net.ErrorJSON(c, errors.New("订单已经在拣货中"))
		return
	}

	tx := db.Begin()

	//批次相关
	batchId, err := SaveBatch(tx, userInfo, pickOrder.ShopName, pickOrder.Line, "", "", pickOrder.DistributionType, &pickOrder.PayAt, nil, &pickOrder.DeliveryAt)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, errors.New("批次数据保存失败:"+err.Error()))
		return
	}

	numbers, pickOrderGoodsId, orderGoodsIds, err := SavePrePickPool(tx, userInfo, batchId, form.PickNumber, nil)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	//更新订单
	err = UpdateOrder(tx, numbers, pickOrderGoodsId, orderGoodsIds, batchId)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, errors.New("订单数据更新失败:"+err.Error()))
		return
	}

	err = UpdateBatchPickNums(tx, batchId)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, errors.New("批次拣货池数量更新失败:"+err.Error()))
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

func GetUserInfo(c *gin.Context) *middlewares.CustomClaims {
	claims, ok := c.Get("claims")

	if !ok {
		return nil
	}

	return claims.(*middlewares.CustomClaims)
}

func UpdateOrder(tx *gorm.DB, numbers []string, pickOrderGoodsId []int, orderGoodsIds []int, batchId int) error {

	result := tx.Model(&model.PickOrder{}).Where("number in (?)", numbers).Update("order_type", model.PickOrderPickingOrderType)

	if result.Error != nil {
		return result.Error
	}

	result = tx.Model(&model.PickOrderGoods{}).Where("id in (?)", pickOrderGoodsId).Updates(map[string]interface{}{"status": 1, "batch_id": batchId})

	if result.Error != nil {
		return result.Error
	}

	//查询orderGoods id

	tx.Model(&model.OrderGoods{}).Where("id in (?)", orderGoodsIds).Updates(map[string]interface{}{"batch_id": batchId})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func SaveBatch(tx *gorm.DB, userInfo *middlewares.CustomClaims, batchName, line, sku, goods string, deliveryMethod int, payEndTime, deliveryStartTime, deliveryEndTime *model.MyTime) (int, error) {

	now := time.Now()

	//创建批次
	batches := model.Batch{
		WarehouseId:     userInfo.WarehouseId,
		BatchName:       batchName,
		DeliveryEndTime: deliveryEndTime,
		UserName:        userInfo.Name,
		Line:            line,
		DeliveryMethod:  deliveryMethod,
		EndTime:         (*model.MyTime)(&now),
		Status:          0,
		Sort:            0,
	}

	result := tx.Save(&batches)

	if result.Error != nil {
		return 0, result.Error
	}

	//批次创建条件
	condition := model.BatchCondition{
		BatchId:         batches.Id,
		WarehouseId:     userInfo.WarehouseId,
		PayEndTime:      payEndTime,
		DeliveryEndTime: deliveryEndTime,
		Line:            line,
		DeliveryMethod:  deliveryMethod,
		Sku:             sku,
		Goods:           goods,
	}

	if deliveryStartTime != nil {
		condition.DeliveryStartTime = deliveryStartTime
	}

	result = tx.Save(&condition)

	if result.Error != nil {
		return 0, result.Error
	}

	return batches.Id, nil
}

func SavePrePickPool(tx *gorm.DB, userInfo *middlewares.CustomClaims, batchId int, pickNumber string, form *req.CreateBatchForm) ([]string, []int, []int, error) {

	//返回给调用方更新订单和订单商品状态
	var (
		numbers          []string
		pickOrderGoodsId []int
		orderGoodsId     []int
	)

	//缓存中的线路数据
	lineCacheMp, errCache := cache.GetShopLine()

	if errCache != nil {
		return numbers, pickOrderGoodsId, orderGoodsId, errors.New("线路缓存获取失败")
	}

	mp, err := cache.GetClassification()

	if err != nil {
		return numbers, pickOrderGoodsId, orderGoodsId, err
	}

	var orderAndGoods []rsp.OrderAndGoods

	//筛选条件
	localDb := global.DB.
		Table("t_pick_order_goods og").
		Select("og.*,o.shop_id,o.shop_name,o.shop_code,o.line,o.distribution_type,o.order_remark").
		Joins("left join t_pick_order o on og.pick_order_id = o.id")

	if pickNumber != "" {
		localDb = localDb.Where("o.pick_number = ?", pickNumber)
	}

	if form != nil {
		localDb = localDb.Where("o.distribution_type = ? and o.pay_at <= ? and o.delivery_at <= ? ", form.DistributionType, form.PayTime, form.DeliveryEndTime)

		if form.DeliveryStartTime != "" {
			localDb = localDb.Where("o.delivery_at >= ?", form.DeliveryStartTime)
		}

		if len(form.Sku) > 0 {
			localDb = localDb.Where("og.sku in (?)", form.Sku)
		}

		if len(form.Lines) > 0 {
			localDb = localDb.Where("o.line in (?) ", form.Lines)
		}
	}

	result := localDb.Where("og.status = 0").Find(&orderAndGoods)

	if result.Error != nil {
		return numbers, pickOrderGoodsId, orderGoodsId, result.Error
	}

	var (
		prePicks      []model.PrePick
		prePickGoods  []*model.PrePickGoods
		prePickRemark []*model.PrePickRemark
	)

	//订单相关数据 -店铺数 订单数 商品数
	goodsNum := 0                              //商品数
	shopNumMp := make(map[int]struct{}, 0)     //店铺
	orderNumMp := make(map[string]struct{}, 0) //订单

	for _, og := range orderAndGoods {
		//拣货单单商品表id
		pickOrderGoodsId = append(pickOrderGoodsId, og.Id)
		//
		orderGoodsId = append(orderGoodsId, og.OrderGoodsId)
		//商品总数量
		goodsNum += og.LackCount
		//商品类型
		goodsType, mpOk := mp[og.GoodsType]

		if !mpOk {
			return numbers, pickOrderGoodsId, orderGoodsId, errors.New("商品类型:" + og.GoodsType + "数据未同步")
		}
		//线路
		cacheMpLine, cacheMpOk := lineCacheMp[og.ShopId]

		if !cacheMpOk {
			return numbers, pickOrderGoodsId, orderGoodsId, errors.New("店铺:" + og.ShopName + "线路未同步，请先同步")
		}

		//店铺mp 订单号去重使用
		_, orderOk := orderNumMp[og.Number]

		if !orderOk {
			//订单号去重
			numbers = append(numbers, og.Number)
		}

		//店铺数
		orderNumMp[og.Number] = struct{}{}

		needNum := og.LackCount

		//如果欠货数量大于限发数量，需拣货数量为限货数
		if og.LackCount > og.LimitNum {
			needNum = og.LimitNum
		}

		prePickGoods = append(prePickGoods, &model.PrePickGoods{
			WarehouseId:      userInfo.WarehouseId,
			BatchId:          batchId,
			OrderGoodsId:     og.Id,
			Number:           og.Number,
			PrePickId:        0, //后续逻辑变更
			ShopId:           og.ShopId,
			DistributionType: og.DistributionType,
			Sku:              og.Sku,
			GoodsName:        og.GoodsName,
			GoodsType:        goodsType,
			GoodsSpe:         og.GoodsSpe,
			Shelves:          og.Shelves,
			DiscountPrice:    og.DiscountPrice,
			Unit:             og.GoodsUnit,
			NeedNum:          needNum,
			CloseNum:         og.CloseCount,
			OutCount:         0,
			NeedOutNum:       og.LackCount,
		})

		if og.GoodsRemark != "" || og.OrderRemark != "" {
			prePickRemark = append(prePickRemark, &model.PrePickRemark{
				WarehouseId:  userInfo.WarehouseId,
				BatchId:      batchId,
				OrderGoodsId: og.Id,
				ShopId:       og.ShopId,
				PrePickId:    0,
				Number:       og.Number,
				OrderRemark:  og.OrderRemark,
				GoodsRemark:  og.GoodsRemark,
				ShopName:     og.ShopName,
				Line:         cacheMpLine,
			})
		}

		_, shopMpOk := shopNumMp[og.ShopId]

		if shopMpOk {
			continue
		}

		shopNumMp[og.ShopId] = struct{}{}

		prePicks = append(prePicks, model.PrePick{
			WarehouseId: userInfo.WarehouseId,
			BatchId:     batchId,
			ShopId:      og.ShopId,
			ShopCode:    og.ShopCode,
			ShopName:    og.ShopName,
			Line:        cacheMpLine,
			Status:      0,
		})
	}

	//预拣池数量
	prePickNum := len(prePicks)

	if prePickNum == 0 {
		return numbers, pickOrderGoodsId, orderGoodsId, ecode.NoOrderFound
	}

	result = tx.Save(&prePicks)

	if result.Error != nil {
		return numbers, pickOrderGoodsId, orderGoodsId, result.Error
	}

	shopMap := make(map[int]int, 0)

	for _, pick := range prePicks {
		shopMap[pick.ShopId] = pick.Id
	}

	for k, good := range prePickGoods {
		val, shopMapOk := shopMap[good.ShopId]
		if !shopMapOk {
			return numbers, pickOrderGoodsId, orderGoodsId, ecode.MapKeyNotExist
		}
		prePickGoods[k].PrePickId = val
	}

	result = tx.Save(&prePickGoods)

	if result.Error != nil {
		return numbers, pickOrderGoodsId, orderGoodsId, result.Error
	}

	if len(prePickRemark) > 0 {
		for k, remark := range prePickRemark {
			val, shopMapOk := shopMap[remark.ShopId]
			if !shopMapOk {
				return numbers, pickOrderGoodsId, orderGoodsId, ecode.MapKeyNotExist
			}
			prePickRemark[k].PrePickId = val
		}

		result = tx.Save(&prePickRemark)

		if result.Error != nil {
			return numbers, pickOrderGoodsId, orderGoodsId, result.Error
		}
	}

	shopNum := len(shopNumMp)
	orderNum := len(orderNumMp)

	updates := map[string]interface{}{}

	updates["goods_num"] = goodsNum
	updates["shop_num"] = shopNum
	updates["order_num"] = orderNum

	result = tx.Model(&model.Batch{}).
		Where("id = ?", batchId).
		Updates(map[string]interface{}{
			"goods_num":    goodsNum,
			"shop_num":     shopNum,
			"order_num":    orderNum,
			"pre_pick_num": prePickNum,
		})

	if result.Error != nil {
		return numbers, pickOrderGoodsId, orderGoodsId, result.Error
	}

	return numbers, pickOrderGoodsId, orderGoodsId, nil
}

// 结束拣货批次
func EndBatch(c *gin.Context) {
	var form req.EndBatchForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		batches       model.Batch
		pickGoods     []model.PickGoods
		pick          []model.Pick
		orderAndGoods []rsp.OrderAndGoods
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
	result = db.Model(&model.Batch{}).Where("id = ?", batches.Id).Updates(map[string]interface{}{"status": 1})

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Table("t_pick_order_goods og").
		Select("og.*,o.shop_id,o.shop_name,o.shop_code,o.line,o.distribution_type,o.order_remark,o.pay_at,o.province,o.city,o.district,o.shop_type,o.latest_picking_time").
		Joins("left join t_pick_order o on og.pick_order_id = o.id").
		Where("batch_id = ? ", form.Id).
		Find(&orderAndGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//查询批次下全部订单
	result = db.Model(&model.PickGoods{}).Where("batch_id = ? ", form.Id).Find(&pickGoods)
	if result.Error != nil {
		global.Logger["err"].Infof("批次结束成功，但推送u8拣货数据查询失败:%s", result.Error.Error())
		xsq_net.ErrorJSON(c, errors.New("批次结束成功，但推送u8拣货数据查询失败"))
		return
	}

	result = db.Model(&model.Pick{}).Where("batch_id = ?", form.Id).Find(&pick)

	if result.Error != nil {
		global.Logger["err"].Infof("批次结束成功，但推送u8拣货数据查询失败:%s", result.Error.Error())
		xsq_net.ErrorJSON(c, errors.New("批次结束成功，但推送u8拣货主表数据查询失败"))
		return
	}

	//拣货表数据map
	mpPick := make(map[int]model.Pick, 0)

	for _, p := range pick {
		mpPick[p.Id] = p
	}

	tx := global.DB.Begin()

	err := YongYouLog(tx, pickGoods, orderAndGoods, form.Id)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	//这里会删数据，要放在推u8之后处理，失败重试要加上这里的逻辑
	err = UpdateCompleteOrder(tx, form.Id)
	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

func UpdateCompleteOrder(tx *gorm.DB, batchId int) error {
	db := global.DB

	var (
		order      []model.Order
		orderGoods []model.OrderGoods
		isSendMQ   = true
	)
	//根据批次拿order会导致完成订单有遗漏
	//OrderGoods 的批次id是否被更新了---不会，批次里的单在被更新为欠货之前一个商品只能被一个批次拿走，如果不是应该优化批次拿数据逻辑
	result := db.Model(&model.OrderGoods{}).Where("batch_id = ?", batchId).Find(&orderGoods)

	if result.Error != nil {
		return result.Error
	}

	//当前批次所有订单号
	var numbers []string

	for _, good := range orderGoods {
		numbers = append(numbers, good.Number)
	}

	numbers = slice.UniqueStringSlice(numbers)

	//根据当前批次的所有订单号，查询订单
	result = db.Model(&model.Order{}).Where("number in (?)", numbers).Find(&order)

	if result.Error != nil {
		return result.Error
	}

	var (
		//完成订单map
		//key => number
		completeMp          = make(map[string]interface{}, 0)
		completeNumbers     []string
		deleteIds           []int    //待删除订单表id
		deleteNumbers       []string //删除订单商品表
		completeOrder       = make([]model.CompleteOrder, 0)
		completeOrderDetail = make([]model.CompleteOrderDetail, 0)
		lackNumbers         []string //待更新为欠货订单表number
	)

	for _, o := range order {
		//TODO [这个字段删除了，逻辑要改] 还有欠货
		//if o.UnPicked > 0 {
		//	lackNumbers = append(lackNumbers, o.Number)
		//	continue
		//}

		deleteIds = append(deleteIds, o.Id)

		completeMp[o.Number] = struct{}{}

		//完成订单
		completeNumbers = append(completeNumbers, o.Number)

		deleteNumbers = append(deleteNumbers, o.Number)

		//完成订单
		completeOrder = append(completeOrder, model.CompleteOrder{
			Number:         o.Number,
			OrderRemark:    o.OrderRemark,
			ShopId:         o.ShopId,
			ShopName:       o.ShopName,
			ShopType:       o.ShopType,
			ShopCode:       o.ShopCode,
			Line:           o.Line,
			DeliveryMethod: o.DistributionType,
			Province:       o.Province,
			City:           o.City,
			District:       o.District,
			PickTime:       o.LatestPickingTime,
			PayAt:          o.PayAt,
		})
	}

	if len(completeNumbers) > 0 {
		//如果有完成订单重新查询订单商品，因为这批商品可能是多个批次拣的，根据批次查的商品不全
		result = db.Model(&model.OrderGoods{}).Where("number in (?)", completeNumbers).Find(&orderGoods)

		if result.Error != nil {
			return result.Error
		}

		for _, og := range orderGoods {
			//完成订单map中不存在订单的跳过
			_, ok := completeMp[og.Number]

			if !ok {
				continue
			}
			//完成订单详情
			completeOrderDetail = append(completeOrderDetail, model.CompleteOrderDetail{
				Number:          og.Number,
				GoodsName:       og.GoodsName,
				Sku:             og.Sku,
				GoodsSpe:        og.GoodsSpe,
				GoodsType:       og.GoodsType,
				Shelves:         og.Shelves,
				PayCount:        og.PayCount,
				CloseCount:      og.CloseCount,
				ReviewCount:     og.OutCount,
				GoodsRemark:     og.GoodsRemark,
				DeliveryOrderNo: og.DeliveryOrderNo,
			})
		}
	}

	if len(deleteIds) > 0 {
		result = tx.Delete(&model.Order{}, "id in (?)", deleteIds)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	if len(deleteNumbers) > 0 {
		deleteNumbers = slice.UniqueStringSlice(deleteNumbers)

		result = tx.Delete(&model.OrderGoods{}, "number in (?)", deleteNumbers)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	// 欠货的单 拣货池还有未完成的，不更新为欠货
	if len(lackNumbers) > 0 {
		var (
			pickAndGoods   []model.PickAndGoods
			pendingNumbers []string
			diffSlice      []string
		)

		//获取欠货的订单number是否有在拣货池中未复核完成的数据，如果有，过滤掉欠货的订单number
		result = db.Table("t_pick_goods pg").
			Select("p.id as pick_id,p.status,pg.number,p.pick_user").
			Joins("left join t_pick p on pg.pick_id = p.id").
			Where("number in (?)", lackNumbers).
			Find(&pickAndGoods)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		//获取拣货id，根据拣货id查出 拣货单中 未复核完成的订单，不更新为欠货，
		//且 有未复核完成的订单 不发送到mq中，完成后再发送到mq中
		for _, p := range pickAndGoods {
			//已经被接单，且未完成复核
			if p.Status < model.ReviewCompletedStatus && p.PickUser != "" {
				pendingNumbers = append(pendingNumbers, p.Number)
				isSendMQ = false
			}
		}

		pendingNumbers = slice.UniqueStringSlice(pendingNumbers)

		diffSlice = slice.StrDiff(lackNumbers, pendingNumbers) // 在 lackNumbers 不在 pendingNumbers 中的

		if len(diffSlice) > 0 {
			//更新为欠货
			result = tx.Model(&model.Order{}).Where("number in (?)", diffSlice).Updates(map[string]interface{}{
				"order_type": model.LackOrderType,
			})

			if result.Error != nil {
				tx.Rollback()
				return result.Error
			}
		}

	}

	//保存完成订单
	if len(completeOrder) > 0 {
		result = tx.Save(&completeOrder)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	//保存完成订单详情
	if len(completeOrderDetail) > 0 {
		result = tx.Save(&completeOrderDetail)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	//更新pickOrder为已完成
	result = tx.Model(&model.PickOrder{}).Where("number in (?)", numbers).Updates(map[string]interface{}{
		"order_type": model.PickOrderCompleteOrderType,
	})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if isSendMQ {
		//mq 存入 批次id
		err := dao.SyncBatch(batchId)
		if err != nil {
			tx.Rollback()
			return errors.New("写入mq失败")
		}
	}

	return nil
}

// 编辑批次
func EditBatch(c *gin.Context) {
	var form req.EditBatchForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result := global.DB.Model(&model.Batch{}).Where("id = ?", form.Id).Updates(map[string]interface{}{"batch_name": form.BatchName})

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 推送u8 日志记录生成
func YongYouLog(tx *gorm.DB, pickGoods []model.PickGoods, orderAndGoods []rsp.OrderAndGoods, batchId int) error {
	mpOrderAndGoods := make(map[int]rsp.OrderAndGoods, 0)

	for _, order := range orderAndGoods {
		_, ok := mpOrderAndGoods[order.Id]
		if ok {
			continue
		}
		mpOrderAndGoods[order.Id] = order
	}

	mpPgv := make(map[string]req.PickGoodsView, 0)

	for _, good := range pickGoods {
		order, ogOk := mpOrderAndGoods[good.OrderGoodsId]
		if !ogOk {
			continue
		}

		//以拣货id和订单编号的纬度来推u8
		mpPgvKey := fmt.Sprintf("%v%v", good.PickId, good.Number)

		pgv, ok := mpPgv[mpPgvKey]

		if !ok {
			pgv = req.PickGoodsView{}
		}

		pgv.PickId = good.PickId
		pgv.SaleNumber = order.Number
		pgv.ShopId = int64(order.ShopId)
		pgv.ShopName = order.ShopName
		pgv.Date = order.PayAt
		pgv.Remark = order.OrderRemark
		pgv.DeliveryType = order.DistributionType //配送方式
		pgv.Line = order.Line
		pgv.List = append(pgv.List, req.PickGoods{
			GoodsName:    good.GoodsName,
			Sku:          good.Sku,
			Price:        int64(order.DiscountPrice),
			GoodsSpe:     good.GoodsSpe,
			Shelves:      good.Shelves,
			RealOutCount: good.ReviewNum,
			SlaveCode:    order.SaleCode,
			GoodsUnit:    order.GoodsUnit,
			SlaveUnit:    order.SaleUnit,
		})

		mpPgv[mpPgvKey] = pgv
	}

	var stockLogs = make([]model.StockLog, 0)

	for _, view := range mpPgv {
		//推送u8
		xml := dao.GenU8Xml(view, view.ShopId, view.ShopName, "05") //店铺属性中获 HouseCode

		stockLogs = append(stockLogs, model.StockLog{
			Number:      view.SaleNumber,
			BatchId:     batchId,
			PickId:      view.PickId,
			Status:      model.StockLogCreatedStatus, //已创建
			RequestXml:  xml,
			ResponseXml: "",
			ShopName:    view.ShopName,
		})
	}

	if len(stockLogs) > 0 {
		result := tx.Save(&stockLogs)
		if result.Error != nil {
			return result.Error
		}

		for _, log := range stockLogs {
			dao.YongYouProducer(log.Id)
		}
	}

	return nil
}

// 批次出库订单和商品明细
func GetBatchOrderAndGoods(c *gin.Context) {
	var form req.GetBatchOrderAndGoodsForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		batch          model.Batch
		pickOrder      []model.PickOrder
		pickOrderGoods []model.PickOrderGoods
		data           rsp.GetBatchOrderAndGoodsRsp
		mp             = make(map[string][]rsp.OutGoods)
		numbers        = make([]string, 0)
	)

	db := global.DB

	result := db.First(&batch, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//状态:0:进行中,1:已结束,2:暂停
	if batch.Status != 1 {
		xsq_net.ErrorJSON(c, errors.New("批次未结束"))
		return
	}

	result = db.Model(&model.PickOrderGoods{}).Where("batch_id = ?", form.Id).Find(&pickOrderGoods)

	totalGoodsNum := 0

	for _, good := range pickOrderGoods {
		//出库为0的不推送
		if good.OutCount == 0 {
			continue
		}

		totalGoodsNum++

		//编号 ，查询订单
		numbers = append(numbers, good.Number)

		_, ok := mp[good.Number]

		if !ok {
			mp[good.Number] = make([]rsp.OutGoods, 0)
		}

		mp[good.Number] = append(mp[good.Number], rsp.OutGoods{
			Id:            good.OrderGoodsId,
			Name:          good.GoodsName,
			Sku:           good.Sku,
			GoodsType:     good.GoodsType,
			GoodsSpe:      good.GoodsSpe,
			DiscountPrice: good.DiscountPrice,
			GoodsUnit:     good.GoodsUnit,
			SaleUnit:      good.SaleUnit,
			SaleCode:      good.SaleCode,
			OutCount:      good.OutCount,
			OutAt:         good.UpdateTime.Format(timeutil.TimeFormat),
			Number:        good.Number,
			CkNumber:      strings.Join(good.DeliveryOrderNo, ","),
		})

	}

	numbers = slice.UniqueStringSlice(numbers)

	result = db.Model(&model.PickOrder{}).Where("number in (?)", numbers).Find(&pickOrder)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.OutOrder, 0, len(pickOrder))

	for _, order := range pickOrder {
		goodsInfo, ok := mp[order.Number]

		if !ok {
			xsq_net.ErrorJSON(c, ecode.DataQueryError)
			return
		}

		list = append(list, rsp.OutOrder{
			DistributionType: order.DistributionType,
			PayAt:            order.PayAt,
			OrderId:          order.OrderId,
			GoodsInfo:        goodsInfo,
		})
	}

	data.Count = totalGoodsNum

	data.List = list

	xsq_net.SucJson(c, data)
}

// 当前批次是否有接单
func IsPick(c *gin.Context) {
	var (
		form   req.EndBatchForm
		pick   model.Pick
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

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var batches model.Batch

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
	result = db.Model(&model.Batch{}).Where("id = ? and status = ?", form.Id, form.Status).Update("status", updateStatus)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 获取批次列表
func GetBatchList(c *gin.Context) {
	var form req.GetBatchListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, res := dao.BatchList(global.DB, form)
	if err != nil {
		return
	}

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

	result := global.DB.Model(&model.Batch{}).
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
		form         req.GetBaseForm
		batch        model.Batch
		outboundTask model.OutboundTask
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	result := db.First(&batch, form.BatchId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.First(&outboundTask, batch.TaskId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	ret := rsp.GetBaseRsp{
		CreateTime:        timeutil.FormatToDateTime(batch.CreateTime),
		PayEndTime:        batch.PayEndTime,
		DeliveryStartTime: batch.DeliveryStartTime,
		DeliveryEndTime:   batch.DeliveryEndTime,
		DeliveryMethod:    batch.DeliveryMethod,
		Line:              batch.Line,
		Goods:             outboundTask.GoodsName,
		Status:            batch.Status,
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
		prePicks []model.PrePick
		//prePickGoods []batch.PrePickGoods
		prePickIds []int
	)

	db := global.DB

	result := db.Where("batch_id = ?", form.BatchId).Where(model.PrePick{ShopId: form.ShopId, Line: form.Line}).Where("status = 0").Find(&prePicks)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	db.Where("batch_id = ?", form.BatchId).Where(model.PrePick{ShopId: form.ShopId, Line: form.Line}).Where("status = 0").Scopes(model.Paginate(form.Page, form.Size)).Find(&prePicks)

	for _, pick := range prePicks {
		prePickIds = append(prePickIds, pick.Id)
	}

	retCount := []rsp.Ret{}

	result = db.Model(&model.PrePickGoods{}).
		Select("SUM(out_count) as out_c, SUM(need_num) AS need_c, shop_id, goods_type").
		Where("pre_pick_id in (?)", prePickIds).
		Where("status = 0"). //状态:0:未处理,1:已进入拣货池
		Group("shop_id, goods_type").
		Find(&retCount)

	typeMap := make(map[int]map[string]rsp.PickCount, 0)

	for _, r := range retCount {
		mp, ok := typeMap[r.ShopId]

		if !ok {
			mp = make(map[string]rsp.PickCount, 0)
		}

		mp[r.GoodsType] = rsp.PickCount{
			WaitingPick: r.NeedC,
			PickedCount: r.OutC,
		}

		typeMap[r.ShopId] = mp
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
		prePick       model.PrePick
		prePickGoods  []model.PrePickGoods
		prePickRemark []model.PrePickRemark
	)

	db := global.DB

	result := db.First(&prePick, form.PrePickId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.TaskName = prePick.ShopName
	res.Line = prePick.Line

	result = db.Where("pre_pick_id = ? and status = 0", form.PrePickId).Find(&prePickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	prePickGoodsSkuMp := make(map[string]rsp.MergePrePickGoods, 0)

	goodsNum := 0

	orderNumMp := make(map[string]struct{}, 0)

	//相同sku合并处理
	for _, goods := range prePickGoods {

		orderNumMp[goods.Number] = struct{}{}

		goodsNum += goods.NeedNum

		val, ok := prePickGoodsSkuMp[goods.Sku]

		paramsId := rsp.ParamsId{
			PickGoodsId:  goods.Id,
			OrderGoodsId: goods.OrderGoodsId,
		}

		if !ok {

			prePickGoodsSkuMp[goods.Sku] = rsp.MergePrePickGoods{
				Id:        goods.Id,
				Sku:       goods.Sku,
				GoodsName: goods.GoodsName,
				GoodsType: goods.GoodsType,
				GoodsSpe:  goods.GoodsSpe,
				Shelves:   goods.Shelves,
				NeedNum:   goods.NeedNum,
				Unit:      goods.Unit,
				ParamsId:  []rsp.ParamsId{paramsId},
			}
		} else {
			val.NeedNum += val.NeedNum
			val.ParamsId = append(val.ParamsId, paramsId)
			prePickGoodsSkuMp[goods.Sku] = val
		}
	}

	//订单数
	res.OrderNum = len(orderNumMp)

	//商品数
	res.GoodsNum = goodsNum

	goodsMap := make(map[string][]rsp.MergePrePickGoods, 0)

	for _, goods := range prePickGoodsSkuMp {

		goodsMap[goods.GoodsType] = append(goodsMap[goods.GoodsType], rsp.MergePrePickGoods{
			Id:        goods.Id,
			Sku:       goods.Sku,
			GoodsName: goods.GoodsName,
			GoodsType: goods.GoodsType,
			GoodsSpe:  goods.GoodsSpe,
			Shelves:   goods.Shelves,
			NeedNum:   goods.NeedNum,
			Unit:      goods.Unit,
			ParamsId:  goods.ParamsId,
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

// 置顶
func Topping(c *gin.Context) {
	var form req.ToppingForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	val, err := cache.Incr(constant.BATCH_TOPPING)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	result := global.DB.Model(model.Batch{}).Where("id = ?", form.Id).Update("sort", val)

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

	result := db.Model(&model.PrePick{}).Select("id").Where("batch_id = ? and status = 0", form.BatchId).Count(&count)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.Pick{}).
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

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err = c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	form.WarehouseId = c.GetInt("warehouseId")
	//类型：1:常规批次,2:快递批次,
	form.Typ = 1

	err = dao.BatchPick(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

// todo 废弃
func UpdateBatchPickNums(tx *gorm.DB, batchId int) error {

	//更新批次 预拣货单 拣货单 复核单 数
	var (
		prePickNum int64 //预拣货单
		pickNum    int
		reviewNum  int
	)

	result := tx.Model(&model.PrePick{}).Select("id").Where("batch_id = ? and status = 0", batchId).Count(&prePickNum)

	if result.Error != nil {
		return result.Error
	}

	type Count struct {
		Count  int `json:"count"`
		Status int `json:"status"`
	}

	var ct []Count

	result = tx.Model(&model.Pick{}).
		Select("count(id) as count,status").
		Where("batch_id = ? and status in (0,1)", batchId).
		Find(&ct)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	for _, c := range ct {
		switch c.Status {
		case 0:
			pickNum += c.Count
			break
		case 1:
			reviewNum += c.Count
			break
		}
	}

	var batch model.Batch

	result = tx.First(&batch, batchId)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Model(&model.Batch{}).Where("id = ? and version = ?", batchId, batch.Version).Updates(map[string]interface{}{
		"pre_pick_num":      prePickNum,
		"pick_num":          pickNum,
		"recheck_sheet_num": reviewNum,
		"version":           gorm.Expr("version + ?", 1),
	})

	if result.Error != nil {
		tx.Rollback()
		return errors.New("更新批次拣货单数等失败，请重试.错误:" + result.Error.Error())
	}

	return nil
}

// 合并拣货
func MergePick(c *gin.Context) {
	var form req.MergePickForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	form.WarehouseId = c.GetInt("warehouseId")

	form.Typ = 1

	err := dao.MergePick(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
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

	printCh := dao.GetPrintJobMap(form.HouseCode)

	//通道中没有任务
	if printCh == nil {
		xsq_net.SucJson(c, nil)
		return
	}

	var (
		pick          model.Pick
		pickGoods     []model.PickGoods
		orderAndGoods []rsp.OrderAndGoods
	)

	db := global.DB

	result := db.Model(&model.Pick{}).Where("delivery_no = ?", printCh.DeliveryOrderNo).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.PickGoods{}).Where("pick_id = ? and shop_id = ? and review_num > 0", pick.Id, printCh.ShopId).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	length := len(pickGoods) //有多少条pickGoods就有多少条OrderInfo数据，map数也是

	orderGoodsIds := make([]int, 0, length)

	goodsMp := make(map[int]model.PickGoods, length)

	for _, good := range pickGoods {
		orderGoodsIds = append(orderGoodsIds, good.OrderGoodsId)

		goodsMp[good.OrderGoodsId] = good
	}

	result = db.Table("t_pick_order_goods og").
		Select("og.*,o.shop_id,o.shop_name,o.shop_code,o.line,o.distribution_type,o.order_remark,o.pay_at,o.province,o.city,o.district,o.shop_type,o.latest_picking_time,o.house_code,o.consignee_name,o.consignee_tel").
		Joins("left join t_pick_order o on og.pick_order_id = o.id").
		Where("og.id in (?)", orderGoodsIds).
		Scan(&orderAndGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if len(orderAndGoods) <= 0 {
		xsq_net.ErrorJSON(c, ecode.OrderDataNotFound)
		return
	}

	packages := pick.Num

	//合并单不打印，ShopCode为空说明是合并单，合并单不会有店编
	if pick.ShopCode == "" {
		packages = 0
	}

	item := rsp.PrintCallGetRsp{
		ShopName:    pick.ShopName,
		JHNumber:    strconv.Itoa(pick.Id),
		PickName:    pick.PickUser, //拣货人
		ShopType:    orderAndGoods[0].ShopType,
		CheckName:   pick.ReviewUser,                                             //复核员
		HouseName:   TransferHouse(orderAndGoods[0].HouseCode),                   //TransferHouse(info.HouseCode)
		Delivery:    TransferDistributionType(orderAndGoods[0].DistributionType), //TransferDistributionType(info.DistributionType),
		OrderRemark: orderAndGoods[0].OrderRemark,
		Consignee:   orderAndGoods[0].ConsigneeName, //info.ConsigneeName
		Shop_code:   pick.ShopCode,
		Packages:    packages,
		Phone:       orderAndGoods[0].ConsigneeTel, //info.ConsigneeTel,
		PriType:     printCh.Type,
	}

	if orderAndGoods[0].ShopCode != "" {
		item.ShopName = orderAndGoods[0].ShopCode + "--" + orderAndGoods[0].ShopName
	}

	item2Mp := make(map[string]rsp.CallGetGoodsView, 0)

	for _, info := range orderAndGoods {

		pgs, ok := goodsMp[info.Id]

		if !ok {
			continue
		}

		item2val, item2ok := item2Mp[info.Number]

		if !item2ok {
			item2val = rsp.CallGetGoodsView{
				SaleNumber:  info.Number,
				Date:        info.PayAt,
				OrderRemark: info.OrderRemark,
			}
		}

		item3 := rsp.CallGetGoods{
			GoodsName:    info.GoodsName,
			GoodsSpe:     info.GoodsSpe,
			GoodsCount:   info.PayCount,
			RealOutCount: pgs.ReviewNum,
			GoodsUnit:    info.GoodsUnit,
			Price:        int64(info.DiscountPrice) * int64(pgs.ReviewNum),
			LackCount:    info.PayCount - pgs.ReviewNum,
		}
		item2val.List = append(item2val.List, item3)

		item2Mp[info.Number] = item2val
	}

	for _, item2 := range item2Mp {
		item.GoodsList = append(item.GoodsList, item2)
	}

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

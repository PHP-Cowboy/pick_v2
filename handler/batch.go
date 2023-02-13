package handler

import (
	"errors"
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
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"strconv"
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

	rdsKey := c.Request.URL.Path + strconv.Itoa(form.TaskId)

	err := cache.AntiRepeatedClick(rdsKey, 30)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, ecode.GetContextUserInfoFailed)
		return
	}

	form.Typ = 1 // 常规批次

	tx := global.DB.Begin()

	err = dao.CreateBatchByTask(tx, form, userInfo)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	//执行完成后删除锁定时间
	_, _ = cache.Del(rdsKey)

	tx.Commit()

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

	if form.Typ == 0 {
		//form.Typ = 1
		err := errors.New("typ 异常")
		xsq_net.ErrorJSON(c, err)
		return
	}

	tx := global.DB.Begin()

	err := dao.CreateBatch(tx, form, userInfo)

	if err != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, err)
		return
	}

	tx.Commit()

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

	err, res := dao.CentralizedPickDetail(global.DB, form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, res)
}

func GetUserInfo(c *gin.Context) *middlewares.CustomClaims {
	claims, ok := c.Get("claims")

	if !ok {
		return nil
	}

	return claims.(*middlewares.CustomClaims)
}

// 结束拣货批次
func EndBatch(c *gin.Context) {
	var form req.EndBatchForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := dao.EndBatch(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

// 后台拣货结束批次
func AdminEndBatch(c *gin.Context) {
	var form req.EndBatchForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := dao.AdminEndBatch(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
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

// 批次出库订单和商品明细
func GetBatchOrderAndGoods(c *gin.Context) {
	var form req.GetBatchOrderAndGoodsForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, res := dao.GetBatchOrderAndGoods(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, res)
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
		form      req.GetBatchPoolNumForm
		batchPool []rsp.BatchPoolNum
		res       rsp.GetBatchPoolNumRsp
		ongoing   int
		suspend   int
		finished  int
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result := global.DB.Model(&model.Batch{}).
		Select("count(id) as count, status").
		Where(model.Batch{Typ: form.Typ}).
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
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, res := dao.GetPrePickDetail(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

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

// 批量变更批次状态为 暂停||进行中
func BatchChangeBatch(c *gin.Context) {
	var form req.BatchChangeBatchForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	err := dao.BatchChangeBatch(db, *form.Status)

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
		pick      model.Pick
		pickGoods []model.PickGoods
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

	err, orderJoinGoods := model.GetOrderGoodsJoinOrderByIds(db, orderGoodsIds)
	if err != nil {
		return
	}

	if len(orderJoinGoods) <= 0 {
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
		ShopType:    orderJoinGoods[0].ShopType,
		CheckName:   pick.ReviewUser,                                              //复核员
		HouseName:   TransferHouse(orderJoinGoods[0].HouseCode),                   //TransferHouse(info.HouseCode)
		Delivery:    TransferDistributionType(orderJoinGoods[0].DistributionType), //TransferDistributionType(info.DistributionType),
		OrderRemark: orderJoinGoods[0].OrderRemark,
		Consignee:   orderJoinGoods[0].ConsigneeName, //info.ConsigneeName
		Shop_code:   pick.ShopCode,
		Packages:    packages,
		Phone:       orderJoinGoods[0].ConsigneeTel, //info.ConsigneeTel,
		PriType:     printCh.Type,
	}

	if orderJoinGoods[0].ShopCode != "" {
		item.ShopName = orderJoinGoods[0].ShopCode + "--" + orderJoinGoods[0].ShopName
	}

	item2Mp := make(map[string]rsp.CallGetGoodsView, 0)

	for _, info := range orderJoinGoods {

		pgs, ok := goodsMp[info.OrderGoodsId]

		if !ok {
			continue
		}

		item2val, item2ok := item2Mp[info.Number]

		if !item2ok {
			item2val = rsp.CallGetGoodsView{
				SaleNumber:  info.Number,
				Date:        timeutil.FormatToDateTime(time.Time(*info.PayAt)),
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

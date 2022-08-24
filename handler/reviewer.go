package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pick_v2/common/constant"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model/batch"
	"pick_v2/model/order"
	"pick_v2/utils/ecode"
	"pick_v2/utils/str_util"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"strconv"
	"time"
)

// 复核列表 通过状态区分是否已完成
func ReviewList(c *gin.Context) {
	var (
		form       req.ReviewListReq
		res        rsp.ReviewListRsp
		pick       []batch.Pick
		pickRemark []batch.PickRemark
		//pickListModel []rsp.PickListModel
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	//result := db.Model(&batch.Pick{}).Where("status = ?", form.Status).Where(batch.Pick{PickUser: form.Name}).Find(&pick)

	localDb := db.Model(&batch.Pick{}).Where("status = ?", form.Status)

	claims, ok := c.Get("claims")

	if !ok {
		xsq_net.ErrorJSON(c, errors.New("claims 获取失败"))
		return
	}

	userInfo := claims.(*middlewares.CustomClaims)

	//1:待复核,2:复核完成
	if form.Status == 1 {
		localDb.Where("review_user = ? or review_user = ''", userInfo.Name)
	} else {
		localDb.Where("review_user = ? ", userInfo.Name)
	}

	result := localDb.Where(batch.Pick{PickUser: form.Name}).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	//拣货ids
	pickIds := make([]int, 0, len(pick))
	for _, p := range pick {
		pickIds = append(pickIds, p.Id)
	}

	//拣货ids 的订单备注
	result = db.Where("pick_id in (?)", pickIds).Find(&pickRemark)
	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//构建pickId 对应的订单 是否有备注map
	remarkMp := make(map[int]struct{}, 0) //key 存在即为有
	for _, remark := range pickRemark {
		remarkMp[remark.PickId] = struct{}{}
	}

	isRemark := false

	list := make([]rsp.Pick, 0, form.Size)

	for _, p := range pick {

		_, ok := remarkMp[p.Id]
		if ok { //拣货id在拣货备注中存在，即为有备注
			isRemark = true
		}

		reviewTime := ""

		if p.ReviewTime != nil {
			reviewTime = p.ReviewTime.Format(timeutil.TimeFormat)
		}

		list = append(list, rsp.Pick{
			Id:             p.Id,
			TaskName:       p.TaskName,
			ShopCode:       p.ShopCode,
			ShopName:       p.ShopName,
			ShopNum:        p.ShopNum,
			OrderNum:       p.OrderNum,
			NeedNum:        p.NeedNum,
			PickUser:       p.PickUser,
			TakeOrdersTime: p.TakeOrdersTime.Format(timeutil.TimeFormat),
			IsRemark:       isRemark,
			PickNum:        p.PickNum,
			ReviewNum:      p.ReviewNum,
			Num:            p.Num,
			ReviewUser:     p.ReviewUser,
			ReviewTime:     reviewTime,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// 复核明细
func ReviewDetail(c *gin.Context) {
	var (
		form req.ReviewDetailReq
		res  rsp.ReviewDetailRsp
		pick batch.Pick
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	result := db.First(&pick, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//复核员为空，接单复核
	if pick.ReviewUser == "" {
		claims, ok := c.Get("claims")

		if !ok {
			xsq_net.ErrorJSON(c, errors.New("claims 获取失败"))
			return
		}

		userInfo := claims.(*middlewares.CustomClaims)

		result = db.Model(&batch.Pick{}).
			Where("id = ? and version = ?", pick.Id, pick.Version).
			Updates(map[string]interface{}{"review_user": userInfo.Name, "version": pick.Version + 1})

		// todo 需要模拟测试
		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		//更新复核单数量
		result = db.Model(batch.Batch{}).Where("id = ?", pick.BatchId).Update("recheck_sheet_num", gorm.Expr("recheck_sheet_num + ?", 1))

		if result.Error != nil {
			xsq_net.ErrorJSON(c, ecode.DataSaveError)
			return
		}
	}

	var (
		pickGoods  []batch.PickGoods
		pickRemark []batch.PickRemark
	)

	res.TaskName = pick.TaskName
	res.ShopCode = pick.ShopCode
	res.PickUser = pick.PickUser
	res.TakeOrdersTime = pick.TakeOrdersTime.Format(timeutil.TimeFormat)
	res.ReviewUser = pick.ReviewUser

	var reviewTime string

	if pick.ReviewTime != nil {
		reviewTime = pick.ReviewTime.Format(timeutil.TimeFormat)
	}
	res.ReviewTime = reviewTime

	result = db.Where("pick_id = ?", form.Id).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	goodsMap := make(map[string][]rsp.PickGoods, 0)

	needTotal := 0
	completeTotal := 0
	reviewTotal := 0
	for _, goods := range pickGoods {
		completeTotal += goods.CompleteNum
		needTotal += goods.NeedNum
		reviewTotal += goods.ReviewNum
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
	res.ReviewTotal = reviewTotal

	res.Goods = goodsMap

	result = db.Where("pick_id = ?", form.Id).Find(&pickRemark)

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

// 确认出库
func ConfirmDelivery(c *gin.Context) {
	var (
		form req.ConfirmDeliveryReq
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		pick      batch.Pick
		pickGoods []batch.PickGoods
		batches   batch.Batch
	)

	db := global.DB

	//根据id获取拣货数据
	result := db.First(&pick, form.Id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			xsq_net.ErrorJSON(c, ecode.DataNotExist)
			return
		}
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if pick.Status == 2 {
		xsq_net.ErrorJSON(c, ecode.OrderHasBeenReviewedAndCompleted)
		return
	}

	claims, ok := c.Get("claims")

	if !ok {
		xsq_net.ErrorJSON(c, errors.New("claims 获取失败"))
		return
	}

	userInfo := claims.(*middlewares.CustomClaims)

	if pick.ReviewUser != userInfo.Name {
		xsq_net.ErrorJSON(c, ecode.DataNotExist)
		return
	}

	//复核拣货数量数据 构造map 最终构造成db的save参数
	pickGoodsMap := make(map[int]int, len(form.CompleteReview))

	for _, cp := range form.CompleteReview {
		pickGoodsMap[cp.Id] = cp.ReviewNum
	}

	//获取拣货商品数据
	result = db.Where("pick_id = ?", form.Id).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//构造打印 chan 结构体数据
	printChMp := make(map[int]struct{}, 0)

	//构造更新 orderInfo 表完成出库数据
	orderInfoIds := []int{}
	for k, pg := range pickGoods {
		_, printChOk := printChMp[pg.ShopId]

		if !printChOk {
			printChMp[pg.ShopId] = struct{}{}
		}

		num, pgOk := pickGoodsMap[pg.Id]

		if !pgOk {
			continue
		}

		pickGoods[k].ReviewNum = num

		if pg.NeedNum == num { //需拣和复核数一致，即为完成
			orderInfoIds = append(orderInfoIds, pg.OrderInfoId)
		}

	}

	now := time.Now()

	dateStr := now.Format(timeutil.DateIntFormat)

	//redis
	redis := global.Redis

	//出库单号key
	redisKey := constant.DELIVERY_ORDER_NO + dateStr

	val, err := redis.Do(context.Background(), "incr", redisKey).Result()
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	number := strconv.Itoa(int(val.(int64)))

	//填充出库单号编码
	deliveryOrderNo := str_util.StrPad(number, 3, "0", 0)

	tx := db.Begin()

	if len(orderInfoIds) > 0 { //有完成的
		//更新orderInfo表数据为完成状态
		result = tx.Model(&order.OrderInfo{}).
			Where("id in (?)", orderInfoIds).
			Updates(map[string]interface{}{
				"delivery_order_no": deliveryOrderNo,
				"pick_time":         now.Format(timeutil.TimeFormat),
			}) //
		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	}

	//更新主表
	result = tx.Model(&batch.Pick{}).Where("id = ?", pick.Id).Updates(map[string]interface{}{"status": 2, "review_time": &now, "num": form.Num, "delivery_order_no": deliveryOrderNo})

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//更新拣货商品数据
	result = tx.Save(&pickGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//拆单 -打印
	for shopId, _ := range printChMp {
		AddPrintJobMap(constant.JH_HUOSE_CODE, &global.PrintCh{
			DeliveryOrderNo: deliveryOrderNo,
			ShopId:          shopId,
		})
	}

	result = db.First(&batches, pick.BatchId)
	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if batches.Status == 1 {
		err = PushU8(pickGoods)
		if err != nil {
			tx.Commit() // u8推送失败不能影响仓库出货，只提示，业务继续
			xsq_net.ErrorJSON(c, err)
			return
		}
	}

	tx.Commit()

	xsq_net.Success(c)
}

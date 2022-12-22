package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/model"
	"pick_v2/utils/slice"
)

// 关闭订单状态数量统计
func CloseOrderCount(db *gorm.DB) (err error, res rsp.CloseOrderCountRsp) {
	var (
		countList []model.CountCloseOrder
	)

	err, countList = model.CountCloseOrderStatus(db)
	if err != nil {
		return
	}

	for _, cl := range countList {
		switch cl.Status {
		//进行中
		case model.CloseOrderStatusPending:
			res.PendingNum = cl.Count
			break

		//已完成
		case model.CloseOrderStatusComplete:
			res.CompleteNum = cl.Count
			break

		//异常
		case model.CloseOrderStatusException:
			res.ExceptionNum = cl.Count
			break
		}
	}

	return
}

// 关闭订单列表
func CloseOrderList(db *gorm.DB, form req.CloseOrderListForm) (err error, res rsp.CloseOrderListRsp) {
	var (
		count       int64
		closeOrders []model.CloseOrder
	)

	cond := model.CloseOrder{Number: form.Number, Status: form.Status}

	err, count = model.CountCloseOrderByCond(db, cond)

	if err != nil {
		return
	}

	list := make([]rsp.CloseOrderList, 0, form.Size)

	res.Total = count

	if count == 0 {
		res.List = list
		return
	}

	err, closeOrders = model.GetCloseOrderPageList(db, cond, form.Page, form.Size)

	if err != nil {
		return
	}

	for _, order := range closeOrders {
		list = append(list, rsp.CloseOrderList{
			Id:               order.Id,
			Number:           order.Number,
			PayAt:            order.PayAt,
			PayTotal:         order.PayTotal,
			NeedCloseTotal:   order.NeedCloseTotal,
			ShopName:         order.ShopName,
			ShopType:         order.ShopType,
			DistributionType: order.DistributionType,
			Province:         order.Province,
			City:             order.City,
			District:         order.District,
			OrderRemark:      order.OrderRemark,
			Status:           order.Status,
		})
	}

	res.List = list

	return
}

// 关闭订单详情
func CloseOrderDetail(db *gorm.DB, form req.CloseOrderDetailForm) (err error, res rsp.CloseOrderDetailRsp) {
	var (
		closeOrder model.CloseOrder
		closeGoods []model.CloseGoods
	)

	err, closeOrder = model.GetCloseOrderByPk(db, form.Id)

	if err != nil {
		return
	}

	err, closeGoods = model.GetCloseGoodsListByCond(db, model.CloseGoods{CloseOrderId: closeOrder.Id})

	if err != nil {
		return
	}

	list := make([]rsp.CloseGoodsList, 0, len(closeGoods))

	for _, good := range closeGoods {
		list = append(list, rsp.CloseGoodsList{
			GoodsName:      good.GoodsName,
			GoodsSpe:       good.GoodsSpe,
			PayCount:       good.PayCount,
			CloseCount:     good.CloseCount,
			OutCount:       0,
			NeedCloseCount: good.NeedCloseCount,
			GoodsRemark:    good.GoodsRemark,
		})
	}

	res = rsp.CloseOrderDetailRsp{
		Number:           closeOrder.Number,
		ShopName:         closeOrder.ShopName,
		DistributionType: closeOrder.DistributionType,
		District:         closeOrder.District,
		Status:           closeOrder.Status,
		OrderRemark:      closeOrder.OrderRemark,
		List:             list,
	}

	return
}

// 关闭订单&&详情列表
func CloseOrderAndGoodsList(db *gorm.DB) (err error, res []rsp.CloseOrderAndGoodsList) {

	var (
		closeOrders   []model.CloseOrder
		closeOrderIds []int
		closeGoods    []model.CloseGoods
		closeGoodsMp  = make(map[int][]rsp.CloseGoodsList, 0)
	)

	err, closeOrders = model.GetCloseOrderList(db, model.CloseOrder{Status: model.CloseOrderStatusPending})

	if err != nil {
		return
	}

	for _, order := range closeOrders {
		closeOrderIds = append(closeOrderIds, order.Id)
	}

	err, closeGoods = model.GetCloseGoodsListByCloseOrderIds(db, closeOrderIds)

	if err != nil {
		return
	}

	for _, good := range closeGoods {
		closeGoodsMpVal, closeGoodsMpOk := closeGoodsMp[good.CloseOrderId]

		if !closeGoodsMpOk {
			closeGoodsMpVal = []rsp.CloseGoodsList{}
		}

		closeGoodsMpVal = append(closeGoodsMpVal, rsp.CloseGoodsList{
			GoodsName:      good.GoodsName,
			GoodsSpe:       good.GoodsSpe,
			PayCount:       good.PayCount,
			CloseCount:     good.CloseCount,
			OutCount:       0,
			NeedCloseCount: good.NeedCloseCount,
			GoodsRemark:    good.GoodsRemark,
		})

		closeGoodsMp[good.CloseOrderId] = closeGoodsMpVal
	}

	for _, order := range closeOrders {
		list, closeGoodsMpOk := closeGoodsMp[order.Id]

		if !closeGoodsMpOk {
			list = make([]rsp.CloseGoodsList, 0, 0)
		}

		res = append(res, rsp.CloseOrderAndGoodsList{
			Number:      order.Number,
			ShopName:    order.ShopName,
			OrderRemark: order.OrderRemark,
			List:        list,
		})
	}

	return
}

// 关闭预拣池
func ClosePrePick(tx *gorm.DB, closeGoodsMp map[int]int, number string, taskId int) (err error, prePickGoodsIds []int) {
	var (
		prePickGoodsJoinPrePick []model.PrePickGoodsJoinPrePick
		prePickGoodsUpdate      []model.PrePickGoods
		prePickIds              []int
		notCloseAllPrePickIds   []int
		batchId                 int
		batch                   model.Batch
	)

	//根据任务ID查询全部预拣池任务和预拣池商品数据
	err, prePickGoodsJoinPrePick = model.GetPrePickGoodsJoinPrePickListByTaskId(tx, taskId, number)

	if err != nil {
		return
	}

	//一个订单只会在一个进行中的任务里，且也只会在这个任务下的某一个批次中
	batchId = prePickGoodsJoinPrePick[0].BatchId

	err, batch = model.GetBatchByPk(tx, batchId)
	if err != nil {
		return
	}

	//不是暂停中的批次不处理[已结束的不能处理，会导致u8数据错误，进行中是处理前要把所有的批次暂停]
	if batch.Status != model.BatchSuspendStatus {
		err = errors.New(fmt.Sprintf("批次状态异常:%d", batch.Status))
		return
	}

	for _, good := range prePickGoodsJoinPrePick {
		prePickIds = append(prePickIds, good.PrePickId)

		closeCount, closeGoodsMpOk := closeGoodsMp[good.OrderGoodsId]

		if !closeGoodsMpOk {
			//没被关闭的
			notCloseAllPrePickIds = append(notCloseAllPrePickIds, good.PrePickId)
			continue
		}

		updatePrePickGoodsData := model.PrePickGoods{
			Base:     model.Base{Id: good.PrePickGoodsId},
			NeedNum:  good.NeedNum - closeCount,
			CloseNum: good.CloseNum + closeCount,
		}

		prePickGoodsIds = append(prePickGoodsIds, good.PrePickGoodsId)

		//需拣如果小于零是有问题的，说明关闭数量大于需拣数量
		if updatePrePickGoodsData.NeedNum < 0 {
			err = errors.New(fmt.Sprintf(
				"需拣变更err:t_pre_pick_goods.id:%d,need_num:%d,closeCount:%d",
				good.PrePickGoodsId,
				good.NeedNum,
				closeCount,
			))
			return
		}

		if updatePrePickGoodsData.NeedNum == 0 {
			updatePrePickGoodsData.Status = model.PrePickGoodsStatusClose
		} else {
			//没有被全部关闭的
			notCloseAllPrePickIds = append(notCloseAllPrePickIds, good.PrePickId)
		}
		//需要被更新的预拣池商品数据
		prePickGoodsUpdate = append(prePickGoodsUpdate, updatePrePickGoodsData)

	}

	err = model.PrePickGoodsReplaceSave(tx, prePickGoodsUpdate, []string{"need_num", "close_num", "status"})

	if err != nil {
		return
	}

	prePickIds = slice.UniqueSlice(prePickIds)
	notCloseAllPrePickIds = slice.UniqueSlice(notCloseAllPrePickIds)

	prePickIds = slice.Diff(prePickIds, notCloseAllPrePickIds)

	// 如果预拣池全部被关闭，则更新预拣池状态
	if len(prePickIds) > 0 {
		err = model.UpdatePrePickStatusByIds(tx, prePickIds, model.PrePickStatusClose)

		if err != nil {
			return
		}
	}

	return
}

// 关闭集中拣货 这一版先不做
func CloseCentralizedPick(tx *gorm.DB) (err error) {
	return
}

// 关闭拣货池
func ClosePick(tx *gorm.DB, closeGoodsMp map[int]int, prePickGoodsIds []int) (err error) {
	var (
		pickGoods       []model.PickGoods
		pickGoodsUpdate []model.PickGoods
		pickIds         []int
	)

	err, pickGoods = model.GetPickGoodsByOrderGoodsIds(tx, prePickGoodsIds)

	if err != nil {
		return
	}

	//一个品可能被拣多次，但只有结束后才能进入新的任务、批次。
	//更新拣货池id最大的 map[order_goods_id]id
	//TODO 第二次拣货，只进入了预拣池时，这样更新的是第一次拣货的数据 预拣池返回map[order_good_id]pre_pick_id
	pickCloseMp := make(map[int]int, 0)

	for _, good := range pickGoods {

		id, pickCloseMpOk := pickCloseMp[good.OrderGoodsId]

		if !pickCloseMpOk || good.Id > id {
			pickCloseMp[good.OrderGoodsId] = good.Id
		}

	}

	for _, good := range pickGoods {
		id, pickCloseMpOk := pickCloseMp[good.OrderGoodsId]

		if !pickCloseMpOk {
			err = errors.New("预拣池商品map异常")
			return
		}
		//订单商品id与对应的预拣池id不匹配，则为历史预拣池数据，跳过
		if id != good.Id {
			continue
		}

		closeCount, closeGoodsMpOk := closeGoodsMp[good.OrderGoodsId]

		if !closeGoodsMpOk {
			err = errors.New("订单商品map异常")
			return
		}

		updatePickGoodsData := model.PickGoods{
			Base:     model.Base{Id: id},
			NeedNum:  good.NeedNum - closeCount,
			CloseNum: good.CloseNum + closeCount,
		}

		if updatePickGoodsData.NeedNum <= 0 {
			updatePickGoodsData.Status = model.PickGoodsStatusClosed
		}
		//需要被更新的预拣池商品数据
		pickGoodsUpdate = append(pickGoodsUpdate, updatePickGoodsData)

		pickIds = append(pickIds, good.PickId)
	}

	err = model.PickGoodsReplaceSave(tx, &pickGoodsUpdate, []string{"need_num", "close_num", "complete_num", "status"})

	if err != nil {
		return
	}

	// 如果预拣池全部被关闭，则更新预拣池状态
	return
}

// 关单处理
func CloseOrderExecNew(db *gorm.DB, form req.CloseOrder) (err error) {
	var (
		closeOrders   []model.CloseOrder
		closeGoods    []model.CloseGoods
		orderGoodsIds []int
		closeGoodsMp  = make(map[int]int, 0)
	)

	//校验是否所有批次全部暂停
	//查询是否有进行中的批次
	err, _ = model.GetBatchFirstByStatus(db, model.BatchOngoingStatus)

	//err 不是未找到数据
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	//关闭订单数据
	err, closeOrders = model.GetCloseOrderByNumbers(db, form.Number)

	if err != nil {
		return
	}

	//关闭订单商品数据
	err, closeGoods = model.GetCloseGoodsListByNumbers(db, form.Number)

	if err != nil {
		return
	}

	for _, good := range closeGoods {
		orderGoodsIds = append(orderGoodsIds, good.OrderGoodsId)

		closeGoodsMp[good.OrderGoodsId] = good.CloseCount
	}

	for _, co := range closeOrders {
		tx := db.Begin()

		err = CloseOrderHandle(tx, co.Number, co.Typ, closeGoodsMp, orderGoodsIds)

		if err != nil {
			tx.Rollback()
			continue
		}

		//更新关单任务状态
		err = model.UpdateCloseOrderByPk(tx, co.Id, map[string]interface{}{"status": model.CloseOrderStatusComplete})

		if err != nil {
			tx.Rollback()
			continue
		}

		err = SendMsgQueue("close_order_result", []string{co.Number})

		if err != nil {
			tx.Rollback()
			continue
		}

		tx.Commit()

	}

	return
}

// 关单处理
func CloseOrderHandle(tx *gorm.DB, number string, typ int, closeGoodsMp map[int]int, orderGoodsIds []int) (err error) {
	var (
		isCommit bool //用于判断是否提交事务，是否还需要继续执行后续流程
		taskId   int
	)

	err, isCommit, taskId = OrderDataHandle(tx, number, typ, closeGoodsMp)

	if err != nil || isCommit {
		return
	}

	//批次数据处理
	err, isCommit = BatchDataHandle(tx, closeGoodsMp, number, taskId)

	if err != nil || isCommit {
		return
	}

	return
}

// 订单数据处理(包括订单相关表和任务订单相关表)
func OrderDataHandle(tx *gorm.DB, number string, typ int, closeGoodsMp map[int]int) (err error, isCommit bool, taskId int) {
	//关闭订单&&订单商品逻辑
	err, isCommit = CloseGoodsAndOrderLogic(tx, number, typ, closeGoodsMp)

	if err != nil || isCommit {
		return
	}

	//关闭出库任务订单&&订单商品
	err, isCommit, taskId = CloseTaskLogic(tx, number, typ, closeGoodsMp)

	if err != nil || isCommit {
		return
	}

	return
}

// 批次数据处理
func BatchDataHandle(tx *gorm.DB, closeGoodsMp map[int]int, number string, taskId int) (err error, isCommit bool) {
	var prePickGoodsIds []int

	err, prePickGoodsIds = ClosePrePick(tx, closeGoodsMp, number, taskId)
	if err != nil {
		return
	}

	err = ClosePick(tx, closeGoodsMp, prePickGoodsIds)

	if err != nil {
		return
	}
	return
}

// 关闭订单以及订单商品逻辑
func CloseGoodsAndOrderLogic(tx *gorm.DB, number string, typ int, closeGoodsMp map[int]int) (err error, isCommit bool) {

	var (
		order      model.Order
		orderGoods []model.OrderGoods
	)

	err, order = model.GetOrderByNumber(tx, number)

	if err != nil {
		return
	}

	//全单关闭
	if typ == model.CloseOrderTypAll {
		//全单关闭时更新订单类型为关闭
		err = model.UpdateOrderByNumbers(tx, []string{number}, map[string]interface{}{"order_type": model.CloseOrderType})
		if err != nil {
			return
		}
	}

	//根据订单号查询全部订单商品数据
	err, orderGoods = model.GetOrderGoodsListByNumbers(tx, []string{number})

	if err != nil {
		return
	}

	//订单商品关闭
	for i, og := range orderGoods {
		closeCount, closeGoodsMpOk := closeGoodsMp[og.Id]

		if !closeGoodsMpOk {
			err = errors.New("订单商品map异常")
			return
		}

		if og.LackCount < closeCount {
			err = errors.New("欠货数量小于关闭数量")
			return
		}

		orderGoods[i].CloseCount += closeCount
		orderGoods[i].LackCount -= closeCount
	}

	err = model.OrderGoodsReplaceSave(tx, &orderGoods, []string{"update_time", "lack_count", "close_count"})

	if err != nil {
		return
	}

	//如果是新订单，则处理完订单后可以直接提交，不做后续逻辑处理
	if order.OrderType == model.NewOrderType {
		isCommit = true
	}

	return
}

// 关闭任务订单以及商品逻辑
func CloseTaskLogic(tx *gorm.DB, number string, typ int, closeGoodsMp map[int]int) (err error, isCommit bool, taskId int) {

	var (
		outboundOrder model.OutboundOrder
		closeGoods    []model.CloseGoods
		outboundGoods []model.OutboundGoods
	)

	//获取订单所在的最新任务
	err, outboundOrder = model.GetOutboundOrderByNumberFirstSortByTaskId(tx, number)

	if err != nil {
		return
	}

	if outboundOrder.OrderType == model.OutboundOrderTypeNew {
		isCommit = true
	}

	taskId = outboundOrder.TaskId

	//全单关闭
	if typ == model.CloseOrderTypAll {

		//更新任务订单为关闭
		err = model.OutboundOrderBatchUpdate(
			tx,
			&model.OutboundOrder{
				TaskId: taskId,
				Number: number,
			},
			map[string]interface{}{"order_type": model.OutboundOrderTypeClose},
		)

		if err != nil {
			return
		}

	}

	//关闭订单商品数据
	err, closeGoods = model.GetCloseGoodsListByNumbers(tx, []string{number})

	if err != nil {
		return
	}

	//构造任务商品查询条件
	skus := make([]string, 0, len(closeGoods))

	for _, cg := range closeGoods {
		skus = append(skus, cg.Sku)
	}

	err, outboundGoods = model.GetOutboundGoodsListByPks(tx, taskId, number, skus)

	if err != nil {
		return
	}

	for i, good := range outboundGoods {
		closeCount, closeGoodsMpOk := closeGoodsMp[good.OrderGoodsId]

		if !closeGoodsMpOk {
			err = errors.New("订单商品map异常")
			return
		}

		if good.LackCount < closeCount {
			err = errors.New(fmt.Sprintf("欠货数量小于关闭数量,taskId:%d,number:%s,sku:%s", taskId, number, good.Sku))
			return
		}

		outboundGoods[i].CloseCount += closeCount
		outboundGoods[i].LackCount -= closeCount
	}

	err = model.OutboundGoodsReplaceSave(tx, &outboundGoods, []string{"lack_count", "out_count"})

	if err != nil {
		return
	}
	return
}

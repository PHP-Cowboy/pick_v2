package dao

import (
	"errors"
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/model"
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

	list := make([]rsp.OrderGoodsList, 0, len(closeGoods))

	for _, good := range closeGoods {
		list = append(list, rsp.OrderGoodsList{
			GoodsName:      good.GoodsName,
			GoodsSpe:       good.GoodsSpe,
			PayCount:       good.PayCount,
			CloseCount:     good.CloseCount,
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

// 关闭订单列表&&详情

func CloseOrderExec(db *gorm.DB, form req.CloseOrder) (err error) {
	var (
		outboundOrders      []model.OutboundOrder
		closeOrders         []model.CloseOrder
		closeGoods          []model.CloseGoods
		closeOrderNumberTyp []rsp.CloseOrderNumberTyp
		orderGoodsIds       []int
		closeGoodsMp        = make(map[int]int, 0)
	)

	//校验是否所有批次全部暂停
	//查询是否有进行中的批次
	err, _ = model.GetBatchFirst(db, model.Batch{Status: model.BatchOngoingStatus})

	//err 不是未找到数据
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	//关闭订单数据
	err, closeOrders = model.GetCloseOrderByNumbers(db, form.Number)

	if err != nil {
		return
	}

	for _, order := range closeOrders {
		closeOrderNumberTyp = append(closeOrderNumberTyp, rsp.CloseOrderNumberTyp{
			Number: order.Number,
			Typ:    order.Typ,
		})
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

	//MQ 中 处理新订单关闭 [包括订单中的新订单和任务中的新订单]
	err, outboundOrders = model.GetOutboundOrderByNumbers(db, form.Number)

	if err != nil {
		return
	}

	//map[number]taskId
	numberTaskMp := make(map[string]int, 0)

	for _, order := range outboundOrders {
		//任务中的新订单，已在关闭订单消息中处理掉

		taskId, numberTaskMpOk := numberTaskMp[order.Number]
		//利用主键自增，找到订单号对应的最大批次
		//订单每次只能在一个任务中，任务结束后变为欠货单才能重新进入下一个任务中
		if !numberTaskMpOk || order.TaskId > taskId {
			numberTaskMp[order.Number] = order.TaskId
		}
	}

	taskOrderCond := [][]interface{}{}

	//构造任务订单查询条件数据
	for no, taskId := range numberTaskMp {
		cond := []interface{}{taskId, no}
		taskOrderCond = append(taskOrderCond, cond)
	}

	//构造任务商品查询条件
	taskOrderSkuCond := [][]interface{}{}
	for _, cg := range closeGoods {
		taskId, numberTaskMpOk := numberTaskMp[cg.Number]

		if !numberTaskMpOk {
			err = errors.New("关闭商品数据有误")
			return
		}

		cond := []interface{}{taskId, cg.Number, cg.Sku}
		taskOrderSkuCond = append(taskOrderSkuCond, cond)
	}

	tx := db.Begin()

	err = CloseOrderAndGoods(tx, closeOrderNumberTyp, orderGoodsIds, closeGoodsMp, taskOrderCond, taskOrderSkuCond)
	if err != nil {
		tx.Rollback()
		return
	}

	err = ClosePrePick(tx, orderGoodsIds, closeGoodsMp)
	if err != nil {
		tx.Rollback()
		return
	}

	err = ClosePick(tx, orderGoodsIds, closeGoodsMp)
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return
}

// 订单&&订单商品关闭
func CloseOrderAndGoods(
	tx *gorm.DB,
	closeOrderNumberTyp []rsp.CloseOrderNumberTyp,
	orderGoodsIds []int,
	closeGoodsMp map[int]int,
	taskOrderCond,
	taskOrderSkuCond [][]interface{},
) (
	err error,
) {
	if len(orderGoodsIds) == 0 {
		err = errors.New("被关闭商品不能为空")
		return
	}

	var (
		orderGoods []model.OrderGoods
	)

	err, orderGoods = model.GetOrderGoodsListByIds(tx, orderGoodsIds)

	if err != nil {
		return
	}

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

	//全单关闭订单号
	var (
		closeOrderNumbers   []string
		closeOrderNumbersMp = make(map[string]struct{}, 0)
	)

	for _, nt := range closeOrderNumberTyp {
		if nt.Typ == 2 {
			closeOrderNumbers = append(closeOrderNumbers, nt.Number)
			closeOrderNumbersMp[nt.Number] = struct{}{}
		}
	}

	//全单关闭
	if len(closeOrderNumbers) == 0 {
		err = model.UpdateOrderByNumbers(tx, closeOrderNumbers, map[string]interface{}{"order_type": model.CloseOrderType})
		if err != nil {
			return
		}
	}

	//不需要处理任务数据
	if taskOrderCond == nil {
		return
	}

	var (
		outboundOrder       []model.OutboundOrder
		outboundOrderUpdate []model.OutboundOrder
	)

	err, outboundOrder = model.GetOutboundOrderInMultiColumn(tx, taskOrderCond)
	if err != nil {
		return
	}

	if len(outboundOrder) == 0 {
		err = errors.New("任务订单未找到")
		return
	}

	for _, order := range outboundOrder {
		_, closeOrderNumbersMpOk := closeOrderNumbersMp[order.Number]

		if closeOrderNumbersMpOk {
			outboundOrderUpdate = append(outboundOrderUpdate, order)
		}
	}

	if len(outboundOrderUpdate) > 0 {
		err = model.OutboundOrderReplaceSave(tx, outboundOrderUpdate, []string{"order_type"})

		if err != nil {
			return
		}
	}

	var (
		outboundGoods []model.OutboundGoods
	)

	err, outboundGoods = model.GetOutboundGoodsInMultiColumn(tx, taskOrderSkuCond)

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
			err = errors.New("欠货数量小于关闭数量")
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

// 关闭预拣池
func ClosePrePick(
	tx *gorm.DB,
	orderGoodsIds []int,
	closeGoodsMp map[int]int,
) (
	err error,
) {
	var (
		prePickGoods       []model.PrePickGoods
		prePickGoodsUpdate []model.PrePickGoods
		prePickIds         []int
	)

	err, prePickGoods = model.GetPrePickGoodsByOrderGoodsIds(tx, orderGoodsIds)

	if err != nil {
		return
	}

	//一个品可能被拣多次，但只有结束后才能进入新的任务、批次。
	//更新预拣池id最大的 map[order_goods_id]id
	prePickCloseMp := make(map[int]int, 0)

	for _, good := range prePickGoods {

		id, prePickCloseMpOk := prePickCloseMp[good.OrderGoodsId]

		if !prePickCloseMpOk || good.Id > id {
			prePickCloseMp[good.OrderGoodsId] = good.Id
		}

	}

	for _, good := range prePickGoods {
		id, prePickCloseMpOk := prePickCloseMp[good.OrderGoodsId]

		if !prePickCloseMpOk {
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

		updatePrePickGoodsData := model.PrePickGoods{
			Base:     model.Base{Id: id},
			NeedNum:  good.NeedNum - closeCount,
			CloseNum: good.CloseNum + closeCount,
		}

		if updatePrePickGoodsData.NeedNum <= 0 {
			updatePrePickGoodsData.Status = model.PrePickGoodsStatusClose
		}
		//需要被更新的预拣池商品数据
		prePickGoodsUpdate = append(prePickGoodsUpdate, updatePrePickGoodsData)

		prePickIds = append(prePickIds, good.PrePickId)
	}

	err = model.PrePickGoodsReplaceSave(tx, prePickGoodsUpdate, []string{"need_num", "close_num", "status"})

	if err != nil {
		return
	}

	//todo 如果预拣池全部被关闭，则更新预拣池状态

	return
}

// 关闭集中拣货 这一版先不做
func CloseCentralizedPick(tx *gorm.DB) (err error) {
	return
}

// 关闭拣货池
func ClosePick(
	tx *gorm.DB,
	orderGoodsIds []int,
	closeGoodsMp map[int]int,
) (
	err error,
) {
	var (
		pickGoods       []model.PickGoods
		pickGoodsUpdate []model.PickGoods
		pickIds         []int
	)

	err, pickGoods = model.GetPickGoodsByOrderGoodsIds(tx, orderGoodsIds)

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

	//todo 如果预拣池全部被关闭，则更新预拣池状态
	return
}

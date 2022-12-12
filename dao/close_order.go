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
		outboundOrders []model.OutboundOrder
		closeOrders    []model.CloseOrder
		closeGoods     []model.CloseGoods
		taskNumberNew  []string //任务中的新订单
	)

	//校验是否所有批次全部暂停
	//查询是否有进行中的批次
	err, _ = model.GetBatchFirst(db, model.Batch{Status: model.BatchOngoingStatus})

	//err 不是未找到数据
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	//MQ 中 处理新订单关闭 [包括订单中的新订单和任务中的新订单][关闭订单任务加异常状态]
	err, outboundOrders = model.GetOutboundOrderByNumbers(db, form.Number)

	if err != nil {
		return
	}

	//关闭订单数据
	err, closeOrders = model.GetCloseOrderByNumbers(db, form.Number)

	if err != nil {
		return
	}

	closeOrderMp := make(map[string]int, 0)

	for _, order := range closeOrders {
		closeOrderMp[order.Number] = order.Typ
	}

	//关闭订单商品数据
	err, closeGoods = model.GetCloseGoodsListByNumbers(db, form.Number)

	if err != nil {
		return
	}

	closeGoodsMp := make(map[string][]model.CloseGoods, 0)

	for _, good := range closeGoods {
		closeGoodsMp[good.Number] = append(closeGoodsMp[good.Number], good)
	}

	//map[number]taskId
	numberTaskMp := make(map[string]int, 0)

	for _, order := range outboundOrders {
		//任务中的新订单，直接关闭
		if order.OrderType == model.OutboundOrderTypeNew {
			taskNumberNew = append(taskNumberNew, order.Number)
		}

		taskId, numberTaskMpOk := numberTaskMp[order.Number]

		//利用主键自增，找到订单号对应的最大批次
		//订单每次只能在一个任务中，任务结束后变为欠货单才能重新进入下一个任务中
		if numberTaskMpOk && order.TaskId > taskId {
			numberTaskMp[order.Number] = order.TaskId
		}
	}

	//构造新订单查询更新数据

	return
}

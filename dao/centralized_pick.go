package dao

import (
	"errors"
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/model"
	"pick_v2/utils/timeutil"
	"time"
)

// 生成集中拣货
func CreateCentralizedPick(db *gorm.DB, outboundGoodsJoinOrder []model.OutboundGoodsJoinOrder, batchId int) error {

	mpCentralized := make(map[string]model.CentralizedPick, 0)

	//按sku归集数据
	for _, order := range outboundGoodsJoinOrder {
		cp, mpCentralizedOk := mpCentralized[order.Sku]

		if !mpCentralizedOk {
			cp = model.CentralizedPick{
				BatchId:        batchId,
				Sku:            order.Sku,
				TaskId:         order.TaskId,
				GoodsName:      order.GoodsName,
				GoodsType:      order.GoodsType,
				GoodsSpe:       order.GoodsSpe,
				PickUser:       "",
				TakeOrdersTime: nil,
				GoodsRemark:    order.GoodsRemark,
				GoodsUnit:      order.GoodsUnit,
				Shelves:        order.Shelves,
			}
		}

		cp.NeedNum += order.LackCount

		mpCentralized[order.Sku] = cp
	}

	//集中拣货数据构造
	centralizedPicks := make([]model.CentralizedPick, 0, len(mpCentralized))

	for _, pick := range mpCentralized {
		centralizedPicks = append(centralizedPicks, pick)
	}

	//集中拣货数据保存
	err := model.CentralizedPickSave(db, &centralizedPicks)

	if err != nil {
		return err
	}

	return nil
}

// 集中拣货列表
func CentralizedPickList(db *gorm.DB, form req.CentralizedPickListForm) (err error, res rsp.CentralizedPickListRsp) {

	err, total, centralizedPickList := model.GetCentralizedPickPageList(db, form.BatchId, form.Page, form.Size, form.IsRemark, form.GoodsName, form.GoodsType)
	if err != nil {
		return err, res
	}

	res.Total = total

	list := make([]*rsp.CentralizedPickList, 0, len(centralizedPickList))

	for _, pick := range centralizedPickList {
		list = append(list, &rsp.CentralizedPickList{
			TaskName:    pick.GoodsName,
			GoodsName:   pick.GoodsName,
			GoodsSpe:    pick.GoodsSpe,
			NeedNum:     pick.NeedNum,
			CompleteNum: pick.CompleteNum,
			PickUser:    pick.PickUser,
			GoodsRemark: pick.GoodsRemark,
		})
	}

	res.List = list
	return
}

// 集中拣货详情
func CentralizedPickDetail(db *gorm.DB, form req.CentralizedPickDetailForm) (err error, list []rsp.CentralizedPickDetailList) {

	err, prePickGoodsList := model.GetPrePickGoodsAndRemark(db, form.BatchId, form.Sku)

	if err != nil {
		return err, nil
	}

	var (
		numbers []string
		pickMp  = make(map[string]rsp.CentralizedPickDetailGoodsInfo, 0)
	)

	for _, l := range prePickGoodsList {
		numbers = append(numbers, l.Number)

		pickMp[l.Number] = rsp.CentralizedPickDetailGoodsInfo{
			GoodsName:   l.GoodsName,
			NeedNum:     l.NeedNum,
			GoodsRemark: l.GoodsRemark,
			GoodsUnit:   l.Unit,
		}
	}

	err, orderList := model.GetOrderListByNumbers(db, numbers)

	if err != nil {
		return err, nil
	}

	list = make([]rsp.CentralizedPickDetailList, 0, len(orderList))

	for _, ol := range orderList {
		goodsInfo, pickOk := pickMp[ol.Number]

		if !pickOk {
			return errors.New("商品数据异常"), nil
		}

		list = append(list, rsp.CentralizedPickDetailList{
			ShopName:    ol.ShopName,
			ShopCode:    ol.ShopCode,
			Line:        ol.Line,
			GoodsName:   goodsInfo.GoodsName,
			NeedNum:     goodsInfo.NeedNum,
			GoodsRemark: goodsInfo.GoodsRemark,
			GoodsUnit:   goodsInfo.GoodsUnit,
			OrderRemark: ol.OrderRemark,
		})
	}

	return err, list
}

// 集中拣货剩余数量统计
func CentralizedPickRemainingQuantity(db *gorm.DB, userName string) (error, int64) {
	var (
		count    int64
		batchIds []int
	)

	//快递批次进行中或暂停的单数量
	err, batchList := model.GetBatchListByTyp(db, model.ExpressDeliveryBatchTyp)

	if err != nil {
		return err, 0
	}

	for _, b := range batchList {
		batchIds = append(batchIds, b.Id)
	}

	if len(batchIds) > 0 {
		err, count = model.CountCentralizedPickByBatchAndUser(db, batchIds, userName)
		if err != nil {
			return err, 0
		}
	}

	return nil, count
}

// 集中拣货接单
func ConcentratedPickReceivingOrders(db *gorm.DB, form req.ConcentratedPickReceivingOrdersForm, userName string) (err error, res rsp.ReceivingOrdersRsp) {

	var (
		pick []model.CentralizedPick
	)

	// 先查询是否有当前拣货员被分配的任务或已经接单且未完成拣货的数据,如果被分配多条，第一按批次优先级，第二按拣货池优先级 优先拣货
	result := db.Model(&model.CentralizedPick{}).Where("batch_id = ? and pick_user = ? and status = 0", form.BatchId, userName).Find(&pick)

	if result.Error != nil {
		return result.Error, res
	}

	now := time.Now()

	if result.RowsAffected > 0 { //有分配的拣货任务
		res = ConcentratedPick(db, pick)
	} else { //没有分配的拣货任务
		//查询未被接单的拣货池数据
		result = db.Model(&model.CentralizedPick{}).Where("batch_id = ? and pick_user = '' and status = 0", form.BatchId).Find(&pick)

		if result.Error != nil {
			return result.Error, res
		}

		//拣货池有未接单的数据
		if result.RowsAffected > 0 {
			res = ConcentratedPick(db, pick)
		} else {
			return errors.New("暂无拣货单"), res
		}
	}

	//更新拣货池 + version 防并发
	result = db.Model(&model.CentralizedPick{}).
		Where("id = ? and version = ?", res.Id, res.Version).
		Updates(map[string]interface{}{
			"pick_user":        userName,
			"take_orders_time": &now,
			"version":          gorm.Expr("version + ?", 1),
		})

	return nil, res

}

// 拣货员拣货单接单分配逻辑
func ConcentratedPick(db *gorm.DB, pick []model.CentralizedPick) (res rsp.ReceivingOrdersRsp) {

	res.Id = pick[0].Id
	res.BatchId = pick[0].BatchId
	res.Version = pick[0].Version
	res.TakeOrdersTime = pick[0].TakeOrdersTime
	res.Sku = pick[0].Sku

	return
}

// 快递拣货列表
func CentralizedAndSecondaryList(db *gorm.DB, userName string) (err error, list []rsp.CentralizedAndSecondary) {
	//1.获取用户已接单但未完成的批次id
	err, pickList := model.GetPickListByPickUserAndNotReviewCompleted(db, userName, model.ExpressDeliveryBatchTyp)
	if err != nil {
		return
	}

	var batchIds []int

	for _, pick := range pickList {
		batchIds = append(batchIds, pick.BatchId)
	}

	//2.获取进行中的批次 或 用户已接单但未完成的批次
	err, batchList := model.GetBatchListByIdsOrPending(db, batchIds, model.ExpressDeliveryBatchTyp)
	if err != nil {
		return
	}

	//重置batchIds
	batchIds = []int{}

	for _, batch := range batchList {
		batchIds = append(batchIds, batch.Id)
	}

	//统计集中拣货相关数量
	err, countCentralizedPickNums := model.CountCentralizedPickByBatch(db, batchIds)

	if err != nil {
		return
	}

	//统计二次分拣相关数量
	err, CountPickNums := model.CountPickByBatch(db, batchIds)
	if err != nil {
		return
	}

	var (
		countCentralizedPickMp = make(map[int]model.CountPickNums, 0)
		countPickMp            = make(map[int]model.CountPickNums, 0)
	)

	for _, num := range countCentralizedPickNums {
		countCentralizedPickMp[num.BatchId] = num
	}

	for _, num := range CountPickNums {
		countPickMp[num.BatchId] = num
	}

	list = make([]rsp.CentralizedAndSecondary, 0, len(batchList))

	for _, batch := range batchList {

		countCentralizedPick, countCentralizedPickOk := countCentralizedPickMp[batch.Id]

		if !countCentralizedPickOk {
			countCentralizedPick.SumNeedNum = 0
			countCentralizedPick.SumPickNum = 0
			countCentralizedPick.CountNeedNum = 0
			countCentralizedPick.CountPickNum = 0
		}

		countPick, countPickOk := countPickMp[batch.Id]

		if !countPickOk {
			countPick.SumNeedNum = 0
			countPick.SumCompleteNum = 0
			countPick.CountPickNum = 0
			countPick.CountCompleteNum = 0
		}

		centralized := rsp.Centralized{
			SumNeedNum:   countCentralizedPick.SumNeedNum,
			SumPickNum:   countCentralizedPick.SumPickNum,
			CountNeedNum: countCentralizedPick.CountNeedNum,
			CountPickNum: countCentralizedPick.CountPickNum,
		}

		secondary := rsp.Secondary{
			SumNeedNum:       countPick.SumNeedNum,
			SumCompleteNum:   countPick.SumCompleteNum,
			CountNeedNum:     countPick.CountPickNum,
			CountCompleteNum: countPick.CountCompleteNum,
		}

		list = append(list, rsp.CentralizedAndSecondary{
			BatchId:     batch.Id,
			BatchName:   batch.BatchName,
			CreateTime:  timeutil.FormatToDateTimeMinute(batch.CreateTime),
			Centralized: centralized,
			Secondary:   secondary,
		})
	}

	return
}

// 集中拣货详情-PDA 拣货使用
func CentralizedPickDetailPDA(db *gorm.DB, form req.CentralizedPickDetailPDAForm) (err error, res rsp.CentralizedPickDetailPDARsp) {
	err, first := model.GetCentralizedPickById(db, form.Id)

	if err != nil {
		return err, rsp.CentralizedPickDetailPDARsp{}
	}

	err, prePickGoodsList := model.GetPrePickGoodsAndRemark(db, first.BatchId, first.Sku)

	if err != nil {
		return err, res
	}

	var (
		numbers []string
	)

	for _, l := range prePickGoodsList {
		numbers = append(numbers, l.Number)
	}

	err, orderList := model.GetOrderListByNumbers(db, numbers)

	if err != nil {
		return err, res
	}

	numberList := make([]rsp.RemarkList, 0, len(orderList))

	for _, ol := range orderList {

		if ol.OrderRemark == "" {
			continue
		}

		numberList = append(numberList, rsp.RemarkList{
			Number: ol.Number,
			Remark: ol.OrderRemark,
		})
	}

	res = rsp.CentralizedPickDetailPDARsp{
		Id:          first.Id,
		Shelves:     first.Shelves,
		GoodsName:   first.GoodsName,
		GoodsSpe:    first.GoodsSpe,
		NeedNum:     first.NeedNum,
		GoodsUnit:   first.GoodsUnit,
		CompleteNum: first.CompleteNum,
		RemarkList:  numberList,
	}

	return err, res
}

func CompleteConcentratedPick(db *gorm.DB, form req.CompleteConcentratedPickForm) (err error) {
	//更新集中拣货状态为已完成，拣货数，类型
	err = model.UpdateCentralizedPickById(db, form.Id, map[string]interface{}{"status": model.CentralizedPickStatusCompleted, "complete_num": form.CompleteNum})
	return err
}

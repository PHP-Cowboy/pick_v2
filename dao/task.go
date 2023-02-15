package dao

import (
	"errors"
	"gorm.io/gorm"
	"math"
	"pick_v2/forms/req"
	"pick_v2/model"
	"pick_v2/utils/slice"
	"strconv"
	"strings"
)

// 验证是否有任务被接单
func CheckPickIsReceived(db *gorm.DB, ids []int) (err error, pickList []model.Pick) {
	err, pickList = model.GetPickListByIds(db, ids)

	if err != nil {
		return
	}

	//拣货池状态
	for _, p := range pickList {
		if p.PickUser != "" {
			return errors.New("已有拣货任务被接单，无法取消"), nil
		}
	}

	return
}

// 校验任务是否有订单已被接单拣货
func CheckPickNumberIsReceived(db *gorm.DB, ids []int) (err error) {

	var (
		numbers       []string
		exist         bool
		pickGoodsList []model.PickGoods
	)

	//拣货池根据所选任务查询全部商品
	err, pickGoodsList = model.GetPickGoodsByPickIds(db, ids)

	if err != nil {
		return
	}

	//全部商品的订单编号
	for _, pg := range pickGoodsList {
		numbers = append(numbers, pg.Number)
	}

	numbers = slice.UniqueSlice(numbers)

	//查询订单是否有已拣的
	err, exist = model.GetFirstPickGoodsByNumbers(db, numbers)

	if err != nil {
		return
	}

	if exist {
		return errors.New("该任务订单部分商品已拣，不可关闭")
	}

	return
}

// 拣货池任务返回预拣池
func PickReturnPrePick(db *gorm.DB, pickList []model.Pick) error {

	var (
		ids []int
	)

	for _, pl := range pickList {
		//获取预拣池id
		prePickIds := strings.Split(pl.PrePickIds, ",")

		for _, id := range prePickIds {
			prePickId, err := strconv.Atoi(id)
			if err != nil {
				return err
			}
			ids = append(ids, prePickId)
		}
	}

	tx := db.Begin()

	var (
		prePickIds []int
		pickIds    []int
	)

	//更新预拣池状态
	err := ChangePrePickStatus(tx, prePickIds)
	if err != nil {
		tx.Rollback()
		return err
	}

	//更新拣货池状态
	err = ChangePickStatus(tx, pickIds)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// 更新预拣池状态
func ChangePrePickStatus(db *gorm.DB, prePickIds []int) error {
	//更新预拣池商品状态
	err := model.UpdatePrePickGoodsByPrePickIds(db, prePickIds, map[string]interface{}{"status": model.PrePickGoodsStatusProcessing})
	if err != nil {
		return err
	}

	//更新预拣池任务状态
	err = model.UpdatePrePickByIds(db, prePickIds, map[string]interface{}{"status": model.PrePickStatusUnhandled})
	if err != nil {
		return err
	}

	//更新预拣池备注状态
	err = model.UpdatePrePickRemarkByPrePickIds(db, prePickIds, map[string]interface{}{"status": model.PrePickRemarkStatusUnhandled})
	if err != nil {
		return err
	}

	return nil
}

// 更新拣货池状态
func ChangePickStatus(db *gorm.DB, pickIds []int) error {
	//拣货池商品和备注表没有状态
	err := model.UpdatePickByIds(db, pickIds, map[string]interface{}{"status": model.ReturnPrePickStatus})
	if err != nil {
		return err
	}

	return nil
}

// 取消拣货
func CancelPick(db *gorm.DB, form req.CancelPickForm) error {

	var (
		err      error
		pickList []model.Pick
	)

	//验证是否有任务被接单
	err, pickList = CheckPickIsReceived(db, form.Ids)

	if err != nil {
		return err
	}

	//校验任务是否有订单已被接单拣货
	err = CheckPickNumberIsReceived(db, form.Ids)

	if err != nil {
		return err
	}

	//拣货池任务返回预拣池
	err = PickReturnPrePick(db, pickList)
	if err != nil {
		return err
	}

	return nil
}

// 修改复核数
func ChangeReviewNum(db *gorm.DB, form req.ChangeReviewNumForm) (err error) {

	var (
		mpForm        = make(map[string]int, len(form.SkuReview)) // sku 变更 数量 map
		skuSlice      []string                                    //查询
		batch         model.Batch
		pick          model.Pick
		pickGoods     []model.PickGoods
		outboundGoods []model.OutboundGoods
		orderGoods    []model.OrderGoods
	)

	err, batch = model.GetBatchByPk(db, form.BatchId)

	if err != nil {
		return
	}

	//只允许改进行中的批次数据
	if batch.Status != 0 {
		err = errors.New("只能修改进行中的批次数据")
		return
	}

	err, pick = model.GetPickByPk(db, form.PickId)

	if err != nil {
		return
	}

	//只允许修改复核完成状态的数据
	if pick.Status != 2 {
		err = errors.New("只能修改拣货复核完成的数据")
		return
	}

	for _, review := range form.SkuReview {
		mpForm[review.Sku] = *review.ReviewNum  //sku 修改的复核数量 map
		skuSlice = append(skuSlice, review.Sku) //被修改复核数量的sku切片
	}

	err, pickGoods = model.GetPickGoodsByPickIdAndSku(db, form.PickId, skuSlice)

	if err != nil {
		return
	}

	var (
		skuReviewTotalMp = make(map[string]int) //当前拣货商品表中sku复核总数
		skuNeedTotalMp   = make(map[string]int) //当前拣货商品表中sku需拣总数
	)

	//计算 pick_goods sku 复核总数 需拣总数
	for _, good := range pickGoods {

		//计算 当前拣货商品表中sku复核总数 开始
		reviewNum, reviewOk := skuReviewTotalMp[good.Sku]

		if !reviewOk {
			reviewNum = 0
		}

		reviewNum += good.ReviewNum

		skuReviewTotalMp[good.Sku] = reviewNum
		//计算 当前拣货商品表中sku复核总数 结束

		//计算 当前拣货商品表中sku需拣总数 开始
		needNum, needOk := skuNeedTotalMp[good.Sku]

		if !needOk {
			needNum = 0
		}

		needNum += good.NeedNum

		skuNeedTotalMp[good.Sku] = needNum
		//计算 当前拣货商品表中sku需拣总数 结束
	}

	//复核数 差值 = 原来sku复核总数 - form.Num
	var (
		reviewNumDiffMp    = make(map[string]int, len(form.SkuReview)) //每个sku的复核数差值map
		reviewNumDiffTotal = 0                                         //复核数 总差值
	)

	for _, sr := range form.SkuReview {
		num, sntOk := skuNeedTotalMp[sr.Sku]

		if !sntOk {
			err = errors.New("sku:" + sr.Sku + "需拣数未找到")
			return
		}

		if *sr.ReviewNum > num {
			err = errors.New("修改后的复核数大于需拣数，请核对")
			return
		}

		mpReviewNum, srtOK := skuReviewTotalMp[sr.Sku] //当前拣货商品表中sku复核总数

		if !srtOK {
			err = errors.New("sku:" + sr.Sku + "复核数未找到")
			return
		}

		reviewNumDiffMp[sr.Sku] = *sr.ReviewNum - mpReviewNum

		reviewNumDiffTotal += reviewNumDiffMp[sr.Sku]
	}

	var (
		//numberConsumeMp       = make(map[string]int, 0) //order相关消费复核差值map
		orderGoodsIdConsumeMp = make(map[int]int, 0) //orderGoods相关消费复核差值map
		outboundGoodsIdSlice  []int                  //outbound_goods 表 id
		orderGoodsIdSlice     []int                  //order_goods 表 id 处理后可能和 outboundGoodsIdSlice 是一致的
	)

	//pick_goods pay_at
	for i, good := range pickGoods {
		reviewNumDiff, reviewNumDiffOk := reviewNumDiffMp[good.Sku]

		if !reviewNumDiffOk {
			continue
		}

		if reviewNumDiff == 0 {
			continue
		}

		//numberConsumeNum, numberConsumeOk := numberConsumeMp[good.Number]
		//
		//if !numberConsumeOk {
		//	numberConsumeNum = 0
		//}

		orderGoodsIdConsumeNum, orderGoodsIdConsumeOk := orderGoodsIdConsumeMp[good.OrderGoodsId]

		if !orderGoodsIdConsumeOk {
			orderGoodsIdConsumeNum = 0
		}

		if reviewNumDiff > 0 { //修改后的复核数比原来大 已拣要加 未拣要减
			if good.NeedNum > good.ReviewNum { //需拣大于复核数，还有欠货，复核数可以加
				surplus := good.NeedNum - good.ReviewNum //当前单还剩的可以复核数量
				if reviewNumDiff > surplus {             //修改的复核数差值 比 当前单还剩的复核数量 大 复核数可以加满
					pickGoods[i].ReviewNum = good.NeedNum
					//pickGoods[i].NeedNum = 0 需拣数不会变，拣货和复核时都不修改需拣数
					reviewNumDiff -= surplus
					orderGoodsIdConsumeNum += surplus //这里复核数加了
					//numberConsumeNum += surplus
				} else { //修改的复核数差值 比 当前单还剩的复核数量 小 当前单可以消耗完差值
					pickGoods[i].ReviewNum += reviewNumDiff //复核数可以消耗完，直接把复核数往上加
					orderGoodsIdConsumeNum += reviewNumDiff //这里增加的复核数被这个单消耗完了
					//numberConsumeNum += reviewNumDiff
					//pickGoods[i].NeedNum -= reviewNumDiff
					reviewNumDiff = 0
				}
			} //没有 else 需拣小于等于复核数表明这单拣货完成，reviewNumDiff > 0 已拣加不上去了

		} else { //修改后的复核数比原来小 已拣要扣，未拣要增

			reviewDiffAbs := math.Abs(float64(reviewNumDiff))

			if good.ReviewNum > int(reviewDiffAbs) { //复核数大于 复核数差值的绝对值 差值直接被消耗完了
				pickGoods[i].ReviewNum = good.ReviewNum + reviewNumDiff //reviewNumDiff	小于0 这里的加的是负数
				orderGoodsIdConsumeNum += reviewNumDiff                 //这里扣减的复核数被这个单消耗完了 这里的加的是负数
				//numberConsumeNum += reviewNumDiff
				reviewNumDiff = 0
			} else { //当复核数本来就是0时，其实相当于什么都没做
				pickGoods[i].ReviewNum = 0               //复核数扣到0
				reviewNumDiff += good.ReviewNum          //差值 被消耗掉 good.ReviewNum 复核数个
				orderGoodsIdConsumeNum -= good.ReviewNum //这里把整单复核数扣光了
				//numberConsumeNum -= good.ReviewNum
			}
		}

		//变更map中sku的差值
		reviewNumDiffMp[good.Sku] = reviewNumDiff
		//变更orderGoods相关消费复核差值map
		orderGoodsIdConsumeMp[good.OrderGoodsId] = orderGoodsIdConsumeNum
		//变更 order相关消费复核差值map
		//numberConsumeMp[good.Number] = numberConsumeNum

		outboundGoodsIdSlice = append(outboundGoodsIdSlice, good.OrderGoodsId)
	}

	//修改 t_outbound_goods lack_count out_count
	err, outboundGoods = model.GetOutboundGoodsListByOrderGoodsIdAndTaskId(db, outboundGoodsIdSlice, pick.TaskId)
	if err != nil {
		return
	}

	//order_goods表 id 对应消费数
	var orderGoodsConsumeMp = make(map[int]int, 0)

	for i, good := range outboundGoods {
		consumeNum, consumeOk := orderGoodsIdConsumeMp[good.OrderGoodsId]

		if !consumeOk {
			continue
		}
		outboundGoods[i].LackCount -= consumeNum
		outboundGoods[i].OutCount += consumeNum

		orderGoodsConsumeMp[good.OrderGoodsId] = consumeNum

		orderGoodsIdSlice = append(orderGoodsIdSlice, good.OrderGoodsId)
	}

	err, orderGoods = model.GetOrderGoodsListByIds(db, orderGoodsIdSlice)

	if err != nil {
		return
	}

	for i, og := range orderGoods {
		consumeNum, consumeOk := orderGoodsConsumeMp[og.Id]

		if !consumeOk {
			continue
		}

		orderGoods[i].LackCount -= consumeNum
		orderGoods[i].OutCount += consumeNum
	}

	tx := db.Begin()

	//更新拣货商品表 复核数量 [need_num,review_num not null]
	err = model.PickGoodsReplaceSave(tx, &pickGoods, []string{"need_num", "complete_num", "review_num"})

	if err != nil {
		tx.Rollback()
		return
	}

	//更新出库单商品表数据
	err = model.OutboundGoodsReplaceSave(tx, &outboundGoods, []string{"lack_count", "out_count"})
	if err != nil {
		tx.Rollback()
		return
	}

	//更新订单商品表欠货数量，出库数量
	err = model.OrderGoodsReplaceSave(tx, &orderGoods, []string{"update_time", "lack_count", "out_count"})
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return
}

// 修改已拣数量
func ChangeCompleteNum(db *gorm.DB, form req.ChangeCompleteNumForm) (err error) {

	var (
		mpForm    = make(map[string]int, len(form.SkuComplete)) // sku 变更 数量 map
		skuSlice  []string                                      //查询
		batch     model.Batch
		pick      model.Pick
		pickGoods []model.PickGoods
	)

	err, batch = model.GetBatchByPk(db, form.BatchId)

	if err != nil {
		return
	}

	//只允许改进行中的批次数据
	if batch.Status != 0 {
		err = errors.New("只能修改进行中的批次数据")
		return
	}

	err, pick = model.GetPickByPk(db, form.PickId)

	if err != nil {
		return
	}

	//只允许修改待复核状态的数据
	if pick.Status != model.ToBeReviewedStatus {
		err = errors.New("只能修改待复核的数据")
		return
	}

	for _, review := range form.SkuComplete {
		mpForm[review.Sku] = *review.CompleteNum //sku 修改的完成数量 map
		skuSlice = append(skuSlice, review.Sku)  //被修改完成数量的sku切片
	}

	err, pickGoods = model.GetPickGoodsByPickIdAndSku(db, form.PickId, skuSlice)

	if err != nil {
		return
	}

	var (
		skuCompleteTotalMp = make(map[string]int) //当前拣货商品表中sku完成总数
		skuNeedTotalMp     = make(map[string]int) //当前拣货商品表中sku需拣总数
	)

	//计算 pick_goods sku 完成总数 需拣总数
	for _, good := range pickGoods {

		//计算 当前拣货商品表中sku完成总数 开始
		completeNum, completeOk := skuCompleteTotalMp[good.Sku]

		if !completeOk {
			completeNum = 0
		}

		completeNum += good.CompleteNum

		skuCompleteTotalMp[good.Sku] = completeNum
		//计算 当前拣货商品表中sku完成总数 结束

		//计算 当前拣货商品表中sku需拣总数 开始
		needNum, needOk := skuNeedTotalMp[good.Sku]

		if !needOk {
			needNum = 0
		}

		needNum += good.NeedNum

		skuNeedTotalMp[good.Sku] = needNum
		//计算 当前拣货商品表中sku需拣总数 结束
	}

	//完成数 差值 = 原来sku完成总数 - form.Num
	var (
		completeNumDiffMp    = make(map[string]int, len(form.SkuComplete)) //每个sku的完成数差值map
		completeNumDiffTotal = 0                                           //完成数 总差值
	)

	for _, sr := range form.SkuComplete {
		num, sntOk := skuNeedTotalMp[sr.Sku]

		if !sntOk {
			err = errors.New("sku:" + sr.Sku + "需拣数未找到")
			return
		}

		if *sr.CompleteNum > num {
			err = errors.New("修改后的完成数大于需拣数，请核对")
			return
		}

		mpCompleteNum, srtOK := skuCompleteTotalMp[sr.Sku] //当前拣货商品表中sku完成总数

		if !srtOK {
			err = errors.New("sku:" + sr.Sku + "完成数未找到")
			return
		}

		completeNumDiffMp[sr.Sku] = *sr.CompleteNum - mpCompleteNum

		completeNumDiffTotal += completeNumDiffMp[sr.Sku]
	}

	//pick_goods pay_at
	for i, good := range pickGoods {
		completeNumDiff, completeNumDiffOk := completeNumDiffMp[good.Sku]

		if !completeNumDiffOk || completeNumDiff == 0 {
			continue
		}

		if completeNumDiff > 0 { //修改后的完成数比原来大 已拣要加 未拣要减
			if good.NeedNum > good.CompleteNum { //需拣大于完成数，还有欠货，完成数可以加
				surplus := good.NeedNum - good.CompleteNum //当前单还剩的可以完成数量
				if completeNumDiff > surplus {             //修改的完成数差值 比 当前单还剩的完成数量 大 完成数可以加满
					pickGoods[i].CompleteNum = good.NeedNum
					completeNumDiff -= surplus
				} else { //修改的完成数差值 比 当前单还剩的完成数量 小 当前单可以消耗完差值
					pickGoods[i].CompleteNum += completeNumDiff //完成数可以消耗完，直接把完成数往上加
					completeNumDiff = 0
				}
			} //没有 else 需拣小于等于完成数表明这单拣货完成，completeNumDiff > 0 已拣加不上去了

		} else { //修改后的完成数比原来小 已拣要扣，未拣要增

			reviewDiffAbs := math.Abs(float64(completeNumDiff))

			if good.CompleteNum > int(reviewDiffAbs) { //完成数大于 完成数差值的绝对值 差值直接被消耗完了
				pickGoods[i].CompleteNum = good.CompleteNum + completeNumDiff //completeNumDiff	小于0 这里的加的是负数
				completeNumDiff = 0
			} else { //当完成数本来就是0时，其实相当于什么都没做
				pickGoods[i].CompleteNum = 0        //完成数扣到0
				completeNumDiff += good.CompleteNum //差值 被消耗掉 good.completeNum 完成数个
			}
		}

		//变更map中sku的差值
		completeNumDiffMp[good.Sku] = completeNumDiff

	}

	//更新拣货商品表 完成数量 [need_num,review_num not null]
	err = model.PickGoodsReplaceSave(db, &pickGoods, []string{"need_num", "complete_num", "review_num"})

	if err != nil {
		return
	}

	return
}

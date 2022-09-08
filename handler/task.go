package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pick_v2/common/constant"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"strings"
	"time"
)

// 置顶
func PickTopping(c *gin.Context) {
	var form req.PickToppingForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	val, err := cache.GetIncrByKey(constant.PICK_TOPPING)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	sort := int(val.(int64))

	result := global.DB.Model(model.Pick{}).Where("id = ?", form.Id).Update("sort", sort)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 拣货池列表
func PickList(c *gin.Context) {
	var form req.PickListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	var (
		ids        []int
		res        rsp.PickListRsp
		pick       []model.Pick
		pickGoods  []model.PickGoods
		pickRemark []model.PickRemark
		result     *gorm.DB
	)

	if form.Goods != "" || form.Number != "" || form.ShopId > 0 {
		result = db.Where(model.PickGoods{BatchId: form.BatchId, GoodsName: form.Goods, Number: form.Number, ShopId: form.ShopId}).Find(&pickGoods)
		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, pg := range pickGoods {
			ids = append(ids, pg.Id)
		}
	}

	localDb := db.Table("t_pick")

	if len(ids) > 0 {
		localDb.Where("id in (?)", ids)
	}

	localDb.Where("batch_id = ? and status = ?", form.BatchId, form.Status)

	result = localDb.Find(&pick)

	res.Total = result.RowsAffected

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

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

	list := make([]rsp.Pick, 0)
	for _, p := range pick {
		_, ok := remarkMp[p.Id]
		if ok { //拣货id在拣货备注中存在，即为有备注
			isRemark = true
		}

		takeOrdersTime := ""

		if p.TakeOrdersTime != nil {
			takeOrdersTime = p.TakeOrdersTime.Format(timeutil.TimeFormat)
		}

		reviewTime := ""

		if p.ReviewTime != nil {
			reviewTime = p.ReviewTime.Format(timeutil.TimeFormat)
		}

		list = append(list, rsp.Pick{
			Id:             p.Id,
			ShopCode:       p.ShopCode,
			ShopName:       p.ShopName,
			ShopNum:        p.ShopNum,
			OrderNum:       p.OrderNum,
			NeedNum:        p.NeedNum,
			PickUser:       p.PickUser,
			TakeOrdersTime: takeOrdersTime,
			IsRemark:       isRemark,
			Status:         p.Status,
			UpdateTime:     p.UpdateTime.Format(timeutil.TimeFormat),
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

// 拣货明细
func GetPickDetail(c *gin.Context) {
	var (
		form req.GetPickDetailForm
		res  rsp.GetPickDetailRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		pick       model.Pick
		pickGoods  []model.PickGoods
		pickRemark []model.PickRemark
	)

	db := global.DB

	result := db.Where("id = ?", form.PickId).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.TaskName = pick.TaskName
	res.ShopNum = pick.ShopNum
	res.OrderNum = pick.OrderNum
	res.GoodsNum = pick.NeedNum
	res.PickUser = pick.PickUser

	takeOrdersTime := ""

	if pick.TakeOrdersTime != nil {
		takeOrdersTime = pick.TakeOrdersTime.Format(timeutil.TimeFormat)
	}
	res.TakeOrdersTime = takeOrdersTime

	result = db.Where("pick_id = ?", form.PickId).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	pickGoodsSkuMp := make(map[string]rsp.MergePickGoods, 0)
	//相同sku合并处理
	for _, goods := range pickGoods {
		val, ok := pickGoodsSkuMp[goods.Sku]

		paramsId := rsp.ParamsId{
			PickGoodsId:  goods.Id,
			OrderGoodsId: goods.OrderGoodsId,
		}

		if !ok {

			pickGoodsSkuMp[goods.Sku] = rsp.MergePickGoods{
				Id:          goods.Id,
				Sku:         goods.Sku,
				GoodsName:   goods.GoodsName,
				GoodsType:   goods.GoodsType,
				GoodsSpe:    goods.GoodsSpe,
				Shelves:     goods.Shelves,
				NeedNum:     goods.NeedNum,
				CompleteNum: goods.CompleteNum,
				ReviewNum:   goods.ReviewNum,
				Unit:        goods.Unit,
				ParamsId:    []rsp.ParamsId{paramsId},
			}
		} else {
			val.NeedNum += goods.NeedNum
			val.CompleteNum += goods.CompleteNum
			val.ParamsId = append(val.ParamsId, paramsId)
			pickGoodsSkuMp[goods.Sku] = val
		}
	}

	goodsMap := make(map[string][]rsp.MergePickGoods, 0)

	for _, goods := range pickGoodsSkuMp {

		goodsMap[goods.GoodsType] = append(goodsMap[goods.GoodsType], rsp.MergePickGoods{
			Id:          goods.Id,
			Sku:         goods.Sku,
			GoodsName:   goods.GoodsName,
			GoodsType:   goods.GoodsType,
			GoodsSpe:    goods.GoodsSpe,
			Shelves:     goods.Shelves,
			NeedNum:     goods.NeedNum,
			CompleteNum: goods.CompleteNum,
			ReviewNum:   goods.ReviewNum,
			Unit:        goods.Unit,
			ParamsId:    goods.ParamsId,
		})
	}

	res.Goods = goodsMap

	result = db.Where("pick_id = ?", form.PickId).Find(&pickRemark)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := []rsp.PickRemark{}
	//这里订单号会重复，但是sku是不一致的，待确认
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

// 修改已出库件数
func ChangeNum(c *gin.Context) {
	var form req.ChangeNumReq

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var pick model.Pick

	db := global.DB

	result := db.First(&pick, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if pick.Status != 2 {
		xsq_net.ErrorJSON(c, errors.New("只允许修改已出库的"))
		return
	}

	result = db.Model(&model.Pick{}).Where("id = ?", form.Id).Update("num", form.Num)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 打印
func PushPrint(c *gin.Context) {
	var form req.PrintReq

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		pick      []model.Pick
		pickGoods []model.PickGoods
	)

	db := global.DB

	//复核完成
	result := db.Where("id in (?) and status = 2", form.Ids).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	pickMp := make(map[int]string, 0)

	for _, p := range pick {
		pickMp[p.Id] = strings.Join(p.DeliveryOrderNo, ",")
	}

	result = db.Where("pick_id in (?)", form.Ids).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	shopAndPickGoodsIdMp := make(map[string]struct{})

	printChSlice := make([]global.PrintCh, 0)
	for _, good := range pickGoods {
		deliveryOrderNo, pickOk := pickMp[good.PickId]

		if !pickOk {
			continue
		}

		mpKey := fmt.Sprintf("%d%d", good.PickId, good.ShopId)
		_, ok := shopAndPickGoodsIdMp[mpKey]
		if ok {
			continue
		}
		shopAndPickGoodsIdMp[mpKey] = struct{}{}
		printChSlice = append(printChSlice, global.PrintCh{
			DeliveryOrderNo: deliveryOrderNo,
			ShopId:          good.ShopId,
		})
	}

	for _, ch := range printChSlice {
		AddPrintJobMap(constant.JH_HUOSE_CODE, &global.PrintCh{
			DeliveryOrderNo: ch.DeliveryOrderNo,
			ShopId:          ch.ShopId,
		})
	}

	xsq_net.Success(c)
}

// 指派
func Assign(c *gin.Context) {
	var form req.AssignReq

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		user  model.User
		picks []model.Pick
	)

	db := global.DB

	result := db.First(&user, form.UserId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if user.RoleId != 3 { //不是拣货员
		xsq_net.ErrorJSON(c, ecode.UserNotFound)
		return
	}

	result = db.Where("id in (?)", form.PickIds).Find(&picks)
	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, p := range picks {
		if p.Status != 0 {
			xsq_net.ErrorJSON(c, errors.New("只能分配待拣货的任务"))
			return
		}
	}

	now := time.Now()

	//已有拣货员可以强转
	result = db.Model(&model.Pick{}).Where("id in (?)", form.PickIds).Updates(map[string]interface{}{"pick_user": user.Name, "take_orders_time": &now})

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

//修改复核数
func ChangeReviewNum(c *gin.Context) {
	var form req.ChangeReviewNumForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		batch          model.Batch
		pick           model.Pick
		pickGoods      []model.PickGoods
		pickOrderGoods []model.PickOrderGoods
		order          []model.Order
		orderGoods     []model.OrderGoods
	)

	db := global.DB

	result := db.First(&batch, form.BatchId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//只允许改进行中的批次数据
	if batch.Status != 0 {
		xsq_net.ErrorJSON(c, errors.New("只能修改进行中的批次数据"))
		return
	}

	result = db.First(&pick, form.PickId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//只允许修改复核完成状态的数据
	if pick.Status != 2 {
		xsq_net.ErrorJSON(c, errors.New("只能修改拣货复核完成的数据"))
		return
	}

	result = db.Model(&model.PickGoods{}).Where("pick_id = ? and sku = ?", form.PickId, form.Sku).Find(&pickGoods)

	var (
		reviewTotal     int
		needTotal       int
		numbers         []string
		numberPickNumMp = make(map[string]int) //订单本次复核数map
	)
	//计算 pick_goods sku 复核总数 需拣总数
	for _, good := range pickGoods {
		reviewTotal += good.ReviewNum
		needTotal += good.NeedNum

		numbers = append(numbers, good.Number) //不需要去重，一个拣货任务 订单号和sku本来就是唯一的

		numberPickNumMp[good.Number] = good.ReviewNum
	}

	num := form.Num

	if *num > needTotal {
		xsq_net.ErrorJSON(c, errors.New("修改后的复核数大于需拣数，请核对"))
		return
	}

	//复核数 差值 = 原来sku复核总数 - form.Num
	reviewNumDiff := reviewTotal - *num

	if reviewNumDiff == 0 { //没改，不折腾了
		xsq_net.Success(c)
		return
	}

	//拣货池 复核总数 = 复核总数 - 复核数 差值
	pick.ReviewNum = pick.ReviewNum - reviewNumDiff

	result = db.Model(&model.Order{}).Where("number in (?)", numbers).Order("pay_at asc").Find(&order)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for i, o := range order {
		if reviewNumDiff == 0 { //开始就小于0或大于0
			break
		}

		//订单本次复核数
		reviewNum, ok := numberPickNumMp[o.Number]

		if !ok {
			continue
		}

		//这里 需要 确认下 closeNum逻辑
		//当前单需拣数
		orderNeed := o.PayTotal - o.CloseNum

		if reviewNumDiff > 0 { //修改后的复核数比原来大 已拣要加 未拣要减
			//当前剩余需拣 是否有往上加的空间
			if orderNeed > 0 { //剩余需拣大于0 有加的空间
				//能加几个
				if orderNeed >= reviewNumDiff { //能加满
					order[i].Picked += reviewNumDiff
					order[i].UnPicked -= reviewNumDiff
					reviewNumDiff = 0
				} else {
					order[i].Picked += orderNeed //加到可以加的数量
					order[i].UnPicked -= orderNeed
					reviewNumDiff -= orderNeed
				}
			} //没有else ，没有空间加了，给下一条加

		} else { //修改后的复核数比原来小 已拣要扣，未拣要增
			reviewNumDiff = 0 - reviewNumDiff //不想用绝对值，直接减吧
			//确认出库前 已拣数量 未拣数量
			//历史已拣数量
			beforePicked := o.Picked - reviewNum
			//历史未拣数量
			beforeUnPicked := o.UnPicked + reviewNum
			//历史剩余需拣 = 需拣 - 历史已拣
			surplus := orderNeed - beforePicked
			//历史剩余需拣 大于等于 复核差值
			if surplus >= reviewNumDiff {
				order[i].Picked -= reviewNumDiff
				order[i].UnPicked += reviewNumDiff
				reviewNumDiff = 0
			} else {
				//历史剩余需拣 比 复核差值小
				order[i].Picked = beforePicked           //已拣直接扣回到历史已拣
				order[i].UnPicked = beforeUnPicked       //未拣还原到历史未拣
				reviewNumDiff -= o.Picked - beforePicked //差值 -= (复核数恢复前-即当前记录的拣货数 - 历史已拣)
			}
		}
	}

	result = db.Table("t_order_goods og").
		Select("og.*").
		Joins("left join t_order o on og.number = o.number").
		Where("number in (?) and sku = ?", numbers, form.Sku).
		Order("o.pay_at asc").
		Find(&orderGoods)

	for i, og := range orderGoods {
		if reviewNumDiff == 0 {
			break
		}

		//订单本次复核数
		reviewNum, ok := numberPickNumMp[og.Number]

		if !ok {
			continue
		}

		//这里 需要 确认下 closeNum逻辑
		//当前单需拣数
		orderNeed := og.PayCount - og.CloseCount

		if reviewNumDiff > 0 { //修改后的复核数比原来大 已拣要加 未拣要减
			//当前剩余需拣 是否有往上加的空间
			if orderNeed > 0 { //剩余需拣大于0 有加的空间
				//能加几个
				if orderNeed >= reviewNumDiff { //能加满
					orderGoods[i].OutCount += reviewNumDiff
					orderGoods[i].LackCount -= reviewNumDiff
					reviewNumDiff = 0
				} else {
					orderGoods[i].OutCount += orderNeed //加到可以加的数量
					orderGoods[i].LackCount -= orderNeed
					reviewNumDiff -= orderNeed
				}
			} //没有else ，没有空间加了，给下一条加

		} else { //修改后的复核数比原来小 已拣要扣，未拣要增
			reviewNumDiff = 0 - reviewNumDiff //不想用绝对值，直接减吧
			//确认出库前 已拣数量 未拣数量
			//历史已拣数量
			beforePicked := og.OutCount - reviewNum
			//历史未拣数量
			beforeUnPicked := og.LackCount + reviewNum
			//历史剩余需拣 = 需拣 - 历史已拣
			surplus := orderNeed - beforePicked
			//历史剩余需拣 大于等于 复核差值
			if surplus >= reviewNumDiff {
				orderGoods[i].OutCount -= reviewNumDiff
				orderGoods[i].LackCount += reviewNumDiff
				reviewNumDiff = 0
			} else {
				//历史剩余需拣 比 复核差值小
				orderGoods[i].OutCount = beforePicked       //已拣直接扣回到历史已拣
				orderGoods[i].LackCount = beforeUnPicked    //未拣还原到历史未拣
				reviewNumDiff -= og.OutCount - beforePicked //差值 -= (复核数恢复前-即当前记录的拣货数 - 历史已拣)
			}
		}

		//if reviewNumDiff > 0 {
		//	if reviewNumDiff <= og.OutCount {
		//		orderGoods[i].OutCount -= reviewNumDiff
		//		orderGoods[i].LackCount += reviewNumDiff
		//	}else {
		//		//修改后的复核数比原来大
		//
		//		//最多只能加到 下单数
		//
		//		//已拣要加 未拣要减
		//		orderGoods[i].OutCount += reviewNumDiff
		//		orderGoods[i].LackCount -= reviewNumDiff
		//	}
		//}

	}

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//pick_order_goods
	result = db.Table("t_pick_order_goods og").
		Select("og.*").
		Joins("left join t_pick_order o on og.number = o.number").
		Where("number in (?) and sku = ?", numbers, form.Sku).
		Order("o.pay_at asc").
		Find(&pickOrderGoods)

	//todo 订单可以多次生成拣货单 这个逻辑有问题
	for i, og := range pickOrderGoods {
		if reviewNumDiff == 0 {
			break
		}

		//订单本次复核数
		reviewNum, ok := numberPickNumMp[og.Number]

		if !ok {
			continue
		}

		//这里 需要 确认下 closeNum逻辑
		//当前单需拣数
		orderNeed := og.PayCount - og.CloseCount

		if reviewNumDiff > 0 { //修改后的复核数比原来大 已拣要加 未拣要减
			//当前剩余需拣 是否有往上加的空间
			if orderNeed > 0 { //剩余需拣大于0 有加的空间
				//能加几个
				if orderNeed >= reviewNumDiff { //能加满
					pickOrderGoods[i].OutCount += reviewNumDiff
					pickOrderGoods[i].LackCount -= reviewNumDiff
					reviewNumDiff = 0
				} else {
					pickOrderGoods[i].OutCount += orderNeed //加到可以加的数量
					pickOrderGoods[i].LackCount -= orderNeed
					reviewNumDiff -= orderNeed
				}
			} //没有else ，没有空间加了，给下一条加

		} else { //修改后的复核数比原来小 已拣要扣，未拣要增
			reviewNumDiff = 0 - reviewNumDiff //不想用绝对值，直接减吧
			//确认出库前 已拣数量 未拣数量
			//历史已拣数量
			beforePicked := og.OutCount - reviewNum
			//历史未拣数量
			beforeUnPicked := og.LackCount + reviewNum
			//历史剩余需拣 = 需拣 - 历史已拣
			surplus := orderNeed - beforePicked
			//历史剩余需拣 大于等于 复核差值
			if surplus >= reviewNumDiff {
				pickOrderGoods[i].OutCount -= reviewNumDiff
				pickOrderGoods[i].LackCount += reviewNumDiff
				reviewNumDiff = 0
			} else {
				//历史剩余需拣 比 复核差值小
				pickOrderGoods[i].OutCount = beforePicked    //已拣直接扣回到历史已拣
				pickOrderGoods[i].LackCount = beforeUnPicked //未拣还原到历史未拣
				reviewNumDiff -= og.OutCount - beforePicked  //差值 -= (复核数恢复前-即当前记录的拣货数 - 历史已拣)
			}
		}

	}

	tx := db.Begin()

	result = tx.
		Select("id", "update_time", "shop_code", "shop_name", "line", "review_num"). //这里要改下表结构
		Save(&pick)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.
		Select("id", "update_time", "picked", "un_picked").
		Save(&order)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.
		Select("id", "lack_count", "out_count").
		Save(&orderGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.
		Select("id", "lack_count", "out_count").
		Save(&pickOrderGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//pick_goods 咋改 以拣货单为突破口，向上 向下更新 订单 和 pick_goods
}

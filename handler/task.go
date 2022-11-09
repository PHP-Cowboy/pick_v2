package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"math"
	"pick_v2/common/constant"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/slice"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"sort"
	"strings"
	"time"
)

// 置顶
func PickTopping(c *gin.Context) {
	var form req.PickToppingForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	val, err := cache.Incr(constant.PICK_TOPPING)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	result := global.DB.Model(model.Pick{}).Where("id = ?", form.Id).Update("sort", val)

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

	result = localDb.Order("sort desc").Scopes(model.Paginate(form.Page, form.Size)).Find(&pick)

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

		list = append(list, rsp.Pick{
			Id:             p.Id,
			ShopCode:       p.ShopCode,
			ShopName:       p.ShopName,
			ShopNum:        p.ShopNum,
			OrderNum:       p.OrderNum,
			NeedNum:        p.NeedNum,
			PickUser:       p.PickUser,
			TakeOrdersTime: p.TakeOrdersTime,
			IsRemark:       isRemark,
			Status:         p.Status,
			UpdateTime:     p.UpdateTime.Format(timeutil.TimeFormat),
			PickNum:        p.PickNum,
			ReviewNum:      p.ReviewNum,
			Num:            p.Num,
			PrintNum:       p.PrintNum,
			ReviewUser:     p.ReviewUser,
			ReviewTime:     p.ReviewTime,
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

	res.BatchId = pick.BatchId
	res.PickId = pick.Id
	res.TaskName = pick.TaskName
	res.ShopCode = pick.ShopCode
	res.ShopNum = pick.ShopNum
	res.OrderNum = pick.OrderNum
	res.GoodsNum = pick.NeedNum
	res.PickUser = pick.PickUser
	res.TakeOrdersTime = pick.TakeOrdersTime

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
			val.ReviewNum += goods.ReviewNum
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

	//按货架号排序
	for s, goods := range goodsMap {

		ret := rsp.MyMergePickGoods(goods)

		sort.Sort(ret)

		goodsMap[s] = ret
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

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
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

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
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

	for i, p := range pick {
		pickMp[p.Id] = strings.Join(p.DeliveryOrderNo, ",")
		pick[i].PrintNum += 1
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
			Type:            form.Type, //打印类型
		})
	}

	result = db.Select("id", "update_time", "shop_code", "shop_name", "line", "print_num").Save(&pick)

	xsq_net.Success(c)
}

// 指派
func Assign(c *gin.Context) {
	var form req.AssignReq

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
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

// 修改复核数
func ChangeReviewNum(c *gin.Context) {
	var form req.ChangeReviewNumForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		mpForm         = make(map[string]int, len(form.SkuReview)) // sku 变更 数量 map
		skuSlice       []string                                    //查询
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

	for _, review := range form.SkuReview {
		mpForm[review.Sku] = *review.Num        //sku 修改的复核数量 map
		skuSlice = append(skuSlice, review.Sku) //被修改复核数量的sku切片
	}

	result = db.Model(&model.PickGoods{}).Where("pick_id = ? and sku in (?)", form.PickId, skuSlice).Find(&pickGoods)

	var (
		skuReviewTotalMp = make(map[string]int) //当前拣货商品表中sku复核总数
		skuNeedTotalMp   = make(map[string]int) //当前拣货商品表中sku需拣总数
		numbers          []string               //订单编号切片
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
		//计算 当前拣货商品表中sku需拣总数 开始

		numbers = append(numbers, good.Number) //订单编号切片
	}

	//订单编号去重
	numbers = slice.UniqueStringSlice(numbers)

	//复核数 差值 = 原来sku复核总数 - form.Num
	var (
		reviewNumDiffMp    = make(map[string]int, len(form.SkuReview)) //每个sku的复核数差值map
		reviewNumDiffTotal = 0                                         //复核数 总差值
	)

	for _, sr := range form.SkuReview {
		num, sntOk := skuNeedTotalMp[sr.Sku]

		if !sntOk {
			xsq_net.ErrorJSON(c, errors.New("sku:"+sr.Sku+"需拣数未找到"))
			return
		}

		if *sr.Num > num {
			xsq_net.ErrorJSON(c, errors.New("修改后的复核数大于需拣数，请核对"))
			return
		}

		mpReviewNum, srtOK := skuReviewTotalMp[sr.Sku] //当前拣货商品表中sku复核总数

		if !srtOK {
			xsq_net.ErrorJSON(c, errors.New("sku:"+sr.Sku+"复核数未找到"))
			return
		}

		reviewNumDiffMp[sr.Sku] = *sr.Num - mpReviewNum

		reviewNumDiffTotal += reviewNumDiffMp[sr.Sku]
	}

	//拣货池 复核总数 = 复核总数 - 复核数 差值
	pick.ReviewNum = pick.ReviewNum + reviewNumDiffTotal

	var (
		numberConsumeMp       = make(map[string]int, 0) //order相关消费复核差值map
		orderGoodsIdConsumeMp = make(map[int]int, 0)    //orderGoods相关消费复核差值map
		pickOrderGoodsIdSlice []int                     //pick_order_goods 表 id map
		orderGoodsIdSlice     []int
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

		numberConsumeNum, numberConsumeOk := numberConsumeMp[good.Number]

		if !numberConsumeOk {
			numberConsumeNum = 0
		}

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
					numberConsumeNum += surplus
				} else { //修改的复核数差值 比 当前单还剩的复核数量 小 当前单可以消耗完差值
					pickGoods[i].ReviewNum += reviewNumDiff //复核数可以消耗完，直接把复核数往上加
					orderGoodsIdConsumeNum += reviewNumDiff //这里增加的复核数被这个单消耗完了
					numberConsumeNum += reviewNumDiff
					//pickGoods[i].NeedNum -= reviewNumDiff
					reviewNumDiff = 0
				}
			} //没有 else 需拣小于等于复核数表明这单拣货完成，reviewNumDiff > 0 已拣加不上去了

		} else { //修改后的复核数比原来小 已拣要扣，未拣要增

			reviewDiffAbs := math.Abs(float64(reviewNumDiff))

			if good.ReviewNum > int(reviewDiffAbs) { //复核数大于 复核数差值的绝对值 差值直接被消耗完了
				pickGoods[i].ReviewNum = good.ReviewNum + reviewNumDiff //reviewNumDiff	小于0 这里的加的是负数
				orderGoodsIdConsumeNum += reviewNumDiff                 //这里扣减的复核数被这个单消耗完了 这里的加的是负数
				numberConsumeNum += reviewNumDiff
				reviewNumDiff = 0
			} else { //当复核数本来就是0时，其实相当于什么都没做
				pickGoods[i].ReviewNum = 0               //复核数扣到0
				reviewNumDiff += good.ReviewNum          //差值 被消耗掉 good.ReviewNum 复核数个
				orderGoodsIdConsumeNum -= good.ReviewNum //这里把整单复核数扣光了
				numberConsumeNum -= good.ReviewNum
			}
		}

		//变更map中sku的差值
		reviewNumDiffMp[good.Sku] = reviewNumDiff
		//变更orderGoods相关消费复核差值map
		orderGoodsIdConsumeMp[good.OrderGoodsId] = orderGoodsIdConsumeNum
		//变更 order相关消费复核差值map
		numberConsumeMp[good.Number] = numberConsumeNum

		pickOrderGoodsIdSlice = append(pickOrderGoodsIdSlice, good.OrderGoodsId)
	}

	//pick_order_goods
	result = db.Model(&model.PickOrderGoods{}).Where("id in (?)", pickOrderGoodsIdSlice).Find(&pickOrderGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//order_goods表 id 对应消费数
	var orderGoodsConsumeMp = make(map[int]int, 0)

	for i, good := range pickOrderGoods {
		consumeNum, consumeOk := orderGoodsIdConsumeMp[good.Id]

		if !consumeOk {
			continue
		}
		pickOrderGoods[i].LackCount -= consumeNum
		pickOrderGoods[i].OutCount += consumeNum

		orderGoodsConsumeMp[good.OrderGoodsId] = consumeNum

		orderGoodsIdSlice = append(orderGoodsIdSlice, good.OrderGoodsId)
	}

	result = db.Model(&model.OrderGoods{}).Where("id in (?)", orderGoodsIdSlice).Find(&orderGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
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

	result = db.Model(&model.Order{}).Where("number in (?)", numbers).Find(&order)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for i, o := range order {
		consumeNum, consumeOk := numberConsumeMp[o.Number]

		if !consumeOk {
			continue
		}

		order[i].Picked += consumeNum
		order[i].UnPicked -= consumeNum
	}

	tx := db.Begin()

	result = tx.Select("id", "update_time", "shop_code", "shop_name", "line", "review_num").
		Save(&pick)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Select("id", "update_time", "need_num", "complete_num", "review_num").
		Save(&pickGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Select("id", "update_time", "lack_count", "out_count").
		Save(&pickOrderGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Select("id", "update_time", "lack_count", "out_count").
		Save(&orderGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Select("id", "update_time", "shop_id", "shop_name", "shop_type", "shop_code", "house_code", "line", "picked", "un_picked").
		Save(&order)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

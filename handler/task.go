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
	"pick_v2/model/batch"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
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

	result := global.DB.Model(batch.Pick{}).Where("id = ?", form.Id).Update("sort", sort)

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
		pick       []batch.Pick
		pickGoods  []batch.PickGoods
		pickRemark []batch.PickRemark
		result     *gorm.DB
	)

	if form.Goods != "" || form.Number != "" || form.ShopId > 0 {
		result = db.Where(batch.PickGoods{BatchId: form.BatchId, GoodsName: form.Goods, Number: form.Number, ShopId: form.ShopId}).Find(&pickGoods)
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
		pick       batch.Pick
		pickGoods  []batch.PickGoods
		pickRemark []batch.PickRemark
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

	var pick batch.Pick

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

	result = db.Model(&batch.Pick{}).Where("id = ?", form.Id).Update("num", form.Num)

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
		pick      []batch.Pick
		pickGoods []batch.PickGoods
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
		pickMp[p.Id] = p.DeliveryOrderNo
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
		picks []batch.Pick
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
	result = db.Model(&batch.Pick{}).Where("id in (?)", form.PickIds).Updates(map[string]interface{}{"pick_user": user.Name, "take_orders_time": &now})

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

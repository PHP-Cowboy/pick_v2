package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"pick_v2/dao"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
	"sort"
)

// 复核列表 通过状态区分是否已完成
func ReviewList(c *gin.Context) {
	var (
		form       req.ReviewListReq
		res        rsp.ReviewListRsp
		pick       []model.Pick
		pickRemark []model.PickRemark
		//pickListModel []rsp.PickListModel
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	//result := db.Model(&batch.Pick{}).Where("status = ?", form.Status).Where(batch.Pick{PickUser: form.Name}).Find(&pick)

	localDb := db.Model(&model.Pick{}).Where("status = ?", form.Status)

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
		localDb.Where("review_user = ? ", userInfo.Name).Order("review_time desc")
	}

	result := localDb.Where(model.Pick{PickUser: form.Name}).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	list := make([]rsp.Pick, 0, form.Size)

	if len(pick) > 0 {
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

		for _, p := range pick {

			_, remarkMpOk := remarkMp[p.Id]
			if remarkMpOk { //拣货id在拣货备注中存在，即为有备注
				isRemark = true
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
				TakeOrdersTime: p.TakeOrdersTime,
				IsRemark:       isRemark,
				PickNum:        p.PickNum,
				ReviewNum:      p.ReviewNum,
				Num:            p.Num,
				ReviewUser:     p.ReviewUser,
				ReviewTime:     p.ReviewTime,
			})
		}
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// 复核明细
func ReviewDetail(c *gin.Context) {
	var (
		form req.ReviewDetailReq
		res  rsp.ReviewDetailRsp
		pick model.Pick
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

		result = db.Model(&model.Pick{}).
			Where("id = ? and version = ?", pick.Id, pick.Version).
			Updates(map[string]interface{}{"review_user": userInfo.Name, "version": pick.Version + 1})

		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	}

	var (
		pickGoods  []model.PickGoods
		pickRemark []model.PickRemark
	)

	res.TaskName = pick.TaskName
	res.ShopCode = pick.ShopCode
	res.PickUser = pick.PickUser
	res.TakeOrdersTime = pick.TakeOrdersTime
	res.ReviewUser = pick.ReviewUser
	res.ReviewTime = pick.ReviewTime

	result = db.Where("pick_id = ?", form.Id).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	pickGoodsSkuMp := make(map[string]rsp.MergePickGoods, 0)

	goodsMap := make(map[string][]rsp.MergePickGoods, 0)

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

	needTotal := 0
	completeTotal := 0
	reviewTotal := 0

	for _, goods := range pickGoodsSkuMp {
		completeTotal += goods.CompleteNum
		needTotal += goods.NeedNum
		reviewTotal += goods.ReviewNum

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

	res.ShopCode = pick.ShopCode
	res.OutTotal = completeTotal
	res.UnselectedTotal = needTotal - completeTotal
	res.ReviewTotal = reviewTotal

	//按货架号排序
	for s, goods := range goodsMap {

		ret := rsp.MyMergePickGoods(goods)

		sort.Sort(ret)

		goodsMap[s] = ret
	}

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

	form.UserName = userInfo.Name

	err := dao.ConfirmDelivery(global.DB, form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

package handler

import (
	"github.com/gin-gonic/gin"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/model/batch"
	"pick_v2/model/order"
	"pick_v2/utils/ecode"
	"pick_v2/utils/helper"
	"pick_v2/utils/xsq_net"
	"strconv"
)

//生成拣货批次
func CreateBatch(c *gin.Context) {
	var form req.CreateBatchForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	tx := global.DB.Begin()

	condition := batch.BatchCondition{
		WarehouseId:       form.WarehouseId,
		PayEndTime:        form.PayEndTime,
		DeliveryStartTime: form.DeliveryStartTime,
		DeliveryEndTime:   form.DeliveryEndTime,
		Line:              form.Lines,
		DeliveryMethod:    form.DeType,
		Sku:               form.Sku,
	}

	//筛选条件保存
	result := tx.Save(&condition)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	claims, ok := c.Get("claims")

	if !ok {
		xsq_net.ErrorJSON(c, ecode.DataNotExist)
		return
	}

	userInfo := claims.(*middlewares.CustomClaims)

	//批次数据
	batches := batch.Batch{
		WarehouseId:     form.WarehouseId,
		BatchName:       form.Lines + helper.GetDeliveryMethod(form.DeType),
		DeliveryEndTime: form.DeliveryEndTime,
		ShopNum:         0,
		OrderNum:        0,
		UserName:        userInfo.Name,
		Line:            form.Lines,
		DeliveryMethod:  form.DeType,
		EndTime:         form.PayEndTime,
		Status:          0,
		PickNum:         0,
		RecheckSheetNum: 0,
		Sort:            0,
	}

	result = tx.Save(&batches)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	form.BatchNumber = strconv.Itoa(batches.Id)

	goodsRes, err := RequestGoodsList(form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	if len(goodsRes.Data.List) <= 0 {
		xsq_net.ErrorJSON(c, ecode.NoOrderFound)
		return
	}

	var (
		orders        []order.OrderInfo
		prePicks      []batch.PrePick
		prePickGoods  []*batch.PrePickGoods
		prePickRemark []*batch.PrePickRemark
		shopMap       = make(map[int]int, 0)
	)

	//订单相关数据
	for _, goods := range goodsRes.Data.List {
		orders = append(orders, order.OrderInfo{
			ShopId:           goods.ShopId,
			ShopName:         goods.ShopName,
			ShopType:         goods.ShopType,
			ShopCode:         goods.ShopCode,
			HouseCode:        goods.HouseCode,
			Line:             goods.Line,
			Number:           goods.Number,
			Status:           goods.Status,
			DeliveryAt:       goods.DeliveryAt,
			DistributionType: goods.DistributionType,
			OrderRemark:      goods.OrderRemark,
			Province:         goods.Province,
			City:             goods.City,
			District:         goods.District,
			Address:          goods.Address,
			ConsigneeName:    goods.ConsigneeName,
			ConsigneeTel:     goods.ConsigneeTel,
			Name:             goods.Name,
			Sku:              goods.Sku,
			GoodsSpe:         goods.GoodsSpe,
			GoodsType:        goods.GoodsType,
			Shelves:          goods.Shelves,
			OriginalPrice:    goods.OriginalPrice,
			DiscountPrice:    int(goods.DiscountPrice * 100),
			GoodsUnit:        goods.GoodsUnit,
			SaleUnit:         goods.SaleUnit,
			SaleCode:         goods.SaleCode,
			PayCount:         goods.PayCount,
			CloseCount:       goods.CloseCount,
			OutCount:         goods.OutCount,
			GoodsRemark:      goods.GoodsRemark,
			PickStatus:       goods.PickStatus,
			PayAt:            goods.PayAt,
			LackCount:        goods.LackCount,
		})

		prePickGoods = append(prePickGoods, &batch.PrePickGoods{
			WarehouseId: form.WarehouseId,
			BatchId:     batches.Id,
			PrePickId:   0,
			ShopId:      goods.ShopId,
			GoodsName:   goods.Name,
			GoodsType:   goods.GoodsType,
			GoodsSpe:    goods.GoodsSpe,
			Shelves:     goods.Shelves,
			NeedNum:     0,
			CloseNum:    goods.CloseCount,
			OutCount:    goods.OutCount,
			NeedOutNum:  0,
		})

		if goods.GoodsRemark != "" || goods.OrderRemark != "" {
			prePickRemark = append(prePickRemark, &batch.PrePickRemark{
				WarehouseId: form.WarehouseId,
				BatchId:     batches.Id,
				ShopId:      goods.ShopId,
				Number:      goods.Number,
				OrderRemark: goods.OrderRemark,
				GoodsRemark: goods.GoodsRemark,
				ShopName:    goods.ShopName,
				Line:        goods.Line,
				PrePickId:   0,
			})
		}

		_, ok := shopMap[goods.ShopId]
		if ok {
			continue
		}
		shopMap[goods.ShopId] = 0
		prePicks = append(prePicks, batch.PrePick{
			WarehouseId: form.WarehouseId,
			BatchId:     batches.Id,
			ShopId:      goods.ShopId,
			ShopCode:    goods.ShopCode,
			ShopName:    goods.ShopName,
			Line:        goods.Line,
			Status:      0,
		})

	}

	result = tx.Save(&orders)
	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = tx.Save(&prePicks)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, pick := range prePicks {
		shopMap[pick.ShopId] = pick.Id
	}

	for k, good := range prePickGoods {
		val, ok := shopMap[good.ShopId]
		if !ok {
			xsq_net.ErrorJSON(c, ecode.MapKeyNotExist)
			return
		}
		prePickGoods[k].PrePickId = val
	}

	result = tx.Save(&prePickGoods)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if len(prePickRemark) > 0 {
		for k, remark := range prePickRemark {
			val, ok := shopMap[remark.ShopId]
			if !ok {
				xsq_net.ErrorJSON(c, ecode.MapKeyNotExist)
				return
			}
			prePickRemark[k].PrePickId = val
		}

		result = tx.Save(&prePickRemark)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	}

	tx.Commit()

	xsq_net.Success(c)
}

//获取批次列表
func GetBatchList(c *gin.Context) {
	var (
		form req.GetBatchListForm
		res  rsp.GetBatchListRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var batches []batch.Batch

	db := global.DB

	result := db.Where(map[string]interface{}{"status": form.Status}).Find(&batches)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	db.Where(map[string]interface{}{"status": form.Status}).Scopes(model.Paginate(form.Page, form.Size)).Find(&batches)

	list := make([]*rsp.Batch, 0, len(batches))
	for _, b := range batches {
		list = append(list, &rsp.Batch{
			BatchName:         b.BatchName,
			DeliveryStartTime: b.DeliveryStartTime,
			DeliveryEndTime:   b.DeliveryEndTime,
			ShopNum:           b.ShopNum,
			OrderNum:          b.OrderNum,
			UserName:          b.UserName,
			Line:              b.Line,
			DeliveryMethod:    b.DeliveryMethod,
			EndTime:           b.EndTime,
			Status:            b.Status,
			PickNum:           b.PickNum,
			RecheckSheetNum:   b.RecheckSheetNum,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

type Ret struct {
	OutC      int
	NeedC     int
	ShopId    int
	GoodsType string
}

//预拣池列表
func GetPrePickList(c *gin.Context) {
	var (
		form req.GetPrePickListForm
		res  rsp.GetPrePickListRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		prePicks []batch.PrePick
		//prePickGoods []batch.PrePickGoods
		prePickIds []int
	)

	db := global.DB

	result := db.Find(&prePicks)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	db.Scopes(model.Paginate(form.Page, form.Size)).Find(&prePicks)

	for _, pick := range prePicks {
		prePickIds = append(prePickIds, pick.Id)
	}

	retCount := []Ret{}

	result = db.Model(&batch.PrePickGoods{}).
		Select("SUM(out_count) as outC, SUM(need_num) AS needC, shop_id, goods_type").
		Where("pre_pick_id in (?)", prePickIds).
		Group("shop_id, goods_type").
		Find(&retCount)

	typeMap := make(map[int]map[string]rsp.PickCount, 0)

	for _, r := range retCount {
		_, ok := typeMap[r.ShopId]
		if !ok {
			countMap := make(map[string]rsp.PickCount, 0)
			typeMap[r.ShopId] = countMap
			countMap[r.GoodsType] = rsp.PickCount{
				WaitingPick: r.NeedC,
				PickedCount: r.OutC,
			}
			typeMap[r.ShopId][r.GoodsType] = countMap[r.GoodsType]
		}
	}

	for _, pick := range prePicks {
		res.List = append(res.List, &rsp.PrePick{
			ShopCode:     pick.ShopCode,
			ShopName:     pick.ShopName,
			Line:         pick.Line,
			CategoryInfo: typeMap[pick.ShopId],
		})
	}

	xsq_net.SucJson(c, res)

}

//预拣货明细
func GetPrePickDetail(c *gin.Context) {
	var (
		form req.GetPrePickDetailForm
		res  rsp.GetPrePickDetailRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		prePickGoods  []batch.PrePickGoods
		prePickRemark []batch.PrePickRemark
	)

	db := global.DB

	result := db.Where("pre_pick_id = ?", form.PrePickId).Find(&prePickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	goodsMap := make(map[string][]rsp.PrePickGoods, 0)

	for _, goods := range prePickGoods {
		goodsMap[goods.GoodsType] = append(goodsMap[goods.GoodsType], rsp.PrePickGoods{
			GoodsName:  goods.GoodsName,
			GoodsSpe:   goods.GoodsSpe,
			Shelves:    goods.Shelves,
			NeedNum:    goods.NeedNum,
			CloseNum:   goods.CloseNum,
			OutCount:   goods.OutCount,
			NeedOutNum: goods.NeedOutNum,
		})
	}

	res.Goods = goodsMap

	result = db.Where("pre_pick_id = ?", form.PrePickId).Find(&prePickRemark)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := []rsp.Remark{}
	for _, remark := range prePickRemark {
		list = append(list, rsp.Remark{
			Number:      remark.Number,
			OrderRemark: remark.OrderRemark,
			GoodsRemark: remark.GoodsRemark,
		})
	}

	res.RemarkList = list

	xsq_net.SucJson(c, res)
}

//置顶
func Topping(c *gin.Context) {
	//redis
}

//暂停

//结束

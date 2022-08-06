package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pick_v2/common/constant"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/model/batch"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
)

//停止拣货
func StopPick() {

}

//终止拣货
func TerminatePick() {

}

//置顶
func PickTopping(c *gin.Context) {
	var form req.PickToppingForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	//redis
	redis := global.Redis

	redisKey := constant.PICK_TOPPING

	val, err := redis.Do(context.Background(), "incr", redisKey).Result()
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

//拣货池列表
func PickList(c *gin.Context) {
	var form req.PickListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	var (
		ids           []int
		res           rsp.PickListRsp
		pick          []batch.Pick
		pickGoods     []batch.PickGoods
		pickListModel []rsp.PickListModel
		result        *gorm.DB
	)

	if form.Goods != "" || form.Number != "" || form.ShopId > 0 {
		result = db.Where(batch.PickGoods{GoodsName: form.Goods, Number: form.Number, ShopId: form.ShopId}).Find(&pickGoods)
		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, pg := range pickGoods {
			ids = append(ids, pg.Id)
		}
	}

	localDb := db.Table("t_pick p").Where("status = 0")

	if len(ids) > 0 {
		localDb.Where("p.id in (?)", ids)
	}

	result = localDb.Find(&pick)

	res.Total = result.RowsAffected

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = localDb.
		Select("p.id,shop_code,p.shop_name,shop_num,order_num,need_num,pick_user,take_orders_time,order_remark,goods_remark").
		Where("p.status = 0").
		Joins("left join t_pick_remark pr on pr.pick_id = p.id").
		Scopes(model.Paginate(form.Page, form.Size)).
		Scan(&pickListModel)

	isRemark := false

	list := make([]rsp.Pick, 0)
	for _, p := range pickListModel {
		if p.GoodsRemark != "" || p.OrderRemark != "" {
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
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

//拣货明细
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
	res.GoodsNum = 0
	res.PickUser = pick.PickUser
	res.TakeOrdersTime = pick.TakeOrdersTime

	result = db.Where("pick_id = ?", form.PickId).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	goodsMap := make(map[string][]rsp.PickGoods, 0)

	for _, goods := range pickGoods {
		goodsMap[goods.GoodsType] = append(goodsMap[goods.GoodsType], rsp.PickGoods{
			GoodsName: goods.GoodsName,
			GoodsSpe:  goods.GoodsSpe,
			Shelves:   goods.Shelves,
			NeedNum:   goods.NeedNum,
		})
	}

	res.Goods = goodsMap

	result = db.Where("pick_id = ?", form.PickId).Find(&pickRemark)

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

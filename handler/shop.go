package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm/clause"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/request"
	"pick_v2/utils/xsq_net"
	"time"
)

type ShopResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    []model.Shop `json:"data"`
}

// 同步门店
func SyncShop(c *gin.Context) {

	path := "api/v1/remote/shop/list"

	body, err := request.Post(path, nil)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var result ShopResponse

	err = json.Unmarshal(body, &result)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	if result.Code != 200 {
		xsq_net.ErrorJSON(c, errors.New(result.Message))
		return
	}

	var shops []*model.Shop

	db := global.DB

	for _, shop := range result.Data {

		shops = append(shops, &model.Shop{
			Id:        shop.Id,
			ShopId:    shop.ShopId,
			ShopName:  shop.ShopName,
			HouseCode: shop.HouseCode,
			Warehouse: shop.Warehouse,
			Typ:       shop.Typ,
			Province:  shop.Province,
			City:      shop.City,
			District:  shop.District,
			Line:      shop.Line,
			ShopCode:  shop.ShopCode,
			Status:    shop.Status,
			CreateAt:  time.Now(),
			UpdateAt:  time.Now(),
		})
	}

	res := db.Model(&model.Shop{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"shop_id", "shop_name", "house_code", "warehouse", "typ", "province", "city", "district", "shop_code", "status"}),
		}).
		CreateInBatches(&shops, model.BatchSize)

	if res.Error != nil {
		xsq_net.ErrorJSON(c, res.Error)
	}

	if err = cache.SetShopLine(); err != nil {
		xsq_net.ErrorJSON(c, errors.New("店铺更新成功，但缓存更新失败"))
		return
	}

	xsq_net.Success(c)
}

// 门店列表
func ShopList(c *gin.Context) {
	var form req.ShopListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		res   rsp.ShopListRsp
		shops []model.Shop
	)

	db := global.DB.Model(model.Shop{})

	if form.ShopName != "" {
		db.Where("shop_name like ?", "%"+form.ShopName+"%")
	}

	if form.LineStatus == 1 { //已设置路线
		db.Where("line != ''")
	} else if form.LineStatus == 2 { //未设置路线
		db.Where("line = ''")
	}

	if form.Line != "" {
		db.Where("line = ?", form.Line)
	}

	ret := db.Find(&shops)

	if ret.Error != nil {
		xsq_net.ErrorJSON(c, ret.Error)
		return
	}

	res.Total = ret.RowsAffected

	if form.IsPage {
		ret = db.Scopes(model.Paginate(form.Page, form.Size)).Find(&shops)

		if ret.Error != nil {
			xsq_net.ErrorJSON(c, ret.Error)
			return
		}
	}

	list := make([]*rsp.Shop, 0, form.Size)

	for _, shop := range shops {
		list = append(list, &rsp.Shop{
			Id:       shop.Id,
			ShopId:   shop.ShopId,
			ShopName: shop.ShopName,
			ShopCode: shop.ShopCode,
			Line:     shop.Line,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// 批量设置线路
func BatchSetLine(c *gin.Context) {
	var form req.BatchSetLineForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	ret := global.DB.Model(model.Shop{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"line": form.Line})

	if ret.Error != nil {
		xsq_net.ErrorJSON(c, ret.Error)
		return
	}

	if err := cache.SetShopLine(); err != nil {
		xsq_net.ErrorJSON(c, errors.New("线路设置成功，但缓存更新失败"))
		return
	}

	xsq_net.Success(c)
}

// 线路名称列表
func LineList(c *gin.Context) {
	var (
		res []*rsp.LineListRsp
	)

	db := global.DB

	ret := db.Raw("SELECT DISTINCT(`line`) AS line FROM `t_shop` WHERE line != ''")

	if ret.Error != nil {
		xsq_net.ErrorJSON(c, ret.Error)
		return
	}

	ret.Scan(&res)

	xsq_net.SucJson(c, res)
}

// 批量设置配送方式
func BatchSetDistributionType(c *gin.Context) {
	var form req.BatchSetDistributionTypeForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	var shopIds []int

	result := db.Model(&model.Shop{}).Select("shop_id").Where("id in (?)", form.Ids).Find(&shopIds)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	type Res struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
		Msg  string      `json:"msg"`
	}

	var (
		res Res
		mp  = make(map[string]interface{}, 8)
	)

	mp["shop_ids"] = shopIds
	mp["delivery_id"] = form.DistributionType

	body, err := request.Post("api/v1/remote/update/shop", mp)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	err = json.Unmarshal(body, &res)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	if res.Code != 200 {
		xsq_net.ErrorJSON(c, errors.New(res.Msg))
		return
	}

	result = db.Model(model.Shop{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"distribution_type": form.DistributionType})

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

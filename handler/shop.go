package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
func SyncU8Shop(c *gin.Context) {

	type Shop struct {
		ShopId    int    `gorm:"column:cCusCode" json:"shop_id"`         //门店ID
		ShopName  string `gorm:"column:cCusName" json:"shop_name"`       //门店名称
		Warehouse string `gorm:"column:column:cWhName" json:"warehouse"` //所属仓库
		HouseCode string `gorm:"column:cCusWhCode" json:"house_code"`    //仓库编码
		Typ       string `gorm:"column:cCCName3" json:"typ"`             //门店类型 直营店 加盟店
		Province  string `gorm:"column:cCCName2" json:"province"`        //省
		City      string `gorm:"column:cCCName1" json:"city"`            //市
		District  string `gorm:"column:cCCName" json:"district"`         //区
		Line      string `gorm:"column:cCusDefine3" json:"line"`         //线路
		ShopCode  string `gorm:"column:cCusDefine4" json:"shop_code"`    //门店编码
	}

	var shopList []Shop

	result := global.SqlServer.Table("view_customer_base").Select("cCusCode,cCusName,cWhName,cCCName3,cCCName2,cCCName1,cCCName,cCusWhCode,cCusDefine3,cCusDefine4").Find(&shopList)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	var (
		shopMp = make(map[int]model.Shop, len(shopList))
		shop   = make([]model.Shop, 0, len(shopList))
		dbShop []model.Shop
	)

	db := global.DB

	result = db.Find(&dbShop)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//构造 shop表 map[shop_id]shop
	for _, s := range dbShop {
		shopMp[s.ShopId] = s
	}

	//step1:循环用友店铺数据
	//step2:如果map中存在用友shop_id更新，没有新增到map中
	//step3:循环map构造shop表数据
	for _, sl := range shopList {

		v, ok := shopMp[sl.ShopId]

		val := model.Shop{
			ShopId:    sl.ShopId,
			ShopName:  sl.ShopName,
			HouseCode: sl.HouseCode,
			Warehouse: sl.Warehouse,
			Typ:       sl.Typ,
			Province:  sl.Province,
			City:      sl.City,
			District:  sl.District,
			ShopCode:  sl.ShopCode,
		}
		//step2:如果map中存在用友shop_id更新，没有新增到map中
		if ok {
			val.Id = v.Id
			val.Line = v.Line
		}

		shop = append(shop, val)
	}

	result = db.Save(&shop)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if err := cache.SetShopLine(); err != nil {
		xsq_net.ErrorJSON(c, errors.New("店铺更新成功，但缓存更新失败"))
		return
	}

	xsq_net.Success(c)
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

	//更新时不能更新线路，新增时保存线路
	maxShopId := 0

	res := db.Raw("SELECT id FROM t_shop ORDER BY id DESC LIMIT 1")
	if res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			xsq_net.ErrorJSON(c, res.Error)
			return
		}
	} else {
		res.Scan(&maxShopId)
	}

	for _, shop := range result.Data {
		//接口查到的id比当前存储的最大id大时，保存数据
		if shop.Id <= maxShopId {
			continue
		}

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

	//有需要保存的店铺数据
	if len(shops) > 0 {
		res = db.Save(&shops)

		if res.Error != nil {
			xsq_net.ErrorJSON(c, res.Error)
		}

		if err = cache.SetShopLine(); err != nil {
			xsq_net.ErrorJSON(c, errors.New("店铺更新成功，但缓存更新失败"))
			return
		}
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

	if err := c.ShouldBind(&form); err != nil {
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

	if err := c.ShouldBind(&form); err != nil {
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

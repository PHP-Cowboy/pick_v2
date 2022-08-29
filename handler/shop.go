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
	"pick_v2/model/other"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/request"
	"pick_v2/utils/xsq_net"
	"time"
)

type ShopResponse struct {
	Code int          `json:"code"`
	Data []other.Shop `json:"data"`
}

// 同步门店
func SyncShop(c *gin.Context) {

	path := "api/v1/remote/shop/list"

	body, err := request.Get(path)
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

	var shops []*other.Shop

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

		shops = append(shops, &other.Shop{
			Id:        shop.Id,
			ShopId:    shop.ShopId,
			ShopName:  shop.ShopName,
			HouseCode: shop.HouseCode,
			Warehouse: shop.Warehouse,
			Typ:       shop.Typ,
			Province:  shop.Province,
			City:      shop.City,
			District:  shop.District,
			//Line:      shop.Line,
			ShopCode: shop.ShopCode,
			Status:   shop.Status,
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
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
		shops []other.Shop
	)

	db := global.DB.Model(other.Shop{})

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

	ret := global.DB.Model(other.Shop{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"line": form.Line})

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

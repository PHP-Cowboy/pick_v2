package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/request"
	"pick_v2/utils/xsq_net"
)

type ClassResponse struct {
	Code int                    `json:"code"`
	Data []model.Classification `json:"data"`
}

// 同步分类
func SyncClassification(c *gin.Context) {
	url := "api/v1/remote/goods/typ/list"

	body, err := request.Post(url, nil)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var result ClassResponse

	err = json.Unmarshal(body, &result)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	db := global.DB

	var (
		class    []model.Classification
		saveData []model.Classification
	)

	ret := db.Find(&class)

	if ret.Error != nil {
		xsq_net.ErrorJSON(c, ret.Error)
		return
	}

	set := make(map[string]struct{}, len(class))

	//数据库表中分类数据存入map
	for _, c := range class {
		set[c.GoodsClass] = struct{}{}
	}

	//检验是否在map中，不在的存入数据库表
	for _, d := range result.Data {
		_, ok := set[d.GoodsClass]
		if !ok {
			saveData = append(saveData, d)
		}
	}

	if len(saveData) > 0 {
		ret = db.Save(&saveData)
		if ret.Error != nil {
			xsq_net.ErrorJSON(c, ret.Error)
			return
		}

		if err = cache.SetClassification(); err != nil {
			xsq_net.ErrorJSON(c, errors.New("分类更新成功，但缓存更新失败"))
			return
		}
	}

	xsq_net.Success(c)
}

// 分类列表
func ClassList(c *gin.Context) {
	var form req.ClassListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		res     rsp.ClassListRsp
		classes []model.Classification
	)

	db := global.DB.Model(model.Classification{})

	if form.ClassStatus == 1 { //已设置路线
		db.Where("goods_class != ''")
	} else if form.ClassStatus == 2 { //未设置路线
		db.Where("goods_class = ''")
	}

	if form.WarehouseClass != "" {
		db.Where("warehouse_class = ?", form.WarehouseClass)
	}

	ret := db.Find(&classes)

	if ret.Error != nil {
		xsq_net.ErrorJSON(c, ret.Error)
		return
	}

	res.Total = ret.RowsAffected

	ret = db.Scopes(model.Paginate(form.Page, form.Size)).Find(&classes)

	list := make([]*rsp.Class, 0, form.Size)

	for _, class := range classes {
		list = append(list, &rsp.Class{
			Id:             class.Id,
			GoodsClass:     class.GoodsClass,
			WarehouseClass: class.WarehouseClass,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// 批量设置分类
func BatchSetClass(c *gin.Context) {
	var form req.BatchSetClassForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	ret := global.DB.Model(model.Classification{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"warehouse_class": form.WarehouseClass})

	if ret.Error != nil {
		xsq_net.ErrorJSON(c, ret.Error)
		return
	}

	if err := cache.SetClassification(); err != nil {
		xsq_net.ErrorJSON(c, errors.New("分类更新成功，但缓存更新失败"))
		return
	}

	xsq_net.Success(c)
}

// 仓库分类名列表
func ClassNameList(c *gin.Context) {
	var (
		res []*rsp.ClassNameListRsp
	)

	db := global.DB

	ret := db.Raw("SELECT DISTINCT(`warehouse_class`) AS warehouse_class FROM `t_classification` WHERE warehouse_class != ''")

	if ret.Error != nil {
		xsq_net.ErrorJSON(c, ret.Error)
		return
	}

	ret.Scan(&res)

	xsq_net.SucJson(c, res)
}

func GoodsClassList(c *gin.Context) {
	var (
		class []model.Classification
		res   []*rsp.GoodsClassListRsp
	)

	ret := global.DB.Find(&class)

	if ret.Error != nil {
		xsq_net.ErrorJSON(c, ret.Error)
		return
	}

	for _, cl := range class {
		res = append(res, &rsp.GoodsClassListRsp{GoodsClass: cl.GoodsClass})
	}

	xsq_net.SucJson(c, res)
}

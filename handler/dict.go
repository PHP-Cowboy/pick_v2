package handler

import (
	"github.com/gin-gonic/gin"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"strings"
)

// 字典类型列表
func DictTypeList(c *gin.Context) {
	var form req.DictTypeListForm

	err := c.ShouldBind(&form)
	if err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		types []model.DictType
		res   rsp.DictTypeListRsp
	)

	db := global.DB

	result := db.Where(model.DictType{Code: form.Code, Name: form.Name}).Find(&types)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	db.Where(model.DictType{Code: form.Code, Name: form.Name}).Scopes(model.Paginate(form.Page, form.Size)).Find(&types)

	list := make([]*rsp.DictType, 0)
	for _, t := range types {
		list = append(list, &rsp.DictType{
			Code:       t.Code,
			Name:       t.Name,
			CreateTime: t.CreateTime.Format(timeutil.TimeFormat),
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// 新增字典类型
func CreateDictType(c *gin.Context) {
	var form req.CreateDictTypeForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	dict := model.DictType{
		Code: strings.ToLower(form.Code),
		Name: form.Name,
	}

	result := db.Save(&dict)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 修改字典类型
func ChangeDictType(c *gin.Context) {
	var form req.ChangeDictTypeForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		dictType model.DictType
	)

	db := global.DB

	result := db.Where(model.DictType{Code: form.Code}).First(&dictType)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		xsq_net.ErrorJSON(c, ecode.DataNotExist)
		return
	}

	result = db.Model(&dictType).Updates(form)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)

}

// 字典数据列表
func DictList(c *gin.Context) {
	var form req.DictListForm

	err := c.ShouldBind(&form)
	if err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		dict []model.Dict
	)

	db := global.DB

	result := db.Where(model.Dict{TypeCode: form.Code}).Find(&dict)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]*rsp.DictListRsp, 0, result.RowsAffected)
	for _, d := range dict {
		list = append(list, &rsp.DictListRsp{
			TypeCode:   d.TypeCode,
			Code:       d.Code,
			Name:       d.Name,
			Value:      d.Value,
			IsEdit:     d.IsEdit,
			CreateTime: d.CreateTime.Format(timeutil.TimeFormat),
		})
	}

	xsq_net.SucJson(c, list)
}

// 新增字典数据
func CreateDict(c *gin.Context) {
	var form req.CreateDictForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	var (
		dictType model.DictType
	)

	result := db.Where(&model.DictType{Code: form.TypeCode}).First(&dictType)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if result.RowsAffected != 1 {
		xsq_net.ErrorJSON(c, ecode.DataNotExist)
		return
	}

	dict := model.Dict{
		Code:     strings.ToLower(form.Code),
		TypeCode: strings.ToLower(form.TypeCode),
		Name:     form.Name,
		Value:    form.Value,
		IsEdit:   form.IsEdit,
	}

	result = db.Save(&dict)
	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 修改字典数据
func ChangeDict(c *gin.Context) {
	var form req.ChangeDictForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		dict model.Dict
	)

	db := global.DB

	result := db.Where(model.Dict{TypeCode: form.TypeCode, Code: form.Code}).First(&dict)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		xsq_net.ErrorJSON(c, ecode.DataNotExist)
		return
	}

	result = db.Model(&dict).Updates(form)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)

}

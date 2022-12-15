package dao

import (
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"strings"
)

// 新增字典数据
func CreateDict(db *gorm.DB, form req.CreateDictForm) (err error) {

	err, _ = model.GetDictTypeByPk(db, form.TypeCode)

	if err != nil {
		return
	}

	err = model.GetDictExistByPk(db, form.TypeCode, form.Code)

	if err != nil {
		return
	}

	dict := model.Dict{
		Code:     strings.ToLower(form.Code),
		TypeCode: strings.ToLower(form.TypeCode),
		Name:     form.Name,
		Value:    form.Value,
		IsEdit:   form.IsEdit,
	}

	err = model.DictSave(db, &dict)

	if err != nil {
		return
	}

	return
}

// 修改字典数据
func ChangeDict(db *gorm.DB, form req.ChangeDictForm) (err error) {
	var (
		dict model.Dict
	)

	err, dict = model.GetDictByPk(db, form.TypeCode, form.Code)

	if err != nil {
		return
	}

	if dict.IsEdit == 0 {
		err = ecode.DataCannotBeModified
		return
	}

	err = db.Model(&dict).Updates(form).Error

	if err != nil {
		return
	}

	return
}

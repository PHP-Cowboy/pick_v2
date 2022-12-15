package model

import (
	"errors"
	"gorm.io/gorm"
	"pick_v2/utils/ecode"
	"time"
)

// 字典表
type Dict struct {
	TypeCode   string     `gorm:"type:varchar(20);not null;primaryKey;comment:字典类型编码"`
	Code       string     `gorm:"type:varchar(50);not null;primaryKey;comment:字典编码"` //这里改成 DictCode
	Name       string     `gorm:"type:varchar(20);not null;comment:字典名称"`
	Value      string     `gorm:"type:varchar(20);not null;comment:字典值"`
	IsEdit     int        `gorm:"type:tinyint;not null;default:0;comment:是否可编辑:0:否,1:是"`
	CreateTime time.Time  `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime time.Time  `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime *time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
}

func DictSave(db *gorm.DB, dict *Dict) (err error) {
	err = db.Model(&Dict{}).Save(dict).Error
	return err
}

func GetDictByPk(db *gorm.DB, typeCode, code string) (err error, dict Dict) {

	err = db.Model(&Dict{}).Where("type_code = ? and code = ?", typeCode, code).First(&dict).Error

	return
}

func GetDictExistByPk(db *gorm.DB, typeCode, code string) (err error) {

	var dict Dict

	err = db.Model(&Dict{}).Where("type_code = ? and code = ?", typeCode, code).First(&dict).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return
	}

	return ecode.DataAlreadyExist
}

func GetDictByTypeCodeAndValue(db *gorm.DB, typeCode, value string) (err error, dict Dict) {
	err = db.Model(&Dict{}).Where("type_code = ? and value = ?", typeCode, value).First(&dict).Error
	return
}

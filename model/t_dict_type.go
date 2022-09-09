package model

import "time"

// 字典类型表
type DictType struct {
	Code       string     `gorm:"type:varchar(50);primaryKey;comment:字典类型编码"` //这里改成 TypeCode
	Name       string     `gorm:"type:varchar(20);not null;comment:字典类型名称"`
	CreateTime time.Time  `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime time.Time  `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime *time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
}

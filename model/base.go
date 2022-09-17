package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

const TableOptions string = "gorm:table_options"

func GetOptions(tableName string) string {
	return "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci comment '" + tableName + "'"
}

func GetUserOptions(tableName string) string {
	return "ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci comment '" + tableName + "'"
}

type Base struct {
	Id         int       `gorm:"primaryKey;type:int(11) unsigned AUTO_INCREMENT;comment:id"`
	CreateTime time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
}

type Creator struct {
	CreatorId int    `gorm:"type:int(11) unsigned;comment:操作人id"`
	Creator   string `gorm:"type:varchar(32);comment:操作人昵称"`
}

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

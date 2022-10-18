package model

import "time"

type InvTask struct {
	OrderNo       string    `gorm:"primaryKey;type:varchar(64);comment:盘点单号"`
	CreateTime    time.Time `gorm:"autoCreateTime;type:datetime;not null;comment:创建时间"`
	UpdateTime    time.Time `gorm:"autoUpdateTime;type:datetime;not null;comment:更新时间"`
	DeleteTime    time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
	TaskDate      *MyTime   `gorm:"type:date;not null;comment:任务日期"`
	TaskName      string    `gorm:"type:varchar(64);comment:任务名称"`
	WarehouseId   int       `gorm:"type:int(11);comment:盘点仓库ID"`
	Warehouse     string    `gorm:"type:varchar(64);comment:盘点仓库"`
	BookNum       int       `gorm:"type:int(11);not null;default:0;comment:账面数量"`
	InventoryNum  int       `gorm:"type:int(11);not null;default:0;comment:盘点数量"`
	ProfitLossNum int       `gorm:"type:int(11);not null;default:0;comment:盈亏数量"`
	Remark        string    `gorm:"type:varchar(255);default:'';comment:备注"`
	Status        int       `gorm:"type:tinyint;default:1;comment:状态:1:进行中,2:结束"`
	Version       int       `gorm:"type:int(11);not null;default:0;comment:版本"`
}

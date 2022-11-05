package model

import "time"

// 自建盘点任务表
type InvTaskSelfBuilt struct {
	Base
	OrderNo    string    `gorm:"index;type:varchar(64);comment:盘点单号"`
	CreateTime time.Time `gorm:"type:datetime;comment:创建时间"` //autoCreateTime 针对 int类型
	UpdateTime time.Time `gorm:"type:datetime;comment:更新时间"`
	DeleteTime time.Time `gorm:"type:datetime;default:null;comment:删除时间"`
	TaskName   string    `gorm:"type:varchar(64);comment:任务名称"`
	Status     int       `gorm:"type:tinyint;default:1;comment:状态:1:进行中,2:结束"`
}

type SelfBuiltJoinTask struct {
	Id         int       `json:"id"`
	CreateTime time.Time `json:"create_time"`
	OrderNo    string    `json:"order_no"`
	TaskName   string    `json:"task_name"`
	Status     int       `json:"status"`
	BookNum    float64   `json:"book_num"`
	Remark     string    `json:"remark"`
}

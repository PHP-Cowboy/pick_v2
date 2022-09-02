package model

// u8 推送日志
type StockLog struct {
	Base
	Number      string `gorm:"type:varchar(64);default:'';comment:订单编号"`
	PickId      int    `gorm:"type:int(11) unsigned;default:0;comment:拣货表id"`
	Status      int    `gorm:"type:tinyint;default:0;comment:状态:0:创建,1:推送成功,2:推送失败"`
	RequestXml  string `gorm:"type:text;comment:请求xml"`
	ResponseXml string `gorm:"type:text;comment:响应xml"`
	Msg         string `gorm:"type:text;comment:信息"`
	ResponseNo  string `gorm:"type:varchar(64);default:'';comment:U8返回单号"`
	ShopName    string `gorm:"type:varchar(64);default:'';not null;comment:店铺名称"`
}

package model

// u8 推送日志
type StockLog struct {
	Base
	Number      string `gorm:"type:varchar(64);default:'';comment:订单编号" json:"number"`
	BatchId     int    `gorm:"type:int(11) unsigned;comment:批次表id" json:"batch_id"`
	Status      int    `gorm:"type:tinyint;default:0;comment:状态:0:创建,1:推送成功,2:推送失败" json:"status"`
	RequestXml  string `gorm:"type:text;comment:请求xml" json:"request_xml"`
	ResponseXml string `gorm:"type:text;comment:响应xml" json:"response_xml"`
	ResponseNo  string `gorm:"type:varchar(64);default:'';comment:U8返回单号" json:"response_no"`
	Msg         string `gorm:"type:text;comment:信息" json:"msg"`
	ShopName    string `gorm:"type:varchar(64);default:'';not null;comment:店铺名称" json:"shop_name"`
}

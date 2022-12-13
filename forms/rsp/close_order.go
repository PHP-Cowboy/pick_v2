package rsp

import "pick_v2/model"

type CloseOrderCountRsp struct {
	PendingNum  int `json:"pending_num"`
	CompleteNum int `json:"complete_num"`
}

type CloseOrderListRsp struct {
	Total int64            `json:"total"`
	List  []CloseOrderList `json:"list"`
}

type CloseOrderList struct {
	Id               int           `json:"id"`
	Number           string        `json:"number"`
	PayAt            *model.MyTime `gorm:"type:datetime;comment:支付时间"`
	PayTotal         int           `json:"pay_total"`
	NeedCloseTotal   int           `json:"need_close_total"`
	ShopName         string        `json:"shop_name"`
	ShopType         string        `json:"shop_type"`
	DistributionType int           `json:"distribution_type"`
	Province         string        `gorm:"type:varchar(64);comment:省"`
	City             string        `gorm:"type:varchar(64);comment:市"`
	District         string        `gorm:"type:varchar(64);comment:区"`
	OrderRemark      string        `gorm:"type:varchar(512);comment:订单备注"`
	Status           int           `gorm:"type:tinyint;default:1;comment:状态:1:处理中,2:已完成"`
}

type CloseOrderDetailRsp struct {
	Number           string           `json:"number"`
	ShopName         string           `json:"shop_name"`
	DistributionType int              `json:"distribution_type"`
	District         string           `json:"district"`
	Status           int              `json:"status"`
	OrderRemark      string           `json:"order_remark"`
	List             []OrderGoodsList `json:"list"`
}

type OrderGoodsList struct {
	GoodsName      string `json:"goods_name"`
	GoodsSpe       string `json:"goods_spe"`
	PayCount       int    `json:"pay_count"`
	CloseCount     int    `json:"close_count"`      //已关闭数量
	NeedCloseCount int    `json:"need_close_count"` //需关闭数量
	GoodsRemark    string `gorm:"type:varchar(255);comment:商品备注"`
}

// 关闭订单 订单号和类型
type CloseOrderNumberTyp struct {
	Number string
	Typ    int
}

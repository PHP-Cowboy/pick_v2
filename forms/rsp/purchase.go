package rsp

import "pick_v2/model"

type OrderRsp struct {
	Code int         `json:"code"`
	Data []OrderInfo `json:"data"`
	Msg  string      `json:"msg"`
}

type OrderInfo struct {
	ShopID           int           `json:"shop_id"`
	ShopName         string        `json:"shop_name"`
	ShopType         string        `json:"shop_type"`
	Number           string        `json:"number"`
	HouseCode        string        `json:"house_code"`
	Line             string        `json:"line"`
	ShopCode         string        `json:"shop_code"`
	DistributionType int           `json:"distribution_type"`
	OrderRemark      string        `json:"order_remark"`
	PayAt            *model.MyTime `json:"pay_at"`
	DeliveryAt       model.MyTime  `json:"delivery_at"`
	Province         string        `json:"province"`
	City             string        `json:"city"`
	District         string        `json:"district"`
	OrderID          int           `json:"order_id"`
	Address          string        `json:"address"`
	ConsigneeName    string        `json:"consignee_name"`
	ConsigneeTel     string        `json:"consignee_tel"`
	GoodsInfo        []GoodsInfo   `json:"goods_info"`
}

type GoodsInfo struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Sku           string `json:"sku"`
	GoodsType     string `json:"goods_type"`
	GoodsSpe      string `json:"goods_spe"`
	Shelves       string `json:"shelves"`
	DiscountPrice int    `json:"discount_price"`
	GoodsUnit     string `json:"goods_unit"`
	SaleUnit      string `json:"sale_unit"`
	SaleCode      string `json:"sale_code"`
	PayCount      int    `json:"pay_count"`
	GoodsRemark   string `json:"goods_remark"`
}

type CloseOrderRsp struct {
	Code int        `json:"code"`
	Data CloseOrder `json:"data"`
	Msg  string     `json:"msg"`
}

type CloseOrder struct {
	Number           string           `json:"number"`            //订单编号
	ShopName         string           `json:"shop_name"`         //店铺名称
	ShopType         string           `json:"shop_type"`         //店铺类型
	DistributionType int              `json:"distribution_type"` //配送方式
	OrderRemark      string           `json:"order_remark"`      //订单备注
	PayAt            *model.MyTime    `json:"pay_at"`            //支付时间
	PayTotal         int              `json:"pay_total"`         //下单总数
	Province         string           `json:"province"`          //省
	City             string           `json:"city"`              //市
	District         string           `json:"district"`          //区
	Typ              int              `json:"typ"`               //1.部分关闭,2.全单关闭,
	Applicant        string           `json:"applicant"`         //申请人
	ApplyTime        *model.MyTime    `json:"apply_time"`        //申请时间
	GoodsInfo        []CloseGoodsInfo `json:"goods_info"`        //被关闭商品明细
}

type CloseGoodsInfo struct {
	ID             int    `json:"id"`               //订单明细表ID
	Name           string `json:"name"`             //商品名称
	Sku            string `json:"sku"`              //sku
	GoodsSpe       string `json:"goods_spe"`        //商品规格
	GoodsUnit      string `json:"goods_unit"`       //商品单位
	PayCount       int    `json:"pay_count"`        //下单数量
	CloseCount     int    `json:"close_count"`      //已关闭数量
	NeedCloseCount int    `json:"need_close_count"` //需关闭数量
	GoodsRemark    string `json:"goods_remark"`     //商品备注
}

type GetBatchOrderAndGoodsRsp struct {
	Count int        `json:"count"`
	List  []OutOrder `json:"list"`
}

type OutOrder struct {
	DistributionType int          `json:"distribution_type"`
	PayAt            model.MyTime `json:"pay_at"`
	OrderId          int          `json:"order_id"`
	GoodsInfo        []OutGoods   `json:"goods_info"`
}

type OutGoods struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Sku           string `json:"sku"`
	GoodsType     string `json:"goods_type"`
	GoodsSpe      string `json:"goods_spe"`
	DiscountPrice int    `json:"discount_price"`
	GoodsUnit     string `json:"goods_unit"`
	SaleUnit      string `json:"sale_unit"`
	SaleCode      string `json:"sale_code"`
	OutCount      int    `json:"out_count"`
	OutAt         string `json:"out_at"`
	Number        string `json:"number"`
	CkNumber      string `json:"ck_number"`
}

type CloseOrderResult struct {
	Id  int `json:"id"`
	Typ int `json:"typ"`
}

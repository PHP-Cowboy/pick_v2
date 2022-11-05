package rsp

import (
	"pick_v2/model"
)

type SyncTaskRsp struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data []SyncTask `json:"data"`
}

type SyncTask struct {
	CcvCode      string  `json:"ccv_code"`
	CInvCode     string  `json:"c_inv_code"`
	IcvQuantity  float64 `json:"icv_quantity"`
	IcvcQuantity float64 `json:"icvc_quantity"`
	CInvName     string  `json:"c_inv_name"`
	CInvStd      string  `json:"c_inv_std"`
	Cate         string  `json:"cate"`
	Dcate        string  `json:"dcate"`
	CComUnitName string  `json:"c_com_unit_name"`
	CWhCode      string  `json:"c_wh_code"`
	CWhName      string  `json:"c_wh_name"`
	CcvMeno      string  `json:"ccv_meno"`
}

type InvNumSum struct {
	Sum     float64 `json:"sum"`
	OrderNo string  `json:"order_no"`
}

type TaskListRsp struct {
	Total int64       `json:"total"`
	Data  []*TaskList `json:"list"`
}

type TaskList struct {
	OrderNo       string        `json:"order_no"`
	TaskName      string        `json:"task_name"`
	TaskDate      *model.MyTime `json:"task_time"`
	WarehouseId   int           `json:"warehouse_id"`
	WarehouseName string        `json:"warehouse_name"`
	BookNum       float64       `json:"book_num"`
	InventoryNum  float64       `json:"inventory_num"`
	ProfitLossNum float64       `json:"profit_loss_num"`
	Remark        string        `json:"remark"`
	Status        int           `json:"status"`
}

type SkuInvNumSum struct {
	Sum float64 `json:"sum"`
	Sku string  `json:"sku"`
}

type RecordListRsp struct {
	Total int64         `json:"total"`
	Data  []*RecordList `json:"list"`
}

type RecordList struct {
	Sku           string  `json:"sku"`
	GoodsName     string  `json:"goods_name"`
	GoodsType     string  `json:"goods_type"`
	GoodsSpe      string  `json:"goods_spe"`
	BookNum       float64 `json:"book_num"`
	InventoryNum  float64 `json:"inventory_num"`
	ProfitLossNum float64 `json:"profit_loss_num"`
	InvType       int     `json:"inv_type"`
}

type InventoryRecordListRsp struct {
	Total int64              `json:"total"`
	Data  []*InventoryRecord `json:"list"`
}

type InventoryRecord struct {
	Id           int     `json:"id"`
	SelfBuiltId  int     `json:"self_built_id"`
	OrderNo      string  `json:"order_no"`
	Sku          string  `json:"sku"`
	UserName     string  `json:"user_name"`
	CreateTime   string  `json:"create_time"`
	InventoryNum float64 `json:"inventory_num"`
	GoodsUnit    string  `json:"goods_unit"`
}

type UserInventoryRecordListRsp struct {
	Total int64                  `json:"total"`
	Data  []*UserInventoryRecord `json:"list"`
}

type UserInventoryRecord struct {
	Id           int     `json:"id"`
	SelfBuiltId  int     `json:"self_built_id"`
	OrderNo      string  `json:"order_no"`
	Sku          string  `json:"sku"`
	UserName     string  `json:"user_name"`
	GoodsName    string  `json:"goods_name"`
	GoodsSpe     string  `json:"goods_spe"`
	InventoryNum float64 `json:"inventory_num"`
	SystemNum    float64 `json:"system_num"`
}

type UserNotInventoryRecord struct {
	Sku       string `json:"sku"`
	GoodsName string `json:"goods_name"`
	GoodsSpe  string `json:"goods_spe"`
}

type SelfBuiltTaskRsp struct {
	Total int64            `json:"total"`
	List  []*SelfBuiltTask `json:"list"`
}

type SelfBuiltTask struct {
	Id            int     `json:"id"`
	CreateTime    string  `json:"create_time"`
	OrderNo       string  `json:"order_no"`
	TaskName      string  `json:"task_name"`
	Status        int     `json:"status"`
	BookNum       float64 `json:"book_num"`
	InventoryNum  float64 `json:"inventory_num"`
	ProfitLossNum float64 `json:"profit_loss_num"`
	Remark        string  `json:"remark"`
}

package rsp

import "pick_v2/model"

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
	BookNum       int           `json:"book_num"`
	InventoryNum  int           `json:"inventory_num"`
	ProfitLossNum int           `json:"profit_loss_num"`
	Remark        string        `json:"remark"`
	Status        int           `json:"status"`
}

type RecordListRsp struct {
	Total int64         `json:"total"`
	Data  []*RecordList `json:"list"`
}

type RecordList struct {
	Sku           string `json:"sku"`
	GoodsName     string `json:"goods_name"`
	GoodsType     string `json:"goods_type"`
	GoodsSpe      string `json:"goods_spe"`
	BookNum       int    `json:"book_num"`
	InventoryNum  int    `json:"inventory_num"`
	ProfitLossNum int    `json:"profit_loss_num"`
}

type InventoryRecordListRsp struct {
	Total int64              `json:"total"`
	Data  []*InventoryRecord `json:"list"`
}

type InventoryRecord struct {
	Id           int    `json:"id"`
	OrderNo      string `json:"order_no"`
	Sku          string `json:"sku"`
	UserName     string `json:"user_name"`
	CreateTime   string `json:"create_time"`
	InventoryNum int    `json:"inventory_num"`
	GoodsUnit    string `json:"goods_unit"`
}

type UserInventoryRecordListRsp struct {
	Total int64                  `json:"total"`
	Data  []*UserInventoryRecord `json:"list"`
}

type UserInventoryRecord struct {
	Id           int    `json:"id"`
	OrderNo      string `json:"order_no"`
	Sku          string `json:"sku"`
	UserName     string `json:"user_name"`
	GoodsName    string `json:"goods_name"`
	GoodsSpe     string `json:"goods_spe"`
	InventoryNum int    `json:"inventory_num"`
}

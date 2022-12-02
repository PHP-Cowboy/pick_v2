package req

type PickGoodsView struct {
	PickId       int         `json:"pick_id"`
	SaleNumber   string      `json:"sale_number"` //销售单号
	ShopId       int64       `json:"shop_id"`
	ShopName     string      `json:"shop_name"`
	HouseCode    string      `json:"houseCode"`
	Date         string      `json:"date"`          //日期
	Remark       string      `json:"remark"`        //备注
	DeliveryType int         `json:"delivery_type"` //配送方式 0-暂无 1-公司配送 2-用户自提 3-三方物流
	Line         string      `json:"line"`          //线路
	List         []PickGoods `json:"list"`
}

type PickGoods struct {
	GoodsName    string `json:"goods_name"`     //
	Sku          string `json:"sku"`            //
	Price        int64  `json:"price"`          //
	GoodsSpe     string `json:"goods_spe"`      //
	Shelves      string `json:"shelves"`        //
	RealOutCount int    `json:"real_out_count"` //
	MasterCode   string `json:"master_code"`    //主计量单位编码
	SlaveCode    string `json:"slave_code"`     //辅计量单位编码 sale_code
	GoodsUnit    string `json:"goods_unit"`     //主计量单位 goods_unit
	SlaveUnit    string `json:"slave_unit"`     //辅计量单位 sale_unit
}

type LogListForm struct {
	Paging
	Status    int    `json:"status" form:"status"`
	StartTime string `json:"start_time" form:"start_time"`
	EndTime   string `json:"end_time" form:"end_time"`
}

type BatchSupplementForm struct {
	Ids []int `json:"ids"`
}

type LogDetailForm struct {
	Number string `json:"number" form:"number" binding:"required"`
	PickId int    `json:"pick_id" form:"pick_id" binding:"required"`
}

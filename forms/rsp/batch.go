package rsp

type GetBatchListRsp struct {
	Total int64    `json:"total"`
	List  []*Batch `json:"list"`
}

type Batch struct {
	Id                int    `json:"id"`
	BatchName         string `json:"batch_name"`
	DeliveryStartTime string `json:"delivery_start_time"`
	DeliveryEndTime   string `json:"delivery_end_time"`
	ShopNum           int    `json:"shop_num"`
	OrderNum          int    `json:"order_num"`
	GoodsNum          int    `json:"goods_num"`
	UserName          string `json:"user_name"`
	Line              string `json:"line"`
	DeliveryMethod    int    `json:"delivery_method"`
	EndTime           string `json:"end_time"`
	Status            int    `json:"status"`
	PickNum           int    `json:"pick_num"`
	RecheckSheetNum   int    `json:"recheck_sheet_num"`
}

type GetPrePickListRsp struct {
	Total int64      `json:"total"`
	List  []*PrePick `json:"list"`
}

type PrePick struct {
	Id           int                  `json:"id"`
	ShopCode     string               `json:"shop_code"`
	ShopName     string               `json:"shop_name"`
	Line         string               `json:"line"`
	CategoryInfo map[string]PickCount `json:"category_info"`
}

type PickCount struct {
	WaitingPick int `json:"waiting_pick"`
	PickedCount int `json:"picked_count"`
}

type Ret struct {
	OutC      int
	NeedC     int
	ShopId    int
	GoodsType string
}

type BatchPoolNum struct {
	Count  int `json:"count"`
	Status int `json:"status"`
}

type GetBatchPoolNumRsp struct {
	Ongoing  int `json:"ongoing"`
	Finished int `json:"finished"`
}

type GetPrePickDetailRsp struct {
	Goods      map[string][]PrePickGoods `json:"goods"`
	RemarkList []Remark                  `json:"remark_list"`
}

type PrePickGoods struct {
	GoodsName   string `json:"goods_name"`
	GoodsSpe    string `json:"goods_spe"`
	Shelves     string `json:"shelves"`
	NeedNum     int    `json:"need_num"`
	CloseNum    int    `json:"close_num"`
	OutCount    int    `json:"out_count"`
	NeedOutNum  int    `json:"need_out_num"`
	GoodsRemark string `json:"goods_remark"`
}

type Remark struct {
	Number      string `json:"number"`
	OrderRemark string `json:"order_remark"`
	GoodsRemark string `json:"goods_remark"`
}

type GetBaseRsp struct {
	CreateTime        string `json:"create_time"`
	PayEndTime        string `json:"pay_end_time"`
	DeliveryStartTime string `json:"delivery_start_time"`
	DeliveryEndTime   string `json:"delivery_end_time"`
	DeliveryMethod    int    `json:"delivery_method"`
	Line              string `json:"line"`
	Goods             string `json:"goods"`
}

type GetPoolNumRsp struct {
	PrePickNum  int64 `json:"pre_pick_num"`
	PickNum     int   `json:"pick_num"`
	ToReviewNum int   `json:"to_review_num"`
	CompleteNum int   `json:"complete_num"`
}

type PoolNumCount struct {
	Count  int `json:"count"`
	Status int `json:"status"`
}

type PrintCallGetRsp struct {
	ShopName    string             `json:"shop_name"`    //门店名称
	JHNumber    string             `json:"jh_number"`    //拣货订单号
	PickName    string             `json:"pick_name"`    //拣货人
	ShopType    string             `json:"shop_type"`    //门店类型
	CheckName   string             `json:"check_name"`   //复核员
	HouseName   string             `json:"house_name"`   //仓库编码
	Delivery    string             `json:"delivery"`     //配送方式
	OrderRemark string             `json:"order_remark"` //订单备注
	Consignee   string             `json:"consignee"`    //收件人
	Shop_code   string             `json:"shop_code"`    //门店编码
	Packages    int                `json:"packages"`     //件数
	Phone       string             `json:"phone"`        //收货电话
	PriType     int                `json:"pri_type"`     // 1-全部打印 2-打印箱单 3-打印出库单 第一次全打，后边的前段选
	GoodsList   []CallGetGoodsView `json:"goods_list"`   //商品信息
}

type CallGetGoodsView struct {
	SaleNumber  string         `json:"sale_number"`
	Date        string         `json:"date"`
	OrderRemark string         `json:"order_remark"`
	List        []CallGetGoods `json:"list"`
}

type CallGetGoods struct {
	GoodsName    string `json:"goods_name"`
	GoodsSpe     string `json:"goods_spe"`
	GoodsCount   int    `json:"sale_count"`     //订量
	RealOutCount int    `json:"real_out_count"` //出量
	GoodsUnit    string `json:"goods_unit"`     //主计量单位
	Price        int64  `json:"price"`
	LackCount    int    `json:"lack_count"` //欠货数量
}

type OutGoodsRsp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

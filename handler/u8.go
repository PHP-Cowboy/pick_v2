package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math"
	"net/http"
	"pick_v2/common/constant"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/labstack/gommon/random"
	"github.com/shopspring/decimal"
)

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

func SendShopXml(xml string) (string, error) {
	global.SugarLogger.Info(xml)
	var (
		requestUrl string
		err        error
		request    *http.Request
		response   *http.Response
		rspBody    []byte
	)
	requestUrl = "http://8.136.191.24/U8EAI/import.asp"

	request, err = http.NewRequest(http.MethodPost, requestUrl, strings.NewReader(xml))
	if err != nil {
		global.SugarLogger.Error("请求失败", err)
		return "", err
	}
	request.Header.Set("Content-Type", "application/xml")
	client := &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		global.SugarLogger.Error("请求用友失败1", err)
		return "", err
	}
	defer response.Body.Close()

	rspBody, err = io.ReadAll(response.Body)
	if err != nil {
		global.SugarLogger.Error("请求用友失败2", err)
		return "", err
	}
	return string(rspBody), nil

}

func GenU8Xml(order PickGoodsView, shopId int64, shopName, houseCode string) string {
	//如果一件商品数量都没有 则直接返回空
	sumCount := 0
	for _, v := range order.List {
		sumCount += v.RealOutCount
	}
	if sumCount == 0 {
		return ""
	}

	tax := "13"
	now := time.Now().Format("2006-01-02")
	id := time.Now().Format("20060102") + random.New().String(uint8(4), random.Numeric)

	doc := etree.NewDocument()
	document := doc.CreateElement("ufinterface")
	//发送方编码
	document.CreateAttr("sender", global.ServerConfig.YongSender)
	//接收方可填U8
	document.CreateAttr("receiver", "u8")
	//单据模版名
	document.CreateAttr("roottag", "consignment")
	//仓库
	document.CreateAttr("warehouse", "05")
	//唯一编号 可空
	document.CreateAttr("docid", "")
	//操作码 此处写死Add
	document.CreateAttr("proc", "ADD")
	//编码是否已转换 Y-是 N-否
	document.CreateAttr("codeexchanged", "N")
	//导出是否需要根据对照表进行转换
	document.CreateAttr("exportneedexch", "N")

	document.CreateAttr("display", "销售发货单")
	document.CreateAttr("family", "销售管理")

	saleorder := document.CreateElement("consignment")
	header := saleorder.CreateElement("header")
	body := saleorder.CreateElement("body")

	/*---------组装请求头报文------------*/
	//id
	header.CreateElement("id").CreateText(id)
	header.CreateElement("code").CreateText(id)

	//单据类型 05 发货单06 委托代销发货单00 委托代销调整单
	header.CreateElement("vouchertype").CreateText("05")
	//销售类型编码
	header.CreateElement("saletype").CreateText("01")
	//日期
	header.CreateElement("date").CreateText(now)
	//部门编号（同步用友组织结构）
	header.CreateElement("deptcode").CreateText("06")
	//职员编号（同步用友组织结构）可空
	header.CreateElement("personcode").CreateText("2051")
	//客户编号（同步用友组织结构）
	header.CreateElement("custcode").CreateText(strconv.FormatInt(shopId, 10))
	//外币名称
	header.CreateElement("currency_name").CreateText("人民币")
	//汇率 （人民币是 该值传1）
	header.CreateElement("currency_rate").CreateText("1")
	//表头税率
	header.CreateElement("taxrate").CreateText(tax)
	header.CreateElement("beginflag").CreateText("0")
	header.CreateElement("returnflag").CreateText("0")
	//制单员
	header.CreateElement("maker").CreateText("沏掌柜拣货")
	header.CreateElement("sale_cons_flag").CreateText("0")
	//散户开票的客户名称
	header.CreateElement("retail_custname").CreateText(shopName)
	//业务类型
	header.CreateElement("operation_type").CreateText("普通销售")
	//验证日期
	header.CreateElement("bcredit").CreateText("否")
	//销售单号
	header.CreateElement("define11").CreateText(order.SaleNumber)

	/*组装请求体报文 (根据商品列表 循环生成)*/
	//batch := &model.Batch{}
	for _, goods := range order.List {
		if goods.RealOutCount == 0 {
			//如果实际出库为0 则直接忽略生成XML
			continue
		}

		rsp := CalculatePrice(goods.Price, goods.RealOutCount)

		entry := body.CreateElement("entry")
		//headid
		entry.CreateElement("headid").CreateText(id)
		//仓库编码
		entry.CreateElement("warehouse_code").CreateText(houseCode)
		//存货编号
		entry.CreateElement("inventory_code").CreateText(goods.Sku)
		//主计量数量
		entry.CreateElement("quantity").CreateText(strconv.Itoa(goods.RealOutCount))
		//主计量单位编码
		entry.CreateElement("ccomunitcode").CreateText(goods.MasterCode)
		//辅助计量单位编码
		entry.CreateElement("unit_code").CreateText(goods.SlaveCode)
		//主计量单位
		entry.CreateElement("cinvm_unit").CreateText(goods.GoodsUnit)
		//辅计量单位
		entry.CreateElement("cinva_unit").CreateText(goods.SlaveUnit)
		//报价
		entry.CreateElement("quotedprice").CreateText(rsp.TaxPrice)
		//单价（原币，无税）
		entry.CreateElement("price").CreateText(rsp.NoTaxPrice)
		//含税单价（原币）
		entry.CreateElement("taxprice").CreateText(rsp.TaxPrice)
		//金额（原币，无税）
		entry.CreateElement("money").CreateText(rsp.TotalNoTaxPrice)
		//税额
		entry.CreateElement("tax").CreateText(rsp.SubTaxPrice)
		//价税合计（原币）
		entry.CreateElement("sum").CreateText(rsp.TotalTaxPrice)
		//单价（本币，无税)
		entry.CreateElement("natprice").CreateText(rsp.NoTaxPrice)
		//金额（本币，无税）
		entry.CreateElement("natmoney").CreateText(rsp.TotalNoTaxPrice)
		//税额（本币）
		entry.CreateElement("nattax").CreateText(rsp.SubTaxPrice)
		//价税合计（本币）
		entry.CreateElement("natsum").CreateText(rsp.TotalTaxPrice)
		entry.CreateElement("natdiscount").CreateText("0")
		entry.CreateElement("backflag").CreateText("正常")
		//打印用存货名称
		entry.CreateElement("inventory_printname").CreateText(goods.GoodsName)
		//税率
		entry.CreateElement("taxrate").CreateText(tax)
		//销售单号
		entry.CreateElement("cordercode").CreateText(order.SaleNumber)
		entry.CreateElement("bsaleprice").CreateText("1")
		entry.CreateElement("bgift").CreateText("0")
		entry.CreateElement("fcusminprice").CreateText("0")
		entry.CreateElement("retailrealamount").CreateText("0")
		entry.CreateElement("retailsettleamount").CreateText("0")

	}

	doc.Indent(2)
	info, _ := doc.WriteToString()
	return info
}

// 用友相关税额计算
type CalculateModel struct {
	TaxPrice        string `json:"tax_price"`          //含税单价
	NoTaxPrice      string `json:"no_tax_price"`       //无税单价
	TotalTaxPrice   string `json:"total_tax_price"`    //含税金额
	TotalNoTaxPrice string `json:"total_no_tax_price"` //无税金额
	SubTaxPrice     string `json:"sub_tax_price"`      //税额
}

func CalculatePrice(price int64, realCount int) (res CalculateModel) {
	//含税单价
	res.TaxPrice = AmountTransfer(price)
	//无税单价
	noTaxPriceInt64 := int64(math.Floor(float64(price)/float64(1.13) + 0.5))
	//无税单价
	res.NoTaxPrice = AmountTransfer(noTaxPriceInt64)
	//含税金额
	res.TotalTaxPrice = AmountTransfer(price * int64(realCount))
	//无税金额
	res.TotalNoTaxPrice = AmountTransfer(int64(math.Floor(float64(price)*float64(realCount)/1.13 + 0.5)))
	//税额
	res.SubTaxPrice = SubKeepNum(res.TotalTaxPrice, res.TotalNoTaxPrice, 2)
	return
}

// 系统金额 由int64 单位分 转 字符串 单位元,且带小数2位
func AmountTransfer(amount int64) string {
	da, _ := decimal.NewFromString(strconv.FormatInt(amount, 10))
	db, _ := decimal.NewFromString("100")
	return da.Div(db).StringFixed(2)
}

// 减法 （a-b）保留指定位小数
func SubKeepNum(a string, b string, num int32) string {
	da, _ := decimal.NewFromString(a)
	db, _ := decimal.NewFromString(b)
	return da.Sub(db).StringFixed(num)
}

// u8推送日志列表
func LogList(c *gin.Context) {
	var form req.LogListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var (
		total    int64
		stockLog []model.StockLog
		res      rsp.LogListRsp
	)

	db := global.DB

	localDb := db.Model(&model.StockLog{}).Where(model.StockLog{Status: form.Status})

	if form.StartTime != "" {
		localDb = localDb.Where("create_time >= ?", form.StartTime)
	}

	if form.EndTime != "" {
		localDb = localDb.Where("create_time <= ?", form.EndTime)
	}

	result := localDb.Count(&total)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = total

	localDb.Scopes(model.Paginate(form.Page, form.Size)).Find(&stockLog)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.LogList, 0, len(stockLog))

	for _, log := range stockLog {
		list = append(list, rsp.LogList{
			Id:          log.Id,
			CreateTime:  log.CreateTime.Format(timeutil.TimeFormat),
			UpdateTime:  log.UpdateTime.Format(timeutil.TimeFormat),
			Number:      log.Number,
			BatchId:     log.BatchId,
			PickId:      log.PickId,
			Status:      log.Status,
			RequestXml:  log.RequestXml,
			ResponseXml: log.ResponseXml,
			ResponseNo:  log.ResponseNo,
			Msg:         log.Msg,
			ShopName:    log.ShopName,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// u8 批量补单
func BatchSupplement(c *gin.Context) {
	var form req.BatchSupplementForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	second, err := cache.GetTtlKey(constant.BATCH_SUPPLEMENT)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	if second > 0 {
		xsq_net.ErrorJSON(c, errors.New(fmt.Sprintf("推送u8处理中，请%v秒后再试", second.Seconds())))
		return
	}

	expire := BaseNum * (len(form.Ids) + 10) //channel 没读取到数据时 等待了 BaseNum * 10 秒

	_, err = cache.Set(constant.BATCH_SUPPLEMENT, "1", expire)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	for _, id := range form.Ids {
		YongYouProducer(id)
	}

	xsq_net.Success(c)
}

// u8 推送
func PushYongYou(id int) {
	var (
		stockLog model.StockLog
		db       = global.DB
	)

	result := db.First(&stockLog, id)

	if result.Error != nil {
		return
	}

	shopXml, err := SendShopXml(stockLog.RequestXml)

	doc := etree.NewDocument()

	if err != nil {
		stockLog.Msg = fmt.Sprintf("SendShopXml err:", err.Error())
	} else {
		xmlErr := doc.ReadFromString(shopXml)

		if xmlErr != nil {
			stockLog.Msg = fmt.Sprintf("解析用友响应错误=", xmlErr.Error())
		} else {
			item := doc.SelectElement("ufinterface").SelectElement("item")
			code := item.SelectAttr("succeed").Value

			if code == "0" { //成功
				stockLog.ResponseNo = item.SelectAttr("u8key").Value
				stockLog.Status = 1
			} else {
				stockLog.Status = 2
			}

			stockLog.Msg = item.SelectAttr("dsc").Value
		}
	}

	stockLog.UpdateTime = time.Now()

	//db.Omit("create_time", "number", "batch_id", "shop_name","delete_time").Save(stockLog)
	db.Select("id", "update_time", "status", "request_xml", "response_xml", "msg").Save(stockLog)
}

// 推送u8拣货详情
func LogDetail(c *gin.Context) {
	var form req.LogDetailForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		res       rsp.LogDetailRsp
		pickGoods []model.PickGoods
		pickOrder model.PickOrder
		pick      model.Pick
	)

	db := global.DB

	result := db.Model(&model.PickOrder{}).Where("number = ?", form.Number).First(&pickOrder)

	if result.Error != nil {

		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.First(&pick, form.PickId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	payAt, payAtErr := time.ParseInLocation(timeutil.TimeZoneFormat, pickOrder.PayAt, time.Local)

	if payAtErr != nil {
		xsq_net.ErrorJSON(c, ecode.DataTransformationError)
		return
	}

	res.ShopName = pickOrder.ShopName
	res.PayAt = payAt.Format(timeutil.TimeFormat)

	res.PickUser = pick.PickUser

	takeOrdersTime := ""

	if pick.TakeOrdersTime != nil {
		takeOrdersTime = pick.TakeOrdersTime.Format(timeutil.TimeFormat)
	}

	res.TakeOrdersTime = takeOrdersTime
	res.ReviewUser = pick.ReviewUser

	reviewTime := ""

	if pick.ReviewTime != nil {
		reviewTime = pick.ReviewTime.Format(timeutil.TimeFormat)
	}

	res.ReviewTime = reviewTime

	result = db.Model(&model.PickGoods{}).Where(model.PickGoods{PickId: form.PickId, Number: form.Number}).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.LogDetail, 0, len(pickGoods))

	for _, pg := range pickGoods {
		list = append(list, rsp.LogDetail{
			Id:               pg.Id,
			UpdateTime:       pg.UpdateTime.Format(timeutil.TimeFormat),
			PickId:           pg.PickId,
			BatchId:          pg.BatchId,
			PrePickGoodsId:   pg.PrePickGoodsId,
			OrderGoodsId:     pg.OrderGoodsId,
			Number:           pg.Number,
			ShopId:           pg.ShopId,
			DistributionType: pg.DistributionType,
			Sku:              pg.Sku,
			GoodsName:        pg.GoodsName,
			GoodsType:        pg.GoodsType,
			GoodsSpe:         pg.GoodsSpe,
			Shelves:          pg.Shelves,
			DiscountPrice:    float64(pg.DiscountPrice) / 100,
			NeedNum:          pg.NeedNum,
			CompleteNum:      pg.CompleteNum,
			ReviewNum:        pg.ReviewNum,
			Unit:             pg.Unit,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

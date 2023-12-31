package handler

import (
	"bytes"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"pick_v2/dao"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/slice"
	"pick_v2/utils/xsq_net"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// 首批物料导出
func FirstMaterial(c *gin.Context) {
	var (
		form req.FirstMaterialExportReq
		orderNumber,
		shopName string
		pick      model.Pick
		pickGoods []model.PickGoods
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	result := db.Model(&model.PickGoods{}).Where("pick_id = ?", form.Id).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.First(&pick, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xFile := excelize.NewFile()
	sheet := xFile.NewSheet("sheet1")
	// 设置单元格的值
	xFile.MergeCell("Sheet1", "A1", "L1")
	xFile.SetCellValue("Sheet1", "A2", "序号")
	xFile.SetCellValue("Sheet1", "B2", "货架号")
	xFile.SetCellValue("Sheet1", "C2", "商品编码")
	xFile.SetCellValue("Sheet1", "D2", "商品名称")
	xFile.SetCellValue("Sheet1", "E2", "商品规格")
	xFile.SetCellValue("Sheet1", "F2", "需出数量")
	xFile.SetCellValue("Sheet1", "G2", "单位")
	xFile.SetCellValue("Sheet1", "H2", "售价")
	xFile.SetCellValue("Sheet1", "I2", "合计")
	xFile.SetCellValue("Sheet1", "J2", "拣货数")
	xFile.SetCellValue("Sheet1", "K2", "复核数")
	xFile.SetCellValue("Sheet1", "L2", "欠货数")
	xFile.SetActiveSheet(sheet)
	//设置指定行高 指定列宽
	xFile.SetRowHeight("Sheet1", 1, 30)
	xFile.SetColWidth("Sheet1", "C", "C", 30)

	startCount := 3
	tOutCount := 0
	for idx, val := range pickGoods {
		item := make([]interface{}, 0)
		item = append(item, idx+1)
		item = append(item, val.Shelves)
		item = append(item, val.Sku)
		item = append(item, val.GoodsName)
		item = append(item, val.GoodsSpe)
		item = append(item, val.NeedNum)
		item = append(item, val.Unit)
		item = append(item, AmountToFloatKeepTwo(val.DiscountPrice))
		item = append(item, AmountToFloatKeepTwo(val.DiscountPrice*val.NeedNum))
		item = append(item, "")
		item = append(item, "")
		item = append(item, "")

		xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d", startCount+idx), &item)
		tOutCount += val.NeedNum
	}
	if len(pickGoods) >= 1 {
		orderNumber = pickGoods[0].Number
	}

	shopName = pick.ShopName

	xFile.SetSheetRow("Sheet1", "A1", &[]interface{}{fmt.Sprintf("门店 : %s  订单编号: %s 配送方式 :首批设备|物料单 需出库数:%d", shopName, orderNumber, tOutCount)})

	var buffer bytes.Buffer
	_ = xFile.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	data, _ := io.ReadAll(content)
	date := time.Now().Format("20060102")
	c.Writer.Header().Add("Content-Type", "application/octet-stream;charset=utf-8")
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fmt.Sprintf("%s-.xlsx", date)+"\"")
	c.Writer.Write(data)
}

// 将int64位金额（分） 转 float 金额（元）保留2位
func AmountToFloatKeepTwo(amount int) float32 {
	if amount == 0 {
		return 0
	}
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(amount)/100), 64)
	return float32(value)
}

// 批次出库导出
func OutboundBatch(c *gin.Context) {
	var (
		form      req.OutboundBatchFormReq
		shopCodes []string
		batch     model.Batch
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	result := db.First(&batch, form.Id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			xsq_net.ErrorJSON(c, ecode.DataNotExist)
			return
		}
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	type PickGoods struct {
		ShopCode  string `json:"shop_code"`
		GoodsName string `json:"goods_name"`
		GoodsSpe  string `json:"goods_spe"`
		Unit      string `json:"unit"`
		Sku       string `json:"sku"`
		ReviewNum int    `json:"review_num"`
	}

	var pickGoods []PickGoods

	result = db.Table("t_pick_goods pg").
		Select("shop_code,goods_name,goods_spe,unit,sku,pg.review_num").
		Where("pg.batch_id = ?", form.Id).
		Joins("left join t_pick p on p.id = pg.pick_id").
		Scan(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	mp := make(map[string]map[string]string, 0)

	mpSum := make(map[string]int, 0)

	for _, good := range pickGoods {
		shopCodes = append(shopCodes, good.ShopCode)
	}

	for _, pg := range pickGoods {

		subMp, ok := mp[pg.Sku]

		if !ok {
			subMp = make(map[string]string, 0)
			mp[pg.Sku] = subMp
		}

		for _, code := range shopCodes {
			_, has := subMp[code]
			if !has {
				subMp[code] = "0"
			}

			if code == pg.ShopCode {
				subMp[code] = strconv.Itoa(pg.ReviewNum)
			}
		}

		subMp["商品名称"] = pg.GoodsName
		subMp["规格"] = pg.GoodsSpe
		subMp["单位"] = pg.Unit

		_, msOk := mpSum[pg.Sku]
		if !msOk {
			mpSum[pg.Sku] = pg.ReviewNum
		} else {
			mpSum[pg.Sku] += pg.ReviewNum
		}

		subMp["总计"] = strconv.Itoa(mpSum[pg.Sku])
	}

	column := []string{"商品名称", "规格", "单位", "总计"}

	shopCodes = slice.UniqueSlice(shopCodes)

	column = append(column, shopCodes...)

	xFile := excelize.NewFile()
	sheet := xFile.NewSheet("sheet1")
	// 设置单元格的值
	xFile.MergeCell("Sheet1", "A1", GetKey(len(column))+"1")

	for i, cn := range column {
		xFile.SetCellValue("Sheet1", GetKey(i)+"2", cn)
	}

	xFile.SetActiveSheet(sheet)
	//设置指定行高 指定列宽
	xFile.SetRowHeight("Sheet1", 1, 30)
	xFile.SetColWidth("Sheet1", "C", "C", 30)

	startCount := 2

	i := 0

	for _, val := range mp {
		i++
		item := make([]interface{}, 0)
		item = append(item, val["商品名称"])
		item = append(item, val["规格"])
		item = append(item, val["单位"])
		item = append(item, val["总计"])

		for _, code := range shopCodes {
			item = append(item, val[code])
		}

		xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d", startCount+i), &item)
	}

	xFile.SetSheetRow("Sheet1", "A1", &[]interface{}{batch.BatchName})

	var buffer bytes.Buffer
	_ = xFile.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	data, _ := io.ReadAll(content)
	date := time.Now().Format("20060102")
	c.Writer.Header().Add("Content-Type", "application/octet-stream;charset=utf-8")
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fmt.Sprintf("%s-.xlsx", date)+"\"")
	c.Writer.Write(data)
}

// 根据下标 获取 对应表头
func GetKey(index int) string {
	colCode := ""
	key := 'A'
	loop := index / 26
	if loop > 0 {
		colCode += GetKey(loop - 1)
	}
	return colCode + string(key+int32(index)%26)
}

// 欠货信息导出
func Lack(c *gin.Context) {
	var form req.LackForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		order      []model.Order
		orderGoods []model.OrderGoods
	)

	db := global.DB

	result := db.Model(&model.Order{}).Where("order_type = ?", model.LackOrderType).Find(&order)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		xsq_net.ErrorJSON(c, ecode.DataNotExist)
		return
	}

	type orderInfo struct {
		Number   string
		ShopId   int
		PayAt    *model.MyTime
		ShopName string
	}

	var (
		numbers []string
		orderMp = make(map[string]orderInfo, len(order))
	)

	for _, o := range order {
		numbers = append(numbers, o.Number)

		orderMp[o.Number] = orderInfo{
			Number:   o.Number,
			ShopId:   o.ShopId,
			PayAt:    o.PayAt,
			ShopName: o.ShopName,
		}
	}

	result = db.Model(&model.OrderGoods{}).Where("number in (?)", numbers).Find(&orderGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xFile := excelize.NewFile()
	sheet := xFile.NewSheet("sheet1")
	// 设置单元格的值
	xFile.MergeCell("Sheet1", "A1", "K1")

	xFile.SetCellValue("Sheet1", "A2", "订单号")
	xFile.SetCellValue("Sheet1", "B2", "客户编码")
	xFile.SetCellValue("Sheet1", "C2", "订单日期 (日)")
	xFile.SetCellValue("Sheet1", "D2", "客户名称")
	xFile.SetCellValue("Sheet1", "E2", "存货编码")
	xFile.SetCellValue("Sheet1", "F2", "存货名称")
	xFile.SetCellValue("Sheet1", "G2", "规格型号")
	xFile.SetCellValue("Sheet1", "H2", "单位")
	xFile.SetCellValue("Sheet1", "I2", "数量")
	xFile.SetCellValue("Sheet1", "J2", "发货数量")
	xFile.SetCellValue("Sheet1", "K2", "欠货数量")
	xFile.SetActiveSheet(sheet)
	//设置指定行高 指定列宽
	xFile.SetRowHeight("Sheet1", 1, 30)
	xFile.SetColWidth("Sheet1", "C", "C", 30)

	startCount := 3
	for idx, val := range orderGoods {
		if val.LackCount <= 0 {
			continue
		}

		o, ok := orderMp[val.Number]

		if !ok {
			continue
		}

		item := make([]interface{}, 0)
		item = append(item, o.Number)
		item = append(item, o.ShopId)
		item = append(item, o.PayAt)
		item = append(item, o.ShopName)
		item = append(item, val.Sku)
		item = append(item, val.GoodsName)
		item = append(item, val.GoodsSpe)
		item = append(item, val.GoodsUnit) //单位
		item = append(item, val.PayCount)
		item = append(item, val.OutCount)
		item = append(item, val.LackCount)

		xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d", startCount+idx), &item)
	}

	xFile.SetSheetRow("Sheet1", "A1", &[]interface{}{fmt.Sprintf("欠货信息")})

	var buffer bytes.Buffer
	_ = xFile.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	data, _ := io.ReadAll(content)
	date := time.Now().Format("20060102")
	c.Writer.Header().Add("Content-Type", "application/octet-stream;charset=utf-8")
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fmt.Sprintf("%s-.xlsx", date)+"\"")
	c.Writer.Write(data)
}

// 批次门店信息
func BatchShop(c *gin.Context) {
	var form req.BatchShopForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		prePicks []model.PrePick
		batch    model.Batch
	)

	db := global.DB

	result := db.Model(&model.PrePick{}).Where("batch_id = ?", form.Id).Order("shop_code ASC").Find(&prePicks)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.Batch{}).First(&batch, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xFile := excelize.NewFile()
	sheet := xFile.NewSheet("sheet1")
	// 设置单元格的值
	xFile.MergeCell("Sheet1", "A1", "H1")
	xFile.SetCellValue("Sheet1", "A2", "序号")
	xFile.SetCellValue("Sheet1", "B2", "门店编码")
	xFile.SetCellValue("Sheet1", "C2", "门店名称")
	xFile.SetCellValue("Sheet1", "D2", "配送位置")
	xFile.SetCellValue("Sheet1", "E2", "数量")
	xFile.SetCellValue("Sheet1", "F2", "冷冻件数")
	xFile.SetCellValue("Sheet1", "G2", "页数")
	xFile.SetCellValue("Sheet1", "H2", "备注")
	xFile.SetActiveSheet(sheet)
	//设置指定行高 指定列宽
	xFile.SetRowHeight("Sheet1", 1, 20)
	xFile.SetColWidth("Sheet1", "C", "C", 40)

	startCount := 3
	for idx, val := range prePicks {
		item := make([]interface{}, 0)
		item = append(item, idx+1)
		item = append(item, val.ShopCode)
		item = append(item, val.ShopName)

		xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d", startCount+idx), &item)
	}

	xFile.SetSheetRow("Sheet1", "A1", &[]interface{}{fmt.Sprintf("批次-%s门店信息", batch.BatchName)})

	var buffer bytes.Buffer
	_ = xFile.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	data, _ := io.ReadAll(content)
	date := time.Now().Format("20060102")
	c.Writer.Header().Add("Content-Type", "application/octet-stream;charset=utf-8")
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fmt.Sprintf("%s-.xlsx", date)+"\"")
	c.Writer.Write(data)
}

// 批次门店物料表
func BatchShopMaterial(c *gin.Context) {
	var form req.BatchShopMaterialForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		batch     model.Batch
		pickGoods []model.PickGoods
	)

	db := global.DB

	result := db.Model(&model.Batch{}).First(&batch, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	var prePickAndGoods []rsp.PrePickAndGoods

	//获取批次预拣池数据
	result = db.Table("t_pre_pick_goods pg").
		Select("shop_code,shop_name,pg.id as pre_pick_goods_id,goods_name,goods_type,goods_spe,unit,need_num").
		Joins("left join t_pre_pick pp on pg.pre_pick_id = pp.id").
		Where("pp.batch_id = ?", form.Id). //t_pre_pick.batch_id 有索引，t_pre_pick_goods batch_id 没有
		Find(&prePickAndGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//获取批次全部拣货池数据
	result = db.Model(&model.PickGoods{}).Where("batch_id = ?", form.Id).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//拣货商品明细表 复核数map
	var pickGoodsMp = make(map[int]int, len(pickGoods))

	for _, good := range pickGoods {
		pickGoodsMp[good.PrePickGoodsId] = good.ReviewNum
	}

	xFile := excelize.NewFile()
	sheet := xFile.NewSheet("sheet1")
	// 设置单元格的值
	xFile.MergeCell("Sheet1", "A1", "I1")
	xFile.SetCellValue("Sheet1", "A2", "序号")
	xFile.SetCellValue("Sheet1", "B2", "门店编码")
	xFile.SetCellValue("Sheet1", "C2", "门店名称")
	xFile.SetCellValue("Sheet1", "D2", "仓库商品分类")
	xFile.SetCellValue("Sheet1", "E2", "商品名称")
	xFile.SetCellValue("Sheet1", "F2", "商品规格")
	xFile.SetCellValue("Sheet1", "G2", "商品单位")
	xFile.SetCellValue("Sheet1", "H2", "需拣数量")
	xFile.SetCellValue("Sheet1", "I2", "已拣数量")
	xFile.SetActiveSheet(sheet)
	//设置指定行高 指定列宽
	xFile.SetRowHeight("Sheet1", 1, 20)
	xFile.SetColWidth("Sheet1", "C", "C", 40)

	startCount := 3
	for idx, val := range prePickAndGoods {

		reviewNum, pgOk := pickGoodsMp[val.PrePickGoodsId]

		//不存在时赋值为0，这种是没有进入拣货池，还在待拣池中
		if !pgOk {
			reviewNum = 0
		}

		item := make([]interface{}, 0)
		item = append(item, idx+1)
		item = append(item, val.ShopCode)
		item = append(item, val.ShopName)
		item = append(item, val.GoodsType)
		item = append(item, val.GoodsName)
		item = append(item, val.GoodsSpe)
		item = append(item, val.Unit)
		item = append(item, val.NeedNum)
		item = append(item, reviewNum)

		xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d", startCount+idx), &item)
	}

	err, methodName := GetDeliveryMethod(db, batch.DeliveryMethod)

	if err != nil {
		methodName = "配送方式"
	}

	xFile.SetSheetRow("Sheet1", "A1", &[]interface{}{fmt.Sprintf("%s-%s-门店物料信息", batch.BatchName, methodName)})

	var buffer bytes.Buffer
	_ = xFile.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	data, _ := io.ReadAll(content)
	date := time.Now().Format("20060102")
	c.Writer.Header().Add("Content-Type", "application/octet-stream;charset=utf-8")
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fmt.Sprintf("%s-.xlsx", date)+"\"")
	c.Writer.Write(data)
}

// 拣货任务导出
func BatchTask(c *gin.Context) {
	var (
		form      req.BatchTaskForm
		shopCodes []string
		pick      model.Pick
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	result := db.First(&pick, form.Id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			xsq_net.ErrorJSON(c, ecode.DataNotExist)
			return
		}
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	err, pickGoods := model.GetPickGoodsJoinOrderByPickIdAndTaskId(db, form.Id, pick.TaskId)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	mp := make(map[string]map[string]string, 0)

	mpSum := make(map[string]int, 0)
	skuCodeSum := make(map[string]int, 0)

	codeSum := make(map[string]int, 0)

	for _, good := range pickGoods {
		shopCodes = append(shopCodes, good.ShopCode)
	}

	shopCodes = slice.UniqueSlice(shopCodes)

	sort.Strings(shopCodes)

	for _, pg := range pickGoods {

		subMp, ok := mp[pg.Sku]

		if !ok {
			subMp = make(map[string]string, 0)
			mp[pg.Sku] = subMp
		}

		num := 0

		if form.Typ == 1 {
			num = pg.NeedNum
		} else if form.Typ == 3 {
			num = pg.ReviewNum
		}

		for _, code := range shopCodes {

			_, has := subMp[code]
			if !has {
				subMp[code] = ""
			}

			if code == pg.ShopCode {
				mpKey := fmt.Sprintf("%s%s", pg.Sku, code)
				skuCodeSumVal, skuCodeSumOk := skuCodeSum[mpKey]

				if !skuCodeSumOk {
					skuCodeSumVal = 0
				}

				skuCodeSumVal += num

				skuCodeSum[mpKey] = skuCodeSumVal

				subMp[code] = strconv.Itoa(skuCodeSumVal)

				//codeSum[code] += skuCodeSumVal
			}
		}

		subMp["商品名称"] = pg.GoodsName
		subMp["规格"] = pg.GoodsSpe
		subMp["单位"] = pg.Unit

		_, msOk := mpSum[pg.Sku]
		if !msOk {
			mpSum[pg.Sku] = num
		} else {
			mpSum[pg.Sku] += num
		}

		subMp["总计"] = strconv.Itoa(mpSum[pg.Sku])

		//codeSum["total"] += mpSum[pg.Sku]

		mp[pg.Sku] = subMp
	}

	column := []string{"商品名称", "规格", "单位", "总计"}

	shopCodes = slice.UniqueSlice(shopCodes)

	column = append(column, shopCodes...)

	total := 0
	codeVal := 0

	for _, sp := range mp {
		total, _ = strconv.Atoi(sp["总计"])

		codeSum["total"] += total

		for _, code := range shopCodes {
			codeVal, _ = strconv.Atoi(sp[code])
			codeSum[code] += codeVal
		}

	}

	xFile := excelize.NewFile()
	sheet := xFile.NewSheet("sheet1")

	xFile.MergeCell("Sheet1", "A1", GetKey(len(column))+"1")
	for i, cn := range column {
		xFile.SetCellValue("Sheet1", GetKey(i)+"2", cn)
	}

	xFile.SetActiveSheet(sheet)
	//设置指定行高 指定列宽
	xFile.SetRowHeight("Sheet1", 1, 30)
	xFile.SetColWidth("Sheet1", "C", "C", 30)

	startCount := 2

	i := 0

	for _, val := range mp {
		if val["总计"] == "0" {
			//总计为0的不输出到Excel
			continue
		}

		i++
		item := make([]interface{}, 0)
		item = append(item, val["商品名称"])
		item = append(item, val["规格"])
		item = append(item, val["单位"])
		item = append(item, val["总计"])

		for _, code := range shopCodes {
			item = append(item, val[code])
		}

		xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d", startCount+i), &item)
	}

	item := make([]interface{}, 0)
	item = append(item, "")
	item = append(item, "")
	item = append(item, "合计")             //单位的位置写入合计
	item = append(item, codeSum["total"]) //总计的和

	for _, code := range shopCodes {
		val, codeSumOk := codeSum[code]

		if !codeSumOk {
			val = 0
		}
		item = append(item, val) //分项的和

		xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d", startCount+i+1), &item)

	}

	xFile.SetSheetRow("Sheet1", "A1", &[]interface{}{pick.TaskName + "拣货单导出"})

	var buffer bytes.Buffer
	_ = xFile.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	data, _ := io.ReadAll(content)
	date := time.Now().Format("20060102")
	c.Writer.Header().Add("Content-Type", "application/octet-stream;charset=utf-8")
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fmt.Sprintf("%s.xlsx", date)+"\"")
	c.Writer.Write(data)
}

// 货品汇总单导出
func GoodsSummaryList(c *gin.Context) {

	var form req.GoodsSummaryListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	form.GoodsTypes = []string{}

	typeStr := ""

	if form.Typ == 1 {
		typeStr = "全部"
		//全部 默认值
	} else if form.Typ == 2 {
		//冻品
		typeStr = "冻品"
		form.GoodsTypes = []string{"冷藏类", "冷冻类"}
	} else {
		//常温
		typeStr = "常温"
		form.GoodsTypes = []string{"常温类"}
	}

	db := global.DB

	err, batch := model.GetBatchByPk(db, form.BatchId)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	codeSum := make(map[string]int, 0)

	err, mp, column, shopCodes, codeSum := dao.GoodsSummaryList(db, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xFile := excelize.NewFile()
	sheet := xFile.NewSheet("sheet1")

	xFile.MergeCell("Sheet1", "A1", GetKey(len(column))+"1")
	for i, cn := range column {
		xFile.SetCellValue("Sheet1", GetKey(i)+"2", cn)
	}

	xFile.SetActiveSheet(sheet)
	//设置指定行高 指定列宽
	xFile.SetRowHeight("Sheet1", 1, 30)
	xFile.SetColWidth("Sheet1", "C", "C", 30)

	startCount := 2

	i := 0

	for _, val := range mp {
		i++
		item := make([]interface{}, 0)
		item = append(item, val["商品名称"])
		item = append(item, val["商品编码"])
		item = append(item, val["货架号"])
		item = append(item, val["规格"])
		item = append(item, val["分类"])
		item = append(item, val["单位"])
		item = append(item, val["总计"])

		for _, code := range shopCodes {
			item = append(item, val[code])
		}

		xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d", startCount+i), &item)
	}

	item := make([]interface{}, 0)
	item = append(item, "")
	item = append(item, "")
	item = append(item, "")
	item = append(item, "")
	item = append(item, "")
	item = append(item, "合计")             //单位的位置写入合计
	item = append(item, codeSum["total"]) //总计的和

	for _, code := range shopCodes {
		val, codeSumOk := codeSum[code]

		if !codeSumOk {
			val = 0
		}
		item = append(item, val) //分项的和

		xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d", startCount+i+1), &item)

	}

	xFile.SetSheetRow("Sheet1", "A1", &[]interface{}{fmt.Sprintf("%s-货品汇总单-%s", batch.BatchName, typeStr)})

	var buffer bytes.Buffer
	_ = xFile.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	data, _ := io.ReadAll(content)
	date := time.Now().Format("20060102")
	c.Writer.Header().Add("Content-Type", "application/octet-stream;charset=utf-8")
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fmt.Sprintf("%s.xlsx", date)+"\"")
	c.Writer.Write(data)
}

// 获取
func GetDeliveryMethod(db *gorm.DB, deliveryMethod int) (err error, methodName string) {
	var dict model.Dict

	err, dict = model.GetDictByTypeCodeAndValue(db, "delivery_method", strconv.Itoa(deliveryMethod))

	if err != nil {
		return
	}

	methodName = dict.Name

	return
}

// 导出店铺地址
func ShopAddress(c *gin.Context) {
	var form req.ShopAddressReq

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, addrListMp := dao.ShopAddress(form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xFile := excelize.NewFile()
	sheet := xFile.NewSheet("sheet1")
	// 设置单元格的值
	xFile.SetCellValue("Sheet1", "A1", "店铺名称")
	xFile.SetCellValue("Sheet1", "B1", "收件人姓名")
	xFile.SetCellValue("Sheet1", "C1", "手机/电话")
	xFile.SetCellValue("Sheet1", "D1", "省")
	xFile.SetCellValue("Sheet1", "E1", "市")
	xFile.SetCellValue("Sheet1", "F1", "区")
	xFile.SetCellValue("Sheet1", "G1", "地址")
	xFile.SetCellValue("Sheet1", "H1", "物品名称")
	xFile.SetCellValue("Sheet1", "I1", "备注")
	xFile.SetActiveSheet(sheet)
	// 指定列宽
	xFile.SetColWidth("Sheet1", "A", "A", 30)
	xFile.SetColWidth("Sheet1", "G", "G", 100)

	startCount := 2

	idx := 0

	for _, val := range addrListMp {
		item := make([]interface{}, 0)
		item = append(item, val["shop_name"])
		item = append(item, val["consignee_name"])
		item = append(item, val["consignee_tel"])
		item = append(item, val["province"])
		item = append(item, val["city"])
		item = append(item, val["district"])
		item = append(item, val["address"])
		item = append(item, "")
		item = append(item, "")

		xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d", startCount+idx), &item)
		idx++
	}

	var buffer bytes.Buffer
	_ = xFile.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	data, _ := io.ReadAll(content)
	date := time.Now().Format("20060102")
	c.Writer.Header().Add("Content-Type", "application/octet-stream;charset=utf-8")
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fmt.Sprintf("%s-.xlsx", date)+"\"")
	c.Writer.Write(data)
}

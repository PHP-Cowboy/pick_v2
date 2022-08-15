package handler

import (
	"bytes"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"pick_v2/forms/req"
	"pick_v2/global"
	"pick_v2/model/batch"
	"pick_v2/utils/ecode"
	"pick_v2/utils/slice"
	"pick_v2/utils/xsq_net"
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
		pick      batch.Pick
		pickGoods []batch.PickGoods
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	result := db.Model(&batch.PickGoods{}).Where("pick_id = ?", form.Id).Find(&pickGoods)

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
		//xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d",startCount+idx), &[]interface{}{1, "11111", "草莓果泥",1,999,9999,0,0,0})
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
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fmt.Sprintf("%s-picking-list.xlsx", date)+"\"")
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

func OutboundBatch(c *gin.Context) {
	var (
		form      req.OutboundBatchFormReq
		shopCodes []string
		batches   batch.Batch
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	result := db.First(&batches, form.Id)

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
			subMp[code] = "0"
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

	shopCodes = slice.UniqueStringSlice(shopCodes)

	column = append(column, shopCodes...)

	//xsq_net.SucJson(c, gin.H{
	//	"mp":     mp,
	//})

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

		//xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d",startCount+idx), &[]interface{}{1, "11111", "草莓果泥",1,999,9999,0,0,0})
		xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d", startCount+i), &item)
	}

	xFile.SetSheetRow("Sheet1", "A1", &[]interface{}{batches.BatchName})

	var buffer bytes.Buffer
	_ = xFile.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	data, _ := ioutil.ReadAll(content)
	date := time.Now().Format("20060102")
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fmt.Sprintf("%s-picking-list.xlsx", date)+"\"")
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

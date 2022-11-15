package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"io"
	"net/http"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"time"
)

// 同步任务
func SyncTask(c *gin.Context) {
	u8 := global.ServerConfig.U8Api
	url := fmt.Sprintf("%s:%d/api/v1/checklist", u8.Url, u8.Port)
	method := "POST"

	client := &http.Client{}
	rq, err := http.NewRequest(method, url, nil)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	rq.Header.Add("x-sign", middlewares.Generate())

	res, err := client.Do(rq)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var result rsp.SyncTaskRsp

	err = json.Unmarshal(body, &result)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	if result.Code != 200 {
		xsq_net.ErrorJSON(c, errors.New(result.Msg))
		return
	}

	var (
		tasks      []model.InvTask                //已同步的盘点任务数据
		tasksMp    = make(map[string]struct{}, 0) //已同步的盘点任务数据map
		task       model.InvTask
		taskRecord []model.InvTaskRecord
	)

	db := global.DB

	dbRes := db.Model(&model.InvTask{}).Find(&tasks)

	if dbRes.Error != nil {
		xsq_net.ErrorJSON(c, dbRes.Error)
		return
	}

	//已同步的盘点任务单号map
	for _, t := range tasks {
		tasksMp[t.OrderNo] = struct{}{}
	}

	now := time.Now()

	var totalBookNum float64

	for _, rd := range result.Data {

		//数据库中已有盘点单，跳过
		_, ok := tasksMp[rd.CcvCode]

		if ok {
			continue
		}

		if task.OrderNo == "" {
			task.OrderNo = rd.CcvCode
			task.TaskDate = (*model.MyTime)(&now)
			task.TaskName = rd.CWhName + now.Format(timeutil.DateNumberFormat)
			task.WarehouseId = 1
			task.Warehouse = rd.CWhName
			task.Remark = rd.CcvMeno
		}

		taskRecord = append(taskRecord, model.InvTaskRecord{
			OrderNo:   rd.CcvCode,
			Sku:       rd.CInvCode,
			GoodsName: rd.CInvName,
			GoodsType: rd.Cate,
			GoodsSpe:  rd.CInvStd,
			GoodsUnit: rd.CComUnitName,
			BookNum:   rd.IcvQuantity,
		})

		totalBookNum += rd.IcvQuantity
	}

	if len(taskRecord) > 0 {
		tx := db.Begin()

		task.BookNum = totalBookNum

		dbRes = tx.Save(&task)

		if dbRes.Error != nil {
			xsq_net.ErrorJSON(c, dbRes.Error)
			return
		}

		dbRes = tx.Save(&taskRecord)

		if dbRes.Error != nil {
			xsq_net.ErrorJSON(c, dbRes.Error)
			return
		}

		tx.Commit()
	}

	xsq_net.Success(c)
}

// 同步商品
func SyncGoods(tx *gorm.DB, selfBuiltId int) error {
	// 自建盘点单，暂不绑定时,同步全部商品
	// 拉取商品数据
	u8 := global.ServerConfig.U8Api
	url := fmt.Sprintf("%s:%d/api/v1/checklist", u8.Url, u8.Port)
	method := "POST"

	client := &http.Client{}
	rq, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}

	rq.Header.Add("x-sign", middlewares.Generate())

	res, err := client.Do(rq)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var resp rsp.SyncTaskRsp

	err = json.Unmarshal(body, &resp)

	if err != nil {
		return err
	}

	if resp.Code != 200 {
		return errors.New(resp.Msg)
	}

	//构造任务商品数据
	list := make([]model.InvTaskRecord, 0, len(resp.Data))

	for _, d := range resp.Data {
		list = append(list, model.InvTaskRecord{
			SelfBuiltId: selfBuiltId,
			OrderNo:     "",
			Sku:         d.CInvCode,
			GoodsName:   d.CInvName,
			GoodsType:   d.Cate,
			GoodsSpe:    d.CInvStd,
			GoodsUnit:   d.CComUnitName,
			BookNum:     0,
		})
	}

	//self_built_id order_no sku inv_type 唯一
	result := tx.Model(&model.InvTaskRecord{}).Save(&list)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// 盘点任务列表
func TaskList(c *gin.Context) {

	var (
		tasks []model.InvTask
	)

	db := global.DB

	result := db.Model(&model.InvTask{}).Where(&model.InvTask{IsBind: 1}).Find(&tasks)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]rsp.TaskList, 0, len(tasks))

	for _, task := range tasks {

		list = append(list, rsp.TaskList{
			OrderNo:  task.OrderNo,
			TaskName: task.TaskName,
		})
	}

	xsq_net.SucJson(c, list)
}

// 结束任务
func ChangeTask(c *gin.Context) {
	var form req.ChangeTaskForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result := global.DB.Model(&model.InvTaskSelfBuilt{}).Where("id = ?", form.Id).Update("status", form.Status)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 导出
func Export(c *gin.Context) {
	var form req.ExportForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		records []model.InvTaskRecord
	)

	db := global.DB

	//只拿 inv_type 为1的，二次盘点的数量根据sku填充
	result := db.Model(&model.InvTaskRecord{}).
		Where(model.InvTaskRecord{SelfBuiltId: form.Id, InvType: 1}).
		Find(&records)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	var (
		invSumMp     = make(map[string]map[int]float64, 0)
		invRecordSum []model.InvRecordSum
	)

	result = db.Model(&model.InvRecordSum{}).Where("self_built_id = ? ", form.Id).Find(&invRecordSum)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	// invSumMp["1010001"][1] = 100.01
	for _, ir := range invRecordSum {

		val, ok := invSumMp[ir.Sku]

		if !ok {
			val = make(map[int]float64, 0)
		}

		val[ir.InvType] = ir.InventoryNum

		invSumMp[ir.Sku] = val
	}

	xFile := excelize.NewFile()
	sheet := xFile.NewSheet("sheet1")

	// 设置单元格的值
	xFile.MergeCell("Sheet1", "A1", "H1")
	xFile.SetCellValue("Sheet1", "A2", "Sku")
	xFile.SetCellValue("Sheet1", "B2", "商品名称")
	xFile.SetCellValue("Sheet1", "C2", "商品分类")
	xFile.SetCellValue("Sheet1", "D2", "账面数量")
	xFile.SetCellValue("Sheet1", "E2", "首次盘点数量")
	xFile.SetCellValue("Sheet1", "F2", "首次盈亏数量")
	xFile.SetCellValue("Sheet1", "G2", "二次盘点数量")
	xFile.SetCellValue("Sheet1", "H2", "二次盈亏数量")

	xFile.SetActiveSheet(sheet)
	//设置指定行高 指定列宽
	xFile.SetRowHeight("Sheet1", 1, 30)
	xFile.SetColWidth("Sheet1", "C", "C", 30)

	startCount := 3
	for idx, val := range records {

		var firstInvNum, secondInvNum float64

		invNum, invSumOk := invSumMp[val.Sku]

		if invSumOk {
			//首次盘点数量
			fNum, fOk := invNum[1]
			if fOk {
				firstInvNum = fNum
			}
			//二次盘点数量
			sNum, sOk := invNum[2]
			if sOk {
				secondInvNum = sNum
			} else {
				//导出时二次盘点如果盘了，就取二次盘点数量，如果没盘就取一次盘点数量
				secondInvNum = firstInvNum
			}
		}

		item := make([]interface{}, 0)
		item = append(item, val.Sku)
		item = append(item, val.GoodsName)
		item = append(item, val.GoodsType)
		item = append(item, val.BookNum)
		item = append(item, firstInvNum)
		item = append(item, firstInvNum-val.BookNum)
		item = append(item, secondInvNum)
		item = append(item, secondInvNum-val.BookNum)

		xFile.SetSheetRow("Sheet1", fmt.Sprintf("A%d", startCount+idx), &item)
	}

	xFile.SetSheetRow("Sheet1", "A1", &[]interface{}{fmt.Sprintf("盘点商品列表")})

	var buffer bytes.Buffer
	_ = xFile.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	data, _ := io.ReadAll(content)
	date := time.Now().Format(timeutil.DateNumberFormat)
	c.Writer.Header().Add("Content-Type", "application/octet-stream;charset=utf-8")
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fmt.Sprintf("%s-.xlsx", date)+"\"")
	c.Writer.Write(data)

	xsq_net.SucJson(c, "")
}

// 任务商品数据列表
func TaskRecordList(c *gin.Context) {
	var form req.TaskRecordListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	localDb := db.Model(&model.InvOrderSkuSum{}).
		Where(model.InvOrderSkuSum{SelfBuiltId: form.SelfBuiltId, InvType: form.InvType, GoodsType: form.GoodsType})

	if form.GoodsName != "" {
		localDb.Where("goods_name like ?", "%"+form.GoodsName+"%")
	}

	if form.IsNeed {
		localDb.Where("book_num != inventory_num")
	}

	var (
		invOrderSkuSum []model.InvOrderSkuSum
		res            rsp.RecordListRsp
	)

	result := localDb.Find(&invOrderSkuSum)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	if form.SortField != "" && form.SortRule != "" {
		localDb.Order(fmt.Sprintf("%s %s", form.SortField, form.SortRule))
	}

	result = localDb.
		Scopes(model.Paginate(form.Page, form.Size)).
		Find(&invOrderSkuSum)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]*rsp.RecordList, 0, len(invOrderSkuSum))

	for _, record := range invOrderSkuSum {

		list = append(list, &rsp.RecordList{
			SelfBuiltId:   record.SelfBuiltId,
			Sku:           record.Sku,
			GoodsName:     record.GoodsName,
			GoodsType:     record.GoodsType,
			GoodsSpe:      record.GoodsSpe,
			BookNum:       record.BookNum,
			InventoryNum:  record.InventoryNum,
			ProfitLossNum: record.InventoryNum - record.BookNum,
			InvType:       record.InvType,
		})
	}

	res.Data = list

	xsq_net.SucJson(c, res)
}

// 任务记录分类列表
func TypeList(c *gin.Context) {
	var form req.TypeListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var goodsType []string
	//todo 新数据有了之后可以优化到其他表去
	result := global.DB.
		Model(&model.InvTaskRecord{}).
		Distinct("goods_type").
		Where("order_no = ?", form.OrderNo).
		Find(&goodsType)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.SucJson(c, goodsType)
}

// 盘库记录
func InventoryRecordList(c *gin.Context) {
	var form req.InventoryRecordListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		records []model.InventoryRecord
		res     rsp.InventoryRecordListRsp
	)

	localDb := global.DB.Model(&model.InventoryRecord{}).
		Where(&model.InventoryRecord{SelfBuiltId: form.SelfBuiltId, Sku: form.Sku, InvType: form.InvType}).
		Where("is_delete = 1")

	result := localDb.Find(&records)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]*rsp.InventoryRecord, 0, len(records))

	//盘库记录总数不提供，从列表传递过来
	for _, record := range records {
		list = append(list, &rsp.InventoryRecord{
			Id:           record.Id,
			SelfBuiltId:  record.SelfBuiltId,
			Sku:          record.Sku,
			CreateTime:   record.CreateTime.Format(timeutil.TimeFormat),
			InventoryNum: record.InventoryNum,
			UserName:     record.UserName,
			GoodsUnit:    record.GoodsUnit,
		})
	}

	res.Total = result.RowsAffected
	res.Data = list

	xsq_net.SucJson(c, res)
}

// 盘库记录删除
func InventoryRecordDelete(c *gin.Context) {
	var form req.InventoryRecordDeleteForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result := global.DB.Model(&model.InventoryRecord{}).
		Where("id = ?", form.Id).
		Update("is_delete", 2)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 已盘商品件数
func InvCount(c *gin.Context) {
	var form req.CountForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, errors.New("获取上下文用户数据失败"))
		return
	}

	var count int64

	result := global.DB.
		Model(&model.InventoryRecord{}).
		Distinct("sku").
		Where("self_built_id = ? and inv_type = ? and user_name = ? and is_delete = 1", form.SelfBuiltId, form.InvType, userInfo.Name).
		Count(&count)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.SucJson(c, gin.H{"count": count})
}

// 未盘商品数量
func NotInvCount(c *gin.Context) {
	var form req.NotInvCountForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, errors.New("获取上下文用户数据失败"))
		return
	}

	var (
		cTaskRec    int64
		cUserInvRec int64
	)

	db := global.DB

	result := db.Model(&model.InvTaskRecord{}).
		Where("self_built_id = ? and inv_type = ?", form.SelfBuiltId, 2).
		Count(&cTaskRec)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.InventoryRecord{}).
		Distinct("sku").
		Where("self_built_id = ? and inv_type = ? and is_delete = 1", form.SelfBuiltId, 2).
		Count(&cUserInvRec)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.SucJson(c, gin.H{"count": cTaskRec - cUserInvRec})
}

// 未盘商品列表
func UserNotInventoryRecordList(c *gin.Context) {
	var form req.UserNotInventoryRecordListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	var (
		invTaskRecords []model.InvTaskRecord
		invRecords     []model.InventoryRecord
		secInv         = 2 //复盘
		resp           = make(map[string][]rsp.UserNotInventoryRecord, 0)
		userInvMp      = make(map[string]struct{}, 0)
	)

	result := db.Model(&model.InvTaskRecord{}).
		Where(&model.InvTaskRecord{SelfBuiltId: form.SelfBuiltId, Sku: form.Sku, InvType: secInv}).
		Order("goods_type desc").
		Order("sku asc").
		Find(&invTaskRecords)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.InventoryRecord{}).
		Where(&model.InventoryRecord{SelfBuiltId: form.SelfBuiltId, Sku: form.Sku, InvType: secInv}).
		Find(&invRecords)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, record := range invRecords {
		userInvMp[record.Sku] = struct{}{}
	}

	for _, tr := range invTaskRecords {
		//用户已复盘商品map，如果sku存在，则跳过
		_, ok := userInvMp[tr.Sku]

		if ok {
			continue
		}

		resp[tr.GoodsType] = append(resp[tr.GoodsType], rsp.UserNotInventoryRecord{
			Sku:       tr.Sku,
			GoodsName: tr.GoodsName,
			GoodsSpe:  tr.GoodsSpe,
		})
	}

	xsq_net.SucJson(c, resp)
}

// 用户已盘商品列表
func UserInventoryRecordList(c *gin.Context) {
	var form req.UserInventoryRecordListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, errors.New("获取上下文用户数据失败"))
		return
	}

	var (
		records []model.InventoryRecord
		res     rsp.UserInventoryRecordListRsp
	)

	db := global.DB

	result := db.Model(&model.InventoryRecord{}).
		Where(&model.InventoryRecord{Sku: form.Sku}).
		Where("self_built_id = ? and inv_type = ? and user_name = ? and is_delete = 1", form.SelfBuiltId, form.InvType, userInfo.Name).
		Find(&records)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	//分页
	result = db.Model(&model.InventoryRecord{}).
		Scopes(model.Paginate(form.Page, form.Size)).
		Where(&model.InventoryRecord{Sku: form.Sku}).
		Where("self_built_id = ? and inv_type = ? and user_name = ? and is_delete = 1", form.SelfBuiltId, form.InvType, userInfo.Name).
		Order("id desc").
		Find(&records)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//系统数量
	var (
		skuSlice     []string
		skuInvNumSum []rsp.SkuInvNumSum
		systemNumMp  = make(map[string]float64, 0)
	)

	for _, rs := range records {
		skuSlice = append(skuSlice, rs.Sku)
	}

	result = db.Model(&model.InventoryRecord{}).
		Select("sum(inventory_num) as sum,sku").
		Where("self_built_id = ? and sku in (?) and inv_type = ?  and is_delete = 1", form.SelfBuiltId, skuSlice, form.InvType).
		Group("sku").
		Find(&skuInvNumSum)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, inv := range skuInvNumSum {
		systemNumMp[inv.Sku] = inv.Sum
	}

	list := make([]*rsp.UserInventoryRecord, 0, len(records))

	for _, record := range records {

		sysNum, sysOk := systemNumMp[record.Sku]

		if !sysOk {
			sysNum = 0
		}

		list = append(list, &rsp.UserInventoryRecord{
			Id:           record.Id,
			SelfBuiltId:  record.SelfBuiltId,
			Sku:          record.Sku,
			UserName:     record.UserName,
			GoodsName:    record.GoodsName,
			GoodsSpe:     record.GoodsSpe,
			InventoryNum: record.InventoryNum,
			SystemNum:    sysNum,
		})
	}

	res.Data = list

	xsq_net.SucJson(c, res)
}

// 修改已盘商品数据
func UpdateInventoryRecord(c *gin.Context) {
	var form req.UpdateInventoryRecordForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		record        model.InventoryRecord
		selfBuiltTask model.InvTaskSelfBuilt
	)

	db := global.DB

	result := db.Model(&model.InventoryRecord{}).First(&record, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.InvTaskSelfBuilt{}).Where("id = ?", record.SelfBuiltId).First(&selfBuiltTask)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	// 验证是否已结束
	if selfBuiltTask.Status == 2 {
		xsq_net.ErrorJSON(c, errors.New("盘点任务已结束"))
		return
	}

	result = db.Model(&model.InventoryRecord{}).
		Where("id = ?", form.Id).
		Update("inventory_num", form.InventoryNum)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 批量盘点
func BatchCreate(c *gin.Context) {
	var form req.BatchCreateForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, errors.New("获取上下文用户数据失败"))
		return
	}

	var (
		selfBuiltTask model.InvTaskSelfBuilt
		taskRecords   []model.InvTaskRecord
	)

	db := global.DB

	result := db.Model(&model.InvTaskSelfBuilt{}).Where("id = ?", form.SelfBuiltId).First(&selfBuiltTask)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	// 验证是否已结束
	if selfBuiltTask.Status == 2 {
		xsq_net.ErrorJSON(c, errors.New("盘点任务已结束"))
		return
	}

	result = db.Model(&model.InvTaskRecord{}).
		Where("order_no = ?", selfBuiltTask.OrderNo).
		Find(&taskRecords)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//sku map 验证是否在盘点列表中
	mp := make(map[string]model.InvTaskRecord, 0)

	for _, record := range taskRecords {
		mp[record.Sku] = record
	}

	records := make([]*model.InventoryRecord, 0, len(form.Records))

	for _, fm := range form.Records {
		mpSku, skuOk := mp[fm.Sku]

		if !skuOk {
			xsq_net.ErrorJSON(c, errors.New("盘点单中没有对应的sku:"+fm.Sku))
			return
		}

		if fm.InventoryNum < 0 {
			xsq_net.ErrorJSON(c, errors.New("盘点数量不能为负,sku:"+fm.Sku))
			return
		}

		records = append(records, &model.InventoryRecord{
			SelfBuiltId:  form.SelfBuiltId,
			Sku:          fm.Sku,
			UserName:     userInfo.Name,
			GoodsName:    mpSku.GoodsName,
			GoodsSpe:     mpSku.GoodsSpe,
			GoodsUnit:    mpSku.GoodsUnit,
			InventoryNum: fm.InventoryNum,
			InvType:      form.InvType,
		})

	}

	result = db.Save(&records)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)

}

// 自建盘点任务
func SelfBuiltTask(c *gin.Context) {
	var form req.InvTaskForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var task model.InvTaskSelfBuilt

	task.OrderNo = form.OrderNo
	task.TaskName = form.TaskName

	tx := global.DB.Begin()

	result := tx.Save(&task)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if form.OrderNo == "" {
		//暂不绑定盘点单
		//同步u8商品  self_built_id 为 自建盘点单id order_no为 空字符串
		if err := SyncGoods(tx, task.Id); err != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, errors.New("商品数据同步失败"))
			return
		}
	} else {
		//直接绑定盘点单
		//将 self_built_id 为 0 order_no 为u8的盘点单号 的数据 的 self_built_id 更新为当前自建盘点单的id
		result = tx.Model(&model.InvTaskRecord{}).
			Where("self_built_id = 0 and order_no = ?", form.OrderNo).
			Update("self_built_id", task.Id)

		if result.Error != nil {
			tx.Rollback()
			xsq_net.ErrorJSON(c, result.Error)
			return
		}
	}

	tx.Commit()

	xsq_net.Success(c)
}

// 修改自建盘点任务
func ChangeSelfBuiltTask(c *gin.Context) {
	var form req.ChangeSelfBuiltTaskForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		task          model.InvTask
		selfBuiltTask model.InvTaskSelfBuilt
	)

	db := global.DB

	result := db.Model(&model.InvTask{}).Where("order_no = ?", form.OrderNo).First(&task)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//验证盘点单是否已经被绑定
	if task.IsBind != 1 {
		xsq_net.ErrorJSON(c, ecode.InvTaskAlreadyBind)
		return
	}

	result = db.First(&selfBuiltTask, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if selfBuiltTask.OrderNo != "" {
		xsq_net.ErrorJSON(c, errors.New("已绑定的暂不允许更换"))
		return
	}

	tx := db.Begin()

	//自建盘点任务绑定u8盘点单
	result = tx.Model(&model.InvTaskSelfBuilt{}).
		Where("id = ?", form.Id).
		Update("order_no", form.OrderNo)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//更新u8盘点单绑定状态
	result = tx.Model(&model.InvTask{}).
		Where("order_no = ?", form.OrderNo).
		Update("is_bind", 2)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//删除暂不绑定时生成的数据
	result = tx.Delete(&model.InvTaskRecord{}, "self_built_id = ?", form.Id)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//更新盘点任务商品记录表
	result = tx.Model(&model.InvTaskRecord{}).
		Where("self_built_id = 0 and order_no = ?", form.OrderNo).
		Update("self_built_id", form.Id)

	if result.Error != nil {
		tx.Rollback()
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	tx.Commit()

	xsq_net.Success(c)
}

// 自建盘点任务列表
func SelfBuiltTaskList(c *gin.Context) {
	var form req.SelfBuiltTaskListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	var (
		tasks []model.SelfBuiltJoinTask
		res   rsp.SelfBuiltTaskRsp
	)

	localDb := db.Table("t_inv_task_self_built sbt").
		Select("sbt.*,it.task_date,it.warehouse,it.book_num,it.remark").
		Joins("left join t_inv_task it on sbt.order_no = it.order_no")

	result := localDb.Find(&tasks)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	result = localDb.Scopes(model.Paginate(form.Page, form.Size)).Find(&tasks)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	taskIds := make([]int, 0, len(tasks))

	for _, ts := range tasks {
		taskIds = append(taskIds, ts.Id)
	}

	var (
		invRecordSums []model.InvRecordSum
		builtSkuMp    = make(map[int]map[string]float64, 0)
		builtNumMp    = make(map[int]float64, 0)
	)

	result = db.Model(&model.InvRecordSum{}).
		Where("self_built_id in (?)", taskIds).
		Order("sku asc, inv_type ASC").
		Find(&invRecordSums)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//sku去重，map覆盖
	for _, sum := range invRecordSums {
		skuVal, skuOk := builtSkuMp[sum.SelfBuiltId]

		if !skuOk {
			skuVal = make(map[string]float64, 0)
		}

		skuVal[sum.Sku] = sum.InventoryNum

		builtSkuMp[sum.SelfBuiltId] = skuVal
	}

	// builtSkuMp[self_built_id][sku]inventory_num
	for i, m := range builtSkuMp {
		for _, n := range m {
			builtNumMp[i] += n
		}
	}

	list := make([]*rsp.SelfBuiltTask, 0, len(tasks))

	for _, task := range tasks {

		invSum, numMpOk := builtNumMp[task.Id]

		if !numMpOk {
			invSum = 0.00
		}

		list = append(list, &rsp.SelfBuiltTask{
			Id:            task.Id,
			CreateTime:    timeutil.FormatToDateTime(task.CreateTime),
			OrderNo:       task.OrderNo,
			TaskName:      task.TaskName,
			Warehouse:     task.Warehouse,
			TaskDate:      timeutil.FormatToDate(task.TaskDate),
			Status:        task.Status,
			BookNum:       task.BookNum,
			InventoryNum:  invSum,
			ProfitLossNum: invSum - task.BookNum,
			Remark:        task.Remark,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// 设置复盘
func SetSecondInventory(c *gin.Context) {
	var form req.SetSecondInventoryForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	db := global.DB

	var (
		invTaskSelfBuilt model.InvTaskSelfBuilt
		taskRecords      []model.InvTaskRecord
	)

	result := db.First(&invTaskSelfBuilt, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if invTaskSelfBuilt.Status != 1 {
		xsq_net.ErrorJSON(c, ecode.InvTaskNotGoing)
		return
	}

	if invTaskSelfBuilt.OrderNo == "" {
		xsq_net.ErrorJSON(c, ecode.UnBindTaskNotAllow)
		return
	}

	result = db.Model(&model.InvTaskRecord{}).
		Where("self_built_id = ? and order_no = ? and sku in (?)", form.Id, invTaskSelfBuilt.OrderNo, form.Sku).
		Find(&taskRecords)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	taskRecordMp := make(map[string]model.InvTaskRecord, 0)

	for _, tr := range taskRecords {
		taskRecordMp[tr.Sku] = tr
	}

	//生成复盘记录数据
	list := make([]model.InvTaskRecord, 0, len(form.Sku))

	for _, s := range form.Sku {
		skuVal, ok := taskRecordMp[s]

		if !ok {
			continue
		}

		list = append(list, model.InvTaskRecord{
			SelfBuiltId: form.Id,
			OrderNo:     invTaskSelfBuilt.OrderNo,
			Sku:         s,
			InvType:     2,
			GoodsName:   skuVal.GoodsName,
			GoodsType:   skuVal.GoodsType,
			GoodsSpe:    skuVal.GoodsSpe,
			GoodsUnit:   skuVal.GoodsUnit,
			BookNum:     skuVal.BookNum,
		})
	}

	result = db.Model(&model.InvTaskRecord{}).Save(&list)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 重新盘点
func InvAgain(c *gin.Context) {
	var form req.InvAgainForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	db := global.DB

	result := db.Model(&model.InventoryRecord{}).
		Where("self_built_id = ? and sku in (?) and inv_type = 2", form.SelfBuiltId, form.Sku).
		Updates(map[string]interface{}{
			"is_delete": 2,
		})

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

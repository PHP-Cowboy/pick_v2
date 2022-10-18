package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/xuri/excelize/v2"
	"io"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"time"
)

// 同步任务
func SyncTask(c *gin.Context) {
	// todo u8接口暂无
	xsq_net.Success(c)
}

// 盘点任务列表
func TaskList(c *gin.Context) {
	var form req.TaskListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		tasks []model.InvTask
		res   rsp.TaskListRsp
	)

	db := global.DB

	result := db.Model(&model.InvTask{}).Where(&model.InvTask{Status: form.Status}).Find(&tasks)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	result = db.Model(&model.InvTask{}).Where(&model.InvTask{Status: form.Status}).Scopes(model.Paginate(form.Page, form.Size)).Find(&tasks)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]*rsp.TaskList, 0, len(tasks))

	for _, task := range tasks {
		list = append(list, &rsp.TaskList{
			OrderNo:       task.OrderNo,
			TaskName:      task.TaskName,
			TaskDate:      task.TaskDate,
			WarehouseId:   task.WarehouseId,
			WarehouseName: task.Warehouse,
			BookNum:       task.BookNum,
			InventoryNum:  task.InventoryNum,
			ProfitLossNum: task.ProfitLossNum,
			Remark:        task.Remark,
			Status:        task.Status,
		})
	}

	res.Data = list

	xsq_net.SucJson(c, res)
}

// 结束任务
func ChangeTask(c *gin.Context) {
	var form req.ChangeTaskForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result := global.DB.Model(&model.InvTask{}).Where("order_no = ?", form.OrderNo).Update("status", form.Status)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 导出
func Export(c *gin.Context) {
	var form req.ExportForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		records []model.InvTaskRecord
	)

	db := global.DB

	result := db.Model(&model.InvTaskRecord{}).
		Where(model.InvTaskRecord{OrderNo: form.OrderNo}).
		Find(&records)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xFile := excelize.NewFile()
	sheet := xFile.NewSheet("sheet1")

	// 设置单元格的值
	xFile.MergeCell("Sheet1", "A1", "F1")
	xFile.SetCellValue("Sheet1", "A2", "Sku")
	xFile.SetCellValue("Sheet1", "B2", "商品名称")
	xFile.SetCellValue("Sheet1", "C2", "商品分类")
	xFile.SetCellValue("Sheet1", "D2", "账面数量")
	xFile.SetCellValue("Sheet1", "E2", "盘点数量")
	xFile.SetCellValue("Sheet1", "F2", "盈亏数量")

	xFile.SetActiveSheet(sheet)
	//设置指定行高 指定列宽
	xFile.SetRowHeight("Sheet1", 1, 30)
	xFile.SetColWidth("Sheet1", "C", "C", 30)

	startCount := 3
	for idx, val := range records {
		item := make([]interface{}, 0)
		item = append(item, val.Sku)
		item = append(item, val.GoodsName)
		item = append(item, val.GoodsType)
		item = append(item, val.BookNum)
		item = append(item, val.InventoryNum)
		item = append(item, val.ProfitLossNum)

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

// 任务记录列表
func TaskRecordList(c *gin.Context) {
	var form req.TaskRecordListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	localDb := db.Model(&model.InvTaskRecord{}).Where(model.InvTaskRecord{OrderNo: form.OrderNo, GoodsType: form.GoodsType})

	if form.GoodsName != "" {
		localDb.Where("goods_name like %?%", form.GoodsName)
	}

	if form.IsNeed {
		localDb.Where("profit_loss_num > 0")
	}

	var (
		records []model.InvTaskRecord
		res     rsp.RecordListRsp
	)

	result := localDb.Find(&records)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	result = localDb.Scopes(model.Paginate(form.Page, form.Size)).Find(&records)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]*rsp.RecordList, 0, len(records))

	for _, record := range records {
		list = append(list, &rsp.RecordList{
			Sku:           record.Sku,
			GoodsName:     record.GoodsName,
			GoodsType:     record.GoodsType,
			BookNum:       record.BookNum,
			InventoryNum:  record.InventoryNum,
			ProfitLossNum: record.ProfitLossNum,
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
		Where(&model.InventoryRecord{Sku: form.Sku, OrderNo: form.OrderNo}).
		Where("is_delete = 1")

	if form.SortRule != "" {
		localDb.Order(fmt.Sprintf("%s %s", form.SortField, form.SortRule))
	}

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

	userName, ok := c.Get("user_name")

	if !ok {
		xsq_net.ErrorJSON(c, errors.New("用户名称不存在"))
		return
	}

	var count int64

	result := global.DB.
		Model(&model.InventoryRecord{}).
		Distinct("sku").
		Group("user_name").
		Where("user_name = ?", userName).
		Count(&count)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.SucJson(c, gin.H{"count": count})
}

// 用户已盘商品列表
func UserInventoryRecordList(c *gin.Context) {
	var form req.UserInventoryRecordListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	userName, ok := c.Get("user_name")

	if !ok {
		xsq_net.ErrorJSON(c, errors.New("用户名称不存在"))
		return
	}

	var (
		records []model.InventoryRecord
		res     rsp.UserInventoryRecordListRsp
	)

	result := global.DB.
		Model(&model.InventoryRecord{}).
		Where(&model.InventoryRecord{Sku: form.Sku}).
		Where("order_no = ? and user_name = ? and is_delete = 1", form.OrderNo, userName).
		Find(&records)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	result = global.DB.
		Model(&model.InventoryRecord{}).
		Scopes(model.Paginate(form.Page, form.Size)).
		Where(&model.InventoryRecord{Sku: form.Sku}).
		Where("order_no = ? and user_name = ? and is_delete = 1", form.OrderNo, userName).
		Find(&records)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	list := make([]*rsp.UserInventoryRecord, 0, len(records))

	for _, record := range records {
		list = append(list, &rsp.UserInventoryRecord{
			Id:           record.Id,
			GoodsName:    record.GoodsName,
			GoodsSpe:     record.GoodsSpe,
			InventoryNum: record.InventoryNum,
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
		record model.InventoryRecord
		task   model.InvTask
	)

	db := global.DB

	result := db.Model(&model.InventoryRecord{}).First(&record)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.Model(&model.InvTask{}).Where("order_no = ?", record.OrderNo).First(&task)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	// 验证是否已结束
	if task.Status == 2 {
		xsq_net.ErrorJSON(c, errors.New("盘点任务已结束"))
		return
	}

	record.InventoryNum = form.InventoryNum

	result = db.Save(&record)

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

	userName, ok := c.Get("user_name")

	if !ok {
		xsq_net.ErrorJSON(c, errors.New("用户名不存在"))
		return
	}

	var taskRecords []model.InvTaskRecord

	db := global.DB

	result := db.Model(&model.InvTaskRecord{}).
		Where("order_no = ?", form.OrderNo).
		Find(&taskRecords)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	//sku map
	mp := make(map[string]model.InvTaskRecord, 0)

	for _, record := range taskRecords {
		mp[record.Sku] = record
	}

	records := make([]*model.InventoryRecord, 0, len(form.Records))

	for _, fm := range form.Records {
		mpSku, skuOk := mp[fm.Sku]

		if !skuOk {
			xsq_net.ErrorJSON(c, errors.New("盘点单中没有对应的sku"))
			return
		}

		records = append(records, &model.InventoryRecord{
			OrderNo:      form.OrderNo,
			Sku:          fm.Sku,
			GoodsName:    mpSku.GoodsName,
			GoodsSpe:     mpSku.GoodsSpe,
			GoodsUnit:    mpSku.GoodsUnit,
			InventoryNum: fm.InventoryNum,
			UserName:     userName.(string),
		})
	}

	result = db.Save(&records)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)

}

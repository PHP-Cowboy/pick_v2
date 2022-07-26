package handler

import (
	"github.com/gin-gonic/gin"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/model/other"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
)

//获取仓库列表
func GetWarehouseList(c *gin.Context) {
	var form req.GetWarehouseListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		warehouses []other.Warehouse
		res        rsp.GetWarehouseListRsp
	)

	db := global.DB

	result := db.Find(&warehouses)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	db.Scopes(model.Paginate(form.Page, form.Size)).Find(&warehouses)

	for _, w := range warehouses {
		res.List = append(res.List, &rsp.WarehouseList{
			Id:            w.Id,
			WarehouseName: w.WarehouseName,
			Abbreviation:  w.Abbreviation,
		})
	}

	xsq_net.SucJson(c, res)
}

//新增仓库
func CreateWarehouse(c *gin.Context) {
	var form req.CreateWarehouseForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		warehouse other.Warehouse
	)

	warehouse.WarehouseName = form.WarehouseName
	warehouse.Abbreviation = form.Abbreviation

	db := global.DB
	result := db.Save(&warehouse)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	warehouseRsp := rsp.CreateWarehouseRsp{
		Id:            warehouse.Id,
		WarehouseName: warehouse.WarehouseName,
		Abbreviation:  warehouse.Abbreviation,
	}

	xsq_net.SucJson(c, warehouseRsp)
}

//修改仓库
func ChangeWarehouse(c *gin.Context) {
	xsq_net.Success(c)
}

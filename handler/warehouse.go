package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
)

// 获取仓库列表
func GetWarehouseList(c *gin.Context) {

	var (
		warehouses []model.Warehouse
		res        []*rsp.WarehouseList
	)

	db := global.DB

	result := db.Find(&warehouses)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, w := range warehouses {
		res = append(res, &rsp.WarehouseList{
			Id:            w.Id,
			WarehouseName: w.WarehouseName,
			Abbreviation:  w.Abbreviation,
		})
	}

	xsq_net.SucJson(c, res)
}

// 新增仓库
func CreateWarehouse(c *gin.Context) {
	var form req.CreateWarehouseForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		warehouse model.Warehouse
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

// 修改仓库
func ChangeWarehouse(c *gin.Context) {
	xsq_net.Success(c)
}

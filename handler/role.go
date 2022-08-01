package handler

import (
	"github.com/gin-gonic/gin"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"time"
)

//创建角色
func CreateRole(ctx *gin.Context) {
	var form req.CreateRoleForm

	err := ctx.ShouldBind(&form)

	if err != nil {
		xsq_net.ErrorJSON(ctx, ecode.ParamInvalid)
		return
	}

	result := global.DB.Create(&model.Role{Name: form.Name})

	if result.Error != nil || result.RowsAffected == 0 {
		xsq_net.ErrorJSON(ctx, result.Error)
		return
	}

	xsq_net.Success(ctx)
}

//修改角色
func ChangeRole(ctx *gin.Context) {
	var form req.ChangeRoleForm

	err := ctx.ShouldBind(&form)

	if err != nil {
		xsq_net.ErrorJSON(ctx, ecode.ParamInvalid)
		return
	}

	db := global.DB

	var role model.Role

	result := db.First(&role, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(ctx, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		xsq_net.ErrorJSON(ctx, ecode.RoleNotFound)
		return
	}

	deleteTime, err := time.ParseInLocation(timeutil.TimeFormat, timeutil.GetDateTime(), time.Local)
	if err != nil {
		xsq_net.ErrorJSON(ctx, ecode.DataTransformationError)
		return
	}

	if form.IsDelete {
		result = db.Model(&role).Updates(model.Role{
			Base: model.Base{
				DeleteTime: deleteTime,
			},
		})
	}

	xsq_net.Success(ctx)
}

//角色列表
func GetRoleList(ctx *gin.Context) {
	var form req.GetRoleListForm

	err := ctx.ShouldBind(&form)

	if err != nil {
		xsq_net.ErrorJSON(ctx, ecode.ParamInvalid)
		return
	}

	db := global.DB

	var (
		roles []model.Role
		res   rsp.GetRoleList
	)

	result := db.Where("delete_time is null").Find(&roles)

	if result.Error != nil {
		xsq_net.ErrorJSON(ctx, result.Error)
		return
	}

	res.Total = result.RowsAffected

	db.Where("delete_time is null").Scopes(model.Paginate(form.Page, form.Size)).Find(&roles)

	for _, role := range roles {
		res.Data = append(res.Data, &rsp.Role{
			Id:         role.Id,
			CreateTime: role.CreateTime.Format(timeutil.TimeFormat),
			Name:       role.Name,
		})
	}

	xsq_net.SucJson(ctx, res)
}

//批量删除角色
func BatchDeleteRole(c *gin.Context) {
	var form req.BatchDeleteRoleForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result := global.DB.Model(model.Role{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"delete_time": time.Now().Format(timeutil.TimeFormat)})

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

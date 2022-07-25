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

	deleteTime, err := time.ParseInLocation(timeutil.DateTime, timeutil.GetDateTime(), time.Local)
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

func GetRoleList(ctx *gin.Context) {
	var form req.GetRoleListForm

	err := ctx.ShouldBind(&form)

	if err != nil {
		xsq_net.ErrorJSON(ctx, ecode.ParamInvalid)
		return
	}

	db := global.DB

	result := db.Find(&model.Role{})

	if result.Error != nil {
		xsq_net.ErrorJSON(ctx, result.Error)
		return
	}

	var (
		roles []model.Role
		res   rsp.GetRoleList
	)

	db.Scopes(model.Paginate(form.Paging.Page, form.Paging.Size)).Find(&roles)

	res.Total = result.RowsAffected

	for _, role := range roles {
		res.Data = append(res.Data, &rsp.Role{
			Id:         role.Id,
			CreateTime: role.CreateTime.Format(timeutil.DateTime),
			Name:       role.Name,
		})
	}

	xsq_net.Success(ctx)
}

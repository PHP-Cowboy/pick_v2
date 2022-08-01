package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func RoleRoute(g *gin.RouterGroup) {
	roleGroup := g.Group("/role", middlewares.JWTAuth(), middlewares.IsSuperAdminAuth())
	{
		//新增角色
		roleGroup.POST("/create", handler.CreateRole)
		//修改角色
		roleGroup.POST("/change", handler.ChangeRole)
		//角色列表
		roleGroup.GET("/list", handler.GetRoleList)
		//批量删除角色
		roleGroup.POST("/batch_delete", handler.BatchDeleteRole)
	}
}

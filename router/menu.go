package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func MenuRoute(g *gin.RouterGroup) {

	menuGroup := g.Group("/menu", middlewares.JWTAuth(), middlewares.IsSuperAdminAuth())
	{
		//菜单列表
		menuGroup.GET("/list", handler.GetMenuList)
		//新增菜单
		menuGroup.POST("/create", handler.CreateMenu)
		//修改菜单
		menuGroup.POST("/change", handler.ChangeMenu)
		//批量删除菜单
		menuGroup.POST("/batch_delete", handler.BatchDeleteMenu)

	}
}

package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func WarehouseRoute(g *gin.RouterGroup) {
	group := g.Group("/warehouse")
	{
		//仓库列表
		group.GET("/list", handler.GetWarehouseList)
	}

	roleGroup := g.Group("/warehouse", middlewares.JWTAuth(), middlewares.IsSuperAdminAuth())
	{
		//新增仓库
		roleGroup.POST("/create", handler.CreateWarehouse)
		//修改仓库
		roleGroup.POST("/change", handler.ChangeWarehouse)

	}
}

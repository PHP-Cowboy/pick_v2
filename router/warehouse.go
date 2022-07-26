package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func WarehouseRoute(g *gin.RouterGroup) {
	roleGroup := g.Group("/warehouse", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//新增仓库
		roleGroup.POST("/create", handler.CreateWarehouse)
		//修改仓库
		roleGroup.POST("/change", handler.ChangeWarehouse)
		//仓库列表
		roleGroup.GET("/list", handler.GetWarehouseList)
	}
}

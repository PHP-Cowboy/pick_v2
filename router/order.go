package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func OrderRoute(g *gin.RouterGroup) {

	menuGroup := g.Group("/order", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//拣货单列表
		menuGroup.GET("/list", handler.PickOrderList)
		//拣货单明细
		menuGroup.GET("/pickOrderDetail", handler.GetPickOrderDetail)

	}
}

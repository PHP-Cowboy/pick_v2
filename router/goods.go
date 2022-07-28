package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func GoodsRoute(g *gin.RouterGroup) {
	roleGroup := g.Group("/goods", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//列表
		roleGroup.GET("/list", handler.GetGoodsList)
		//明细
		roleGroup.GET("/detail", handler.GetOrderDetail)
	}
}

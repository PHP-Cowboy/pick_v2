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
		roleGroup.POST("/list", handler.GetGoodsList)
	}
}

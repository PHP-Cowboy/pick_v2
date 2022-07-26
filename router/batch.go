package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func BatchRoute(g *gin.RouterGroup) {
	roleGroup := g.Group("/batch", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//列表
		roleGroup.POST("/create", handler.CreateBatch)
	}
}

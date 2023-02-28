package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
)

func TestRoute(g *gin.RouterGroup) {
	taskGroup := g.Group("/test")
	{
		//列表
		taskGroup.GET("/sAdd", handler.SAdd)
		taskGroup.POST("/orderRepair", handler.OrderRepair)
	}
}

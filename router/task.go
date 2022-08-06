package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func TaskRoute(g *gin.RouterGroup) {
	taskGroup := g.Group("/task", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//列表
		taskGroup.GET("/list", handler.PickList)
		//详情
		taskGroup.GET("/pick_detail", handler.GetPickDetail)
		//拣货置顶
		taskGroup.POST("/pick_topping", handler.PickTopping)
	}
}

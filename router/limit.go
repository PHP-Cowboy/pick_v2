package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func LimitRoute(g *gin.RouterGroup) {

	limitGroup := g.Group("/limit", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//订单批量限发
		limitGroup.POST("/orderLimit", handler.OrderLimit)
		//任务批量限发
		limitGroup.POST("/taskLimit", handler.TaskLimit)
	}

}

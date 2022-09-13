package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func YongYouRoute(g *gin.RouterGroup) {

	yongYouRouteGroup := g.Group("/yongYou", middlewares.JWTAuth(), middlewares.IsSuperAdminAuth())
	{
		//u8推送日志列表
		yongYouRouteGroup.GET("/list", handler.LogList)
		//批量补单
		yongYouRouteGroup.POST("/batchSupplement", handler.BatchSupplement)
		//推送u8拣货详情
		yongYouRouteGroup.GET("/logDetail", handler.LogDetail)
	}
}

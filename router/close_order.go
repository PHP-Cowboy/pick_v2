package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func CloseOrderRoute(g *gin.RouterGroup) {

	closeOrderGroup := g.Group("/closeOrder", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//关闭订单状态数量统计
		closeOrderGroup.GET("/closeOrderCount", handler.CloseOrderCount)
		//关闭订单列表
		closeOrderGroup.POST("/closeOrderList", handler.CloseOrderList)
		//关闭订单详情
		closeOrderGroup.POST("/closeOrderDetail", handler.CloseOrderDetail)
		//处理关单
		closeOrderGroup.POST("/closeOrderExec", handler.CloseOrderExec)
	}

}

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
		closeOrderGroup.GET("/closeOrderList", handler.CloseOrderList)
		//关闭订单详情
		closeOrderGroup.GET("/closeOrderDetail", handler.CloseOrderDetail)
		//关闭订单&&详情列表
		closeOrderGroup.GET("/closeOrderAndGoodsList", handler.CloseOrderAndGoodsList)
		//关单处理
		closeOrderGroup.POST("/closeOrderExec", handler.CloseOrderExec)
		//异常关单处理
		closeOrderGroup.POST("/closeOrderExecException", handler.CloseOrderExecException)
	}

}

package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
)

func CloseOrderRoute(g *gin.RouterGroup) {

	closeOrderGroup := g.Group("/closeOrder")
	{
		//关闭订单状态数量统计
		closeOrderGroup.GET("/closeOrderCount", handler.CloseOrderCount)
		//关闭订单列表
		closeOrderGroup.GET("/closeOrderList", handler.CloseOrderList)
		//关闭订单详情
		closeOrderGroup.GET("/closeOrderDetail", handler.CloseOrderDetail)
		//关闭订单&&详情列表
		closeOrderGroup.GET("/closeOrderAndGoodsList", handler.CloseOrderAndGoodsList)
		//关闭关单任务
		closeOrderGroup.POST("/closeCloseOrderTask", handler.CloseCloseOrderTask)
		//关单处理
		closeOrderGroup.POST("/closeOrderExec", handler.CloseOrderExec)
		//
		closeOrderGroup.GET("/testMsgQueue", handler.TestMsgQueue)
	}

}

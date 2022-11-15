package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func OutboundRoute(g *gin.RouterGroup) {

	outboundGroup := g.Group("/outbound", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//生成出库任务
		outboundGroup.POST("/createTask", handler.CreateOutboundTask)
		//出库单任务列表
		outboundGroup.GET("/taskList", handler.OutboundTaskList)
		//出库单订单列表
		outboundGroup.GET("/orderList", handler.OutboundOrderList)
		//出库单订单详情
		outboundGroup.GET("/orderDetail", handler.OutboundOrderDetail)
	}

}

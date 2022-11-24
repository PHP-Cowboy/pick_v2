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
		//简化版任务列表
		outboundGroup.GET("/taskListSimple", handler.OutboundTaskListSimple)
		//出库任务状态数量
		outboundGroup.GET("/count", handler.OutboundTaskCount)
		//出库单订单列表
		outboundGroup.GET("/orderList", handler.OutboundOrderList)
		//出库订单数量
		outboundGroup.GET("/orderCount", handler.OutboundOrderCount)
		//出库单订单详情
		outboundGroup.GET("/orderDetail", handler.OutboundOrderDetail)
		//出库任务商品列表
		outboundGroup.GET("/goodsList", handler.OutboundOrderGoodsList)
		//结束任务
		outboundGroup.POST("/endOutboundTask", handler.EndOutboundTask)
		//关闭订单
		outboundGroup.POST("/closeOrder", handler.OutboundTaskCloseOrder)
		//加单
		outboundGroup.POST("/addOrder", handler.OutboundTaskAddOrder)
		//订单出库记录
		outboundGroup.POST("/orderOutboundRecord", handler.OrderOutboundRecord)
	}

}

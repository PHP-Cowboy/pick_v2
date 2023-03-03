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
		//taskGroup.POST("/orderRepair", handler.OrderRepair)
		//关闭订单
		taskGroup.POST("/closeOrder", handler.CloseOrder)
		//收到拉单
		taskGroup.POST("/pullOrder", handler.PullOrder)
		//批次数据同步
		taskGroup.POST("/getBatchPickData", handler.GetBatchPickData)
	}
}

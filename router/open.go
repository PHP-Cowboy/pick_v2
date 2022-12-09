package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func OpenRoute(g *gin.RouterGroup) {
	group := g.Group("open")
	{
		group.POST("/test", handler.Test)
		group.POST("/testCall", handler.TestCall)
	}

	openGroup := g.Group("open", middlewares.SignAuth())
	{
		//关闭订单
		openGroup.POST("/closeOrder", handler.CloseOrder)
		//关闭商品
		openGroup.POST("/closeOrderGoods", handler.CloseOrderGoods)
		//批次出库订单和商品明细
		openGroup.POST("/getBatchOrderAndGoods", handler.GetBatchOrderAndGoods)
	}
}

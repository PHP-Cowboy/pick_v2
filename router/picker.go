package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func PickerRoute(g *gin.RouterGroup) {
	pickerGroup := g.Group("/picker", middlewares.JWTAuth(), middlewares.IsPickerAuth())
	{
		//接单拣货
		pickerGroup.GET("/receiving_orders", handler.ReceivingOrders)
		//完成拣货
		pickerGroup.POST("/complete", handler.CompletePick)
		//拣货记录
		pickerGroup.GET("/record", handler.PickingRecord)
		//拣货记录明细
		pickerGroup.GET("/detail", handler.PickingRecordDetail)
	}
}

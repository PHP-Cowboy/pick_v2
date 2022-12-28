package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func PickerRoute(g *gin.RouterGroup) {
	pickerGroup := g.Group("/picker", middlewares.JWTAuth(), middlewares.IsPickerAuth())
	{
		//快递拣货列表
		pickerGroup.GET("/centralizedAndSecondaryList", handler.CentralizedAndSecondaryList)
		//集中拣货明细
		pickerGroup.GET("/centralizedPickDetailPDA", handler.CentralizedPickDetailPDA)
		//接单拣货
		pickerGroup.POST("/receiving_orders", handler.ReceivingOrders)
		//集中拣货接单
		pickerGroup.POST("/concentratedPickReceivingOrders", handler.ConcentratedPickReceivingOrders)
		//完成拣货
		pickerGroup.POST("/complete", handler.CompletePick)
		//完成集中拣货
		pickerGroup.POST("/completeConcentratedPick", handler.CompleteConcentratedPick)
		//剩余数量
		pickerGroup.GET("/remaining_quantity", handler.RemainingQuantity)
		//集中拣货剩余数量
		pickerGroup.GET("/centralizedPickRemainingQuantity", handler.CentralizedPickRemainingQuantity)
		//拣货记录
		pickerGroup.GET("/record", handler.PickingRecord)
		//拣货记录明细
		pickerGroup.GET("/detail", handler.PickingRecordDetail)
		//关单提醒
		pickerGroup.GET("/customsDeclarationReminder", handler.CustomsDeclarationReminder)
	}
}

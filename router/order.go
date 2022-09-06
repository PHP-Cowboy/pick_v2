package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func OrderRoute(g *gin.RouterGroup) {

	orderGroup := g.Group("/order", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//拣货单列表
		orderGroup.GET("/list", handler.PickOrderList)
		//拣货单统计
		orderGroup.GET("/pickOrderCount", handler.PickOrderCount)
		//拣货单明细
		orderGroup.GET("/pickOrderDetail", handler.GetPickOrderDetail)
		//配送方式明细
		orderGroup.GET("/deliveryMethodInfo", handler.DeliveryMethodInfo)
		//修改配送方式
		orderGroup.POST("/changeDeliveryMethod", handler.ChangeDeliveryMethod)
		//拣货单商品列表
		orderGroup.GET("/orderGoodsList", handler.OrderGoodsList)
		//限发
		orderGroup.POST("/restrictedShipment", handler.RestrictedShipment)
		//批量限发
		orderGroup.POST("/batchRestrictedShipment", handler.BatchRestrictedShipment)
		// 批量设置限发商品数量
		orderGroup.GET("/goodsNum", handler.GoodsNum)
		//限发列表
		orderGroup.GET("/restrictedShipmentList", handler.RestrictedShipmentList)
		//撤销限发
		orderGroup.POST("/revokeRestrictedShipment", handler.RevokeRestrictedShipment)
	}

}

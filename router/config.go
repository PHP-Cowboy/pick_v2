package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func ConfigRoute(g *gin.RouterGroup) {

	configGroup := g.Group("/config", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//配送方式明细
		configGroup.GET("/deliveryMethodInfo", handler.DeliveryMethodInfo)
		//修改配送方式
		configGroup.POST("/changeDeliveryMethod", handler.ChangeDeliveryMethod)
	}

}

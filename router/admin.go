package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
)

// 后台拣货
func AdminRoute(g *gin.RouterGroup) {
	adminGroup := g.Group("/admin")
	{
		//后台拣货列表
		adminGroup.GET("/pickList", handler.AdminPickList)
		//后台拣货数据详情
		adminGroup.GET("/pickDetail", handler.AdminPickDetail)
		//后台拣货数据详情
		adminGroup.GET("/batchShopGoodsList", handler.BatchShopGoodsList)
		//快捷出库
		adminGroup.POST("/quickDelivery", handler.QuickDelivery)
	}
}

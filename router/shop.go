package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func ShopRoute(g *gin.RouterGroup) {

	group := g.Group("/shop", middlewares.JWTAuth())
	{
		//线路名称列表
		group.GET("line_list", handler.LineList)
		//门店列表
		group.GET("/list", handler.ShopList)
	}

	shopGroup := g.Group("/shop", middlewares.JWTAuth(), middlewares.IsSuperAdminAuth())
	{
		//同步门店
		shopGroup.GET("/sync", handler.SyncShop)
		//批量设置线路
		shopGroup.POST("batch_set", handler.BatchSetLine)

	}
}

package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func GoodsRoute(g *gin.RouterGroup) {
	goodsGroup := g.Group("/goods", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//列表
		goodsGroup.GET("/list", handler.GetGoodsList)
		//明细
		goodsGroup.GET("/detail", handler.GetOrderDetail)
		//商品列表
		goodsGroup.GET("/commodity_list", handler.CommodityList)
	}
}

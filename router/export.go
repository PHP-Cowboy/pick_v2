package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func ExportRoute(g *gin.RouterGroup) {
	exportGroup := g.Group("/export", middlewares.JWTAuth())
	{
		//首批物料导出
		exportGroup.GET("/first_material", handler.FirstMaterial)
		//批次出库导出
		exportGroup.GET("/outbound_batch", handler.OutboundBatch)
	}
}

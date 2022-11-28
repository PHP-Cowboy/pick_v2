package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
)

func ExportRoute(g *gin.RouterGroup) {
	exportGroup := g.Group("/export")
	{
		//首批物料导出
		exportGroup.GET("/first_material", handler.FirstMaterial)
		//批次出库导出
		exportGroup.GET("/outbound_batch", handler.OutboundBatch)
		//欠货信息导出
		exportGroup.GET("/lack", handler.Lack)
		//批次门店信息
		exportGroup.GET("/batch_shop", handler.BatchShop)
		//批次门店物料表
		exportGroup.GET("/batchShopMaterial", handler.BatchShopMaterial)
		//拣货任务导出
		exportGroup.GET("/batchTask", handler.BatchTask)
		//货品汇总单
		exportGroup.GET("/goodsSummaryList", handler.GoodsSummaryList)
	}
}

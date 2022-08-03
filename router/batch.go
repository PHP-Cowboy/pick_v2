package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func BatchRoute(g *gin.RouterGroup) {
	batchGroup := g.Group("/batch", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//列表
		batchGroup.GET("/list", handler.GetBatchList)
		//创建批次
		batchGroup.POST("/create", handler.CreateBatch)
		//预拣池基础信息
		batchGroup.GET("/base", handler.GetBase)
		//预拣池列表
		batchGroup.GET("/pre_pick_list", handler.GetPrePickList)
		//预拣货明细
		batchGroup.GET("/pre_pick_detail", handler.GetPrePickDetail)
	}
}

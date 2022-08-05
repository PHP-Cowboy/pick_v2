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
		//批次池内单数量
		batchGroup.GET("/pool_num", handler.GetPoolNum)
		//批次置顶
		batchGroup.POST("/topping", handler.Topping)
		//批量拣货
		batchGroup.POST("/batch_pick", handler.BatchPick)
		//合并拣货
		batchGroup.POST("/merge_pick", handler.MergePick)
	}
}

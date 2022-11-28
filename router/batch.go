package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func BatchRoute(g *gin.RouterGroup) {
	group := g.Group("/batch")
	{
		//打印
		group.GET("/print_call", handler.PrintCallGet)
	}

	batchGroup := g.Group("/batch", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//列表
		batchGroup.GET("/list", handler.GetBatchList)
		//新创建批次
		batchGroup.POST("/newBatch", handler.NewBatch)
		//创建快递批次
		batchGroup.POST("/courierBatch", handler.CourierBatch)
		//集中拣货列表
		batchGroup.GET("/centralizedPickList", handler.CentralizedPickList)
		//集中拣货详情
		batchGroup.GET("/centralizedPickDetail", handler.CentralizedPickDetail)
		//创建批次
		batchGroup.POST("/create", handler.CreateBatch)
		//根据订单生成批次
		batchGroup.POST("/create_by_order", handler.CreateByOrder)
		//停止拣货
		batchGroup.POST("/change", handler.ChangeBatch)
		//当前批次是否有接单
		batchGroup.GET("/is_pick", handler.IsPick)
		//结束批次
		batchGroup.POST("/end", handler.EndBatch)
		//编辑批次
		batchGroup.POST("/edit", handler.EditBatch)
		//批次池数量
		batchGroup.GET("/batch_pool_num", handler.GetBatchPoolNum)
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

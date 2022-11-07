package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func InvTaskRoute(g *gin.RouterGroup) {

	group := g.Group("/invTask", middlewares.JWTAuth())
	{
		//盘点任务列表
		group.GET("/list", handler.TaskList)
		//自建盘点任务列表
		group.GET("/selfBuiltTaskList", handler.SelfBuiltTaskList)
	}

	taskGroup := g.Group("/invTask", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//同步任务
		taskGroup.POST("/syncTask", handler.SyncTask)
		//变更任务
		taskGroup.POST("/changeTask", handler.ChangeTask)
		//导出
		taskGroup.GET("/export", handler.Export)
		//自建盘点任务
		taskGroup.POST("/selfBuiltTask", handler.SelfBuiltTask)
		//修改自建盘点任务
		taskGroup.POST("/changeSelfBuiltTask", handler.ChangeSelfBuiltTask)
		//设置复盘
		taskGroup.POST("/setSecondInventory", handler.SetSecondInventory)
		//自建盘点任务列表
		taskGroup.POST("/invAgain", handler.InvAgain)
	}

	rGroup := group.Group("/record", middlewares.JWTAuth())
	{
		//任务商品数据列表
		rGroup.GET("/list", handler.TaskRecordList)
		//盘库记录
		rGroup.GET("/inventoryRecord", handler.InventoryRecordList)
		//任务记录分类列表
		rGroup.GET("/typeList", handler.TypeList)
		//盘库记录删除
		rGroup.POST("/inventoryRecordDelete", handler.InventoryRecordDelete)
	}

	recordGroup := group.Group("/record", middlewares.JWTAuth(), middlewares.IsInvAuth())
	{
		//已盘商品件数
		recordGroup.GET("/count", handler.InvCount)
		//用户已盘商品列表
		recordGroup.GET("/userInventoryRecordList", handler.UserInventoryRecordList)
		//修改已盘商品数据
		recordGroup.POST("/updateInventoryRecord", handler.UpdateInventoryRecord)
		//批量盘点
		recordGroup.POST("/batchCreate", handler.BatchCreate)
		//未盘商品数量
		recordGroup.GET("/notInvCount", handler.NotInvCount)
		//未盘商品列表
		recordGroup.GET("/userNotInventoryRecordList", handler.UserNotInventoryRecordList)
	}
}

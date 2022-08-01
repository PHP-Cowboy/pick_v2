package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func ClassRoute(g *gin.RouterGroup) {
	syncGroup := g.Group("/class", middlewares.JWTAuth(), middlewares.IsSuperAdminAuth())
	{
		//同步分类
		syncGroup.GET("/sync", handler.SyncClassification)
		//分类列表
		syncGroup.GET("/list", handler.ClassList)
		//批量设置分类
		syncGroup.POST("/batch_set", handler.BatchSetClass)
		//分类名称列表
		syncGroup.GET("/class_name_list", handler.ClassNameList)
	}
}

package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func U8Route(g *gin.RouterGroup) {

	u8Group := g.Group("/user", middlewares.JWTAuth(), middlewares.IsSuperAdminAuth())
	{
		//用户列表
		u8Group.GET("/list", handler.LogList)
		//批量补单
		u8Group.POST("/batchSupplement", handler.BatchSupplement)
	}
}

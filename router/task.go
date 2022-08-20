package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func TaskRoute(g *gin.RouterGroup) {
	taskGroup := g.Group("/task", middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//列表
		taskGroup.GET("/list", handler.PickList)
		//详情
		taskGroup.GET("/pick_detail", handler.GetPickDetail)
		//修改已出库件数
		taskGroup.POST("/change_num", handler.ChangeNum)
		//拣货置顶
		taskGroup.POST("/pick_topping", handler.PickTopping)
		//打印
		taskGroup.POST("/push_print", handler.PushPrint)
		//指派
		taskGroup.POST("/assign", handler.Assign)
	}
}

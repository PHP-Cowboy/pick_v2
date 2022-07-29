package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func UserRoute(g *gin.RouterGroup) {
	//登录
	group := g.Group("/user")
	{
		group.POST("/login", handler.Login)
	}

	userGroup := g.Group("/user", middlewares.JWTAuth(), middlewares.IsSuperAdminAuth())
	{
		//用户列表
		userGroup.GET("/list", handler.GetUserList)
		//新增用户
		userGroup.POST("/create", handler.CreateUser)
		//修改密码
		userGroup.POST("/change", handler.ChangeUser)
		//获取仓库用户数
		userGroup.POST("/getWarehouseUserCount", handler.GetWarehouseUserCountList)
	}
}

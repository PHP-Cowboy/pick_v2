package initialize

import (
	"github.com/gin-gonic/gin"
	"pick_v2/middlewares"
	"pick_v2/router"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	//跨域
	r.Use(middlewares.Cors())

	group := r.Group("/v2")
	//用户
	router.UserRoute(group)
	//角色
	router.RoleRoute(group)
	//订单商品
	router.GoodsRoute(group)
	//批次
	router.BatchRoute(group)
	//仓库
	router.WarehouseRoute(group)

	return r
}

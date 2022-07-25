package initialize

import (
	"github.com/gin-gonic/gin"
	"pick_v2/middlewares"
	"pick_v2/router"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.Cors())

	group := r.Group("/v2")

	router.UserRoute(group)

	router.BaseRoute(group)

	router.RoleRoute(group)

	router.GoodsRoute(group)

	return r
}

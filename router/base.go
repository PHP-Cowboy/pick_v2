package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
)

func BaseRoute(g *gin.RouterGroup) {
	userGroup := g.Group("/base")
	{
		userGroup.GET("/captcha", handler.GenerateCaptcha)
	}
}

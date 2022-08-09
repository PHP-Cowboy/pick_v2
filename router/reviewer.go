package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func ReviewerRoute(g *gin.RouterGroup) {
	reviewGroup := g.Group("/review", middlewares.JWTAuth(), middlewares.IsReviewerAuth())
	{
		//复核列表
		reviewGroup.GET("/review_list", handler.ReviewList)
		//复核明细
		reviewGroup.GET("/review_detail", handler.ReviewDetail)
		//确认出库
		reviewGroup.POST("/confirm_delivery", handler.ConfirmDelivery)
	}
}

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
	//出库单
	router.OutboundRoute(group)
	//限发
	router.LimitRoute(group)
	//批次
	router.BatchRoute(group)
	//任务
	router.TaskRoute(group)
	//仓库
	router.WarehouseRoute(group)
	//字典
	router.DictRoute(group)
	//店铺
	router.ShopRoute(group)
	//数据同步
	router.ClassRoute(group)
	//拣货员
	router.PickerRoute(group)
	//复核员
	router.ReviewerRoute(group)
	//菜单
	router.MenuRoute(group)
	//导出
	router.ExportRoute(group)
	//开放数据
	router.OpenRoute(group)
	//u8
	router.YongYouRoute(group)
	//盘点
	router.InvTaskRoute(group)
	//配置
	router.ConfigRoute(group)

	return r
}

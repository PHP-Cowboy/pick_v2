package router

import (
	"github.com/gin-gonic/gin"
	"pick_v2/handler"
	"pick_v2/middlewares"
)

func DictRoute(g *gin.RouterGroup) {
	roleGroup := g.Group("/dict", middlewares.JWTAuth(), middlewares.IsSuperOrWarehouseAdminAuth())
	{
		//字典类型列表
		roleGroup.GET("/dict_type_list", handler.DictTypeList)
		//新增字典类型
		roleGroup.POST("/create_dict_type", handler.CreateDictType)
		//修改字典类型
		roleGroup.POST("/change_dict_type", handler.ChangeDictType)
		//删除字典类型
		roleGroup.POST("/delete_dict_type", handler.DeleteDictType)
		//字典数据列表
		roleGroup.GET("/dict_list", handler.DictList)
		//新增字典数据
		roleGroup.POST("/create_dict", handler.CreateDict)
		//修改字典数据
		roleGroup.POST("/change_dict", handler.ChangeDict)
		//删除字典
		roleGroup.POST("/delete_dict", handler.DeleteDict)
	}
}

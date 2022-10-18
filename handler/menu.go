package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"pick_v2/forms/req"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
	"time"
)

func GetMenuList(c *gin.Context) {
	var menu []Menu

	db := global.DB

	result := db.Model(&model.Menu{}).Where("delete_time is null").Order("sort desc").Find(&menu)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	treeList := GetMenuTree(menu, 0)

	xsq_net.SucJson(c, treeList)
}

type Menu struct {
	Id    int    `json:"id"`
	Type  int    `json:"type"`
	Title string `json:"title"`
	Path  string `json:"path"`
	PId   int    `json:"p_id"`
	Child []Menu `json:"child"`
}

type TreeList struct {
	Id    int        `json:"id"`
	Type  int        `json:"type"`
	Title string     `json:"title"`
	Path  string     `json:"path"`
	PId   int        `json:"p_id"`
	Child []TreeList `json:"child"`
}

func GetMenuTree(menuList []Menu, pId int) (treeList []TreeList) {

	for _, v := range menuList {
		if v.PId == pId {
			child := GetMenuTree(menuList, v.Id)
			node := TreeList{
				Id:    v.Id,
				Type:  v.Type,
				Title: v.Title,
				Path:  v.Path,
				PId:   v.PId,
			}
			node.Child = child
			treeList = append(treeList, node)
		}
	}

	return treeList
}

func CreateMenu(c *gin.Context) {
	var (
		form req.CreateMenuReq
		menu model.Menu
	)

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	menu.Type = form.Type
	menu.Title = form.Title
	menu.Path = form.Path
	menu.PId = form.PId
	menu.Sort = form.Sort

	result := global.DB.Save(&menu)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

func ChangeMenu(c *gin.Context) {
	var (
		form req.ChangeMenuReq
	)

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result := global.DB.Model(&model.Menu{}).
		Where("id = ?", form.Id).
		Updates(map[string]interface{}{
			"type":  form.Type,
			"title": form.Title,
			"path":  form.Path,
			"p_id":  form.PId,
			"sort":  form.Sort,
		})

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

func BatchDeleteMenu(c *gin.Context) {
	var form req.BatchDeleteMenuReq

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	now := time.Now()

	result := global.DB.Model(model.Menu{}).Where("id in (?)", form.Ids).Update("delete_time", &now)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

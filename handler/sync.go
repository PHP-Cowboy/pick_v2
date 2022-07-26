package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"

	"pick_v2/utils/request"
	"pick_v2/utils/xsq_net"
)

//同步门店
func SyncShop(c *gin.Context) {

	url := ""

	body, err := request.Request(url, "")
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var result interface{}

	err = json.Unmarshal(body, &result)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, result)
}

//同步分类
func SyncClassification(c *gin.Context) {
	url := ""

	body, err := request.Request(url, "")
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	var result interface{}

	err = json.Unmarshal(body, &result)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, result)
}

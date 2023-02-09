package handler

import (
	"github.com/gin-gonic/gin"
	"pick_v2/utils/cache"
	"pick_v2/utils/xsq_net"
)

func SAdd(c *gin.Context) {

	arr := []string{"1", "2"}

	err := cache.SAdd("setSalesman", arr)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

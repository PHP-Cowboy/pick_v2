package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
)

//获取待拣货订单商品列表
func GetGoodsList(c *gin.Context) {

	var form req.GetGoodsListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	cfg := global.ServerConfig

	url := fmt.Sprintf("%s:%d/api/v1/remote/pick/lack/list", cfg.GoodsApi.Url, cfg.GoodsApi.Port)
	method := "POST"

	client := &http.Client{}

	rq, err := http.NewRequest(method, url, nil)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}
	res, err := client.Do(rq)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	fmt.Println(string(body))

	var result rsp.ApiGoodsListRsp

	err = json.Unmarshal(body, &result)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, result)
}

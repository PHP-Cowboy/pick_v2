package xsq_net

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ErrMsg struct {
	Err  error  `json:"err"`
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type Response struct {
	// 代码
	Code int `json:"code" example:"200"`
	// 数据集
	Data interface{} `json:"data"`
	// 消息
	Msg string `json:"msg"`
}

func ReplyError(c *gin.Context, err error, msg string, code int, req interface{}) {
	zap.S().Info("[API]:", c.Request.URL.RequestURI(), "[ERR]:", err, "[MSG]:", msg, "[REQUEST]:", req)
	c.JSON(http.StatusInternalServerError, ErrMsg{err, msg, code})
}

func ReplyOK(c *gin.Context, data interface{}, msg string) {
	var res Response
	res.Code = 200
	res.Data = data
	if msg != "" {
		res.Msg = msg
	}else {
		res.Msg = "success"
	}
	c.JSON(http.StatusOK, res)
}

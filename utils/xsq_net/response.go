package xsq_net

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"pick_v2/utils/ecode"
)

var Empty = &struct{}{}

// 通过传递Codes方式
func CodeWithJSON(c *gin.Context, data interface{}, code ecode.Codes) {
	c.JSON(http.StatusOK, gin.H{
		"code": code.Code(),
		"msg":  code.Message(),
		"data": data,
	})
}

//未知错误 错误码为ecode.ServerErr，msg：为具体的错误信息
func WithOutECodeErrorJSON(c *gin.Context, data interface{}, err error) {
	c.JSON(http.StatusOK, gin.H{
		"code": ecode.ServerErr,
		"msg":  err.Error(),
		"data": data,
	})
}

// 通过传递error方式，不用关心具体错误
// 如果是业务捕捉的错误码 则返回业务错误码
// 如果不是业务错误码 则返回具体的错误码
func JSON(c *gin.Context, data interface{}, err error) {
	if err == nil {
		SucJson(c, data)
		return
	}
	e, ok := err.(ecode.Codes)
	if ok {
		CodeWithJSON(c, data, ecode.Cause(e))
		return
	}
	WithOutECodeErrorJSON(c, data, err)
}

// 错误业务JSON 不需关注输出数据
func ErrorJSON(c *gin.Context, err error) {
	if err == ecode.ParamInvalid {
		var body []byte
		if cb, ok := c.Get(gin.BodyBytesKey); ok {
			if cbb, ok := cb.([]byte); ok {
				body = cbb
			}
		}
		if len(body) > 0 {
			zap.S().Errorf("paramsInvalid, url:%s,params:%s", c.Request.URL.Path, string(body))
		}
	}
	JSON(c, Empty, err)
}

// 成功业务JSON
func SucJson(c *gin.Context, data interface{}) {
	CodeWithJSON(c, data, ecode.Success)
}

// 成功业务JSON 只需要状态码
func Success(c *gin.Context) {
	CodeWithJSON(c, Empty, ecode.Success)
}

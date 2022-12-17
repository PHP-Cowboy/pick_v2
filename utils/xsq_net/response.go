package xsq_net

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pick_v2/global"
	"pick_v2/utils/ecode"
	"runtime"
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

// 未知错误 错误码为ecode.ServerErr，msg：为具体的错误信息
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
	e, ok := err.(ecode.Codes)
	if ok {
		CodeWithJSON(c, data, ecode.Cause(e))
		return
	}
	WithOutECodeErrorJSON(c, data, err)
}

// 错误业务JSON 不需关注输出数据
func ErrorJSON(c *gin.Context, err error) {
	_, file, line, _ := runtime.Caller(1)

	WriteErrLog(c, err, file, line)
	JSON(c, Empty, err)
}

// 成功业务JSON
func SucJson(c *gin.Context, data interface{}) {
	log.Print(c.Request.URL.RawQuery)
	WriteInfoLog(c, data)
	CodeWithJSON(c, data, ecode.Success)
}

// 成功业务JSON 只需要状态码
func Success(c *gin.Context) {
	WriteInfoLog(c, "{}")
	CodeWithJSON(c, Empty, ecode.Success)
}

func WriteErrLog(c *gin.Context, err error, file string, line int) {

	l, ok := global.Logger["err"]

	if !ok {
		panic("err日志加载失败")
	}

	params, pErr := Params(c)

	if pErr != nil {
		panic(pErr.Error())
	}

	l.Infof("file:%s,line:%d,url:%s,params:%s,err:%s", file, line, c.Request.URL.Path, params, err.Error())

}

func WriteInfoLog(c *gin.Context, data interface{}) {

	l, ok := global.Logger["info"]

	if !ok {
		panic("info日志加载失败")
	}

	params, pErr := Params(c)
	if pErr != nil {
		panic(pErr.Error())
	}

	j, err := json.Marshal(data)
	if err != nil {
		return
	}

	l.Infof("url:%s,params:%s,rsp:%v", c.Request.URL.Path, params, string(j))

}

func Params(c *gin.Context) (params interface{}, err error) {

	switch c.Request.Method {
	case "GET":
		params, err = ParamsGet(c)
		break
	case "POST":
		params, err = ParamsPost(c)
		break
	}

	return
}

func ParamsGet(c *gin.Context) (params interface{}, err error) {
	//rawQuery := c.Request.URL.RawQuery

	query := c.Request.URL.Query()

	params, err = json.Marshal(query)

	if err != nil {
		return nil, err
	}

	return
}

func ParamsPost(c *gin.Context) (params interface{}, err error) {

	var body []byte
	if cb, ok := c.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}

	return body, nil
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

func GenerateCaptcha(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)

	c := base64Captcha.NewCaptcha(driver, store)

	id, b64s, err := c.Generate()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":   id,
		"data": b64s,
	})
}

//func Verify(id string, value string) bool {
//	return store.Verify(id, value, true)
//}

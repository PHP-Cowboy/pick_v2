package middlewares

import (
	"crypto/sha512"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/gin-gonic/gin"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
)

var key = "9sWBFw96W1Vf7Bb4"

func GetOptions() *password.Options {
	return &password.Options{16, 16, 16, sha512.New}
}

func SignAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sign := c.Request.Header.Get("x-sign")
		salt := c.Request.Header.Get("x-salt")
		if sign == "" {
			xsq_net.ErrorJSON(c, ecode.IllegalRequest)
			c.Abort()
			return
		}

		options := GetOptions()

		if !password.Verify(key, salt, sign, options) {
			xsq_net.ErrorJSON(c, ecode.CommunalSignInvalid)
			c.Abort()
			return
		}

		c.Next()
	}
}

func Generate() (salt string, sign string) {
	options := GetOptions()

	return password.Encode(key, options)
}

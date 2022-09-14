package middlewares

import (
	"crypto/sha512"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/gin-gonic/gin"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
	"strings"
)

var key = "9sWBFw96W1Vf7Bb4"

func GetOptions() *password.Options {
	return &password.Options{16, 16, 16, sha512.New}
}

func SignAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sign := c.Request.Header.Get("x-sign")
		if sign == "" {
			xsq_net.ErrorJSON(c, ecode.IllegalRequest)
			c.Abort()
			return
		}

		options := GetOptions()

		signSlice := strings.Split(sign, "$")

		if len(signSlice) != 2 {
			xsq_net.ErrorJSON(c, ecode.CommunalSignInvalid)
			c.Abort()
			return
		}

		if !password.Verify(key, signSlice[0], signSlice[1], options) {
			xsq_net.ErrorJSON(c, ecode.CommunalSignInvalid)
			c.Abort()
			return
		}

		c.Next()
	}
}

func Generate() string {
	options := GetOptions()

	salt, encode := password.Encode(key, options)

	return salt + "$" + encode
}

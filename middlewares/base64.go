package middlewares

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
	"strconv"
	"strings"
	"time"
)

//整个字符串不能有$
var salt = "9sWBFw96W1Vf7Bb4"

func GetBase64EncodeStr() string {
	//加密字符串
	encodingStr := fmt.Sprintf("%s$%d", salt, time.Now().Unix())

	return base64.StdEncoding.EncodeToString([]byte(encodingStr))
}

func Base64SignAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sign := c.Request.Header.Get("x-sign")
		if sign == "" {
			xsq_net.ErrorJSON(c, ecode.IllegalRequest)
			c.Abort()
			return
		}

		decode, err := base64.StdEncoding.DecodeString(sign)
		if err != nil {
			xsq_net.ErrorJSON(c, ecode.IllegalRequest)
			c.Abort()
			return
		}

		decodeString := string(decode)

		signSlice := strings.Split(decodeString, "$")

		if len(signSlice) != 2 || signSlice[0] != salt {
			xsq_net.ErrorJSON(c, ecode.CommunalSignInvalid)
			c.Abort()
			return
		}

		unix, err := strconv.ParseInt(signSlice[1], 10, 64)
		if err != nil {
			xsq_net.ErrorJSON(c, ecode.CommunalSignInvalid)
			c.Abort()
			return
		}

		if time.Now().Unix()-unix > 10 {
			xsq_net.ErrorJSON(c, ecode.CommunalSignInvalid)
			c.Abort()
			return
		}

		c.Next()
	}
}

package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"

)

func IsAdminAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "claims获取失败",
			})
			c.Abort()
			return
		}

		userInfo := claims.(*CustomClaims)

		if userInfo.AuthorityId != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "无权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}

}
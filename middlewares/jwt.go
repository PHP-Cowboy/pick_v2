package middlewares

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"pick_v2/common/constant"
	"pick_v2/global"
	"pick_v2/utils/ecode"
	"pick_v2/utils/xsq_net"
)

type CustomClaims struct {
	ID          int
	Account     string
	Name        string
	WarehouseId int
	AuthorityId int
	jwt.StandardClaims
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中
		//不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			xsq_net.ErrorJSON(c, ecode.UserNotLogin)
			c.Abort()
			return
		}
		j := NewJwt()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				xsq_net.ErrorJSON(c, ecode.TokenExpired)
				c.Abort()
				return
			}

			xsq_net.ErrorJSON(c, ecode.UserNotLogin)
			c.Abort()
			return
		}
		//校验redis中存储的是否和传入的一致
		redisKey := constant.LOGIN_PREFIX + claims.Account
		val, err := global.Redis.Get(context.Background(), redisKey).Result()
		if err != nil {
			xsq_net.ErrorJSON(c, ecode.RedisFailedToGetData)
			c.Abort()
			return
		}

		//和redis中的不一致，则返还token失效，前端退出
		if token != val {
			xsq_net.ErrorJSON(c, ecode.TokenExpired)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("warehouseId", claims.WarehouseId)
		c.Next()
	}
}

type Jwt struct {
	Key []byte `json:"key"`
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJwt() *Jwt {
	return &Jwt{Key: []byte(global.ServerConfig.JwtInfo.SigningKey)}
}

// 创建token
func (j *Jwt) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Key)
}

// 解析token
func (j *Jwt) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.Key, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}

// 更新token
func (j *Jwt) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.Key, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

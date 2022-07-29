package initialize

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"pick_v2/global"
)

func InitRedis() {
	redisInfo := global.ServerConfig.RedisInfo

	global.Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisInfo.Host, redisInfo.Port),
		Password: redisInfo.Password, // no password set
		DB:       0,                  // use default DB
	})
}

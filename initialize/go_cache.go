package initialize

import (
	"github.com/patrickmn/go-cache"
	"pick_v2/global"
	"time"
)

func InitGoCache() {
	//创建一个默认过期时间为24小时的缓存
	//每12小时清洗一次过期物品
	global.GoCache = cache.New(24*time.Hour, 12*time.Hour)
}

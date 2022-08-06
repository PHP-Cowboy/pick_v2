package cache

import (
	"github.com/patrickmn/go-cache"
	"pick_v2/global"
	"pick_v2/model/other"
	"time"
)

var goCache *cache.Cache

func init() {
	//创建一个默认过期时间为24小时的缓存
	//每12小时清洗一次过期物品
	goCache = cache.New(24*time.Hour, 12*time.Hour)
}

func GetClassification() (map[string]string, error) {

	goodsClassMap, ok := goCache.Get("goodsClassMap")

	mp := make(map[string]string, 0)

	//ok = false

	if !ok {
		var class []other.Classification

		result := global.DB.Find(&class)

		if result.Error != nil {
			return nil, result.Error
		}

		for _, cl := range class {
			mp[cl.GoodsClass] = cl.WarehouseClass
		}

		goCache.Set("goodsClassMap", mp, cache.DefaultExpiration)

	} else {
		mp = goodsClassMap.(map[string]string)
	}

	return mp, nil
}

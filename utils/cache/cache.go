package cache

import (
	"github.com/patrickmn/go-cache"
	"pick_v2/global"
	"pick_v2/model"
	"time"
)

var goCache *cache.Cache

func init() {
	//创建一个默认过期时间为24小时的缓存
	//每12小时清洗一次过期物品
	goCache = cache.New(24*time.Hour, 12*time.Hour)
}

// 获取分类缓存
func GetClassification() (map[string]string, error) {

	goodsClassMap, ok := goCache.Get("goodsClassMap")

	mp := make(map[string]string, 0)

	if !ok {
		var class []model.Classification

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

// 更新分类缓存
func SetClassification() error {

	var class []model.Classification

	result := global.DB.Find(&class)

	if result.Error != nil {
		return result.Error
	}

	mp := make(map[string]string, 0)

	for _, cl := range class {
		mp[cl.GoodsClass] = cl.WarehouseClass
	}

	goCache.Set("goodsClassMap", mp, cache.DefaultExpiration)

	return nil
}

// 获取店铺线路
func GetShopLine() (map[int]string, error) {
	shopLineMap, ok := goCache.Get("shopLine")

	mp := make(map[int]string, 0)

	if !ok {
		var shops []model.Shop

		result := global.DB.Select("shop_id,line").Find(&shops)

		if result.Error != nil {
			return nil, result.Error
		}

		for _, shop := range shops {
			mp[shop.ShopId] = shop.Line
		}

		goCache.Set("shopLine", mp, cache.DefaultExpiration)

	} else {
		mp = shopLineMap.(map[int]string)
	}

	return mp, nil
}

// 更新店铺线路缓存
func SetShopLine() error {
	var shops []model.Shop

	result := global.DB.Find(&shops)

	if result.Error != nil {
		return result.Error
	}

	mp := make(map[int]string, 0)

	for _, shop := range shops {
		mp[shop.ShopId] = shop.Line
	}

	goCache.Set("shopLine", mp, cache.DefaultExpiration)
	return nil
}

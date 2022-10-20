package cache

import (
	"context"
	"pick_v2/global"
	"pick_v2/utils/str_util"
	"pick_v2/utils/timeutil"
	"strconv"
	"time"
)

func GetIncrNumberByKey(key string, padLength int) (string, error) {

	dateNumber := time.Now().Format(timeutil.DateNumberFormat)
	//rds key
	redisKey := key + dateNumber

	val, err := global.Redis.Do(context.Background(), "incr", redisKey).Result()
	if err != nil {
		return "", err
	}

	//设置过期时间
	err = Expire(redisKey, 24*60*60)
	if err != nil {
		return "", err
	}

	number := strconv.Itoa(int(val.(int64)))

	No := dateNumber + str_util.StrPad(number, padLength, "0", 0)

	return No, nil
}

func GetIncrByKey(key string) (interface{}, error) {

	val, err := global.Redis.Do(context.Background(), "incr", key).Result()

	if err != nil {
		return val, err
	}

	return val, nil
}

func TTL(key string) (time.Duration, error) {
	return global.Redis.TTL(context.Background(), key).Result()
}

func Set(key, val string, second int) (string, error) {
	return global.Redis.Set(context.Background(), key, val, time.Duration(second)*time.Second).Result()
}

func Get(key string) (string, error) {
	return global.Redis.Get(context.Background(), key).Result()
}

func SetNx(key string, val string) error {
	return global.Redis.SetNX(context.Background(), key, val, 0).Err()
}

// 删除
func Del(key string) error {
	return global.Redis.Del(context.Background(), key).Err()
}

// 设置过期时间
func Expire(key string, second int) error {
	return global.Redis.Expire(context.Background(), key, time.Duration(second)).Err()
}

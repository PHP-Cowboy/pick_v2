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

	//ok, err := global.Redis.Expire(context.Background(), "key", 24*time.Hour).Result()
	//
	//if err != nil {
	//	return "", errors.New("设置过期返回false")
	//}
	//
	//if !ok {
	//	return "", nil
	//}

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

func GetTtlKey(key string) (time.Duration, error) {

	return global.Redis.TTL(context.Background(), key).Result()
}

func SetTtlKey(key string, second int) (bool, error) {
	return global.Redis.Expire(context.Background(), key, time.Duration(second)*time.Second).Result()
}

func Set(key, val string, second int) (string, error) {
	return global.Redis.Set(context.Background(), key, val, time.Duration(second)*time.Second).Result()
}

package cache

import (
	"context"
	"pick_v2/global"
	"pick_v2/utils/str_util"
	"pick_v2/utils/timeutil"
	"strconv"
	"time"
)

func GetIncrNumberByKey(key string) (string, error) {

	//rds key
	redisKey := key + time.Now().Format(timeutil.DateNumberFormat)

	val, err := global.Redis.Do(context.Background(), "incr", redisKey).Result()
	if err != nil {
		return "", err
	}

	number := strconv.Itoa(int(val.(int64)))

	return str_util.StrPad(number, 3, "0", 0), nil
}

func GetIncrByKey(key string) (interface{}, error) {

	val, err := global.Redis.Do(context.Background(), "incr", key).Result()

	if err != nil {
		return val, err
	}

	return val, nil
}

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

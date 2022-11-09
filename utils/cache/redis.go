package cache

import (
	"context"
	"errors"
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

	val, err := Incr(redisKey)
	if err != nil {
		return "", err
	}

	//设置过期时间
	err = Expire(redisKey, 24*time.Hour)

	number := strconv.Itoa(int(val))

	No := dateNumber + str_util.StrPad(number, padLength, "0", 0)

	return No, nil
}

func Incr(key string) (int64, error) {
	return global.Redis.Incr(context.Background(), key).Result()
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
	ok, err := global.Redis.SetNX(context.Background(), key, val, 0).Result()

	if !ok {
		return errors.New("设置过期时间失败")
	}

	return err
}

// DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
// 返回值: 被删除 key 的数量。
func Del(key string) (int64, error) {
	return global.Redis.Del(context.Background(), key).Result()
}

// 设置过期时间
func Expire(key string, d time.Duration) error {
	ok, err := global.Redis.Expire(context.Background(), key, d).Result()

	if !ok {
		return errors.New("设置过期时间失败")
	}

	return err
}

func ExpireAt(key string, second int) error {
	ok, err := global.Redis.ExpireAt(context.Background(), key, time.Now().Add(time.Duration(second))).Result()

	if !ok {
		return errors.New("设置过期时间失败")
	}

	return err
}

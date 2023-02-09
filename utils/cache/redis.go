package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
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

func Get(key string) (val string, err error) {
	val, err = global.Redis.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return
}

func SetNx(key string, val string) error {
	ok, err := global.Redis.SetNX(context.Background(), key, val, 0).Result()

	if !ok {
		return errors.New("设置失败")
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

/*
Redis Sadd 命令将一个或多个成员元素加入到集合中，已经存在于集合的成员元素将被忽略。
假如集合 key 不存在，则创建一个只包含添加的元素作成员的集合。
当集合 key 不是集合类型时，返回一个错误。
注意：在 Redis2.4 版本以前， SADD 只接受单个成员值。
*/
func SAdd(key string, members []string) (err error) {
	_, err = global.Redis.SAdd(context.Background(), key, members).Result()

	return
}

/*
Redis Srem 命令用于移除集合中的一个或多个成员元素，不存在的成员元素会被忽略。
当 key 不是集合类型，返回一个错误。
在 Redis 2.4 版本以前， SREM 只接受单个成员值。
*/
func SRem(key string, members []string) (err error) {
	_, err = global.Redis.SRem(context.Background(), key, members).Result()

	return
}

/*
Redis Smembers 命令返回集合中的所有的成员。 不存在的集合 key 被视为空集合。
*/
func SMembers(key string) (err error) {
	_, err = global.Redis.SMembers(context.Background(), key).Result()

	return
}

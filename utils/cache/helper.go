package cache

import (
	"errors"
)

// redis防重复点击
func AntiRepeatedClick(key string, d int) (err error) {
	var value string

	value, err = Get(key)
	if err != nil {
		return
	}

	if value != "" {
		err = errors.New("处理中，请稍后重试")
		return
	}

	_, err = Set(key, "1", d)

	if err != nil {
		return
	}

	return
}

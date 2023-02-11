package slice

import (
	"fmt"
	"strings"
)

// 泛型版 slice 去重
func UniqueSlice[T comparable](arr []T) []T {
	ret := make([]T, 0)

	mp := make(map[T]struct{}, 0)

	for _, t := range arr {
		_, ok := mp[t]
		if ok {
			continue
		}

		mp[t] = struct{}{}
		ret = append(ret, t)
	}

	return ret
}

func DelVal[T comparable](arr []T, val T) (ret []T) {

	ret = make([]T, 0, len(arr))

	for _, t := range arr {
		if t != val {
			ret = append(ret, t)
		}
	}

	return
}

// map 转 slice
func MapToSlice[T comparable](mp map[T]struct{}) (arr []T) {
	for id, _ := range mp {
		arr = append(arr, id)
	}
	return arr
}

// slice 转 string
func SliceToString[T comparable](arr []T, separator string) (res string) {
	for _, t := range arr {
		res += fmt.Sprintf("%s%s", t, separator)
	}

	return strings.TrimRight(res, separator)
}

func StringToSlice(str string, separator string) (res []string) {
	return strings.Split(str, separator)
}

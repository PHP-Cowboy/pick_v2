package slice

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

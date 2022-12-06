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

// slice 转 map
func SliceToMap[T comparable](arr []T) (mp map[T]struct{}) {
	mp = make(map[T]struct{}, 0)
	for _, t := range arr {
		mp[t] = struct{}{}
	}
	return mp
}

// map 转 slice
func MapToSlice[T comparable](mp map[T]struct{}) (arr []T) {
	for id, _ := range mp {
		arr = append(arr, id)
	}
	return arr
}

// 在 a 不在 b 中的
// a := []string{"5", "2", "3", "4"}
// b := []string{"0", "1", "2", "3"}
// c => [5 4]
func StrDiff(a []string, b []string) []string {

	var c []string

	// map[string]struct{}{}创建了一个key类型为String值类型为空struct的map，Equal -> make(map[string]struct{})
	temp := map[string]struct{}{}

	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			c = append(c, val)
		}
	}

	return c
}

func UniqueIntSlice(arr []int) []int {
	ret := make([]int, 0)

	mp := make(map[int]struct{}, 0)

	for _, s := range arr {
		_, ok := mp[s]
		if ok {
			continue
		}

		mp[s] = struct{}{}
		ret = append(ret, s)
	}

	return ret
}

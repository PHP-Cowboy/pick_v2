package slice

func UniqueStringSlice(arr []string) []string {
	ret := make([]string, 0)

	mp := make(map[string]struct{}, 0)

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

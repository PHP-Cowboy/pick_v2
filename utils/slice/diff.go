package slice

// 在 a 不在 b 中的
// a := []T{"5", "2", "3", "4"}
// b := []T{"0", "1", "2", "3"}
// c => [5 4]
func Diff[T comparable](a []T, b []T) (c []T) {

	// map[string]struct{}{}创建了一个key类型为String值类型为空struct的map，Equal -> make(map[string]struct{})
	temp := map[T]struct{}{}

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

	return
}

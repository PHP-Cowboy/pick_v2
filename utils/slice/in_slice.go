package slice

//字符串切片转map 空结构体 struct{} 不占用内存空间
func StrSliceToMap(strSlice []string) map[string]struct{} {
	set := make(map[string]struct{}, len(strSlice))

	for _, s := range strSlice {
		set[s] = struct{}{}
	}

	return set
}

func InMap(m map[string]struct{}, s string) bool {
	_, ok := m[s]
	return ok
}

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

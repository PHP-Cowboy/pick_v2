package slice

// 返回在 arr1 但不在 arr2 中的数据
func DiffInt(arr1, arr2 []int) []int {
	mp := make(map[int]struct{}, 0)

	arr := []int{}

	for _, a := range arr1 {
		mp[a] = struct{}{}
	}

	for _, i := range arr2 {
		_, ok := mp[i]

		if ok {
			continue
		}

		arr = append(arr, i)
	}

	return arr
}

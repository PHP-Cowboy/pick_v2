package slice

import "reflect"

// slice 转 map
func SliceToMap[T comparable](arr []T) (mp map[T]struct{}) {
	mp = make(map[T]struct{}, 0)
	for _, t := range arr {
		mp[t] = struct{}{}
	}
	return mp
}

func InMap(m map[string]struct{}, s string) bool {
	_, ok := m[s]
	return ok
}

// 判断某一个值是否含在切片之中
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		length := s.Len()

		for i := 0; i < length; i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

// 判断某一个值是否含在切片之中,并返回找到的第一个值的索引
func inArray[T comparable](val T, array []T) (exists bool, index int) {
	exists = false
	index = -1

	for i, t := range array {
		if t == val {
			exists = true
			index = i
			break
		}
	}

	return
}

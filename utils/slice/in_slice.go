package slice

import "reflect"

// 字符串切片转map 空结构体 struct{} 不占用内存空间
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

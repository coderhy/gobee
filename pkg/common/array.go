package common

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// ArrayFlip 简单反转(key => val)
func ArrayFlip(m map[string]string) map[string]string {
	// func ArrayFlip(m interface{}) map[string]interface{} {
	// n := make(map[string]interface{})
	var n = make(map[string]string)
	for i, v := range m {
		// for i, v := range m.(map[string]string) {
		n[v] = i
	}
	return n
}

// ArrayMerge 合并数组array_merge()
func ArrayMerge(dataMap ...[]interface{}) []interface{} {
	n := 0
	for k := range dataMap {
		n += len(dataMap[k])
	}
	result := make([]interface{}, 0, n)
	for k := range dataMap {
		result = append(result, dataMap[k]...)
	}
	return result
}

// ArrayColumn array_column()
func ArrayColumn(data []map[string]interface{}, val string, key string) (result interface{}) {
	if val != "nil" && key != "nil" {
		result := map[interface{}]interface{}{}
		for k := range data {
			kVal, keyOk := data[k][key]
			vVal, valOk := data[k][val]
			if keyOk && valOk {
				result[kVal] = vVal
			}
		}
		return result
	} else if val != "nil" && key == "nil" {
		result := []interface{}{}
		for k := range data {
			if vVal, valOk := data[k][val]; valOk {
				result = append(result, vVal)
			}
		}
		return result
	} else if val == "nil" && key != "nil" {
		result := map[interface{}]interface{}{}
		for k := range data {
			if kVal, keyOk := data[k][key]; keyOk {
				result[kVal] = data[k]
			}
		}
		return result
	}
	return result
}

// ArrayColumn_bak array_column()
func ArrayColumn_bak(input map[string]map[string]interface{}, columnKey string) []interface{} {
	columns := make([]interface{}, 0, len(input))
	for k := range input {
		if v, ok := input[k][columnKey]; ok {
			columns = append(columns, v)
		}
	}
	return columns
}

// ArrayPop array_pop()
// Pop the element off the end of slice
func ArrayPop(s *[]interface{}) interface{} {
	if len(*s) == 0 {
		return nil
	}
	ep := len(*s) - 1
	e := (*s)[ep]
	*s = (*s)[:ep]
	return e
}

// ArrayUnshift array_unshift()
// Prepend one or more elements to the beginning of a slice
func ArrayUnshift(s *[]interface{}, elements ...interface{}) int {
	*s = append(elements, *s...)
	return len(*s)
}

// ArrayShift array_shift()
// Shift an element off the beginning of slice
func ArrayShift(s *[]interface{}) interface{} {
	if len(*s) == 0 {
		return nil
	}
	f := (*s)[0]
	*s = (*s)[1:]
	return f
}

// ArrayPush array_push()
// Push one or more elements onto the end of slice
func ArrayPush(s *[]interface{}, elements ...interface{}) int {
	*s = append(*s, elements...)
	return len(*s)
}

//KeyArray key是否在数组中
func KeyArray(key string, args []string) bool {
	sort.Strings(args)
	//只要不存在i肯定是len(args)+1的值，否则就是其中的一个key序列值
	i := sort.SearchStrings(args, key)

	result := false
	if i < len(args) {
		if args[i] == key {
			result = true
		}
	}
	return result
}

// InArray 值是否在数组中
// func InArray(s interface{}, d map[string]string) bool {
func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		return false
	}
	return false
}

//Empty 是否为空
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// Explode 字符串拆分成数组
func Explode(delimiter, str string) []string {
	return strings.Split(str, delimiter)
}

// Implode 拆分数组成字符串
func Implode(list interface{}, seq string) string {
	listValue := reflect.Indirect(reflect.ValueOf(list))
	if listValue.Kind() != reflect.Slice {
		return ""
	}
	count := listValue.Len()
	listStr := make([]string, 0, count)
	for i := 0; i < count; i++ {
		v := listValue.Index(i)
		if str, err := GetValue(v); err == nil {
			listStr = append(listStr, str)
		}
	}
	return strings.Join(listStr, seq)
}

// GetValue 获取值
func GetValue(value reflect.Value) (res string, err error) {
	switch value.Kind() {
	case reflect.Ptr:
		res, err = GetValue(value.Elem())
	default:
		res = fmt.Sprint(value.Interface())
	}
	return
}

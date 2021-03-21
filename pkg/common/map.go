package common

import (
	"reflect"
)

//MapStringMerge 合并多个map[string] 键相同时 向前覆盖
func MapStringMerge(dataMap ...map[string]interface{}) map[string]interface{} {
	var result = make(map[string]interface{})
	for i := range dataMap {
		for k := range dataMap[i] {
			if _, ok := result[k]; !ok {
				result[k] = dataMap[i][k]
			}
		}
	}
	return result
}

//MapIntMerge 合并多个map[int] 键相同时 叠加
func MapIntMerge(dataMap ...map[int]interface{}) []interface{} {
	var result = []interface{}{}
	for i := range dataMap {
		for k := range dataMap[i] {
			result = append(result, dataMap[i][k])
		}
	}
	return result
}

//StructToMap struct转换到map
func StructToMap(obj interface{}) (data map[string]interface{}) {

	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	// 获取到obj对应的类别
	if val.Kind() != reflect.Struct {
		return data
	}

	data = make(map[string]interface{})
	num := val.NumField()
	// 遍历结构体的所有字段
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i)
		val := val.Field(i)
		data[tagVal.Name] = val.Interface()
	}
	return data
}

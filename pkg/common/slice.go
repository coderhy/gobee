package common

import (
	"reflect"
	"strconv"
)

//SliceMapStringColumnMapString 获取[]map 中某一列信息指定键值 类似php array_column()
//data 要获取的数据集
//val 作为返回结果中的值
//key 作为返回结果中的键
//val = nil key != nil时 data[i][key]作为返回数据的主键，值为data[i]
//val != nil key != nil时 data[i][key]作为返回数据的主键，值为data[val]
//示例：res:=SliceMapStringColumnMapString(data,"username","userid")
func SliceMapStringColumnMapString(data []map[string]interface{}, val string, key string) (result map[string]interface{}) {
	if val != "nil" && key != "nil" {
		result := map[string]interface{}{}
		for k := range data {
			kVal, keyOk := data[k][key]
			vVal, valOk := data[k][val]
			if keyOk && valOk {
				switch reflect.TypeOf(kVal).String() {
				case "string":
					result[kVal.(string)] = vVal
				case "int":
					result[strconv.Itoa(kVal.(int))] = vVal
				case "int64":
					result[strconv.Itoa(int(kVal.(int64)))] = vVal
				default:
					continue
				}
			}
		}
		return result
	} else if val == "nil" && key != "nil" {
		result := map[string]interface{}{}
		for k := range data {
			if kVal, keyOk := data[k][key]; keyOk {
				switch reflect.TypeOf(kVal).String() {
				case "string":
					result[kVal.(string)] = data[k]
				case "int":
					result[strconv.Itoa(kVal.(int))] = data[k]
				case "int64":
					result[strconv.Itoa(int(kVal.(int64)))] = data[k]
				default:
					continue
				}
			}
		}
		return result
	}
	return result
}

//SliceMapStringColumnSlice 获取[]map 中指定列信息 类似php array_column()
//data 要获取的数据集
//val 作为返回结果中值的键名
//示例：res:=SliceMapStringColumnSlice(data,"userid")
func SliceMapStringColumnSlice(data []map[string]interface{}, val string) (result []interface{}) {
	for k := range data {
		if vVal, valOk := data[k][val]; valOk {
			result = append(result, vVal)
		}
	}
	return result
}

//SliceMapStringColumnSliceString 获取[]map 中指定列信息 类似php array_column()
//data 要获取的数据集
//val 作为返回结果中值的键名
//示例：res:=SliceMapStringColumnSliceString(data,"userid")
func SliceMapStringColumnSliceString(data []map[string]interface{}, val string) (result []string) {
	for k := range data {
		if vVal, valOk := data[k][val]; valOk {
			switch reflect.TypeOf(vVal).String() {
			case "string":
				result = append(result, vVal.(string))
			case "int":
				result = append(result, strconv.Itoa(vVal.(int)))
			case "int64":
				result = append(result, strconv.Itoa(int(vVal.(int64))))
			default:
				continue
			}
		}
	}
	return result
}

//SliceMapStringColumnSliceInt 获取[]map 中指定列信息 类似php array_column()
//data 要获取的数据集
//val 作为返回结果中值的键名
//示例：res:=SliceMapStringColumnSliceInt(data,"userid")
func SliceMapStringColumnSliceInt(data []map[string]interface{}, val string) (result []int) {
	for k := range data {
		if vVal, valOk := data[k][val]; valOk {
			switch reflect.TypeOf(vVal).String() {
			case "string":
				v, _ := strconv.Atoi(vVal.(string))
				result = append(result, v)
			case "int":
				result = append(result, vVal.(int))
			case "int64":
				result = append(result, int(vVal.(int64)))
			default:
				continue
			}
		}
	}
	return result
}

//SliceMapStringColumnSliceInt64 获取[]map 中指定列信息 类似php array_column()
//data 要获取的数据集
//val 作为返回结果中值的键名
//示例：res:=SliceMapStringColumnSliceInt64(data,"userid")
func SliceMapStringColumnSliceInt64(data []map[string]interface{}, val string) (result []int64) {
	for k := range data {
		if vVal, valOk := data[k][val]; valOk {
			switch reflect.TypeOf(vVal).String() {
			case "string":
				v, _ := strconv.Atoi(vVal.(string))
				result = append(result, int64(v))
			case "int":
				result = append(result, int64(vVal.(int)))
			case "int64":
				result = append(result, vVal.(int64))
			default:
				continue
			}
		}
	}
	return result
}

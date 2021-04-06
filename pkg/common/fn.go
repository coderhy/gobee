package common

import (
	"strconv"
	"unsafe"
)

// If 三元表达式、三目运算 a > b ? c : d
// 用例 result : = If( 3 > 4,5,9).(int)
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// FormatDatabaseType 格式化数据类型转换
func FormatDatabaseType(val []byte, databaseType string) interface{} {
	var sVal interface{}
	switch databaseType {
	case "BIGINT":
		sVal, _ = strconv.ParseInt(string(val), 10, 64)
		// sVal, _ = strconv.Atoi(string(val)) //兼容为int
	case "INT", "TINYINT", "SMALLINT", "MEDIUMINT":
		sVal, _ = strconv.Atoi(string(val))
	case "FLOAT", "DECIMAL":
		sVal, _ = strconv.ParseFloat(string(val), 64)
	default:
		sVal = string(val)
	}
	return sVal

}

// IsLittleEndian 判断系统是大端存储还是小端存储
func IsLittleEndian() bool {
	var i int32 = 0x01020304
	// 下面这两句是为了将int32类型的指针转换为byte类型的指针
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb // 取得pb位置对应的值

	// 由于b是byte类型的,最多保存8位,那么只能取得开始的8位
	// 小端: 04 (03 02 01)
	// 大端: 01 (02 03 04)
	return (b == 0x04)
}

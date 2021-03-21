package common

import (
	"bytes"
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

// Trim trim()
func Trim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimSpace(str)
	}
	return strings.Trim(str, characterMask[0])
}

// Strtoupper 转换成大写
func Strtoupper(str string) string {
	return strings.ToUpper(str)
}

// Strtolower 换成小写
func Strtolower(str string) string {
	return strings.ToLower(str)
}

// Ucfirst 首字符大写
func Ucfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}
	return ""
}

// Lcfirst 首字符小写
func Lcfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}
	return ""
}

//StrReplace 字符串替换
func StrReplace(dataOne []string, dataTwo []string, sourceData string) string {
	for i, toReplace := range dataOne {
		r := strings.NewReplacer(toReplace, dataTwo[i])
		sourceData = r.Replace(sourceData)
	}
	return sourceData
}

// Addslashes addslashes()
func Addslashes(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// Stripslashes stripslashes()
func Stripslashes(str string) string {
	var buf bytes.Buffer
	l, skip := len(str), false
	for i, char := range str {
		if skip {
			skip = false
		} else if char == '\\' {
			if i+1 < l && str[i+1] == '\\' {
				skip = true
			}
			continue
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

//UnsafeToString 避免内存copy方式
func UnsafeToString(bytes []byte) *string {
	hdr := &reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&bytes[0])),
		Len:  len(bytes),
	}
	return (*string)(unsafe.Pointer(hdr))
}

// UnescapeJSONMarshal 反转义处理
func UnescapeJSONMarshal(jsonRaw interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	//带有缩进的格式化
	encoder.SetIndent("", "  ")
	err := encoder.Encode(jsonRaw)
	return buffer.Bytes(), err
}

// StringToByte 字符串转[]byte
func StringToByte(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// ByteToString []byte转string
func ByteToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: bh.Data,
		Len:  bh.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))
}

//GetTaskConsoleOtherArgs 获取任务调度其他命令行参数为map
func GetTaskConsoleOtherArgs(params []string) map[string]string {
	ret := map[string]string{}
	if !Empty(params) {
		for k := range params {
			temp := strings.Split(params[k], "=")
			if len(temp) == 2 {
				ret[temp[0]] = temp[1]
			}
		}
	}
	return ret
}

//过滤特殊字符
func FilterSpecialWord(text string) string {
	if Empty(text) {
		return ""
	}
	result := strings.Replace(text, "`", "", -1) //对字符`特殊处理，因为pattern中无法拼出来
	pattern := `(\/|~|!|@|#|\$|%|\^|&|\*|\(|\)|\_|\+|\{|\}|:|<|>|\?|\[|\]|,|\.|;|'|"|-|=|\\)`
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}

	result = reg.ReplaceAllString(result, "") //将特殊字符替换成空

	return result
}

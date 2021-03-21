package common

import (
	"bytes"
)

// StringToBufferMerge 字符串拼接合并
// 用例：StringToBufferMerge("abc","哈!","=") 返回 abc=哈！
func StringToBufferMerge(start, end string, delimiter string) string {
	var buffer bytes.Buffer
	buffer.WriteString(start)
	buffer.WriteString(delimiter)
	buffer.WriteString(end)
	return buffer.String()
}

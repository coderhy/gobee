package e

//用户错误码 10开头 后三位自定义
const (
	ERROR_NOT_EXIST_USER = 100001
)

//UserMsg 错误码信息
var UserMsg = map[int]string{
	ERROR_NOT_EXIST_USER: "不存在此用户",
}

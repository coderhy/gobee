package e

//公共错误码 0～999
const (
	SUCCESS      = 0
	ERROR        = -1
	ERROR_SYSTEM = 999
)

//CommonMsg 错误码信息
var CommonMsg = map[int]string{
	SUCCESS:      "成功",
	ERROR:        "失败",
	ERROR_SYSTEM: "系统错误", //例如添加失败、修改失败 等错误操作返回系统错误
}

//MsgMap 错误码字典
var MsgMap = []map[int]string{
	CommonMsg,
	UserMsg,
}

//GetMsg 返回错误码
func GetMsg(code int) string {
	for k := range MsgMap {
		msg, ok := MsgMap[k][code]
		if ok {
			return msg
		}
	}
	return CommonMsg[ERROR]
}

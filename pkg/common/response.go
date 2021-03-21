package common

import (
	"encoding/json"
)

//Response 返回json
//@param int code 信息代码
//@param string msg 信息内容
//@param data 数据
func Response(code int, msg string, data interface{}) string {

	returnMap := map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	}

	result, _ := json.Marshal(returnMap)
	return string(result)
}

/*
ReturnMapMsg 返回map
@param code string 信息代码
@param msg string 信息内容
@param data []string 数据
*/
func ReturnMapMsg(code int, msg string, data interface{}) interface{} {

	returnMap := make(map[string]interface{})
	returnMap["code"] = code
	returnMap["msg"] = msg
	returnMap["data"] = data //[]string{"dev", "test", "pro"}
	return returnMap
}

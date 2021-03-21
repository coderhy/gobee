package common

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

// Md5 md5加密
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

// HmacSha256 加密哈希算法
func HmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

/*
BuildSign 加密字符生成（适用键值对相同，顺序不同时生成相同加密字符）
* @param args
* @return string
*/
func BuildSign(args interface{}) string {

	dataType := reflect.TypeOf(args).String()
	// if dataType == "string" {
	// 	return args.(string)
	// }
	validParams := []string{}
	switch dataType {
	case "string":
		return args.(string)
	case "map[string]string":
		for k, v := range args.(map[string]string) {
			if len(v) > 0 {
				validParams = append(validParams, k+"="+v)
			}
		}
	case "map[string]interface {}":
		for k, v := range args.(map[string]interface{}) {
			switch reflect.TypeOf(v).String() {
			case "[]interface {}":
				jsonData, _ := json.Marshal(v)
				validParams = append(validParams, k+"="+string(jsonData))
			case "[]string":
				jsonData, _ := json.Marshal(v)
				validParams = append(validParams, k+"="+string(jsonData))
			case "string":
				validParams = append(validParams, k+"="+v.(string))
			case "int":
				validParams = append(validParams, k+"="+strconv.Itoa(v.(int)))
			case "float64":
				temp := strconv.FormatFloat(v.(float64), 'f', -1, 64) //float64
				validParams = append(validParams, k+"="+temp)
			}
		}
	}
	if len(validParams) == 0 {
		return ""
	}

	sort.Strings(validParams)
	return Md5(Implode(validParams, "&"))
}

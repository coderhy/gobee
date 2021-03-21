package common

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// DingDingPush 钉钉推送 https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq
func DingDingPush(accessToken string, secret string, content string) (string, error) {

	if len(accessToken) == 0 || len(secret) == 0 {
		log.Println("dingDingPush error: accessToken or secret is empty")
		return "", errors.New("dingDingPush error: accessToken or secret is empty")
	}
	URL := "https://oapi.dingtalk.com/robot/send?access_token=" + accessToken

	timestamp := time.Now().UnixNano() / 1e6 //毫秒
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	sign := HmacSha256(stringToSign, secret)
	URL = fmt.Sprintf("%s&timestamp=%d&sign=%s", URL, timestamp, sign)

	//设置post数据
	// message := "推送警报1111 \n 日志推送2222 "
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"title":   "rebot",
			"content": content,
		},
		"at": map[string]interface{}{
			"isAtAll": true,
		},
	}
	// params, _ := json.Marshal(data)
	params, _ := UnescapeJSONMarshal(data)
	body := bytes.NewBuffer(params)
	contentType := "application/json;charset=utf-8"
	resp, err := http.Post(URL, contentType, body)

	if err != nil {
		// log.Println("Post failed:", err)
		return "", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read failed:", result, err)
		// return result, err
		return "", err

	}
	return string(result), err
}

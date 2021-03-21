package common

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
)

// CreateToken 生成token
func CreateToken(args map[string]interface{}, secret string) (string, error) {

	if len(args) == 0 {
		return "", errors.New("数据为空")
	}

	data := jwt.MapClaims{}
	for key, val := range args {
		data[key] = val
	}
	encryption := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	// encryption := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"uuid":   uuid,
	// 	"expire": time.Now().Unix(), //time.Now().Add(time.Minute * 15).Unix(),
	// })
	token, err := encryption.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParseToken 解析token
func ParseToken(token string, secret string) (map[string]interface{}, error) {
	if token == "" {
		return make(map[string]interface{}), errors.New("token不能为空")
	}

	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return make(map[string]interface{}), err
	}
	// return claim.Claims.(jwt.MapClaims)["uid"].(string), nil
	return claim.Claims.(jwt.MapClaims), nil
}

// Cors 定义全局的CORS中间件(cross-origin sharing standard 跨源资源共享)
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

//限流文案 demo
// params := map[string]interface{}{
// 	"username": "",
// 	"age":      "",
// }
// limitData, _ := json.Marshal(params)
// result := string(limitData)
// router.GET("/", utils.CurrentLimiting(1, result, "json"), V1.Index)

// CurrentLimiting rate-limit 限流中间件
func CurrentLimiting(limit float64, result string, contentType string) gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(limit, nil)
	if contentType == "json" {
		lmt.SetMessageContentType("application/json") //内容类型
		// data, _ := json.Marshal(result)
		// result = string(data)
	}
	// lmt.SetMessage("服务繁忙，请稍后再试...")
	lmt.SetMessage(result)

	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}

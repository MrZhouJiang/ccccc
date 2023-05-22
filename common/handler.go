package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	CODOON_REQUEST_ID   = "codoon_request_id"
	CODOON_SERVICE_CODE = "codoon_service_code"
	CODOON_CLIENT_IP    = "codoon_client_ip"
)

type EndpointFunc func(*gin.Context) (interface{}, error)

//定义返回
func homeEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS"})
}

/*func wrapHandler(handler EndpointFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理跨域请求
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		Cors(c)

		//从header 获取token,
		var token = c.GetHeader("token")
		//解析token,如果token不存在，或者 解析错误,则返回 没权限
		userId, err := util.TokenHandle(token)
		if err != nil {
			fmt.Errorf(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"code": "ERROR", "message": "token is invalid"})
			return
		}
		//解析后 userId 放入上下文中,后面的service服务，可以获取该用户信息
		c.Set("userId", userId)
		data, err := handler(c)
		if err != nil {
			var systemErr *common.Error
			if errors.As(err, &systemErr) {
				c.JSON(systemErr.Status, common.NewErrorResult(systemErr.Code, systemErr.Msg))
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"code": "ERROR", "message": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "data": data})
	}
}
*/
func Cors(context *gin.Context) {
	method := context.Request.Method
	// 必须，接受指定域的请求，可以使用*不加以限制，但不安全
	//context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Origin", context.GetHeader("Origin"))
	fmt.Println(context.GetHeader("Origin"))
	// 必须，设置服务器支持的所有跨域请求的方法
	context.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	// 服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段
	context.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
	// 可选，设置XMLHttpRequest的响应对象能拿到的额外字段
	context.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
	// 可选，是否允许后续请求携带认证信息Cookir，该值只能是true，不需要则不设置
	context.Header("Access-Control-Allow-Credentials", "true")
	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		context.AbortWithStatus(http.StatusNoContent)
		return
	}
	context.Next()

}

// gin请求日志
func GinLogger() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		fmt.Println(latency, clientIP, method, statusCode)

	}
}

func CommonParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get request_id, and ensure set request_id into header
		values := c.Request.URL.Query()
		reqId := values.Get(CODOON_REQUEST_ID)
		if reqId == "" {
			reqId = c.Request.Header.Get(CODOON_REQUEST_ID)
		} else {
			c.Request.Header.Set(CODOON_REQUEST_ID, reqId)
		}

		c.Writer.Header().Set("RequestId", reqId)
		c.Writer.Header().Set("ClientIP", c.ClientIP())
	}
}

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get request_id, and ensure set request_id into header
		values := c.Request.URL.Query()
		reqId := values.Get(CODOON_REQUEST_ID)
		if reqId == "" {
			reqId = c.Request.Header.Get(CODOON_REQUEST_ID)
		} else {
			c.Request.Header.Set(CODOON_REQUEST_ID, reqId)
		}

		c.Writer.Header().Set("RequestId", reqId)
		c.Writer.Header().Set("ClientIP", c.ClientIP())
	}
}

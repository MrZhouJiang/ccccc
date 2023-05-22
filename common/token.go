package common

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func getParamToken(c *gin.Context) (string, error) {
	var token string
	//  不再从form中解析token会导致body丢失
	token = c.DefaultQuery("bearer_token", "")

	if len(token) == 0 {
		return "", errors.New("empty token from param bearer_token")
	}
	return token, nil
}

func TokenChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			reqId := c.Request.Header.Get(CODOON_REQUEST_ID)
			fmt.Errorf("token checker: %s %d", reqId, time.Now().Sub(start)/1e6)
		}()

		c.Request.Header.Del(USER_ID)
		// 20190326 增加白名单接口也尝试获取user_id(如果header中包含有效token)
		// 是否强制验证Token
		should_verify_token := shouldVerifyToken(c)
		token, err := getHeaderToken(c)
		if err != nil {
			// 直解析Query参数,form造成body丢失
			token, err = getParamToken(c)
			if err != nil {
				fmt.Errorf("[token:%s] getToken failed:%v", token, err)
				if should_verify_token {
					abortInvalidToken(c, gin.H{
						"error":             "invalid_request",
						"error_code":        1001,
						"error_description": "No authentication credentials provided.",
					})
				}
				return
			}
		}
		if userId, err := verifyToken(c, token); err != nil {
			fmt.Errorf("[token:%s] verify error:%+v", token, err)
			if should_verify_token {
				abortInvalidToken(c, gin.H{
					"error":             "invalid_token",
					"error_code":        1002,
					"error_description": "Verify token failed",
				})
			}
		} else {
			// c.Request.Header.Set(USER_ID, userId)
			// be compatible with go 1.5.2 when migrate to 1.7.3
			c.Request.Header[USER_ID] = []string{userId}
			fmt.Errorf("TokenChecker OK, setUserId:[%s]", userId)
		}
	}

}

func Core() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}

}
func shouldVerifyToken(c *gin.Context) bool {
	path := c.Request.URL.Path
	white := isWhiteUri(path)
	fmt.Errorf("isWhiteUri(%s):%v", path, white)
	return !white
}

func isWhiteUri(path string) bool {

	return true
}

const (
	AUTH_BEARER = "bearer"
	USER_ID     = "user_id"
)

func getHeaderToken(c *gin.Context) (string, error) {
	authStr := c.Request.Header.Get("Authorization")
	auth := strings.Split(authStr, " ")
	if len(auth) < 2 {
		return "", errors.New(fmt.Sprintf("invalid authStr [%s]", authStr))
	}
	if strings.ToLower(auth[0]) != AUTH_BEARER {
		return "", errors.New(fmt.Sprintf("auth type [%s] not supported", auth[0]))
	}
	token := strings.Trim(auth[len(auth)-1], " ")
	if len(token) == 0 {
		return "", errors.New("empty token")
	}

	return token, nil
}

func abortInvalidToken(c *gin.Context, rsp gin.H) {
	reqId := c.Request.URL.Query().Get(CODOON_REQUEST_ID)
	fmt.Errorf("abortInvalidToken [req_id:%s] %s %s %#v, rsp:%+v",
		reqId, c.Request.Method, c.Request.URL.Path, c.Request.Header, rsp)
	c.Writer.Header().Add("WWW-Authenticate", `OAuth realm="users"`) // avoid android throw exceptions
	c.JSON(http.StatusUnauthorized, rsp)
	c.Abort()
}

func verifyToken(c *gin.Context, token string) (string, error) {
	return "", nil
}

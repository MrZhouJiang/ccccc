package util

import (
	"ccccc/stringutil"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

var mySigningKey = []byte("asfasfdafasdfdasfa.")

//TokenHandle 解析token

func TokenHandle(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	data := token.Claims.(jwt.MapClaims)["data"]
	var tokenData = new(TokenData)
	if data != nil {
		var info = data.(string)
		err = json.Unmarshal([]byte(info), &tokenData)
	}

	return tokenData.UserId, err

}

//CreateToken is created token
func CreateToken(data TokenData) (string, error) {
	dataByte, err := json.Marshal(data)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": string(dataByte),
		"exp":  time.Now().Unix() + 1000*5,
		"iss":  "ibc_business",
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Println("加密后的token字符串", tokenString)
	return tokenString, nil
}

type TokenData struct {
	UserId   string
	Age      int32
	NickName string
	Name     string
	Phone    string
}

func ReturnCompFunc(c *gin.Context, obj interface{}) {
	GinRsp(c, 200, obj)
}
func GetRequestData(c *gin.Context) string {
	var requestData string
	method := c.Request.Method
	if method == "GET" || method == "DELETE" {
		requestData = c.Request.RequestURI
	} else {
		c.Request.ParseForm()
		requestData = fmt.Sprintf("%s [%s]", c.Request.RequestURI, c.Request.Form.Encode())
	}
	return requestData
}
func GinRsp(c *gin.Context, statusCode int, obj interface{}) {
	requestData := GetRequestData(c)
	objData := fmt.Sprintf("%+v", obj)

	clientIP := c.ClientIP()

	reqId := c.Request.Header.Get("codoon_request_id")
	userId := c.Request.Header.Get("user_id")
	fmt.Errorf("[GIN-RSP] %s [ip:%s] [req_id:%s] [user_id:%s] [rsp:%s]",
		stringutil.Cuts(requestData, 1024),
		clientIP,
		reqId,
		userId,
		stringutil.Cuts(objData, 1024),
	)
	c.JSON(statusCode, obj)
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {

	x1, x2, t := 0, 0, 0
	temp := []int{}

	for l1 != nil || l2 != nil {
		if l1 != nil {
			x1 = l1.Val
			l1 = l1.Next
		} else {
			x1 = 0
		}
		if l2 != nil {
			x2 = l2.Val
			l2 = l2.Next
		} else {
			x2 = 0
		}

		temp = append(temp, (x1+x2+t)%10)

		if x1+x2+t >= 10 {
			t = 1
		} else {
			t = 0
		}
	}
	if t == 1 {
		temp = append(temp, t)
	}

	h := &ListNode{Next: nil}
	th := h
	for i := 0; i < len(temp); i++ {
		nm := ListNode{Val: temp[i]}
		th.Next = &nm
		th = th.Next

	}

	return h.Next

}

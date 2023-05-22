package stringutil

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math"
	"strconv"
	"strings"
)

//截取前n个
func Cuts(s string, n int) string {
	if len(s) > n {
		return s[:n]
	} else {
		return s
	}
}

// 截取前 n 个自然长度字符
func CutsRune(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	} else {
		return s
	}
}

//字符串是否在slice中
func StrIn(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

// 字符串s是否包含slice中任一前缀
func StrPrefixAny(s string, arr []string) bool {
	for _, a := range arr {
		if strings.Index(s, a) == 0 {
			return true
		}
	}
	return false
}

//两个字符串中间值
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		return ""
	}

	n += len(start)

	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

func EmptyDefault(s string, d string) string {
	if s == "" {
		return d
	} else {
		return s
	}
}

func IntArrayToString(intarr []int) string {
	var buffer bytes.Buffer

	for index, value := range intarr {
		str := strconv.Itoa(value)
		if index != 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(str)
	}

	return buffer.String()
}

//字节混淆
func RotateStr(str string) string {
	bs := []byte(str)
	for idx := range bs {
		if (idx+1)%2 == 0 && idx > 0 {
			b := bs[idx-1]
			bs[idx-1] = bs[idx]
			bs[idx] = b
		}
	}

	return string(bs)
}

func RandStr(length int) string {
	if length == 0 {
		return ""
	}

	newLen := math.Ceil(float64(length) / 2)
	buf := make([]byte, int(newLen))
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Printf("gen rand str err: %s\n", err.Error())
		return ""
	}

	out := fmt.Sprintf("%x", buf)
	return out[:int(length)]
}

//ConvertNick 将昵称中间部分按*隐藏起来
func ConvertNick(word string) string {
	var nick string
	rune_nick := []rune(word)
	if len(rune_nick) > 2 {
		nick = string(rune_nick[0])
		for i := 1; i <= len(rune_nick)-2; i++ {
			nick += "*"
		}
		nick += string(rune_nick[len(rune_nick)-1])
	} else if len(rune_nick) == 0 {

	} else if len(rune_nick) <= 2 {
		nick = string(rune_nick[0]) + "*"
	}
	return nick
}

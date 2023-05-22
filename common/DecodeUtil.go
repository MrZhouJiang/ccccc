package common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// 16位随机字符串
var bytes1 = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// 密钥，实际应用中应该从环境变量或文件中获取
const MySecret string = "abc&1*~#^2^#s0^=)^^7%b34"

// Base64编码和解码方法
func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// 加密方法
func encrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}

	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes1)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	return Encode(cipherText), nil
}

//加密
func Encrypt(text string) (string, error) {
	return encrypt(text, MySecret)

}

// 解密方法
func decrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}

	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes1)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}

//解密
func Decrypt(text string) (string, error) {

	return decrypt(text, MySecret)
}

/*func main() {
	StringToEncrypt := "123456"

	// 对StringToEncrypt变量值进行加密
	encText, err := Encrypt(StringToEncrypt, MySecret)
	if err != nil {
		fmt.Println("error encrypting your classified text: ", err)
	}
	fmt.Println(encText)

	// 对密文进行解密
	decText, err := Decrypt(encText, MySecret)
	if err != nil {
		fmt.Println("error decrypting your encrypted text: ", err)
	}

	fmt.Println(decText)
}
*/

package str

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// StringTool 字符串工具
type StringTool struct {
	string string
}

// NetStringTool 实例化工具类
func NetStringTool(str string) *StringTool {
	return &StringTool{string: str}
}

// GetBetween 获取文本中间字符
func (stringTool *StringTool) GetBetween(start string, end string) string {
	n := strings.Index(stringTool.string, start)
	if n == -1 {
		n = 0
	}
	stringTool.string = string([]byte(stringTool.string)[n:])
	m := strings.Index(stringTool.string, end)
	if m == -1 {
		m = len(stringTool.string)
	}
	stringTool.string = string([]byte(stringTool.string)[len(start):m])
	return stringTool.string
}

// Md5 获取字符串md5值
func (stringTool *StringTool) Md5() string {
	h := md5.New()
	h.Write([]byte(stringTool.string))
	return hex.EncodeToString(h.Sum(nil))
}

// AesCbcEncrypt aes cbc加密
func (stringTool *StringTool) AesCbcEncrypt(key string) string {
	plainByte := []byte(stringTool.string)
	keyByte := []byte(key)
	if len(plainByte)%aes.BlockSize != 0 {
		return ""
	}
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return ""
	}
	cipherByte := make([]byte, aes.BlockSize+len(plainByte))
	iv := cipherByte[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherByte[aes.BlockSize:], plainByte)
	return fmt.Sprintf("%x\n", cipherByte)
}

// AesCbcDecrypt aes cbc解密
func (stringTool *StringTool) AesCbcDecrypt(key string) string {
	cipherByte, _ := hex.DecodeString(stringTool.string)
	keyByte := []byte(key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return ""
	}
	if len(cipherByte) < aes.BlockSize {
		return ""
	}
	iv := cipherByte[:aes.BlockSize]
	cipherByte = cipherByte[aes.BlockSize:]
	if len(cipherByte)%aes.BlockSize != 0 {
		return ""
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherByte, cipherByte)
	return string(cipherByte[:])
}

func (stringTool *StringTool) AesEncrypt(key string) string {
	origData := []byte(stringTool.string)
	k := []byte(key)
	block, _ := aes.NewCipher(k)
	blockSize := aes.BlockSize
	origData = stringTool.pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}

func (stringTool *StringTool) AesDecrypt(key string) string {
	crytedByte, _ := base64.StdEncoding.DecodeString(stringTool.string)
	k := []byte(key)
	block, _ := aes.NewCipher(k)
	blockSize := aes.BlockSize
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	orig := make([]byte, len(crytedByte))
	blockMode.CryptBlocks(orig, crytedByte)
	orig = stringTool.pKCS7UnPadding(orig)
	return string(orig)
}

func (StringTool) pKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func (StringTool) pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

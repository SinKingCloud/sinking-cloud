package str

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	rand2 "math/rand"
	"strings"

	"github.com/axgle/mahonia"
	"golang.org/x/crypto/bcrypt"
)

// StringTool 字符串工具
type StringTool struct {
	string string
}

// NewStringTool 实例化工具类
func NewStringTool(str string) *StringTool {
	return &StringTool{string: str}
}

// NewStrTool 实例化工具类
func NewStrTool() *StringTool {
	return &StringTool{}
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

func (*StringTool) pKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func (*StringTool) pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// GetPassword 获取密码
func (stringTool *StringTool) GetPassword() string {
	hash, err := bcrypt.GenerateFromPassword([]byte(stringTool.string), bcrypt.DefaultCost) //加密处理
	if err != nil {
		return ""
	}
	return string(hash)
}

// CheckPassword 比对密码
func (stringTool *StringTool) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(stringTool.string), []byte(password)) //验证（对比）
	if err != nil {
		return false
	}
	return true
}

// GenValidateCode 生成指定位数验证码
func (stringTool *StringTool) GenValidateCode(width int) string {
	numeric := [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, err := fmt.Fprintf(&sb, "%d", numeric[rand2.Intn(r)])
		if err != nil {
			return ""
		}
	}
	return sb.String()
}

func (stringTool *StringTool) RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand2.Int63n(max-min) + min
}

// Show 显示首尾数据字符个数
func (stringTool *StringTool) Show(left int, right int, str string) string {
	if stringTool.string == "" {
		return ""
	}
	nameRune := []rune(stringTool.string)
	lens := len(nameRune)
	if lens-right < 0 || left > lens {
		return ""
	}
	leftStr := nameRune[0:left]
	rightStr := nameRune[lens-right : lens]
	if len(leftStr)+len(rightStr) >= lens {
		return string(leftStr) + string(rightStr)
	}
	temp := ""
	for i := 0; i < lens-len(leftStr)-len(rightStr); i++ {
		temp += str
	}
	return string(leftStr) + temp + string(rightStr)
}

// ConvertToString 编码转换
func (stringTool *StringTool) ConvertToString(srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(stringTool.string)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// ToMap 转map
func (stringTool *StringTool) ToMap() map[string]string {
	var d map[string]string
	e := json.Unmarshal([]byte(stringTool.string), &d)
	if e == nil {
		return d
	}
	return nil
}

// ToDevice 解析User-Agent，返回设备型号和大类
func (stringTool *StringTool) ToDevice() (name string, dtype string) {
	if stringTool.string == "" {
		return "Unknown", "Unknown"
	}
	uaLower := strings.ToLower(stringTool.string)

	// iOS
	if strings.Contains(uaLower, "iphone") {
		ident := ""
		for _, token := range strings.Fields(uaLower) {
			if strings.HasPrefix(token, "iphone") && strings.Count(token, ",") == 1 {
				if v, ok := map[string]string{
					"iphone16,1": "15 Pro", "iphone16,2": "15 Pro Max",
					"iphone15,4": "15", "iphone15,5": "15 Plus",
					"iphone14,7": "14", "iphone14,8": "14 Plus",
					"iphone15,2": "14 Pro", "iphone15,3": "14 Pro Max",
					"iphone14,4": "13 mini", "iphone14,5": "13",
					"iphone14,2": "13 Pro", "iphone14,3": "13 Pro Max",
					"iphone13,1": "12 mini", "iphone13,2": "12",
					"iphone13,3": "12 Pro", "iphone13,4": "12 Pro Max",
					"iphone12,8": "SE (2nd)",
					"iphone12,1": "11", "iphone12,3": "11 Pro", "iphone12,5": "11 Pro Max",
					"iphone11,8": "XR", "iphone11,2": "XS", "iphone11,4": "XS Max",
					"iphone11,6": "XS Max", "iphone11,15": "11", "iphone11,16": "11 Pro", "iphone11,17": "11 Pro Max",
					"iphone10,3": "X", "iphone10,6": "X",
					"iphone10,1": "8", "iphone10,4": "8", "iphone10,2": "8 Plus", "iphone10,5": "8 Plus",
					"iphone9,1": "7", "iphone9,3": "7", "iphone9,2": "7 Plus", "iphone9,4": "7 Plus",
					"iphone8,1": "6s", "iphone8,2": "6s Plus", "iphone8,4": "SE (1st)",
					"iphone7,1": "6 Plus", "iphone7,2": "6",
					"iphone6,1": "5s", "iphone6,2": "5s",
					"iphone5,4": "5c", "iphone5,3": "5c", "iphone5,1": "5", "iphone5,2": "5",
					"iphone4,1": "4S",
					"iphone3,1": "4", "iphone3,2": "4", "iphone3,3": "4",
				}[token]; ok {
					ident = v
					break
				}
			}
		}
		return "iPhone " + ident, "Mobile"
	}
	if strings.Contains(uaLower, "ipad") {
		ident := ""
		for _, token := range strings.Fields(uaLower) {
			if strings.HasPrefix(token, "ipad") && strings.Count(token, ",") == 1 {
				if v, ok := map[string]string{
					"ipad14,1": "10th", "ipad14,2": "10th",
					"ipad13,18": "10th", "ipad13,19": "10th",
					"ipad13,16": "Air (5th)", "ipad13,17": "Air (5th)",
					"ipad14,3": "Pro 11-inch (4th)", "ipad14,4": "Pro 11-inch (4th)",
					"ipad14,5": "Pro 12.9-inch (6th)", "ipad14,6": "Pro 12.9-inch (6th)",
					"ipad13,4": "Pro 11-inch (3rd)", "ipad13,5": "Pro 11-inch (3rd)",
					"ipad13,6": "Pro 11-inch (3rd)", "ipad13,7": "Pro 11-inch (3rd)",
					"ipad13,8": "Pro 12.9-inch (5th)", "ipad13,9": "Pro 12.9-inch (5th)",
					"ipad13,10": "Pro 12.9-inch (5th)", "ipad13,11": "Pro 12.9-inch (5th)",
					"ipad12,1": "Pro 12.9-inch (4th)", "ipad12,2": "Pro 12.9-inch (4th)",
					"ipad13,1": "Air (4th)", "ipad13,2": "Air (4th)",
					"ipad11,6": "Pro 11-inch", "ipad11,7": "Pro 11-inch",
					"ipad11,3": "Air (3rd)",
					"ipad6,7":  "Pro 12.9-inch", "ipad6,8": "Pro 12.9-inch",
					"ipad6,3": "Pro 9.7-inch", "ipad6,4": "Pro 9.7-inch",
					"ipad6,11": "9.7-inch", "ipad6,12": "9.7-inch",
					"ipad7,1": "Pro 12.9-inch (2nd)", "ipad7,2": "Pro 12.9-inch (2nd)",
					"ipad7,3": "Pro 10.5-inch", "ipad7,4": "Pro 10.5-inch",
					"ipad7,5": "9.7-inch", "ipad7,6": "9.7-inch",
					"ipad7,11": "10.2-inch", "ipad7,12": "10.2-inch",
					"ipad8,1": "Pro 11-inch", "ipad8,2": "Pro 11-inch",
					"ipad8,3": "Pro 11-inch", "ipad8,4": "Pro 11-inch",
					"ipad8,5": "Pro 12.9-inch (3rd)", "ipad8,6": "Pro 12.9-inch (3rd)",
					"ipad8,7": "Pro 12.9-inch (3rd)", "ipad8,8": "Pro 12.9-inch (3rd)",
					"ipad11,1": "Mini (5th)", "ipad11,2": "Mini (5th)", "ipad11,4": "Mini (5th)",
					"ipad5,1": "Mini (4th)", "ipad5,2": "Mini (4th)",
					"ipad5,3": "Air", "ipad5,4": "Air",
					"ipad2,1": "2", "ipad2,2": "2", "ipad2,3": "2", "ipad2,4": "2",
					"ipad2,5": "Mini", "ipad2,6": "Mini", "ipad2,7": "Mini",
					"ipad3,1": "3", "ipad3,2": "3", "ipad3,3": "3",
					"ipad3,4": "4", "ipad3,5": "4", "ipad3,6": "4",
				}[token]; ok {
					ident = v
					break
				}
			}
		}
		return "iPad " + ident, "Tablet"
	}
	if strings.Contains(uaLower, "ipod") {
		return "iPod", "Mobile"
	}

	// Android
	if strings.Contains(uaLower, "android") {
		brand := ""
		for kw, b := range map[string]string{
			"pixel": "Google Pixel", "google": "Google",
			"samsung": "Samsung", "xiaomi": "Xiaomi", "redmi": "Redmi", "poco": "POCO",
			"huawei": "Huawei", "honor": "Honor",
			"oppo": "OPPO", "vivo": "vivo", "realme": "realme", "oneplus": "OnePlus",
			"meizu": "Meizu", "nubia": "Nubia", "lenovo": "Lenovo", "motorola": "Motorola",
			"zte": "ZTE", "sony": "Sony", "nokia": "Nokia",
			"lg": "LG", "htc": "HTC", "asus": "ASUS",
			"tcl": "TCL", "alcatel": "Alcatel", "blackberry": "BlackBerry",
			"sharp": "Sharp", "fujitsu": "Fujitsu", "panasonic": "Panasonic",
			"infinix": "Infinix", "tecno": "TECNO", "itel": "itel",
		} {
			if strings.Contains(uaLower, kw) {
				brand = b
				break
			}
		}
		model := ""
		if idx := strings.Index(uaLower, "linux; android "); idx != -1 {
			rest := uaLower[idx+15:]
			semi := strings.Index(rest, ";")
			if semi == -1 || semi > 80 {
				semi = strings.Index(rest, ")")
			}
			if semi != -1 {
				model = strings.TrimSpace(rest[:semi])
				model = strings.SplitN(model, " build/", 2)[0]
				model = strings.Join(strings.Fields(model), " ")
			}
		}
		dt := "Tablet"
		if strings.Contains(uaLower, "mobile") {
			dt = "Mobile"
		}
		if model != "" {
			if brand != "" {
				return brand + " " + model, dt
			}
			return model, dt
		}
		return "Android " + dt, dt
	}

	// Windows
	if strings.Contains(uaLower, "windows phone") || strings.Contains(uaLower, "iemobile") || strings.Contains(uaLower, "wpdesktop") {
		return "Windows Phone", "Mobile"
	}
	if strings.Contains(uaLower, "windows") {
		return "Windows PC", "Desktop"
	}

	// Mac / Linux / Chrome OS
	if strings.Contains(uaLower, "macintosh") || strings.Contains(uaLower, "mac os x") {
		return "Mac", "Desktop"
	}
	if strings.Contains(uaLower, "cros") {
		return "Chrome OS", "Desktop"
	}
	if strings.Contains(uaLower, "linux") {
		return "Linux PC", "Desktop"
	}
	if strings.Contains(uaLower, "freebsd") {
		return "FreeBSD", "Desktop"
	}

	// 游戏主机
	if strings.Contains(uaLower, "playstation") {
		return "PlayStation", "Desktop"
	}
	if strings.Contains(uaLower, "xbox") {
		return "Xbox", "Desktop"
	}
	if strings.Contains(uaLower, "nintendo") || strings.Contains(uaLower, "wiiu") || strings.Contains(uaLower, "switch") {
		return "Nintendo", "Desktop"
	}

	// 智能电视
	if strings.Contains(uaLower, "smarttv") || strings.Contains(uaLower, "hbb-tv") || strings.Contains(uaLower, "googletv") {
		return "Smart TV", "Desktop"
	}

	// 爬虫
	if strings.Contains(uaLower, "bot") || strings.Contains(uaLower, "crawler") || strings.Contains(uaLower, "spider") ||
		strings.Contains(uaLower, "slurp") || strings.Contains(uaLower, "googlebot") || strings.Contains(uaLower, "bingbot") ||
		strings.Contains(uaLower, "yandex") || strings.Contains(uaLower, "baiduspider") {
		return "Bot", "Bot"
	}

	return "Unknown", "Unknown"
}

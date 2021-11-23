package jwt

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// MyClaims jwt载体
type MyClaims struct {
	User interface{}
	jwt.StandardClaims
}

// getKey 获取加密key
func getKey() []byte {
	return []byte(util.JwtTokenKey)
}

// GetToken 生成token
func GetToken(user interface{}) string {
	setClaim := MyClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(util.JwtTokenExpire * time.Second).Unix(), //过期时间（时间戳）
		},
	}
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, setClaim) //生成token
	token, err := reqClaim.SignedString(getKey())                   //转换为字符串
	if err != nil {
		return ""
	}
	return str.NetStringTool(token).AesCbcEncrypt(util.JwtTokenKey)
}

// CheckToken 验证token
func CheckToken(token string) *MyClaims {
	token = str.NetStringTool(token).AesCbcDecrypt(util.JwtTokenKey)
	setToken, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getKey(), nil
	})
	if err != nil {
		return nil
	}
	if key, success := setToken.Claims.(*MyClaims); setToken.Valid && success {
		return key
	} else {
		return nil
	}
}

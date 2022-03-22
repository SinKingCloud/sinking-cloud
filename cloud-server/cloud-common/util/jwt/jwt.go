package jwt

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/golang-jwt/jwt"
	"time"
)

// MyClaims jwt载体
type MyClaims struct {
	WebUser    *WebUser
	SystemUser *SystemUser
	jwt.StandardClaims
}

// SystemUser 系统用户表
type SystemUser struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	RoleId     int          `gorm:"column:role_id" json:"role_id"`
	LoginToken string       `gorm:"column:login_token" json:"login_token"`
	LoginIp    string       `gorm:"column:login_ip" json:"login_ip"`
	LoginTime  str.DateTime `gorm:"column:login_time" json:"login_time"`
	Status     int          `gorm:"column:status" json:"status"`
}

// WebUser 会员用户表
type WebUser struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	LoginToken string       `gorm:"column:login_token" json:"login_token"`
	LoginIp    string       `gorm:"column:login_ip" json:"login_ip"`
	LoginTime  str.DateTime `gorm:"column:login_time" json:"login_time"`
	Status     int          `gorm:"column:status" json:"status"`
	RealAuth   int          `gorm:"column:real_auth" json:"real_auth"`
}

// getKey 获取加密key
func getKey() []byte {
	return []byte(util.JwtTokenKey)
}

// GetLoginToken 获取token
func GetLoginToken(token string) map[string]string {
	temp := make(map[string]string)
	_ = json.Unmarshal([]byte(token), &temp)
	return temp
}

// GetToken 生成token user 用户信息
func GetToken(webUser *WebUser, systemUser *SystemUser) string {
	setClaim := MyClaims{
		WebUser:    webUser,
		SystemUser: systemUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(util.JwtTokenExpire * time.Second).Unix(), //过期时间（时间戳）
		},
	}
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, setClaim) //生成token
	token, err := reqClaim.SignedString(getKey())                   //转换为字符串
	if err != nil {
		return ""
	}
	return token
}

// CheckToken 验证token
func CheckToken(token string) *MyClaims {
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

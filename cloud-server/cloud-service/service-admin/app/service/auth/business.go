package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	constant2 "service-admin/app/constant"
	"service-admin/app/model"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	ip2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/ip"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Login 账号登录
func (s *service) Login(account string, password string, code string, userAgent string, ip string) (string, error) {
	var id int64
	if password != "" {
		u, err := s.repositoryUser.FindByIdentifier(account)
		if err != nil {
			return "", errors.New("系统繁忙,请稍后再试")
		}
		if u == nil || u.Id <= 0 {
			return "", errors.New("账号不存在")
		}
		if !str.NewStringTool(u.Password).CheckPassword(password) {
			return "", errors.New("密码错误")
		}
		id = u.Id
	} else {
		check := s.authService.CheckSms
		if !regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(account) {
			if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(account) {
				return "", errors.New("手机号或邮箱格式不正确")
			}
			check = s.authService.CheckMail
		}
		if !regexp.MustCompile(`^\d{6}$`).MatchString(code) {
			return "", errors.New("验证码格式不正确")
		}
		if !check(account, code) {
			return "", errors.New("验证码错误")
		}
		u, err := s.repositoryUser.FindByIdentifier(account)
		if err != nil {
			return "", errors.New("系统繁忙,请稍后再试")
		}
		if u == nil || u.Id <= 0 {
			return "", errors.New("用户不存在")
		}
		id = u.Id
	}
	if id <= 0 {
		return "", errors.New("获取用户信息错误")
	}
	return s.genLoginToken(id, userAgent, ip)
}

// Logout 注销登录
func (s *service) Logout(token string) error {
	key := jwt.CheckToken(token)
	if key == nil || key.User.Id <= 0 || key.User.Token == "" {
		return errors.New("未获取到登陆信息")
	}
	s.cache.Delete(constant2.CacheNameWithSession + key.User.Token)
	err := s.repositorySession.DeleteByToken(key.User.Token)
	if err != nil {
		return errors.New("注销登录失败")
	}
	return nil
}

// GetSessionInfo 获取session详情
func (s *service) GetSessionInfo(token string) *model.Session {
	cacheKey := constant2.CacheNameWithSession + token
	if data, exists := s.cache.Get(cacheKey); exists && data != "" {
		var session *model.Session
		if err := json.Unmarshal([]byte(data), &session); err == nil {
			return session
		}
	}
	session, err := s.repositorySession.FindByToken(token)
	if err != nil || session == nil || session.Id <= 0 {
		return nil
	}
	if data, err := json.Marshal(session); err == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant2.CacheTimeWithSession)
	}
	return session
}

// GetTokenInfo 根据token获取信息
func (s *service) GetTokenInfo(token string) (*jwt.User, error) {
	key := jwt.CheckToken(token)
	if key == nil || key.User == nil || key.User.Id <= 0 || key.User.Token == "" {
		return nil, errors.New("登陆超时,请重新登陆")
	}
	if s.GetSessionInfo(key.User.Token) == nil {
		return nil, errors.New("登陆超时,请重新登陆")
	}
	return key.User, nil
}

// GenLoginToken 创建会话
func (s *service) genLoginToken(userId int64, userAgent string, ip string) (string, error) {
	token := str.NewStringTool(str.GetSnowWorkIns().GetUuid()).Md5()
	expireTime := str.DateTime(time.Now().Add(constant.JwtExpireTime * time.Second))
	location := "未知"
	ipInfo, err := ip2.Query(ip)
	if err == nil && ipInfo != nil {
		location = strings.ReplaceAll(fmt.Sprintf("%s %s %s %s", ipInfo.Country, ipInfo.Province, ipInfo.City, ipInfo.Isp), "  ", "")
	}
	device, types := str.NewStringTool(userAgent).ToDevice()
	if err = s.repositorySession.Create(&model.Session{
		UserId:     userId,
		Token:      token,
		Device:     fmt.Sprintf("%s %s", device, types),
		Ip:         ip,
		Location:   location,
		ExpireTime: expireTime,
	}); err != nil {
		return "", err
	}
	jwtToken := jwt.GetToken(&jwt.User{
		Id:        userId,
		Token:     token,
		Ip:        ip,
		LoginTime: str.DateTime(time.Now()),
	}, constant.JwtExpireTime)
	return jwtToken, nil
}

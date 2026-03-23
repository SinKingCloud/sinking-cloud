package auth

import (
	"errors"
	"fmt"
	"regexp"
	constant2 "service-user/app/constant"
	"service-user/app/enum/real_auth"
	"service-user/app/enum/user"
	"service-user/app/model"
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
	if password == "" {
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
		if !regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(account) {
			return "", errors.New("手机号格式不正确")
		}
		if !regexp.MustCompile(`^\d{6}$`).MatchString(code) {
			return "", errors.New("验证码格式不正确")
		}
		if !s.authService.CheckSms(account, code) {
			return "", errors.New("验证码错误")
		}
		u, err := s.repositoryUser.FindByIdentifier(account)
		if err != nil {
			return "", errors.New("系统繁忙,请稍后再试")
		}
		if u == nil || u.Id <= 0 {
			lock := constant2.LockUserCreate + account
			if !s.cache.Lock(lock, constant2.LockTimeUserCreate) {
				return "", errors.New("系统繁忙,请稍后重试")
			}
			defer s.cache.UnLock(lock)
			num, err2 := s.repositoryUser.CountByPhone(account)
			if err2 != nil || num != 0 {
				return "", errors.New("该手机号已被占用")
			}
			u = &model.User{
				Phone:          account,
				RealAuthStatus: int(real_auth.UnAuth),
				Status:         int(user.StatusEnable),
			}
			err2 = s.repositoryUser.Create(u)
			if err2 != nil {
				return "", err
			}
			id = u.Id
		} else {
			id = u.Id
		}
	}
	if id <= 0 {
		return "", errors.New("获取用户信息错误")
	}
	return s.genLoginToken(id, userAgent, ip)
}

// Logout 注销登录
func (s *service) Logout(token string) error {
	key := jwt.CheckToken(token)
	if key == nil {
		return errors.New("未获取到登陆信息")
	}
	err := s.repositorySession.DeleteByToken(token)
	if err != nil {
		return errors.New("注销登录失败")
	}
	return nil
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
	if err = s.repositorySession.Create(&model.Session{
		UserId:     userId,
		Token:      token,
		Device:     userAgent,
		Ip:         ip,
		Location:   location,
		ExpireTime: expireTime,
	}); err != nil {
		return "", err
	}
	jwtToken := jwt.GetToken(&jwt.User{
		Id:         userId,
		LoginToken: token,
		LoginIp:    ip,
		LoginTime:  str.DateTime(time.Now()),
	}, constant.JwtExpireTime)
	return jwtToken, nil
}

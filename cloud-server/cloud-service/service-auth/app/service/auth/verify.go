package auth

import (
	"errors"
	"service-auth/app/enum/session"
	"time"
)

// VerifyUserToken 验证网站用户token有效性
func (s *service) VerifyUserToken(token string) error {
	return s.verifyToken(token, session.UserTypeWeb)
}

// VerifyAdminToken 验证后台用户token有效性
func (s *service) VerifyAdminToken(token string) error {
	return s.verifyToken(token, session.UserTypeAdmin)
}

// verifyToken 验证token
func (s *service) verifyToken(token string, userType session.UserType) error {
	info, err := s.repositorySession.FindByToken(token)
	if err != nil {
		return errors.New("token不存在")
	}
	if info.UserType != int(userType) {
		return errors.New("用户类型不匹配")
	}
	if time.Now().After(time.Time(info.ExpireTime)) {
		return errors.New("token已过期")
	}
	return nil
}

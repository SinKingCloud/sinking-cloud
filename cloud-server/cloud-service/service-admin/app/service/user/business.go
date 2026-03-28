package user

import (
	"encoding/json"
	"errors"
	"service-admin/app/constant"
	"service-admin/app/model"
	"service-admin/app/repository/user"
	"strconv"
)

// GetUserInfo 获取用户信息
func (s *service) GetUserInfo(userId int64, refresh bool) (user *model.User, err error) {
	cacheKey := constant.CacheNameWithUserInfo + strconv.FormatInt(userId, 10)
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			if err = json.Unmarshal([]byte(data), &user); err == nil {
				return user, nil
			}
		}
	}
	user, err = s.repositoryUser.FindById(userId)
	if err != nil {
		return nil, err
	}
	if data, err2 := json.Marshal(user); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithUserInfo)
	}
	return user, nil
}

// UpdateById 更新用户信息
func (s *service) UpdateById(id int64, data *user.UpdateUser) error {
	if data == nil {
		return errors.New("参数不能为空")
	}
	err := s.repositoryUser.UpdateById(id, data)
	if err == nil {
		_, _ = s.GetUserInfo(id, true)
	}
	return err
}

// SetUserVerifyStatus 设置用户安全验证状态
func (s *service) SetUserVerifyStatus(userId int64, allow bool, minute int) {
	key := constant.CacheNameWithUserVerify + strconv.FormatInt(userId, 10)
	value := "0"
	if allow {
		value = "1"
	}
	s.cache.SetWithExpire(key, value, int64(minute*60*1000))
}

// CheckUserVerifyStatus 判断用户安全验证状态
func (s *service) CheckUserVerifyStatus(userId int64, clear bool) bool {
	key := constant.CacheNameWithUserVerify + strconv.FormatInt(userId, 10)
	num, exists := s.cache.Get(key)
	if exists && num == "1" {
		if clear {
			s.cache.Delete(key)
		}
		return true
	}
	return false
}

// CheckPhoneIsUse 判断手机号是否被使用
func (s *service) CheckPhoneIsUse(phone string) error {
	num, err := s.repositoryUser.CountByPhone(phone)
	if err != nil {
		return errors.New("获取统计数量失败")
	}
	if num > 0 {
		return errors.New("该手机号已被使用")
	}
	return nil
}

// CheckEmailIsUse 判断邮箱是否被使用
func (s *service) CheckEmailIsUse(email string) error {
	num, err := s.repositoryUser.CountByEmail(email)
	if err != nil {
		return errors.New("获取统计数量失败")
	}
	if num > 0 {
		return errors.New("该邮箱已被使用")
	}
	return nil
}

// CheckAccountIsUse 判断账号是否被使用
func (s *service) CheckAccountIsUse(email string) error {
	num, err := s.repositoryUser.CountByAccount(email)
	if err != nil {
		return errors.New("获取统计数量失败")
	}
	if num > 0 {
		return errors.New("该账号已被使用")
	}
	return nil
}

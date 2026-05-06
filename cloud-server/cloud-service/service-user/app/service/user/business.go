package user

import (
	"encoding/json"
	"errors"
	"service-user/app/constant"
	"service-user/app/enum/real_auth_status"
	"service-user/app/enum/user_status"
	"service-user/app/model"
	"service-user/app/repository/user"
	"strconv"
)

// GetInfo 获取用户信息
func (s *service) GetInfo(userId int64, refresh bool) (user *model.User, err error) {
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
	if user == nil || user.Id <= 0 {
		return nil, errors.New("未查询到用户信息")
	}
	if data, err2 := json.Marshal(user); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithUserInfo)
	}
	return user, nil
}

// UpdateById 更新用户信息
func (s *service) UpdateById(id int64, data *user.UpdateUser) error {
	if id <= 0 || data == nil {
		return errors.New("参数错误")
	}
	info, err := s.GetInfo(id, false)
	if err != nil {
		return errors.Join(errors.New("查询用户失败"), err)
	}
	if info == nil || info.Id <= 0 {
		return errors.New("该用户不存在")
	}
	hasUpdateData := false
	if data.Account != nil {
		account, ok := data.Account.(string)
		if !ok {
			return errors.New("账号数据不合法")
		}
		if account == "" {
			return errors.New("账号不能为空")
		}
		if account == info.Account {
			data.Account = nil
		} else {
			if err = s.checkAccountIsUse(account); err != nil {
				return err
			}
			hasUpdateData = true
		}
	}
	if data.Phone != nil {
		phone, ok := data.Phone.(string)
		if !ok {
			return errors.New("手机号数据不合法")
		}
		if phone == "" {
			return errors.New("手机号不能为空")
		}
		if phone == info.Phone {
			data.Phone = nil
		} else {
			if err = s.checkPhoneIsUse(phone); err != nil {
				return err
			}
			hasUpdateData = true
		}
	}
	if data.Email != nil {
		email, ok := data.Email.(string)
		if !ok {
			return errors.New("邮箱数据不合法")
		}
		if email == "" {
			return errors.New("邮箱不能为空")
		}
		if email == info.Email {
			data.Email = nil
		} else {
			if err = s.checkEmailIsUse(email); err != nil {
				return err
			}
			hasUpdateData = true
		}
	}
	if data.Password != nil {
		password, ok := data.Password.(string)
		if !ok {
			return errors.New("密码数据不合法")
		}
		if password == "" {
			return errors.New("密码不能为空")
		}
		hasUpdateData = true
	}
	if data.Avatar != nil {
		avatar, ok := data.Avatar.(string)
		if !ok {
			return errors.New("头像数据不合法")
		}
		if avatar == info.Avatar {
			data.Avatar = nil
		} else {
			hasUpdateData = true
		}
	}
	if data.NickName != nil {
		nickName, ok := data.NickName.(string)
		if !ok {
			return errors.New("昵称数据不合法")
		}
		if nickName == info.NickName {
			data.NickName = nil
		} else {
			hasUpdateData = true
		}
	}
	if data.RealAuthStatus != nil {
		realAuthStatus, ok := data.RealAuthStatus.(int)
		if !ok {
			return errors.New("状态数据不合法")
		}
		if _, ok = real_auth_status.Map()[realAuthStatus]; !ok {
			return errors.New("状态数据不合法")
		}
		if realAuthStatus == info.RealAuthStatus {
			data.RealAuthStatus = nil
		} else {
			hasUpdateData = true
		}
	}
	if data.Status != nil {
		status, ok := data.Status.(int)
		if !ok {
			return errors.New("状态数据不合法")
		}
		if _, ok = user_status.Map()[status]; !ok {
			return errors.New("状态数据不合法")
		}
		if status == info.Status {
			data.Status = nil
		} else {
			hasUpdateData = true
		}
	}
	if !hasUpdateData {
		return nil
	}
	err = s.repositoryUser.UpdateById(id, data)
	if err == nil {
		_, _ = s.GetInfo(id, true)
	}
	return err
}

// Select 查询用户列表
func (s *service) Select(where *user.SelectUser, orderByField string, orderByType string, page int, pageSize int) ([]*user.User, int64, error) {
	if where.RealAuthStatus != "" {
		status, err := strconv.Atoi(where.RealAuthStatus)
		if err != nil {
			return nil, 0, errors.New("实名状态参数错误")
		}
		if _, ok := real_auth_status.Map()[status]; !ok {
			return nil, 0, errors.New("实名状态参数不合法")
		}
	}
	if where.Status != "" {
		status, err := strconv.Atoi(where.Status)
		if err != nil {
			return nil, 0, errors.New("状态参数错误")
		}
		if _, ok := user_status.Map()[status]; !ok {
			return nil, 0, errors.New("状态参数不合法")
		}
	}
	return s.repositoryUser.Select(where, orderByField, orderByType, page, pageSize)
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

// checkPhoneIsUse 判断手机号是否被使用
func (s *service) checkPhoneIsUse(phone string) error {
	exists, err := s.repositoryUser.ExistsByPhone(phone)
	if err != nil {
		return errors.New("获取结果失败")
	}
	if exists {
		return errors.New("该手机号已被使用")
	}
	return nil
}

// checkEmailIsUse 判断邮箱是否被使用
func (s *service) checkEmailIsUse(email string) error {
	exists, err := s.repositoryUser.ExistsByEmail(email)
	if err != nil {
		return errors.New("获取结果失败")
	}
	if exists {
		return errors.New("该邮箱已被使用")
	}
	return nil
}

// checkAccountIsUse 判断账号是否被使用
func (s *service) checkAccountIsUse(email string) error {
	exists, err := s.repositoryUser.ExistsByAccount(email)
	if err != nil {
		return errors.New("获取结果失败")
	}
	if exists {
		return errors.New("该账号已被使用")
	}
	return nil
}

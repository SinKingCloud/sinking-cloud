package user

import (
	"encoding/json"
	"errors"
	"service-user/app/constant"
	"service-user/app/model"
	"service-user/app/repository/user"
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

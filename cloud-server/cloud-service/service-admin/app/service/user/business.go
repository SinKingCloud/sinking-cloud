package user

import (
	"encoding/json"
	"errors"
	"service-admin/app/constant"
	"service-admin/app/enum/user_status"
	"service-admin/app/model"
	"service-admin/app/repository/user"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"gorm.io/gorm"
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

// Select 查询用户列表
func (s *service) Select(where *user.SelectUser, queryPage *page.Query) (*page.Result[*User], error) {
	if where != nil && where.Status != nil {
		if _, ok := user_status.Map()[*where.Status]; !ok {
			return nil, errors.New("状态参数不合法")
		}
	}
	result, err := s.repositoryUser.Select(where, queryPage)
	if err != nil {
		return nil, err
	}
	list := result.List
	userIds := make([]int64, 0, len(list))
	for _, item := range list {
		userIds = append(userIds, item.Id)
	}
	items, err := s.repositoryUserRole.SelectByUserIds(userIds)
	if err != nil {
		return nil, err
	}
	roleMap := make(map[int64][]int64, len(userIds))
	for _, item := range items {
		if item == nil || item.UserId <= 0 || item.RoleId <= 0 {
			continue
		}
		roleIds := roleMap[item.UserId]
		if len(roleIds) > 0 && roleIds[len(roleIds)-1] == item.RoleId {
			continue
		}
		roleMap[item.UserId] = append(roleMap[item.UserId], item.RoleId)
	}
	users := make([]*User, 0, len(result.List))
	for _, item := range result.List {
		users = append(users, &User{User: item, RoleIds: roleMap[item.Id]})
	}
	return &page.Result[*User]{
		Total:            result.Total,
		Page:             result.Page,
		PageSize:         result.PageSize,
		List:             users,
		HasMore:          result.HasMore,
		NextCursorId:     result.NextCursorId,
		NextCursorLastId: result.NextCursorLastId,
	}, nil
}

// Create 创建用户
func (s *service) Create(data *model.User) error {
	if data == nil {
		return errors.New("参数错误")
	}
	if data.Account == "" {
		return errors.New("账号不能为空")
	}
	if data.Password == "" {
		return errors.New("密码不能为空")
	}
	if _, ok := user_status.Map()[data.Status]; !ok {
		return errors.New("状态数据不合法")
	}
	if err := s.checkAccountIsUse(data.Account); err != nil {
		return err
	}
	if err := s.checkPhoneIsUse(data.Phone); err != nil {
		return err
	}
	if err := s.checkEmailIsUse(data.Email); err != nil {
		return err
	}
	return s.repositoryUser.Create(data)
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
		account := *data.Account
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
		phone := *data.Phone
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
		email := *data.Email
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
		password := *data.Password
		if password == "" {
			return errors.New("密码不能为空")
		}
		hasUpdateData = true
	}
	if data.Avatar != nil {
		avatar := *data.Avatar
		if avatar == info.Avatar {
			data.Avatar = nil
		} else {
			hasUpdateData = true
		}
	}
	if data.NickName != nil {
		nickName := *data.NickName
		if nickName == info.NickName {
			data.NickName = nil
		} else {
			hasUpdateData = true
		}
	}
	if data.Status != nil {
		status := *data.Status
		if _, ok := user_status.Map()[status]; !ok {
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

// DeleteById 删除用户
func (s *service) DeleteById(id int64) error {
	if id <= 0 {
		return nil
	}
	lockKey := constant.LockUserOperate + strconv.FormatInt(id, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeUserOperate) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	uidStr := strconv.FormatInt(id, 10)
	configIdentId := constant.UserConfigIdentId + uidStr
	err := s.database.Transaction(func(tx *gorm.DB) error {
		err := s.repositoryUserRole.DeleteByUserId(id, tx)
		if err != nil {
			return errors.Join(errors.New("删除用户角色关联失败"), err)
		}
		err = s.repositorySession.DeleteByUserId(id, tx)
		if err != nil {
			return errors.Join(errors.New("删除用户会话失败"), err)
		}
		err = s.repositoryConfig.DeleteByIdentId(configIdentId, tx)
		if err != nil {
			return errors.Join(errors.New("删除用户配置失败"), err)
		}
		err = s.repositoryUserLog.DeleteByUserId(id, tx)
		if err != nil {
			return errors.Join(errors.New("删除用户操作日志失败"), err)
		}
		err = s.repositoryUser.DeleteById(id, tx)
		if err != nil {
			return errors.Join(errors.New("删除用户失败"), err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.cache.Delete(constant.CacheNameWithUserInfo+uidStr, constant.CacheNameWithUserRoles+uidStr, constant.CacheNameWithUserVerify+uidStr)
	return nil
}

// AssignRole 分配用户角色
func (s *service) AssignRole(userId int64, roleIds []int64) error {
	info, err := s.GetInfo(userId, false)
	if err != nil {
		return errors.Join(errors.New("查询用户失败"), err)
	}
	if info == nil || info.Id <= 0 {
		return errors.New("该用户不存在")
	}
	return s.rbacService.AssignRoleToUser(userId, roleIds)
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

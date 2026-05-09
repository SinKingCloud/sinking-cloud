package user

import (
	"encoding/json"
	"errors"
	"service-user/app/constant"
	"service-user/app/enum/real_auth_status"
	"service-user/app/enum/real_auth_type"
	"service-user/app/enum/user_status"
	"service-user/app/model"
	repositoryRealAuth "service-user/app/repository/real_auth"
	"service-user/app/repository/user"
	"strconv"
	"strings"
	"time"

	commonConstant "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/aliyun"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
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
	if data.RealAuthType != nil {
		realAuthType, ok := data.RealAuthType.(int)
		if !ok {
			return errors.New("实名类型数据不合法")
		}
		if _, ok = real_auth_type.Map()[realAuthType]; !ok {
			return errors.New("实名类型数据不合法")
		}
		if realAuthType == info.RealAuthType {
			data.RealAuthType = nil
		} else {
			hasUpdateData = true
		}
	}
	if data.RealAuthStatus != nil {
		realAuthStatus, ok := data.RealAuthStatus.(int)
		if !ok {
			return errors.New("实名状态数据不合法")
		}
		if _, ok = real_auth_status.Map()[realAuthStatus]; !ok {
			return errors.New("实名状态数据不合法")
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
	realAuthUpdated := data.RealAuthType != nil || data.RealAuthStatus != nil
	err = s.repositoryUser.UpdateById(id, data)
	if err == nil {
		_, _ = s.GetInfo(id, true)
		if realAuthUpdated {
			s.cache.Delete(constant.CacheNameWithRealAuthInfo + strconv.FormatInt(id, 10))
		}
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

// GetRealAuthInfo 获取实名信息。
func (s *service) GetRealAuthInfo(userId int64, refresh bool) (*RealAuthInfo, error) {
	if userId <= 0 {
		return nil, errors.New("用户ID不能为空")
	}
	cacheKey := constant.CacheNameWithRealAuthInfo + strconv.FormatInt(userId, 10)
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			var cacheInfo RealAuthInfo
			if err := json.Unmarshal([]byte(data), &cacheInfo); err == nil {
				return &cacheInfo, nil
			}
			s.cache.Delete(cacheKey)
		}
	}
	userInfo, err := s.GetInfo(userId, false)
	if err != nil {
		return nil, errors.Join(errors.New("获取用户信息失败"), err)
	}
	info := &RealAuthInfo{
		UserId: userId,
		Type:   userInfo.RealAuthType,
		Status: userInfo.RealAuthStatus,
	}
	authInfo, err := s.repositoryRealAuth.FindOrCreateByUserId(userId)
	if err != nil {
		return nil, errors.Join(errors.New("获取实名信息失败"), err)
	}
	if authInfo == nil || authInfo.Id <= 0 {
		return nil, errors.New("获取实名信息失败,实名信息不存在")
	}
	info.Type = authInfo.Type
	info.Status = authInfo.Status
	info.CompanyName = authInfo.CompanyName
	info.CompanyNo = authInfo.CompanyNo
	info.CardName = authInfo.CardName
	info.CardNo = authInfo.CardNo
	info.AuthTime = authInfo.AuthTime
	if strings.TrimSpace(authInfo.Files) != "" {
		var files RealAuthFiles
		if err := json.Unmarshal([]byte(authInfo.Files), &files); err == nil {
			info.Files = &files
			aliConfig := s.config.GetSystemConfig(commonConstant.AliGroup, false)
			domain := strings.TrimRight(strings.TrimSpace(aliConfig[commonConstant.AliFaceDomain]), "/")
			if domain != "" {
				if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
					domain = "https://" + domain
				}
				if info.Files.VideoPath != "" && !strings.HasPrefix(info.Files.VideoPath, "http://") && !strings.HasPrefix(info.Files.VideoPath, "https://") {
					info.Files.VideoPath = domain + "/" + strings.TrimLeft(info.Files.VideoPath, "/")
				}
				for i, path := range info.Files.PicPaths {
					if path == "" || strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
						continue
					}
					info.Files.PicPaths[i] = domain + "/" + strings.TrimLeft(path, "/")
				}
			}
		}
	}
	if data, err := json.Marshal(info); err == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithRealAuthInfo)
	}
	return info, nil
}

// InitRealAuth 发起实名认证。
func (s *service) InitRealAuth(userId int64, authType int, cardName, cardNo, companyName, companyNo, env string) (string, string, error) {
	if userId <= 0 {
		return "", "", errors.New("用户ID不能为空")
	}
	if _, ok := real_auth_type.Map()[authType]; !ok || authType == real_auth_type.UnAuth {
		return "", "", errors.New("实名类型不合法")
	}
	if cardName == "" {
		return "", "", errors.New("身份证姓名不能为空")
	}
	if cardNo == "" {
		return "", "", errors.New("身份证号码不能为空")
	}
	if env == "" {
		return "", "", errors.New("环境信息不能为空")
	}
	if authType == real_auth_type.Company {
		if companyName == "" {
			return "", "", errors.New("企业名称不能为空")
		}
		if companyNo == "" {
			return "", "", errors.New("企业编码不能为空")
		}
	}
	lockKey := constant.LockRealAuth + strconv.FormatInt(userId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeRealAuth) {
		return "", "", errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	authInfo, err := s.repositoryRealAuth.FindOrCreateByUserId(userId)
	if err != nil {
		return "", "", errors.Join(errors.New("获取实名信息失败"), err)
	}
	if authInfo != nil && authInfo.Status == real_auth_status.Authed {
		return "", "", errors.New("账户已实名成功")
	}
	aliConfig := s.config.GetSystemConfig(commonConstant.AliGroup, false)
	key := strings.TrimSpace(aliConfig[commonConstant.AliKey])
	secret := strings.TrimSpace(aliConfig[commonConstant.AliSecret])
	if key == "" || secret == "" {
		return "", "", errors.New("阿里云配置不完整")
	}
	ali, err := aliyun.New(key, secret)
	if err != nil {
		return "", "", err
	}
	if authType == real_auth_type.Company {
		if err = ali.Company.Verify(companyName, companyNo, cardName, cardNo); err != nil {
			return "", "", err
		}
	}
	scene := strings.TrimSpace(aliConfig[commonConstant.AliFaceScene])
	if scene == "" {
		return "", "", errors.New("实名认证场景ID未配置")
	}
	sceneId, err := strconv.ParseInt(scene, 10, 64)
	if err != nil || sceneId <= 0 {
		return "", "", errors.New("实名认证场景ID配置不合法")
	}
	webConfig := s.config.GetSystemConfig(commonConstant.WebGroup, false)
	returnUrl := strings.TrimRight(strings.TrimSpace(webConfig[commonConstant.WebUrl]), "/")
	callbackUrl := strings.TrimRight(strings.TrimSpace(webConfig[commonConstant.WebApi]), "/")
	callbackUrl = strings.TrimPrefix(callbackUrl, "http://")
	if !strings.HasPrefix(callbackUrl, "https://") {
		callbackUrl = "https://" + callbackUrl
	}
	outerOrderNo := strconv.FormatInt(str.GetSnowWorkIns().GetId(), 10)
	result, err := ali.Face.Init(sceneId, outerOrderNo, cardNo, cardName, returnUrl, callbackUrl, env)
	if err != nil {
		return "", "", err
	}
	process := &realAuthProcess{
		UserId:       userId,
		Type:         authType,
		CompanyName:  companyName,
		CompanyNo:    companyNo,
		CardName:     cardName,
		CardNo:       cardNo,
		CertifyId:    result.CertifyId,
		CertifyUrl:   result.CertifyUrl,
		SceneId:      sceneId,
		AccessKey:    key,
		AccessSecret: secret,
		OuterOrderNo: outerOrderNo,
	}
	processData, err := json.Marshal(process)
	if err != nil {
		return "", "", errors.Join(errors.New("生成实名流程缓存失败"), err)
	}
	if !s.cache.SetWithExpire(constant.CacheNameWithRealAuthProcess+process.CertifyId, string(processData), constant.CacheTimeWithRealAuthProcess) {
		return "", "", errors.New("写入实名流程缓存失败")
	}
	return result.CertifyId, result.CertifyUrl, nil
}

// VerifyRealAuth 查询实名认证状态，认证通过后写入实名记录。
func (s *service) VerifyRealAuth(certifyId string) (*RealAuthInfo, error) {
	certifyId = strings.TrimSpace(certifyId)
	if certifyId == "" {
		return nil, errors.New("认证ID不能为空")
	}
	cacheKey := constant.CacheNameWithRealAuthProcess + certifyId
	data, exists := s.cache.Get(cacheKey)
	if !exists || data == "" {
		return nil, errors.New("未查询到实名认证流程,请重新发起认证")
	}
	var process realAuthProcess
	if err := json.Unmarshal([]byte(data), &process); err != nil {
		s.cache.Delete(cacheKey)
		return nil, errors.New("未查询到实名认证流程,请重新发起认证")
	}
	if process.UserId <= 0 || process.SceneId <= 0 || process.AccessKey == "" || process.AccessSecret == "" {
		s.cache.Delete(cacheKey)
		return nil, errors.New("未查询到实名认证流程,请重新发起认证")
	}
	userId := process.UserId
	lockKey := constant.LockRealAuth + strconv.FormatInt(userId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeRealAuth) {
		return nil, errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)

	userInfo, err := s.GetInfo(userId, false)
	if err != nil {
		return nil, errors.Join(errors.New("获取用户信息失败"), err)
	}
	authInfo, err := s.repositoryRealAuth.FindOrCreateByUserId(userId)
	if err != nil {
		return nil, errors.Join(errors.New("获取实名信息失败"), err)
	}
	if authInfo != nil && authInfo.Status == real_auth_status.Authed {
		if userInfo.RealAuthStatus != real_auth_status.Authed {
			if err = s.repositoryUser.UpdateById(userId, &user.UpdateUser{
				RealAuthType:   authInfo.Type,
				RealAuthStatus: real_auth_status.Authed,
			}); err != nil {
				return nil, errors.Join(errors.New("更新用户实名状态失败"), err)
			}
		}
		s.cache.Delete(
			constant.CacheNameWithUserInfo+strconv.FormatInt(userId, 10),
			constant.CacheNameWithRealAuthInfo+strconv.FormatInt(userId, 10),
			cacheKey,
		)
		return s.GetRealAuthInfo(userId, true)
	}
	ali, err := aliyun.New(process.AccessKey, process.AccessSecret)
	if err != nil {
		return nil, err
	}
	result, err := ali.Face.Query(process.CertifyId, process.SceneId)
	if err != nil {
		return nil, err
	}
	if result == nil || !result.Passed {
		return nil, errors.New("实名认证未通过")
	}
	files := &RealAuthFiles{
		Bucket:    result.Bucket,
		PicPaths:  result.PicPaths,
		VideoPath: result.VideoPath,
	}
	filesData, err := json.Marshal(files)
	if err != nil {
		return nil, errors.Join(errors.New("生成实名资料失败"), err)
	}
	err = s.database.Transaction(func(tx *gorm.DB) error {
		authInfo, err = s.repositoryRealAuth.FindOrCreateByUserId(userId, tx)
		if err != nil {
			return errors.Join(errors.New("获取实名信息失败"), err)
		}
		if authInfo == nil || authInfo.Id <= 0 {
			return errors.New("获取实名信息失败")
		}
		if authInfo.Status == real_auth_status.Authed {
			return nil
		}
		if err = s.repositoryRealAuth.UpdateById(authInfo.Id, &repositoryRealAuth.UpdateRealAuth{
			Type:        process.Type,
			CompanyName: process.CompanyName,
			CompanyNo:   process.CompanyNo,
			CardName:    process.CardName,
			CardNo:      process.CardNo,
			Files:       string(filesData),
			Status:      real_auth_status.Authed,
			AuthTime:    str.DateTime(time.Now()),
		}, tx); err != nil {
			return errors.Join(errors.New("更新实名记录失败"), err)
		}
		if err = s.repositoryUser.UpdateById(userId, &user.UpdateUser{
			RealAuthType:   process.Type,
			RealAuthStatus: real_auth_status.Authed,
		}, tx); err != nil {
			return errors.Join(errors.New("更新用户实名状态失败"), err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.cache.Delete(
		constant.CacheNameWithUserInfo+strconv.FormatInt(userId, 10),
		constant.CacheNameWithRealAuthInfo+strconv.FormatInt(userId, 10),
		cacheKey,
	)
	return s.GetRealAuthInfo(userId, true)
}

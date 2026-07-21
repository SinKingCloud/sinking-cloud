package mind

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"service-space/app/constant"
	"service-space/app/enum/mind_share_status"
	"service-space/app/model"
	mindShareRepository "service-space/app/repository/mind_share"
	"strconv"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// CreateShare 创建导图分享
func (s *service) CreateShare(userId int64, dataId int64, status int, expireTime string, remark string) (*model.MindShare, error) {
	if userId <= 0 || dataId <= 0 {
		return nil, errors.New("参数错误")
	}
	if _, ok := mind_share_status.Map()[status]; !ok {
		return nil, errors.New("分享状态错误")
	}
	expireTime = strings.TrimSpace(expireTime)
	var expireAt str.DateTime
	if expireTime != "" {
		value, err := time.ParseInLocation("2006-01-02 15:04:05", expireTime, time.Local)
		if err != nil {
			return nil, errors.New("分享过期时间格式错误")
		}
		if !value.After(time.Now()) {
			return nil, errors.New("分享过期时间必须晚于当前时间")
		}
		expireAt = str.DateTime(value)
	}
	lockKey := constant.LockMindDataOperate + strconv.FormatInt(dataId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeMindDataOperate) {
		return nil, errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	for attempt := 0; attempt < 5; attempt++ {
		codeBytes := make([]byte, 6)
		if _, err := rand.Read(codeBytes); err != nil {
			return nil, errors.Join(errors.New("生成分享码失败"), err)
		}
		data := &model.MindShare{
			Id:         str.GetSnowWorkIns().GetId(),
			DataId:     dataId,
			Code:       base64.RawURLEncoding.EncodeToString(codeBytes),
			Status:     status,
			ExpireTime: expireAt,
			Remark:     strings.TrimSpace(remark),
		}
		codeExists := false
		var insertErr error
		createErr := s.database.Transaction(func(tx *gorm.DB) error {
			exists, err := s.mindDataRepository.ExistsByIdAndUserId(dataId, userId, tx)
			if err != nil {
				return errors.Join(errors.New("获取导图失败"), err)
			}
			if !exists {
				return errors.New("导图不存在")
			}
			share, err := s.mindShareRepository.FindByCode(data.Code, tx)
			if err != nil {
				return errors.Join(errors.New("检查分享码失败"), err)
			}
			if share != nil && share.Id > 0 {
				codeExists = true
				return nil
			}
			insertErr = s.mindShareRepository.Create(data, tx)
			return insertErr
		})
		if codeExists {
			if createErr != nil {
				return nil, createErr
			}
			continue
		}
		if createErr == nil {
			return data, nil
		}
		if insertErr == nil {
			return nil, createErr
		}
		var share *model.MindShare
		if err := s.database.Transaction(func(tx *gorm.DB) error {
			var err error
			share, err = s.mindShareRepository.FindByCode(data.Code, tx)
			return err
		}); err != nil {
			return nil, errors.Join(errors.New("检查分享码失败"), err)
		}
		if share != nil && share.Id > 0 {
			continue
		}
		return nil, errors.Join(errors.New("创建导图分享失败"), insertErr)
	}
	return nil, errors.New("生成分享码失败")
}

// UpdateShare 更新导图分享
func (s *service) UpdateShare(userId int64, dataId int64, id int64, data *mindShareRepository.UpdateMindShare) error {
	if userId <= 0 || dataId <= 0 || id <= 0 || data == nil {
		return errors.New("参数错误")
	}
	update := &mindShareRepository.UpdateMindShare{}
	if data.Status != nil {
		value, ok := data.Status.(int)
		if !ok {
			return errors.New("分享状态错误")
		}
		if _, ok := mind_share_status.Map()[value]; !ok {
			return errors.New("分享状态错误")
		}
		update.Status = value
	}
	if data.Remark != nil {
		value, ok := data.Remark.(string)
		if !ok {
			return errors.New("分享备注错误")
		}
		update.Remark = strings.TrimSpace(value)
	}
	lockKey := constant.LockMindDataOperate + strconv.FormatInt(dataId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeMindDataOperate) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	exists, err := s.mindDataRepository.ExistsByIdAndUserId(dataId, userId)
	if err != nil {
		return errors.Join(errors.New("获取导图失败"), err)
	}
	if !exists {
		return errors.New("导图不存在")
	}
	share, err := s.mindShareRepository.FindByIdAndDataId(id, dataId)
	if err != nil {
		return errors.Join(errors.New("获取导图分享失败"), err)
	}
	if share == nil || share.Id <= 0 {
		return errors.New("导图分享不存在")
	}
	if err = s.mindShareRepository.UpdateByIdAndDataId(id, dataId, update); err != nil {
		return errors.Join(errors.New("更新导图分享失败"), err)
	}
	return nil
}

// DeleteShare 删除导图分享
func (s *service) DeleteShare(userId int64, dataId int64, id int64) error {
	if userId <= 0 || dataId <= 0 || id <= 0 {
		return errors.New("参数错误")
	}
	lockKey := constant.LockMindDataOperate + strconv.FormatInt(dataId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeMindDataOperate) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	exists, err := s.mindDataRepository.ExistsByIdAndUserId(dataId, userId)
	if err != nil {
		return errors.Join(errors.New("获取导图失败"), err)
	}
	if !exists {
		return errors.New("导图不存在")
	}
	share, err := s.mindShareRepository.FindByIdAndDataId(id, dataId)
	if err != nil {
		return errors.Join(errors.New("获取导图分享失败"), err)
	}
	if share == nil || share.Id <= 0 {
		return errors.New("导图分享不存在")
	}
	if err = s.mindShareRepository.DeleteByIdAndDataId(id, dataId); err != nil {
		return errors.Join(errors.New("删除导图分享失败"), err)
	}
	return nil
}

// SelectSharesByUserId 查询用户导图分享列表
func (s *service) SelectSharesByUserId(userId int64, dataId int64, query *page.Query) ([]*model.MindShare, int64, error) {
	if userId <= 0 || dataId <= 0 {
		return nil, 0, errors.New("参数错误")
	}
	if query == nil {
		return nil, 0, errors.New("分页参数错误")
	}
	exists, err := s.mindDataRepository.ExistsByIdAndUserId(dataId, userId)
	if err != nil {
		return nil, 0, errors.Join(errors.New("获取导图失败"), err)
	}
	if !exists {
		return nil, 0, errors.New("导图不存在")
	}
	list, total, err := s.mindShareRepository.SelectByDataId(dataId, query)
	if err != nil {
		return nil, 0, errors.Join(errors.New("查询导图分享失败"), err)
	}
	return list, total, nil
}

// GetDataByShareCode 通过分享码获取公开导图
func (s *service) GetDataByShareCode(code string) (*SharedData, error) {
	code = strings.TrimSpace(code)
	if len(code) != 8 {
		return nil, errors.New("分享码格式错误")
	}
	share, err := s.mindShareRepository.FindByCode(code)
	if err != nil {
		return nil, errors.Join(errors.New("获取分享信息失败"), err)
	}
	if share == nil || share.Id <= 0 || share.Status != mind_share_status.Enable {
		return nil, errors.New("分享不存在或已失效")
	}
	expireTime := time.Time(share.ExpireTime)
	if !expireTime.IsZero() && !expireTime.After(time.Now()) {
		return nil, errors.New("分享不存在或已失效")
	}
	info, err := s.mindDataRepository.FindById(share.DataId)
	if err != nil {
		return nil, errors.Join(errors.New("获取导图失败"), err)
	}
	if info == nil || info.Id <= 0 {
		return nil, errors.New("分享不存在或已失效")
	}
	return &SharedData{
		Title:      info.Title,
		Describe:   info.Describe,
		Value:      info.Value,
		UpdateTime: info.UpdateTime,
	}, nil
}

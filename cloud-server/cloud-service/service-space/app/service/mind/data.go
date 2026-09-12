package mind

import (
	"bytes"
	"encoding/json"
	"errors"
	"service-space/app/constant"
	"service-space/app/enum/mind_type"
	"service-space/app/enum/mind_version_type"
	"service-space/app/model"
	mindDataRepository "service-space/app/repository/mind_data"
	"strconv"
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// CreateData 创建导图数据
func (s *service) CreateData(userId int64, typeId int64, title string, describe string, value string) (*Data, error) {
	if userId <= 0 {
		return nil, errors.New("用户信息错误")
	}
	if typeId < 0 {
		return nil, errors.New("导图分类参数错误")
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("导图标题不能为空")
	}
	valueText, err := s.normalizeValue(value)
	if err != nil {
		return nil, err
	}
	settings := s.configService.GetSystemConfig(constant.MindGroup, false)
	createMax, _ := strconv.ParseInt(strings.TrimSpace(settings[constant.MindCreateMax]), 10, 64)
	if createMax != 0 && createMax < 10 {
		createMax = 10
	}
	if createMax > 10000 {
		createMax = 10000
	}
	if createMax > 0 {
		lockKey := constant.LockMindDataCreate + strconv.FormatInt(userId, 10)
		if !s.cache.Lock(lockKey, constant.LockTimeMindDataCreate) {
			return nil, errors.New("系统繁忙,请稍候重试")
		}
		defer s.cache.UnLock(lockKey)
	}
	data := &model.MindData{
		Id:       str.GetSnowWorkIns().GetId(),
		UserId:   userId,
		TypeId:   typeId,
		Title:    title,
		Describe: strings.TrimSpace(describe),
		Value:    valueText,
	}
	err = s.database.Transaction(func(tx *gorm.DB) error {
		if _, ok := mind_type.Map()[int(typeId)]; !ok {
			typeInfo, err := s.mindTypeRepository.FindByIdAndUserId(typeId, userId, tx)
			if err != nil {
				return errors.Join(errors.New("获取导图分类失败"), err)
			}
			if typeInfo == nil || typeInfo.Id <= 0 {
				return errors.New("导图分类不存在")
			}
		}
		if createMax > 0 {
			total, err := s.mindDataRepository.CountByUserId(userId, tx)
			if err != nil {
				return errors.Join(errors.New("获取导图数量失败"), err)
			}
			if total >= createMax {
				return errors.New("导图数量已达到上限")
			}
		}
		if err = s.mindDataRepository.Create(data, tx); err != nil {
			return errors.Join(errors.New("创建导图失败"), err)
		}
		version := &model.MindVersion{
			Id:     str.GetSnowWorkIns().GetId(),
			DataId: data.Id,
			Type:   mind_version_type.Manual,
			Value:  data.Value,
			Remark: "初始版本",
		}
		if err = s.mindVersionRepository.Create(version, tx); err != nil {
			return errors.Join(errors.New("创建初始版本失败"), err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Data{
		Id:         data.Id,
		TypeId:     data.TypeId,
		Title:      data.Title,
		Describe:   data.Describe,
		CreateTime: data.CreateTime,
		UpdateTime: data.UpdateTime,
	}, nil
}

// CopyData 复制导图数据
func (s *service) CopyData(userId int64, id int64) (*Data, error) {
	if userId <= 0 || id <= 0 {
		return nil, errors.New("参数错误")
	}
	settings := s.configService.GetSystemConfig(constant.MindGroup, false)
	createMax, _ := strconv.ParseInt(strings.TrimSpace(settings[constant.MindCreateMax]), 10, 64)
	if createMax != 0 && createMax < 10 {
		createMax = 10
	}
	if createMax > 10000 {
		createMax = 10000
	}
	if createMax > 0 {
		createLockKey := constant.LockMindDataCreate + strconv.FormatInt(userId, 10)
		if !s.cache.Lock(createLockKey, constant.LockTimeMindDataCreate) {
			return nil, errors.New("系统繁忙,请稍候重试")
		}
		defer s.cache.UnLock(createLockKey)
	}
	lockKey := constant.LockMindDataOperate + strconv.FormatInt(id, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeMindDataOperate) {
		return nil, errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	var data *model.MindData
	presetTypes := mind_type.Map()
	for attempt := 0; attempt < 3; attempt++ {
		var snapshot *model.MindData
		err := s.database.Transaction(func(tx *gorm.DB) error {
			var findErr error
			snapshot, findErr = s.mindDataRepository.FindByIdAndUserId(id, userId, tx)
			return findErr
		})
		if err != nil {
			return nil, errors.Join(errors.New("获取导图失败"), err)
		}
		if snapshot == nil || snapshot.Id <= 0 {
			return nil, errors.New("导图不存在")
		}
		retry := false
		err = s.database.Transaction(func(tx *gorm.DB) error {
			if _, isPresetType := presetTypes[int(snapshot.TypeId)]; !isPresetType {
				typeInfo, err := s.mindTypeRepository.FindByIdAndUserId(snapshot.TypeId, userId, tx)
				if err != nil {
					return errors.Join(errors.New("获取导图分类失败"), err)
				}
				if typeInfo == nil || typeInfo.Id <= 0 {
					retry = true
					return nil
				}
			}
			info, err := s.mindDataRepository.FindByIdAndUserId(id, userId, tx)
			if err != nil {
				return errors.Join(errors.New("获取导图失败"), err)
			}
			if info == nil || info.Id <= 0 {
				return errors.New("导图不存在")
			}
			if info.TypeId != snapshot.TypeId {
				retry = true
				return nil
			}
			if createMax > 0 {
				total, err := s.mindDataRepository.CountByUserId(userId, tx)
				if err != nil {
					return errors.Join(errors.New("获取导图数量失败"), err)
				}
				if total >= createMax {
					return errors.New("导图数量已达到上限")
				}
			}
			data = &model.MindData{
				Id:       str.GetSnowWorkIns().GetId(),
				UserId:   userId,
				TypeId:   info.TypeId,
				Describe: info.Describe,
				Value:    info.Value,
			}
			const suffix = " 副本"
			title := []rune(strings.TrimSpace(info.Title))
			maxLength := 60 - len([]rune(suffix))
			if len(title) > maxLength {
				title = title[:maxLength]
			}
			data.Title = string(title) + suffix
			if err = s.mindDataRepository.Create(data, tx); err != nil {
				return errors.Join(errors.New("复制导图失败"), err)
			}
			version := &model.MindVersion{
				Id:     str.GetSnowWorkIns().GetId(),
				DataId: data.Id,
				Type:   mind_version_type.Manual,
				Value:  data.Value,
				Remark: "复制创建",
			}
			if err = s.mindVersionRepository.Create(version, tx); err != nil {
				return errors.Join(errors.New("创建初始版本失败"), err)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		if !retry {
			return &Data{
				Id:         data.Id,
				TypeId:     data.TypeId,
				Title:      data.Title,
				Describe:   data.Describe,
				CreateTime: data.CreateTime,
				UpdateTime: data.UpdateTime,
			}, nil
		}
	}
	return nil, errors.New("导图分类已变更，请重试复制")
}

// DeleteData 删除导图数据
func (s *service) DeleteData(userId int64, id int64) error {
	if userId <= 0 || id <= 0 {
		return errors.New("参数错误")
	}
	lockKey := constant.LockMindDataOperate + strconv.FormatInt(id, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeMindDataOperate) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	return s.database.Transaction(func(tx *gorm.DB) error {
		exists, err := s.mindDataRepository.ExistsByIdAndUserId(id, userId, tx)
		if err != nil {
			return errors.Join(errors.New("获取导图失败"), err)
		}
		if !exists {
			return errors.New("导图不存在")
		}
		if err = s.mindShareRepository.DeleteByDataIds([]int64{id}, tx); err != nil {
			return errors.Join(errors.New("删除导图分享失败"), err)
		}
		if err = s.mindVersionRepository.DeleteByDataIds([]int64{id}, tx); err != nil {
			return errors.Join(errors.New("删除导图版本失败"), err)
		}
		if err = s.mindDataRepository.DeleteByIdAndUserId(id, userId, tx); err != nil {
			return errors.Join(errors.New("删除导图失败"), err)
		}
		return nil
	})
}

// GetDataInfoByUserId 查询用户导图数据详情
func (s *service) GetDataInfoByUserId(userId int64, id int64) (*Data, error) {
	if userId <= 0 || id <= 0 {
		return nil, errors.New("参数错误")
	}
	info, err := s.mindDataRepository.FindByIdAndUserId(id, userId)
	if err != nil {
		return nil, errors.Join(errors.New("获取导图失败"), err)
	}
	if info == nil || info.Id <= 0 {
		return nil, errors.New("导图不存在")
	}
	sharedDataIds, err := s.mindShareRepository.SelectValidDataIds([]int64{id})
	if err != nil {
		return nil, errors.Join(errors.New("获取分享信息失败"), err)
	}
	return &Data{
		Id:         info.Id,
		TypeId:     info.TypeId,
		Title:      info.Title,
		Describe:   info.Describe,
		Value:      info.Value,
		Shared:     len(sharedDataIds) > 0,
		CreateTime: info.CreateTime,
		UpdateTime: info.UpdateTime,
	}, nil
}

// SelectDataByUserId 查询用户导图数据列表
func (s *service) SelectDataByUserId(userId int64, where *mindDataRepository.SelectMindData, query *page.Query) (*page.Result[*Data], error) {
	if userId <= 0 {
		return nil, errors.New("用户信息错误")
	}
	if query == nil {
		return nil, errors.New("分页参数错误")
	}
	if where == nil {
		where = &mindDataRepository.SelectMindData{}
	}
	where.UserId = &userId
	if where.TypeId != nil && *where.TypeId <= 0 {
		where.TypeId = nil
	}
	if where.Keyword != nil {
		keyword := strings.TrimSpace(*where.Keyword)
		where.Keyword = nil
		if keyword != "" {
			where.Keyword = &keyword
		}
	}
	data, err := s.mindDataRepository.Select(where, query)
	if err != nil {
		return nil, errors.Join(errors.New("查询导图失败"), err)
	}
	ids := make([]int64, 0, len(data.List))
	for _, item := range data.List {
		if item != nil && item.Id > 0 {
			ids = append(ids, item.Id)
		}
	}
	sharedDataIds, err := s.mindShareRepository.SelectValidDataIds(ids)
	if err != nil {
		return nil, errors.Join(errors.New("获取分享信息失败"), err)
	}
	sharedMap := make(map[int64]bool, len(sharedDataIds))
	for _, dataId := range sharedDataIds {
		sharedMap[dataId] = true
	}
	list := make([]*Data, 0, len(data.List))
	for _, item := range data.List {
		if item == nil {
			list = append(list, nil)
			continue
		}
		list = append(list, &Data{
			Id:         item.Id,
			TypeId:     item.TypeId,
			Title:      item.Title,
			Describe:   item.Describe,
			Shared:     sharedMap[item.Id],
			CreateTime: item.CreateTime,
			UpdateTime: item.UpdateTime,
		})
	}
	return &page.Result[*Data]{
		Total:            data.Total,
		Page:             data.Page,
		PageSize:         data.PageSize,
		List:             list,
		HasMore:          data.HasMore,
		NextCursorId:     data.NextCursorId,
		NextCursorLastId: data.NextCursorLastId,
	}, nil
}

// UpdateData 更新导图数据信息
func (s *service) UpdateData(userId int64, id int64, data *mindDataRepository.UpdateMindData) error {
	if userId <= 0 || id <= 0 || data == nil {
		return errors.New("参数错误")
	}
	update := &mindDataRepository.UpdateMindData{}
	isPresetType := false
	if data.TypeId != nil {
		_, isPresetType = mind_type.Map()[int(*data.TypeId)]
		if !isPresetType && *data.TypeId <= 0 {
			return errors.New("导图分类参数错误")
		}
		update.TypeId = data.TypeId
	}
	if data.Title != nil {
		value := strings.TrimSpace(*data.Title)
		if value == "" {
			return errors.New("导图标题不能为空")
		}
		update.Title = &value
	}
	if data.Describe != nil {
		value := strings.TrimSpace(*data.Describe)
		update.Describe = &value
	}
	if data.TypeId != nil && !isPresetType {
		typeLockKey := constant.LockMindTypeOperate + strconv.FormatInt(userId, 10)
		if !s.cache.Lock(typeLockKey, constant.LockTimeMindTypeOperate) {
			return errors.New("系统繁忙,请稍候重试")
		}
		defer s.cache.UnLock(typeLockKey)
		typeInfo, err := s.mindTypeRepository.FindByIdAndUserId(*data.TypeId, userId)
		if err != nil {
			return errors.Join(errors.New("获取导图分类失败"), err)
		}
		if typeInfo == nil || typeInfo.Id <= 0 {
			return errors.New("导图分类不存在")
		}
	}
	lockKey := constant.LockMindDataOperate + strconv.FormatInt(id, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeMindDataOperate) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	exists, err := s.mindDataRepository.ExistsByIdAndUserId(id, userId)
	if err != nil {
		return errors.Join(errors.New("获取导图失败"), err)
	}
	if !exists {
		return errors.New("导图不存在")
	}
	if err = s.mindDataRepository.UpdateByIdAndUserId(id, userId, update); err != nil {
		return errors.Join(errors.New("更新导图失败"), err)
	}
	return nil
}

// normalizeValue 验证并压缩导图内容
func (s *service) normalizeValue(value string) (string, error) {
	value = strings.TrimSpace(value)
	content := []byte(value)
	if len(content) == 0 || !json.Valid(content) {
		return "", errors.New("导图内容格式错误")
	}
	var object map[string]interface{}
	if err := json.Unmarshal(content, &object); err != nil || object == nil {
		return "", errors.New("导图内容格式错误")
	}
	var buffer bytes.Buffer
	if err := json.Compact(&buffer, content); err != nil {
		return "", errors.New("导图内容格式错误")
	}
	return buffer.String(), nil
}

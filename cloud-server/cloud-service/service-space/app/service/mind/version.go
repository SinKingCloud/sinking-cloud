package mind

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"service-space/app/constant"
	"service-space/app/enum/mind_version_type"
	"service-space/app/model"
	mindDataRepository "service-space/app/repository/mind_data"
	mindVersionRepository "service-space/app/repository/mind_version"
	"strconv"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// SaveData 保存导图数据并生成版本记录
func (s *service) SaveData(userId int64, id int64, versionType int, remark string, value string) (*Version, error) {
	if userId <= 0 || id <= 0 {
		return nil, errors.New("参数错误")
	}
	if _, ok := mind_version_type.Map()[versionType]; !ok {
		return nil, errors.New("版本类型参数错误")
	}
	valueText, err := s.normalizeValue(value)
	if err != nil {
		return nil, err
	}
	remark = strings.TrimSpace(remark)
	if remark == "" {
		if versionType == mind_version_type.Auto {
			remark = "自动保存"
		} else {
			remark = "手动保存"
		}
	}
	settings := s.configService.GetSystemConfig(constant.MindGroup, false)
	versionMax, _ := strconv.Atoi(strings.TrimSpace(settings[constant.MindVersionMax]))
	if versionMax != 0 && versionMax < 10 {
		versionMax = 10
	}
	if versionMax > 10000 {
		versionMax = 10000
	}
	lockKey := constant.LockMindDataOperate + strconv.FormatInt(id, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeMindDataOperate) {
		return nil, errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	var savedVersion *model.MindVersion
	err = s.database.Transaction(func(tx *gorm.DB) error {
		info, err := s.mindDataRepository.FindByIdAndUserId(id, userId, tx)
		if err != nil {
			return errors.Join(errors.New("获取导图失败"), err)
		}
		if info == nil || info.Id <= 0 {
			return errors.New("导图不存在")
		}
		latest, err := s.mindVersionRepository.FindLatestByDataId(id, tx)
		if err != nil {
			return errors.Join(errors.New("获取最新版本失败"), err)
		}
		if versionType == mind_version_type.Auto && s.sameValue(info.Value, valueText) && latest != nil {
			savedVersion = latest
		} else {
			if err = s.mindDataRepository.UpdateByIdAndUserId(id, userId, &mindDataRepository.UpdateMindData{Value: valueText}, tx); err != nil {
				return errors.Join(errors.New("保存导图失败"), err)
			}
			merged := false
			if versionType == mind_version_type.Auto && latest != nil && latest.Type == mind_version_type.Auto {
				snapshotTime := time.Time(latest.CreateTime)
				interval := time.Now().Sub(snapshotTime)
				if interval >= 0 && interval < 10*time.Second {
					if err = s.mindVersionRepository.UpdateByIdAndDataId(latest.Id, id, &mindVersionRepository.UpdateMindVersion{
						Value:  valueText,
						Remark: remark,
					}, tx); err != nil {
						return errors.Join(errors.New("更新自动保存版本失败"), err)
					}
					savedVersion, err = s.mindVersionRepository.FindByIdAndDataId(latest.Id, id, tx)
					if err != nil {
						return errors.Join(errors.New("获取自动保存版本失败"), err)
					}
					merged = true
				}
			}
			if !merged {
				savedVersion = &model.MindVersion{
					Id:     str.GetSnowWorkIns().GetId(),
					DataId: id,
					Type:   versionType,
					Value:  valueText,
					Remark: remark,
				}
				if err = s.mindVersionRepository.Create(savedVersion, tx); err != nil {
					return errors.Join(errors.New("创建导图版本失败"), err)
				}
			}
		}
		if versionMax > 0 {
			if err = s.mindVersionRepository.DeleteOldByDataIdAndType(id, versionType, versionMax, tx); err != nil {
				return errors.Join(errors.New("清理旧版本失败"), err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	result := s.convertVersion(savedVersion)
	if result != nil {
		result.Current = true
	}
	return result, nil
}

// GetVersion 获取导图版本详情
func (s *service) GetVersion(userId int64, dataId int64, id int64) (*Version, error) {
	if userId <= 0 || dataId <= 0 || id <= 0 {
		return nil, errors.New("参数错误")
	}
	exists, err := s.mindDataRepository.ExistsByIdAndUserId(dataId, userId)
	if err != nil {
		return nil, errors.Join(errors.New("获取导图失败"), err)
	}
	if !exists {
		return nil, errors.New("导图不存在")
	}
	version, err := s.mindVersionRepository.FindByIdAndDataId(id, dataId)
	if err != nil {
		return nil, errors.Join(errors.New("获取导图版本失败"), err)
	}
	if version == nil || version.Id <= 0 {
		return nil, errors.New("导图版本不存在")
	}
	result := s.convertVersion(version)
	latestId, err := s.mindVersionRepository.FindLatestIdByDataId(dataId)
	if err != nil {
		return nil, errors.Join(errors.New("获取最新版本失败"), err)
	}
	result.Current = latestId == version.Id
	return result, nil
}

// RestoreVersion 恢复导图版本
func (s *service) RestoreVersion(userId int64, dataId int64, id int64) (*Data, error) {
	if userId <= 0 || dataId <= 0 || id <= 0 {
		return nil, errors.New("参数错误")
	}
	settings := s.configService.GetSystemConfig(constant.MindGroup, false)
	versionMax, _ := strconv.Atoi(strings.TrimSpace(settings[constant.MindVersionMax]))
	if versionMax != 0 && versionMax < 10 {
		versionMax = 10
	}
	if versionMax > 10000 {
		versionMax = 10000
	}
	lockKey := constant.LockMindDataOperate + strconv.FormatInt(dataId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeMindDataOperate) {
		return nil, errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	var restoredData *model.MindData
	var shared bool
	err := s.database.Transaction(func(tx *gorm.DB) error {
		info, err := s.mindDataRepository.FindByIdAndUserId(dataId, userId, tx)
		if err != nil {
			return errors.Join(errors.New("获取导图失败"), err)
		}
		if info == nil || info.Id <= 0 {
			return errors.New("导图不存在")
		}
		version, err := s.mindVersionRepository.FindByIdAndDataId(id, dataId, tx)
		if err != nil {
			return errors.Join(errors.New("获取导图版本失败"), err)
		}
		if version == nil || version.Id <= 0 {
			return errors.New("导图版本不存在")
		}
		if s.sameValue(info.Value, version.Value) {
			return errors.New("当前内容已与该版本一致")
		}
		latest, err := s.mindVersionRepository.FindLatestByDataId(dataId, tx)
		if err != nil {
			return errors.Join(errors.New("获取最新版本失败"), err)
		}
		if latest == nil || !s.sameValue(latest.Value, info.Value) {
			backup := &model.MindVersion{
				Id:     str.GetSnowWorkIns().GetId(),
				DataId: dataId,
				Type:   mind_version_type.Auto,
				Value:  info.Value,
				Remark: "恢复前备份",
			}
			if err = s.mindVersionRepository.Create(backup, tx); err != nil {
				return errors.Join(errors.New("备份当前导图失败"), err)
			}
		}
		if err = s.mindDataRepository.UpdateByIdAndUserId(dataId, userId, &mindDataRepository.UpdateMindData{Value: version.Value}, tx); err != nil {
			return errors.Join(errors.New("恢复导图失败"), err)
		}
		info.Value = version.Value
		info.UpdateTime = str.DateTime(time.Now())
		restoredData = info
		restored := &model.MindVersion{
			Id:     str.GetSnowWorkIns().GetId(),
			DataId: dataId,
			Type:   mind_version_type.Manual,
			Value:  version.Value,
			Remark: "恢复历史版本",
		}
		if err = s.mindVersionRepository.Create(restored, tx); err != nil {
			return errors.Join(errors.New("创建恢复版本失败"), err)
		}
		if versionMax > 0 {
			if err = s.mindVersionRepository.DeleteOldByDataIdAndType(dataId, mind_version_type.Auto, versionMax, tx); err != nil {
				return errors.Join(errors.New("清理自动保存版本失败"), err)
			}
			if err = s.mindVersionRepository.DeleteOldByDataIdAndType(dataId, mind_version_type.Manual, versionMax, tx); err != nil {
				return errors.Join(errors.New("清理手动保存版本失败"), err)
			}
		}
		sharedDataIds, err := s.mindShareRepository.SelectValidDataIds([]int64{dataId})
		if err != nil {
			return errors.Join(errors.New("获取分享信息失败"), err)
		}
		shared = len(sharedDataIds) > 0
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Data{
		Id:         restoredData.Id,
		TypeId:     restoredData.TypeId,
		Title:      restoredData.Title,
		Describe:   restoredData.Describe,
		Value:      restoredData.Value,
		Shared:     shared,
		CreateTime: restoredData.CreateTime,
		UpdateTime: restoredData.UpdateTime,
	}, nil
}

// SelectVersions 查询导图版本列表
func (s *service) SelectVersions(userId int64, dataId int64, versionType int, query *page.Query) ([]*Version, int64, error) {
	if userId <= 0 || dataId <= 0 {
		return nil, 0, errors.New("参数错误")
	}
	if query == nil {
		return nil, 0, errors.New("分页参数错误")
	}
	if _, ok := mind_version_type.Map()[versionType]; !ok {
		return nil, 0, errors.New("版本类型参数错误")
	}
	exists, err := s.mindDataRepository.ExistsByIdAndUserId(dataId, userId)
	if err != nil {
		return nil, 0, errors.Join(errors.New("获取导图失败"), err)
	}
	if !exists {
		return nil, 0, errors.New("导图不存在")
	}
	list, total, err := s.mindVersionRepository.Select(&mindVersionRepository.SelectMindVersion{
		DataId: dataId,
		Type:   strconv.Itoa(versionType),
	}, query)
	if err != nil {
		return nil, 0, errors.Join(errors.New("查询导图版本失败"), err)
	}
	latestId, err := s.mindVersionRepository.FindLatestIdByDataId(dataId)
	if err != nil {
		return nil, 0, errors.Join(errors.New("获取最新版本失败"), err)
	}
	result := make([]*Version, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		result = append(result, &Version{
			Id:         item.Id,
			DataId:     item.DataId,
			Type:       item.Type,
			Remark:     item.Remark,
			Current:    latestId == item.Id,
			CreateTime: item.CreateTime,
			UpdateTime: item.UpdateTime,
		})
	}
	return result, total, nil
}

// sameValue 判断导图内容是否一致
func (s *service) sameValue(left string, right string) bool {
	decode := func(value string) (interface{}, bool) {
		decoder := json.NewDecoder(strings.NewReader(value))
		decoder.UseNumber()
		var data interface{}
		if err := decoder.Decode(&data); err != nil {
			return nil, false
		}
		var extra interface{}
		if err := decoder.Decode(&extra); err != io.EOF {
			return nil, false
		}
		return data, true
	}
	leftData, leftValid := decode(left)
	rightData, rightValid := decode(right)
	if !leftValid || !rightValid {
		return left == right
	}
	return reflect.DeepEqual(leftData, rightData)
}

// convertVersion 转换导图版本返回数据
func (s *service) convertVersion(data *model.MindVersion) *Version {
	if data == nil {
		return nil
	}
	return &Version{
		Id:         data.Id,
		DataId:     data.DataId,
		Type:       data.Type,
		Value:      data.Value,
		Remark:     data.Remark,
		CreateTime: data.CreateTime,
		UpdateTime: data.UpdateTime,
	}
}

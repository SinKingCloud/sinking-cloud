package channel

import (
	"encoding/json"
	"errors"
	"service-driver/app/constant"
	"service-driver/app/enum/channel_status"
	"service-driver/app/enum/channel_type"
	"service-driver/app/model"
	"service-driver/app/repository/channel"
	"strconv"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Select 查询数据
func (s *service) Select(where *channel.SelectChannel, orderByField string, orderByType string, page int, pageSize int) (list []*channel.Channel, total int64, err error) {
	if where.Status != "" {
		status, err := strconv.Atoi(where.Status)
		if err != nil {
			return nil, 0, errors.New("状态参数错误")
		}
		if _, ok := channel_status.Map()[status]; !ok {
			return nil, 0, errors.New("状态参数不合法")
		}
	}
	if where.Type != "" {
		Type, err := strconv.Atoi(where.Type)
		if err != nil {
			return nil, 0, errors.New("类型参数错误")
		}
		if _, ok := channel_type.Map()[Type]; !ok {
			return nil, 0, errors.New("类型参数不合法")
		}
	}
	return s.repositoryChannel.Select(where, orderByField, orderByType, page, pageSize)
}

// GetIdNameMap 获取id name映射
func (s *service) GetIdNameMap(refresh bool) (map[int64]string, error) {
	cacheKey := constant.CacheNameWithEnum + "channel_names"
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			var result map[int64]string
			if err := json.Unmarshal([]byte(data), &result); err == nil {
				return result, nil
			}
		}
	}
	maps, err := s.repositoryChannel.SelectIdNameMap()
	if err != nil {
		return nil, err
	}
	if data, err2 := json.Marshal(maps); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithEnum)
	}
	return maps, nil
}

// DeleteById 通过ID删除数据
func (s *service) DeleteById(id int64) error {
	exists, err := s.repositoryScene.ExistsByChannelId(id)
	if err != nil {
		return errors.New("获取场景使用信息失败")
	}
	if exists {
		return errors.New("该通道下的场景数据仍在使用,请先清理场景数据")
	}
	exists, err = s.repositoryFile.ExistsByChannelId(id)
	if err != nil {
		return errors.New("获取文件使用信息失败")
	}
	if exists {
		return errors.New("该通道下的文件数据仍在使用,请先清理文件数据")
	}
	err = s.repositoryChannel.DeleteById(id)
	if err != nil {
		return errors.Join(errors.New("删除失败"), err)
	}
	return nil
}

// Create 创建数据
func (s *service) Create(name string, types int, status int, config string) error {
	if name == "" {
		return errors.New("通道名称不能为空")
	}
	if _, ok := channel_type.Map()[types]; !ok {
		return errors.New("通道类型不合法")
	}
	if _, ok := channel_status.Map()[status]; !ok {
		return errors.New("通道状态不合法")
	}
	if config != "" {
		value, err := s.checkSetting(config, types)
		if err != nil {
			return errors.Join(errors.New("设置不正确"), err)
		}
		config = value
	} else {
		config = "{}"
	}
	err := s.repositoryChannel.Create(&model.Channel{
		Type:       types,
		Name:       name,
		Config:     config,
		Status:     status,
		CreateTime: str.DateTime(time.Now()),
		UpdateTime: str.DateTime(time.Now()),
	})
	if err != nil {
		return errors.Join(errors.New("创建失败"), err)
	}
	return nil
}

// UpdateById 通过ID更新数据
func (s *service) UpdateById(id int64, data *channel.UpdateChannel) error {
	if id <= 0 {
		return errors.New("通道ID不合法")
	}
	info, err := s.FindById(id, false)
	if err != nil {
		return errors.New("获取通道信息失败")
	}
	update := &channel.UpdateChannel{}
	if data.Name != nil {
		if _, ok := data.Name.(string); !ok {
			return errors.New("通道名称字段不合法")
		}
		update.Name = data.Name
	}
	if data.Status != nil {
		if value, ok := data.Status.(int); ok {
			if _, ok = channel_status.Map()[value]; !ok {
				return errors.New("通道状态字段不合法")
			}
		} else {
			return errors.New("通道状态字段不合法")
		}
		update.Status = data.Status
	}
	if data.Config != nil {
		value, ok := data.Config.(string)
		if !ok {
			return errors.New("通道配置内容不合法")
		}
		value, err = s.checkSetting(value, info.Type)
		if err != nil {
			return errors.Join(errors.New("设置不正确"), err)
		}
		update.Config = value
	}
	err = s.repositoryChannel.UpdateById(id, update)
	if err != nil {
		return errors.Join(errors.New("更新失败"), err)
	}
	return nil
}

// FindById 通过ID查询数据
func (s *service) FindById(id int64, refresh bool) (*Channel, error) {
	cacheKey := constant.CacheNameWithChannel + strconv.FormatInt(id, 10)
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			var result *Channel
			if err := json.Unmarshal([]byte(data), &result); err == nil {
				return result, nil
			}
		}
	}
	info, err := s.repositoryChannel.FindById(id)
	if err != nil {
		return nil, err
	}
	data := &Channel{
		Id:         info.Id,
		Type:       info.Type,
		Name:       info.Name,
		Config:     s.formatSetting(info.Config, info.Type),
		Status:     info.Status,
		CreateTime: info.CreateTime,
		UpdateTime: info.UpdateTime,
	}
	if bytes, err2 := json.Marshal(data); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(bytes), constant.CacheTimeWithChannel)
	}
	return data, nil
}

// formatSetting 格式化设置
func (s *service) formatSetting(str string, types int) interface{} {
	var res interface{}
	switch types {
	case channel_type.Aliyun:
		data := AliyunSetting{}
		_ = json.Unmarshal([]byte(str), &data)
		res = data
		break
	default:
		res = map[string]string{}
		break
	}
	return res
}

// formatSetting 格式化设置
func (s *service) checkSetting(str string, types int) (res string, err error) {
	switch types {
	case channel_type.Aliyun:
		data := AliyunSetting{}
		_ = json.Unmarshal([]byte(str), &data)
		if data.AccessKeySecret == "" {
			err = errors.New("AccessKeySecret不能为空")
		}
		if data.AccessKeyID == "" {
			err = errors.New("AccessKeyID不能为空")
		}
		if data.Endpoint == "" {
			err = errors.New("endpoint不能为空")
		}
		if data.Bucket == "" {
			err = errors.New("bucket不能为空")
		}
		if data.Domain == "" {
			err = errors.New("domain不能为空")
		}
		if data.Path == "" {
			err = errors.New("path不能为空")
		}
	}
	return res, err
}

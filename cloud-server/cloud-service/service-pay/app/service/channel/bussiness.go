package channel

import (
	"encoding/json"
	"errors"
	"service-pay/app/constant"
	"service-pay/app/enum/channel_status"
	"service-pay/app/enum/channel_type"
	"service-pay/app/model"
	"service-pay/app/repository/channel"
	"strconv"
	"strings"
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

// DeleteById 通过ID删除数据
func (s *service) DeleteById(id int64) error {
	if id <= 0 {
		return errors.New("通道ID不合法")
	}
	exists, err := s.orderRepository.CheckChannelHasUnPayOrder(id)
	if err != nil {
		return errors.Join(errors.New("查询未支付订单失败"), err)
	}
	if exists {
		return errors.New("当前通道存在未过期且未支付的订单，不能删除")
	}
	err = s.repositoryChannel.DeleteById(id)
	if err != nil {
		return errors.Join(errors.New("删除失败"), err)
	}
	s.cache.Delete(constant.CacheNameWithChannel + strconv.FormatInt(id, 10))
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
	exists, err := s.orderRepository.CheckChannelHasUnPayOrder(id)
	if err != nil {
		return errors.Join(errors.New("查询未支付订单失败"), err)
	}
	if exists {
		return errors.New("当前通道存在未过期且未支付的订单，不能修改")
	}
	info, err := s.FindById(id, false)
	if err != nil || info == nil {
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
	s.cache.Delete(constant.CacheNameWithChannel + strconv.FormatInt(id, 10))
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
	if info == nil || info.Id <= 0 {
		return nil, errors.New("通道信息不存在")
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
	case channel_type.Alipay:
		data := Alipay{}
		_ = json.Unmarshal([]byte(str), &data)
		res = data
		break
	case channel_type.Epay:
		data := Epay{}
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
	case channel_type.Alipay:
		data := Alipay{}
		_ = json.Unmarshal([]byte(str), &data)
		if data.AppId == "" {
			err = errors.New("AppId不能为空")
		}
		if data.PrivateKey == "" {
			err = errors.New("PrivateKey不能为空")
		}
		if data.PublicKey == "" {
			err = errors.New("PublicKey不能为空")
		}
		bytes, err2 := json.Marshal(data)
		if err2 != nil {
			err = errors.New("配置格式化失败")
		} else {
			res = string(bytes)
		}
		break
	case channel_type.Epay:
		data := Epay{}
		_ = json.Unmarshal([]byte(str), &data)
		if data.Pid == "" {
			err = errors.New("pid不能为空")
		}
		if data.Key == "" {
			err = errors.New("key不能为空")
		}
		if data.Url == "" {
			err = errors.New("url不能为空")
		}
		bytes, err2 := json.Marshal(data)
		if err2 != nil {
			err = errors.New("配置格式化失败")
		} else {
			res = string(bytes)
		}
		break
	}
	return res, err
}

// GetConfig 获取支付配置
func (s *service) GetConfig(refresh bool) *Config {
	if s.configService == nil {
		return &Config{}
	}
	if !refresh {
		if data, exists := s.cache.Get(constant.CacheNameWithPayConfig); exists && data != "" {
			var result *Config
			if err := json.Unmarshal([]byte(data), &result); err == nil {
				return result
			}
		}
	}
	values := s.configService.GetSystemConfig(constant.PayGroup, false)

	minPay, _ := strconv.ParseFloat(strings.TrimSpace(values[constant.PayMin]), 64)
	maxPay, _ := strconv.ParseFloat(strings.TrimSpace(values[constant.PayMax]), 64)
	if minPay < 0 {
		minPay = 0
	}
	if maxPay < 0 {
		maxPay = 0
	}
	alipay, _ := strconv.ParseInt(strings.TrimSpace(values[constant.PayAlipay]), 10, 64)
	wxpay, _ := strconv.ParseInt(strings.TrimSpace(values[constant.PayWxpay]), 10, 64)
	if alipay < 0 {
		alipay = 0
	}
	if wxpay < 0 {
		wxpay = 0
	}
	result := &Config{
		Min:    minPay,
		Max:    maxPay,
		Alipay: alipay,
		Wxpay:  wxpay,
	}
	if bytes, err := json.Marshal(result); err == nil {
		s.cache.SetWithExpire(constant.CacheNameWithPayConfig, string(bytes), constant.CacheTimeWithPayConfig)
	}
	return result
}

// SetConfig 更新支付配置
func (s *service) SetConfig(data *Config) error {
	if s.configService == nil {
		return errors.New("配置服务不存在")
	}
	if data == nil {
		return errors.New("支付配置不能为空")
	}
	if data.Min > 0 && data.Max > 0 && data.Min > data.Max {
		return errors.New("最小充值金额不能大于最大充值金额")
	}
	if data.Alipay != 0 {
		info, err := s.FindById(data.Alipay, false)
		if err != nil {
			return errors.Join(errors.New("支付宝支付通道不存在"), err)
		}
		if info.Status != channel_status.Enable {
			return errors.New("支付宝支付通道已禁用")
		}
	}
	if data.Wxpay != 0 {
		info, err := s.FindById(data.Wxpay, false)
		if err != nil {
			return errors.Join(errors.New("微信支付通道不存在"), err)
		}
		if info.Status != channel_status.Enable {
			return errors.New("微信支付通道已禁用")
		}
	}
	if !s.configService.SetSystemConfigs(map[string]string{
		constant.PayMin:    strconv.FormatFloat(data.Min, 'f', -1, 64),
		constant.PayMax:    strconv.FormatFloat(data.Max, 'f', -1, 64),
		constant.PayAlipay: strconv.FormatInt(data.Alipay, 10),
		constant.PayWxpay:  strconv.FormatInt(data.Wxpay, 10),
	}) {
		return errors.New("更新支付配置失败")
	}
	s.cache.Delete(constant.CacheNameWithPayConfig)
	return nil
}

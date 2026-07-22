package channel

import (
	"encoding/json"
	"errors"
	"net/http"
	"service-pay/app/constant"
	"service-pay/app/enum/channel_plugin"
	"service-pay/app/enum/channel_status"
	"service-pay/app/enum/channel_type"
	"service-pay/app/enum/pay_type"
	"service-pay/app/model"
	"service-pay/app/repository/channel"
	"strconv"
	"strings"
	"time"

	commonConstant "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/pay"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Select 查询数据
func (s *service) Select(where *channel.SelectChannel, queryPage *page.Query) (*page.Result[*channel.Channel], error) {
	if where.Status != "" {
		status, err := strconv.Atoi(where.Status)
		if err != nil {
			return nil, errors.New("状态参数错误")
		}
		if _, ok := channel_status.Map()[status]; !ok {
			return nil, errors.New("状态参数不合法")
		}
	}
	if where.Plugin != "" {
		plugin, err := strconv.Atoi(where.Plugin)
		if err != nil {
			return nil, errors.New("类型参数错误")
		}
		if _, ok := channel_plugin.Map()[plugin]; !ok {
			return nil, errors.New("类型参数不合法")
		}
	}
	return s.repositoryChannel.Select(where, queryPage)
}

// GetIdNameMap 获取id name映射
func (s *service) GetIdNameMap(refresh bool) (map[int64]string, error) {
	cacheKey := constant.CacheNameWithChannelNameEnum
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
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithChannelNameEnum)
	}
	return maps, nil
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
	s.cache.Delete(constant.CacheNameWithChannelNameEnum, constant.CacheNameWithPayConfig, constant.CacheNameWithChannel+strconv.FormatInt(id, 10))
	return nil
}

// Create 创建数据
func (s *service) Create(name string, Type string, plugin int, status int, config string) error {
	if name == "" {
		return errors.New("通道名称不能为空")
	}
	if _, ok := channel_type.Map()[Type]; !ok {
		return errors.New("通道类型不合法")
	}
	if _, ok := channel_plugin.Map()[plugin]; !ok {
		return errors.New("通道插件不合法")
	}
	if _, ok := channel_status.Map()[status]; !ok {
		return errors.New("通道状态不合法")
	}
	if config != "" {
		value, err := s.checkSetting(config, plugin)
		if err != nil {
			return errors.Join(errors.New("设置不正确"), err)
		}
		config = value
	} else {
		config = "{}"
	}
	err := s.repositoryChannel.Create(&model.Channel{
		Type:       Type,
		Plugin:     plugin,
		Name:       name,
		Config:     config,
		Status:     status,
		CreateTime: str.DateTime(time.Now()),
		UpdateTime: str.DateTime(time.Now()),
	})
	if err != nil {
		return errors.Join(errors.New("创建失败"), err)
	}
	s.cache.Delete(constant.CacheNameWithChannelNameEnum)
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
		value, err = s.checkSetting(value, info.Plugin)
		if err != nil {
			return errors.Join(errors.New("设置不正确"), err)
		}
		update.Config = value
	}
	err = s.repositoryChannel.UpdateById(id, update)
	if err != nil {
		return errors.Join(errors.New("更新失败"), err)
	}
	s.cache.Delete(constant.CacheNameWithChannelNameEnum, constant.CacheNameWithPayConfig, constant.CacheNameWithChannel+strconv.FormatInt(id, 10))
	return nil
}

// FindById 通过ID查询数据
func (s *service) FindById(id int64, refresh bool) (*Channel, error) {
	cacheKey := constant.CacheNameWithChannel + strconv.FormatInt(id, 10)
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			var result *model.Channel
			if err := json.Unmarshal([]byte(data), &result); err == nil && result != nil && result.Id > 0 {
				return &Channel{
					Id:         result.Id,
					Type:       result.Type,
					Plugin:     result.Plugin,
					Name:       result.Name,
					Config:     s.formatSetting(result.Config, result.Plugin),
					Status:     result.Status,
					CreateTime: result.CreateTime,
					UpdateTime: result.UpdateTime,
				}, nil
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
	if bytes, err2 := json.Marshal(info); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(bytes), constant.CacheTimeWithChannel)
	}
	return &Channel{
		Id:         info.Id,
		Type:       info.Type,
		Plugin:     info.Plugin,
		Name:       info.Name,
		Config:     s.formatSetting(info.Config, info.Plugin),
		Status:     info.Status,
		CreateTime: info.CreateTime,
		UpdateTime: info.UpdateTime,
	}, nil
}

// GetConfig 获取支付配置
func (s *service) GetConfig(refresh bool) *Config {
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
	balance, _ := strconv.ParseInt(strings.TrimSpace(values[constant.PayBalance]), 10, 64)
	alipay, _ := strconv.ParseInt(strings.TrimSpace(values[constant.PayAlipay]), 10, 64)
	wxpay, _ := strconv.ParseInt(strings.TrimSpace(values[constant.PayWxpay]), 10, 64)
	expire, _ := strconv.ParseInt(strings.TrimSpace(values[constant.PayExpire]), 10, 64)
	if balance < 0 {
		balance = 0
	}
	if alipay < 0 {
		alipay = 0
	}
	if wxpay < 0 {
		wxpay = 0
	}
	if alipay != 0 {
		info, err := s.FindById(alipay, false)
		if err != nil || info == nil || info.Id <= 0 || info.Status != channel_status.Enable {
			alipay = 0
		}
	}
	if wxpay != 0 {
		info, err := s.FindById(wxpay, false)
		if err != nil || info == nil || info.Id <= 0 || info.Status != channel_status.Enable {
			wxpay = 0
		}
	}
	if expire < 5 {
		expire = 5
	}
	result := &Config{
		Min:     minPay,
		Max:     maxPay,
		Balance: balance,
		Alipay:  alipay,
		Wxpay:   wxpay,
		Expire:  expire,
	}
	if bytes, err := json.Marshal(result); err == nil {
		s.cache.SetWithExpire(constant.CacheNameWithPayConfig, string(bytes), constant.CacheTimeWithPayConfig)
	}
	return result
}

// SetConfig 更新支付配置
func (s *service) SetConfig(data *Config) error {
	if data == nil {
		return errors.New("支付配置不能为空")
	}
	if data.Min > 0 && data.Max > 0 && data.Min > data.Max {
		return errors.New("最小充值金额不能大于最大充值金额")
	}
	if data.Expire < 5 {
		data.Expire = 5
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
		constant.PayMin:     strconv.FormatFloat(data.Min, 'f', -1, 64),
		constant.PayMax:     strconv.FormatFloat(data.Max, 'f', -1, 64),
		constant.PayBalance: strconv.FormatInt(data.Balance, 10),
		constant.PayAlipay:  strconv.FormatInt(data.Alipay, 10),
		constant.PayWxpay:   strconv.FormatInt(data.Wxpay, 10),
		constant.PayExpire:  strconv.FormatInt(data.Expire, 10),
	}) {
		return errors.New("更新支付配置失败")
	}
	s.cache.Delete(constant.CacheNameWithPayConfig)
	return nil
}

// GetPayResultPageUrl 获取支付结果页地址
func (s *service) GetPayResultPageUrl(code int) string {
	config := s.configService.GetSystemConfig(commonConstant.WebGroup, false)
	webUrl := strings.TrimRight(strings.TrimSpace(config[commonConstant.WebUrl]), "/")
	return webUrl + "/pay?code=" + strconv.Itoa(code)
}

// GetPayUrl 获取支付跳转链接
func (s *service) GetPayUrl(channelId int64, orderId int64, name string, money float64) (string, error) {
	config := s.configService.GetSystemConfig(commonConstant.WebGroup, false)
	apiUrl := strings.TrimRight(strings.TrimSpace(config[commonConstant.WebApi]), "/")
	if apiUrl == "" {
		return "", errors.New("网站API地址不能为空")
	}
	client, channelInfo, err := s.getPay(channelId)
	if err != nil {
		return "", err
	}
	notifyUrl := apiUrl + strings.ReplaceAll("/"+s.consul.GetName()+"/callback/"+strconv.FormatInt(channelInfo.Id, 10)+"/notify", "//", "/")
	returnUrl := apiUrl + strings.ReplaceAll("/"+s.consul.GetName()+"/callback/"+strconv.FormatInt(channelInfo.Id, 10)+"/return", "//", "/")
	url, err := client.Submit(&pay.Order{
		Id:        orderId,
		Name:      name,
		Money:     money,
		NotifyUrl: notifyUrl,
		ReturnUrl: returnUrl,
		PayType:   channelInfo.Type,
	})
	if err != nil {
		return "", errors.New("获取支付信息失败")
	}
	return url, nil
}

// ValidateCallback 验证回调
func (s *service) ValidateCallback(channelId int64, req *http.Request, successFun func(data *pay.CallbackData) error) (code int, header http.Header, body string, data *pay.CallbackData, err error) {
	client, _, err := s.getPay(channelId)
	if err != nil {
		return 500, nil, "fail", nil, err
	}
	data, err = client.ParseCallback(req)
	if err == nil {
		if data == nil {
			err = errors.New("回调信息不能为空")
			code, header, body = client.FailResponse()
		} else if !data.Status {
			err = errors.New("该订单回调信息支付状态为未支付")
			code, header, body = client.FailResponse()
		} else if err = successFun(data); err != nil {
			code, header, body = client.FailResponse()
		} else {
			code, header, body = client.SuccessResponse()
		}
	} else {
		code, header, body = client.FailResponse()
	}
	return code, header, body, data, err
}

// GetPayChannel 获取支付通道
func (s *service) GetPayChannel(payType string) (channel *Channel, err error) {
	payConfig := s.GetConfig(false)
	var channelID int64
	if payType == pay_type.Alipay {
		channelID = payConfig.Alipay
	} else if payType == pay_type.Wxpay {
		channelID = payConfig.Wxpay
	} else if payType == pay_type.Balance {
		channelID = 0
	}
	if payType != pay_type.Balance && channelID <= 0 {
		return nil, errors.New("未开启该支付方式")
	}
	if payType == pay_type.Balance && payConfig.Balance <= 0 {
		return nil, errors.New("未开启该支付方式")
	}
	channel, err = &Channel{
		Id:     0,
		Type:   pay_type.Balance,
		Name:   "余额支付",
		Status: channel_status.Enable,
	}, nil
	if channelID > 0 {
		channel, err = s.FindById(channelID, false)
		if err != nil || channel == nil || channel.Id <= 0 {
			return nil, errors.Join(errors.New("支付通道不存在"), err)
		}
	}
	return channel, nil
}

// getPay 获取支付对象
func (s *service) getPay(id int64) (client pay.Interface, channel *Channel, err error) {
	info, err := s.FindById(id, false)
	if err != nil {
		return nil, nil, errors.Join(errors.New("获取支付通道信息失败"), err)
	}
	if info.Status != channel_status.Enable {
		return nil, nil, errors.Join(errors.New("获取支付通道信息失败"), errors.New("该通道已被禁用,请重新发起支付"))
	}
	switch info.Plugin {
	case channel_plugin.Alipay:
		config, ok := info.Config.(Alipay)
		if !ok {
			client, err = nil, errors.New("获取支付通道配置信息失败")
		} else {
			client, err = pay.NewAlipayGateway(config.AppId, config.PrivateKey, config.PublicKey)
		}
		break
	case channel_plugin.Epay:
		config, ok := info.Config.(Epay)
		if !ok {
			client, err = nil, errors.New("获取支付通道配置信息失败")
		} else {
			client, err = pay.NewEpayGateway(config.Pid, config.Key, config.Url)
		}
		break
	default:
		client = nil
		err = errors.New("获取支付通道类型失败")
		break
	}
	return client, info, err
}

// formatSetting 格式化设置
func (s *service) formatSetting(str string, plugin int) interface{} {
	var res interface{}
	switch plugin {
	case channel_plugin.Alipay:
		data := Alipay{}
		_ = json.Unmarshal([]byte(str), &data)
		res = data
		break
	case channel_plugin.Epay:
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

// checkSetting 判断设置
func (s *service) checkSetting(str string, plugin int) (res string, err error) {
	switch plugin {
	case channel_plugin.Alipay:
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
	case channel_plugin.Epay:
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

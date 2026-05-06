package auth

import (
	"service-auth/app/constant"
	"strconv"

	constant2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/aliyun"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// CheckSms 判断短信验证码正确性
func (s *service) CheckSms(phone string, code string) bool {
	key := constant.CacheNameWithSms + phone
	key2 := key + "Num"
	s.cache.Inc(key2)
	s.cache.SetExpire(key2, constant.CacheTimeWithSms)
	num, exists := s.cache.Get(key2)
	if exists && num != "" {
		num2, err := strconv.Atoi(key2)
		if err == nil && num2 >= 10 {
			s.cache.Delete(key, key2)
			return false
		}
	}
	data, exists := s.cache.Get(key)
	if exists {
		res := data == code
		if res {
			s.cache.Delete(key, key2)
		}
		return res
	}
	return false
}

// SendSms 发送短信验证码
func (s *service) SendSms(phone string) bool {
	//发送间隔
	key := constant.CacheNameWithSms + phone
	_, exists := s.cache.Get(key)
	if exists {
		return false
	}
	//发送验证码
	configs := s.configService.GetSystemConfig(constant2.AliGroup, false)
	code := strconv.FormatInt(str.NewStrTool().RandInt64(100000, 999999), 10)
	ali, err := aliyun.New(
		configs[constant2.AliKey],
		configs[constant2.AliSecret],
	)
	if err != nil {
		return false
	}
	if ali.Sms.Send(
		phone,
		configs[constant2.AliSmsSign],
		configs[constant2.AliSmsTemplate],
		"{\""+configs[constant2.AliSmsVar]+"\":"+code+"}",
	) {
		s.cache.SetWithExpire(key, code, constant.CacheTimeWithSms)
		return true
	}
	return false
}

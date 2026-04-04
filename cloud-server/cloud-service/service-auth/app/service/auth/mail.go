package auth

import (
	"service-auth/app/constant"
	"strconv"

	constant2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/aliyun"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// CheckMail 判断邮件验证码正确性
func (s *service) CheckMail(mail string, code string) bool {
	key := constant.CacheNameWithMail + mail
	key2 := key + "Num"
	s.cache.Inc(key2)
	s.cache.SetExpire(key2, constant.CacheTimeWithMail)
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

// SendMail 发送短信验证码
func (s *service) SendMail(mail string) bool {
	//发送间隔
	key := constant.CacheNameWithMail + mail
	_, exists := s.cache.Get(key)
	if exists {
		return false
	}
	//发送验证码
	configs := s.configService.GetSystemConfigs([]string{constant2.AliGroup, constant2.WebGroup}, false)
	ali := configs[constant2.AliGroup]
	web := configs[constant2.WebGroup]
	code := strconv.FormatInt(str.NewStrTool().RandInt64(100000, 999999), 10)
	if aliyun.Mail.Send(
		ali[constant2.AliKey],
		ali[constant2.AliSecret],
		ali[constant2.AliMailAccount],
		web[constant2.WebName], mail, "验证码通知", "您的验证码为："+code+"，请勿泄露于他人！",
	) {
		s.cache.SetWithExpire(key, code, constant.CacheTimeWithMail)
		return true
	}
	return false
}

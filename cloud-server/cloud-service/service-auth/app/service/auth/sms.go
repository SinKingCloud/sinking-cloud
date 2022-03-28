package auth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/aliyun"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/constant"
	"strconv"
)

// SmsCheck 判断验证码正确性
func SmsCheck(phone string, code string) bool {
	key := constant.CacheNameWithSmsPhone + phone
	data, exists := util.Cache.Get(key)
	if exists {
		res := data == code
		if res {
			util.Cache.Delete(key)
		}
		return res
	}
	return false
}

// SmsSend 发送短信验证码
func SmsSend(phone string) bool {
	//发送间隔
	key := constant.CacheNameWithExpire + phone
	_, exists := util.Cache.Get(key)
	if exists {
		return false
	}
	//发送验证码
	configs := service.GetSystemConfigServiceIns().GetByCache(false)
	code := strconv.FormatInt(str.NewStrTool().RandInt64(100000, 999999), 10)
	if aliyun.Sms.Send(
		configs[config.AliyunAccessKey],
		configs[config.AliyunAccessSecret],
		phone,
		configs[config.AliyunSmsSignName],
		configs[config.AliyunSmsTemplateCode],
		"{\"code\":"+code+"}",
	) {
		util.Cache.SetWithExpire(key, code, constant.CacheTimeWithExpire)
		util.Cache.SetWithExpire(constant.CacheNameWithSmsPhone+phone, code, constant.CacheTimeWithSmsPhone)
		return true
	}
	return false
}

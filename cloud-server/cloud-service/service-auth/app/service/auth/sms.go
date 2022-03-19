package auth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/constant"
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

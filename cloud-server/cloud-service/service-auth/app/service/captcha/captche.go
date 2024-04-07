package captcha

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/captcha"
)

const (
	prefix    = "Captcha_"
	cacheTime = 600 * 1000
)

// Get 获取验证码
func Get(key string) (*captcha.Image, string) {
	img, code := captcha.NewCaptchaRandCode(4)
	//写入redis
	util.Cache.SetWithExpire(prefix+key, code, cacheTime)
	return img, code
}

// Check 判断验证码是否正确
func Check(key string, code string) bool {
	value, exists := util.Cache.Get(prefix + key)
	if exists {
		if value == code {
			util.Cache.Delete(prefix + key)
			return true
		}
		return false
	}
	return false
}

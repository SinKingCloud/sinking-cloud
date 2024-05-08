package captcha

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/captcha"
)

// Check 判断验证码是否正确
func Check(randStr string, ticket string) bool {
	return captcha.Tencent.Check(randStr, ticket)
}

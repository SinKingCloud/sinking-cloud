package captcha

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/captcha"
)

type ControllerCaptcha struct {
}

// Check 判断验证码正确性
func (ControllerCaptcha) Check(c *sinking_server.Context) {
	data := c.GetRpcData()
	c.SuccessWithData("success", captcha.Check(data["key"].(string), data["code"].(string)))
}

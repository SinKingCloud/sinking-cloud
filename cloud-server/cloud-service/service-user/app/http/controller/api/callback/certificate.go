package callback

import (
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Certificate 实名回调
func (Controller) Certificate(c *context.Context) {
	var form struct {
		CertifyId string `json:"certifyId" default:"" validate:"required" label:"认证ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	_, err := service.User.VerifyRealAuth(form.CertifyId)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("实名成功")
}

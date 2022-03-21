package callback

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/service/rpc"
)

// Face 实人验证回调
func (ControllerCallback) Face(c *sinking_server.Context) {
	type Form struct {
		CertifyId string `form:"certifyId" validate:"required,len=32" label:"认证ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	if rpc.Face.Check(form.CertifyId) {
		c.Success("success")
	} else {
		c.Error("error")
	}
}

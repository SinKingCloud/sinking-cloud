package callback

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/rpc"
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
	if res, userId, cardNo, cardName := rpc.Face.Check(form.CertifyId); res && userId > 0 {
		_ = service.GetWebUserServiceIns().SetRealAuthStatu(userId, res, cardNo, cardName)
		c.Success("success")
	} else {
		c.Error("error")
	}
}

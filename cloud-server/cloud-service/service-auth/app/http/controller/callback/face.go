package callback

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/service/user"
)

// Face 实人验证回调
func (ControllerCallback) Face(c *sinking_server.Context) {
	type Form struct {
		Passed    string `form:"passed" validate:"required" label:"认证结果"`
		CertifyId string `form:"certifyId" validate:"required" label:"认证ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	if form.Passed == "T" {
		//认证成功
		config, err := service.GetWebUserConfSrvIns().FindByValue(form.CertifyId)
		if err == nil {
			err = user.GetIns().UpdateById(config.UserId, map[string]interface{}{
				"real_auth": "1",
			})
			if err == nil {
				service.GetWebUserServiceIns().ClearCache(config.UserId)
			}
		}
	}
	c.Success("success")
}

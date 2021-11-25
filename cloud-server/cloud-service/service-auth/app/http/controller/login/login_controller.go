package login

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
)

func (ControllerLogin) Index(c *sinking_server.Context) {
	user := service.GetWebUserServiceIns().FindById(1)
	token := jwt.GetToken(user, nil)
	c.SuccessWithData("请求成功", token)
}

package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	user2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_user"
)

type ControllerUser struct {
}

func (ControllerUser) FindById(c *sinking_server.Context) {
	data := c.GetRpcData()
	userId := int(data["user_id"].(float64))
	user := user2.GetIns().GetCurrentUserInfo(userId, true)
	if user == nil {
		c.Error("fail")
	} else {
		c.SuccessWithData("success", user)
	}
}

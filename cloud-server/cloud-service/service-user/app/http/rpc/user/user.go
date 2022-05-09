package user

import (
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	user2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/user"
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

func (ControllerUser) Money(c *sinking_server.Context) {
	data := c.GetRpcData()
	userId := int(data["user_id"].(float64))
	mod := int(data["mod"].(float64))
	money := data["money"].(string)
	err := service2.GetWebUserServiceIns().Money(userId, money, mod)
	if err != nil {
		c.Error("fail")
	} else {
		service2.GetWebUserServiceIns().ClearCache(userId)
		c.Success("success")
	}
}

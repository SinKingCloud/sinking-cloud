package oauth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth"
)

type ControllerOauth struct {
}

// AuthTicket 授权ticket
func (ControllerOauth) AuthTicket(c *sinking_server.Context) {
	data := c.GetRpcData()
	userId := int(data["user_id"].(float64))
	ticket := data["ticket"].(string)
	redirectUrl := data["redirect_url"].(string)
	e := web_oauth.AuthTicket(userId, ticket, redirectUrl)
	if e == nil {
		c.Success("success")
	} else {
		c.Error(e.Error())
	}
}

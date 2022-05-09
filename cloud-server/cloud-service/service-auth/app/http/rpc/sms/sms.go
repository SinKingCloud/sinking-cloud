package sms

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/auth"
)

type ControllerSms struct {
}

func (ControllerSms) Check(c *sinking_server.Context) {
	data := c.GetRpcData()
	c.SuccessWithData("success", auth.SmsCheck(data["phone"].(string), data["code"].(string)))
}

func (ControllerSms) Send(c *sinking_server.Context) {
	data := c.GetRpcData()
	c.SuccessWithData("success", auth.SmsSend(data["phone"].(string)))
}

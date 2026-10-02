package person

import (
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Info 账户详情
func (Controller) Info(c *context.Context) {
	u := c.GetUserInfo()
	info, err := service.User.GetInfo(u.Id, false)
	if err != nil || info == nil || info.Id <= 0 {
		c.Error("获取用户信息失败")
		return
	}
	var session sinking_web.H
	s := service.Auth.GetSessionInfo(u.Token)
	if s != nil {
		session = sinking_web.H{
			"device":      s.Device,
			"ip":          s.Ip,
			"location":    s.Location,
			"login_time":  s.CreateTime.String(),
			"expire_time": s.ExpireTime.String(),
		}
	}
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventShow, "查看账户信息", "查看账户信息")
	c.SuccessWithData("获取用户信息成功", sinking_web.H{
		"id":        info.Id,
		"avatar":    info.Avatar,
		"nick_name": info.NickName,
		"account":   info.Account,
		"phone":     str.NewStringTool(info.Phone).Show(3, 4, "*"),
		"email":     str.NewStringTool(info.Email).Show(3, 8, "*"),
		"password":  info.Password != "",
		"session":   session,
	})
}

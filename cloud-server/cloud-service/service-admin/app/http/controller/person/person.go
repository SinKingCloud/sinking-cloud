package person

import (
	"service-admin/app/repository/user"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

type Controller struct {
}

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

// Update 更新信息
func (Controller) Update(c *context.Context) {
	var form struct {
		NickName string `json:"nick_name" default:"" validate:"omitempty,min=1,max=15" label:"用户昵称"`
		Avatar   string `json:"avatar" default:"" validate:"omitempty,uri" label:"用户头像"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &user.UpdateUser{}
	if form.Avatar != "" {
		data.Avatar = form.Avatar
	}
	if form.NickName != "" {
		data.NickName = form.NickName
	}
	if err := service.User.UpdateById(c.GetUserInfo().Id, data); err == nil {
		c.Success("修改成功")
	} else {
		c.Error("修改失败")
	}
}

package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/service/user"
	"strconv"
)

// Update 修改资料
func (ControllerPerson) Update(c *sinking_server.Context) {
	type Form struct {
		Account  string `json:"account" default:"" validate:"omitempty,min=5,max=25" label:"账号"`
		Email    string `json:"email" default:"" validate:"omitempty,email" label:"邮箱"`
		NickName string `json:"nick_name" default:"" validate:"omitempty,min=1,max=15" label:"用户昵称"`
		Avatar   string `json:"avatar" default:"" validate:"omitempty,uri" label:"用户头像"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	data := make(map[string]interface{})
	if form.Account != "" {
		webUser := user.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, true)
		if webUser.Account == "" {
			if num, err := user.GetIns().CountByAccount(form.Account); err != nil || num > 0 {
				c.Error("该账户名已被占用")
				return
			}
			data["account"] = form.Account
		}
	}
	if form.Email != "" {
		data["email"] = form.Email
	}
	if form.NickName != "" {
		data["nick_name"] = form.NickName
	}
	if form.Avatar != "" {
		data["avatar"] = form.Avatar
	}
	if err := user.GetIns().UpdateById(c.GetWebUserInfo().Id, data); err == nil {
		service.GetWebUserServiceIns().ClearCache(c.GetWebUserInfo().Id)
		c.Success("更新资料成功")
	} else {
		c.Error("更新资料失败")
	}
}

// Password 修改密码
func (ControllerPerson) Password(c *sinking_server.Context) {
	type Form struct {
		Code     string `form:"code" json:"code" default:"" validate:"required,len=6" label:"短信验证码"`
		Password string `form:"password" json:"password" default:"" validate:"required,min=6,max=20" label:"密码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	if !rpc.Sms.Check(strconv.FormatInt(c.GetWebUserInfo().Phone, 10), form.Code) {
		c.Error("短信验证码校验失败")
		return
	}
	if err := user.GetIns().UpdateById(c.GetWebUserInfo().Id,
		map[string]interface{}{"password": str.NetStringTool(form.Password).GetPassword()}); err == nil {
		service.GetWebUserServiceIns().ClearCache(c.GetWebUserInfo().Id)
		c.Success("修改密码成功")
	} else {
		c.Error("修改密码失败")
	}
}

// Phone 修改手机号
func (ControllerPerson) Phone(c *sinking_server.Context) {
	type Form struct {
		Phone    string `form:"phone" json:"phone" default:"" validate:"required,len=11" label:"手机号码"`
		Code     string `form:"code" json:"code" default:"" validate:"required,len=6" label:"短信验证码"`
		Password string `form:"password" json:"password" default:"" validate:"required,min=6,max=20" label:"密码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	u := user.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, true)
	if !str.NetStringTool(u.Password).CheckPassword(form.Password) {
		c.Error("密码验证失败")
		return
	}
	if !rpc.Sms.Check(form.Phone, form.Code) {
		c.Error("短信验证码校验失败")
		return
	}
	if err := user.GetIns().UpdateById(c.GetWebUserInfo().Id,
		map[string]interface{}{"phone": form.Phone}); err == nil {
		service.GetWebUserServiceIns().ClearCache(u.Id)
		c.Success("修改绑定手机号成功")
	} else {
		c.Error("修改密码失败")
	}
}

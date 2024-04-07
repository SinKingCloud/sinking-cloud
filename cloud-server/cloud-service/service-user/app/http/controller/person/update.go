package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/user"
	"regexp"
	"strconv"
)

// Update 修改资料
func (ControllerPerson) Update(c *sinking_server.Context) {
	type Form struct {
		Account  string `form:"account" json:"account" default:"" validate:"omitempty,min=5,max=25" label:"账号"`
		Email    string `form:"email" json:"email" default:"" validate:"omitempty,email,max=30" label:"邮箱"`
		NickName string `form:"nick_name" json:"nick_name" default:"" validate:"omitempty,min=1,max=15" label:"用户昵称"`
		Avatar   string `form:"avatar" json:"avatar" default:"" validate:"omitempty,uri" label:"用户头像"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	data := make(map[string]interface{})
	if form.Account != "" {
		webUser := user.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, true)
		if webUser.Account == "" {
			if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{5,25}$", form.Account); !ok {
				c.Error("该账户名不合法")
				return
			}
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
		log.GetIns().CreateByQueue(c, log.EventEdit, "修改资料", "修改账户资料")
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
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	u := user.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	if !rpc.Sms.Check(strconv.FormatInt(u.Phone, 10), form.Code) {
		c.Error("短信验证码校验失败")
		return
	}
	if err := user.GetIns().UpdateById(c.GetWebUserInfo().Id,
		map[string]interface{}{"password": str.NewStringTool(form.Password).GetPassword()}); err == nil {
		log.GetIns().CreateByQueue(c, log.EventEdit, "修改密码", "修改账户密码")
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
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	u := user.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, true)
	if !str.NewStringTool(u.Password).CheckPassword(form.Password) {
		c.Error("密码验证失败")
		return
	}
	if !rpc.Sms.Check(form.Phone, form.Code) {
		c.Error("短信验证码校验失败")
		return
	}
	num, err := user.GetIns().CountByPhone(form.Phone)
	if err != nil || num > 0 {
		c.Error("该手机号已被他人占用")
		return
	}
	if err := user.GetIns().UpdateById(c.GetWebUserInfo().Id,
		map[string]interface{}{"phone": form.Phone}); err == nil {
		log.GetIns().CreateByQueue(c, log.EventEdit, "修改手机", "修改绑定手机")
		c.Success("修改绑定手机号成功")
	} else {
		c.Error("修改密码失败")
	}
}

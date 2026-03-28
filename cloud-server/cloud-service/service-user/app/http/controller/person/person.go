package person

import (
	"fmt"
	"service-user/app/enum/real_auth"
	"service-user/app/repository/user"
	"service-user/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

type Controller struct {
}

// Info 账户详情
func (Controller) Info(c *context.Context) {
	u := c.GetUserInfo()
	info, err := service.User.GetUserInfo(u.Id, false)
	if err != nil {
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
		"real_auth": info.RealAuthStatus == int(real_auth.Authed),
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

// Verify 安全验证
func (Controller) Verify(c *context.Context) {
	var form struct {
		Action string `json:"action" default:"get" validate:"required,oneof=get verify" label:"步骤"`
		Type   string `json:"type" default:"phone" validate:"required,oneof=mail phone" label:"验证方式"`
		Code   string `json:"code" default:"" validate:"required,numeric,len=6" label:"验证码"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	u, err := service.User.GetUserInfo(c.GetUserInfo().Id, false)
	if err != nil || u == nil || u.Id <= 0 {
		c.Error("获取用户信息失败")
		return
	}
	if form.Action == "get" {
		data := sinking_web.H{
			"status": service.User.CheckUserVerifyStatus(u.Id, false),
		}
		if u.Phone != "" {
			data["phone"] = str.NewStringTool(u.Phone).Show(3, 4, "*")
		}
		if u.Email != "" {
			data["mail"] = str.NewStringTool(u.Email).Show(3, 8, "*")
		}
		c.SuccessWithData("获取验证方式成功", data)
		return
	}
	account := u.Email
	check := service2.Auth.CheckMail
	if form.Type == "phone" {
		check = service2.Auth.CheckSms
		account = u.Phone
	}
	if !check(account, form.Code) {
		c.Error("验证码校验失败")
		return
	}
	minute := 10
	service.User.SetUserVerifyStatus(u.Id, true, minute)
	c.Success(fmt.Sprintf("验证成功,此次验证%d分钟内有效", minute))
}

// Password 修改密码
func (Controller) Password(c *context.Context) {
	var form struct {
		Password string `json:"password" default:"" validate:"required,min=6,max=20" label:"新密码"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	u, err := service.User.GetUserInfo(c.GetUserInfo().Id, false)
	if err != nil || u == nil || u.Id <= 0 {
		c.Error("获取用户信息失败")
		return
	}
	if !service.User.CheckUserVerifyStatus(u.Id, false) {
		c.Error("请先完成安全验证")
		return
	}
	if err = service.User.UpdateById(c.GetUserInfo().Id, &user.UpdateUser{
		Password: str.NewStringTool(form.Password).GetPassword(),
	}); err == nil {
		c.Success("修改密码成功")
	} else {
		c.Error("修改密码失败")
	}
}

// Phone 修改手机
func (Controller) Phone(c *context.Context) {
	var form struct {
		Phone string `json:"phone" default:"" validate:"required,numeric,len=11" label:"手机号"`
		Code  string `json:"code" default:"" validate:"required,numeric,len=6" label:"验证码"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	u, err := service.User.GetUserInfo(c.GetUserInfo().Id, false)
	if err != nil || u == nil || u.Id <= 0 {
		c.Error("获取用户信息失败")
		return
	}
	if !service.User.CheckUserVerifyStatus(u.Id, false) {
		c.Error("请先完成安全验证")
		return
	}
	if err = service.User.CheckPhoneIsUse(form.Phone); err != nil {
		c.Error(err.Error())
		return
	}
	if !service2.Auth.CheckSms(form.Phone, form.Code) {
		c.Error("验证码错误")
		return
	}
	if err = service.User.UpdateById(c.GetUserInfo().Id, &user.UpdateUser{
		Phone: form.Phone,
	}); err == nil {
		c.Success("修改手机号成功")
	} else {
		c.Error("修改手机号失败")
	}
}

// Email 修改邮箱
func (Controller) Email(c *context.Context) {
	var form struct {
		Mail string `json:"mail" default:"" validate:"required,email" label:"邮箱"`
		Code string `json:"code" default:"" validate:"required,numeric,len=6" label:"验证码"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	u, err := service.User.GetUserInfo(c.GetUserInfo().Id, false)
	if err != nil || u == nil || u.Id <= 0 {
		c.Error("获取用户信息失败")
		return
	}
	if !service.User.CheckUserVerifyStatus(u.Id, false) {
		c.Error("请先完成安全验证")
		return
	}
	if err = service.User.CheckEmailIsUse(form.Mail); err != nil {
		c.Error(err.Error())
		return
	}
	if !service2.Auth.CheckMail(form.Mail, form.Code) {
		c.Error("验证码错误")
		return
	}
	if err = service.User.UpdateById(c.GetUserInfo().Id, &user.UpdateUser{
		Email: form.Mail,
	}); err == nil {
		c.Success("修改邮箱成功")
	} else {
		c.Error("修改邮箱失败")
	}
}

// Account 修改账号
func (Controller) Account(c *context.Context) {
	var form struct {
		Account string `json:"account" default:"" validate:"required,alphanum,min=5,max=20" label:"账号"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	u, err := service.User.GetUserInfo(c.GetUserInfo().Id, false)
	if err != nil || u == nil || u.Id <= 0 {
		c.Error("获取用户信息失败")
		return
	}
	if u.Account != "" {
		c.Error("您的账户已修改过账号,无需修改")
		return
	}
	if err = service.User.CheckAccountIsUse(form.Account); err != nil {
		c.Error(err.Error())
		return
	}
	if err = service.User.UpdateById(c.GetUserInfo().Id, &user.UpdateUser{
		Account: form.Account,
	}); err == nil {
		c.Success("修改账号成功")
	} else {
		c.Error("修改账号失败")
	}
}

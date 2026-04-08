package security

import (
	"fmt"
	sessionRepository "service-user/app/repository/session"
	"service-user/app/repository/user"
	"service-user/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

type Controller struct {
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

// Session 会话列表
func (Controller) Session(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("create_time", "desc", "create_time,update_time,expire_time")
	var form struct {
		Device          string `json:"device" default:"" validate:"omitempty,max=100" label:"设备"`
		Ip              string `json:"ip" default:"" validate:"omitempty,max=50" label:"IP"`
		Location        string `json:"location" default:"" validate:"omitempty,max=100" label:"登录地点"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"登录起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"登录结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
		ExpireTimeStart string `json:"expire_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"过期起始时间"`
		ExpireTimeEnd   string `json:"expire_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"过期结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &sessionRepository.SelectSession{}
	if form.Device != "" {
		where.Device = form.Device
	}
	if form.Ip != "" {
		where.Ip = form.Ip
	}
	if form.Location != "" {
		where.Location = form.Location
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where.UpdateTimeStart = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where.UpdateTimeEnd = form.UpdateTimeEnd
	}
	if form.ExpireTimeStart != "" {
		where.ExpireTimeStart = form.ExpireTimeStart
	}
	if form.ExpireTimeEnd != "" {
		where.ExpireTimeEnd = form.ExpireTimeEnd
	}
	data, total, err := service.Auth.SelectSessionByUserId(c.GetUserInfo().Id, where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
}

// Logout 注销会话
func (Controller) Logout(c *context.Context) {
	var form struct {
		Tokens []string `json:"tokens" default:"" validate:"required,min=1,max=1000,unique" label:"会话token列表"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Auth.LogoutByUserIdAndTokens(c.GetUserInfo().Id, form.Tokens)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("注销会话成功")
}

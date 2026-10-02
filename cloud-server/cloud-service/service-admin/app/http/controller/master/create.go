package master

import (
	"service-admin/app/model"
	"service-admin/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Create 创建用户
func (Controller) Create(c *context.Context) {
	var form struct {
		Account  string `json:"account" default:"" validate:"required,alphanum,min=5,max=20" label:"账号"`
		Phone    string `json:"phone" default:"" validate:"omitempty,numeric,len=11" label:"手机号"`
		Email    string `json:"email" default:"" validate:"omitempty,email" label:"邮箱"`
		Password string `json:"password" default:"" validate:"required,min=6,max=20" label:"密码"`
		Avatar   string `json:"avatar" default:"" validate:"omitempty,uri" label:"头像"`
		NickName string `json:"nick_name" default:"" validate:"omitempty,min=1,max=15" label:"昵称"`
		Status   int    `json:"status" default:"0" validate:"numeric" label:"状态"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.User.Create(&model.User{
		Account:  form.Account,
		Phone:    form.Phone,
		Email:    form.Email,
		Password: str.NewStringTool(form.Password).GetPassword(),
		Avatar:   form.Avatar,
		NickName: form.NickName,
		Status:   form.Status,
	})
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventCreate, "创建管理员", "创建管理员["+form.Account+"]")
	c.Success("创建数据成功")
}

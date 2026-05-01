package master

import (
	"service-admin/app/repository/user"
	"service-admin/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Update 修改用户
func (Controller) Update(c *context.Context) {
	var form struct {
		Id       int64  `json:"id" default:"" validate:"required,min=1" label:"用户ID"`
		Account  string `json:"account" default:"" validate:"omitempty,alphanum,min=5,max=20" label:"账号"`
		Phone    string `json:"phone" default:"" validate:"omitempty,numeric,len=11" label:"手机号"`
		Email    string `json:"email" default:"" validate:"omitempty,email" label:"邮箱"`
		Password string `json:"password" default:"" validate:"omitempty,min=6,max=20" label:"密码"`
		Avatar   string `json:"avatar" default:"" validate:"omitempty,uri" label:"头像"`
		NickName string `json:"nick_name" default:"" validate:"omitempty,min=1,max=15" label:"昵称"`
		Status   string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &user.UpdateUser{}
	if form.Account != "" {
		data.Account = form.Account
	}
	if form.Phone != "" {
		data.Phone = form.Phone
	}
	if form.Email != "" {
		data.Email = form.Email
	}
	if form.Password != "" {
		data.Password = str.NewStringTool(form.Password).GetPassword()
	}
	if form.Avatar != "" {
		data.Avatar = form.Avatar
	}
	if form.NickName != "" {
		data.NickName = form.NickName
	}
	if form.Status != "" {
		data.Status, _ = strconv.Atoi(form.Status)
	}
	err := service.User.UpdateById(form.Id, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}

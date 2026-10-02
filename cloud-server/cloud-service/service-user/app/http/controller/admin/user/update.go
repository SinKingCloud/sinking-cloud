package user

import (
	"fmt"
	"service-user/app/repository/user"
	"service-user/app/service"
	"strconv"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Update 修改用户
func (Controller) Update(c *context.Context) {
	var form struct {
		Ids            []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"用户ID列表"`
		Account        string  `json:"account" default:"" validate:"omitempty,alphanum,min=5,max=20" label:"账号"`
		Phone          string  `json:"phone" default:"" validate:"omitempty,numeric,len=11" label:"手机号"`
		Email          string  `json:"email" default:"" validate:"omitempty,email" label:"邮箱"`
		Password       string  `json:"password" default:"" validate:"omitempty,min=6,max=20" label:"密码"`
		Avatar         string  `json:"avatar" default:"" validate:"omitempty,uri" label:"头像"`
		NickName       string  `json:"nick_name" default:"" validate:"omitempty,min=1,max=15" label:"昵称"`
		RealAuthType   string  `json:"real_auth_type" default:"" validate:"omitempty,numeric" label:"实名类型"`
		RealAuthStatus string  `json:"real_auth_status" default:"" validate:"omitempty,numeric" label:"实名状态"`
		Status         string  `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &user.UpdateUser{}
	if form.Account != "" {
		data.Account = &form.Account
	}
	if form.Phone != "" {
		data.Phone = &form.Phone
	}
	if form.Email != "" {
		data.Email = &form.Email
	}
	if form.Password != "" {
		password := str.NewStringTool(form.Password).GetPassword()
		data.Password = &password
	}
	if form.Avatar != "" {
		data.Avatar = &form.Avatar
	}
	if form.NickName != "" {
		data.NickName = &form.NickName
	}
	if form.RealAuthType != "" {
		realAuthType, err := strconv.Atoi(form.RealAuthType)
		if err != nil {
			c.Error("实名类型参数错误")
			return
		}
		data.RealAuthType = &realAuthType
	}
	if form.RealAuthStatus != "" {
		realAuthStatus, err := strconv.Atoi(form.RealAuthStatus)
		if err != nil {
			c.Error("实名状态参数错误")
			return
		}
		data.RealAuthStatus = &realAuthStatus
	}
	if form.Status != "" {
		status, err := strconv.Atoi(form.Status)
		if err != nil {
			c.Error("状态参数错误")
			return
		}
		data.Status = &status
	}
	var err error
	if len(form.Ids) == 1 {
		err = service.User.UpdateById(form.Ids[0], data)
	} else {
		err = service.User.UpdateByIds(form.Ids, data)
	}
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventEdit, "修改用户信息", fmt.Sprintf("修改用户%v", form.Ids))
	c.Success("更新数据成功")
}

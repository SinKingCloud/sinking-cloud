package user

import (
	"service-user/app/repository/user"
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// List 用户列表
func (Controller) List(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time,update_time")
	var form struct {
		Account         string `json:"account" default:"" validate:"omitempty,alphanum" label:"账号"`
		Phone           string `json:"phone" default:"" validate:"omitempty,numeric" label:"手机号"`
		Email           string `json:"email" default:"" validate:"omitempty" label:"邮箱"`
		NickName        string `json:"nick_name" default:"" validate:"omitempty,max=15" label:"昵称"`
		RealAuthType    string `json:"real_auth_type" default:"" validate:"omitempty,numeric" label:"实名类型"`
		RealAuthStatus  string `json:"real_auth_status" default:"" validate:"omitempty,numeric" label:"实名状态"`
		Status          string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &user.SelectUser{}
	if form.Account != "" {
		where.Account = form.Account
	}
	if form.Phone != "" {
		where.Phone = form.Phone
	}
	if form.Email != "" {
		where.Email = form.Email
	}
	if form.NickName != "" {
		where.NickName = form.NickName
	}
	if form.RealAuthType != "" {
		where.RealAuthType = form.RealAuthType
	}
	if form.RealAuthStatus != "" {
		where.RealAuthStatus = form.RealAuthStatus
	}
	if form.Status != "" {
		where.Status = form.Status
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
	data, total, err := service.User.Select(where, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, query.Page, query.PageSize, data))
}

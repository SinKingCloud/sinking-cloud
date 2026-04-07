package user

import (
	"errors"
	"service-admin/app/model"
	"service-admin/app/repository/user"
	"service-admin/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

type Controller struct {
}

// Role 分配角色
func (Controller) Role(c *context.Context) {
	var form struct {
		Id      int64   `json:"id" default:"" validate:"required,min=1" label:"用户ID"`
		RoleIds []int64 `json:"role_ids" default:"" validate:"required,unique" label:"角色ID列表"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.User.AssignRole(form.Id, form.RoleIds)
	if err != nil {
		c.Error(err.Error())
	} else {
		c.Success("分配角色成功")
	}
}

// List 用户列表
func (Controller) List(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,create_time,update_time")
	var form struct {
		Account         string `json:"account" default:"" validate:"omitempty,alphanum" label:"账号"`
		Phone           string `json:"phone" default:"" validate:"omitempty,numeric,len=11" label:"手机号"`
		Email           string `json:"email" default:"" validate:"omitempty,email" label:"邮箱"`
		NickName        string `json:"nick_name" default:"" validate:"omitempty,max=15" label:"昵称"`
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
	data, total, err := service.User.Select(where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
}

// Create 创建用户
func (Controller) Create(c *context.Context) {
	var form struct {
		Account  string `json:"account" default:"" validate:"required,alphanum,min=5,max=20" label:"账号"`
		Phone    string `json:"phone" default:"" validate:"omitempty,numeric,len=11" label:"手机号"`
		Email    string `json:"email" default:"" validate:"omitempty,email" label:"邮箱"`
		Password string `json:"password" default:"" validate:"required,min=6,max=20" label:"密码"`
		Avatar   string `json:"avatar" default:"" validate:"omitempty,uri" label:"头像"`
		NickName string `json:"nick_name" default:"" validate:"omitempty,min=1,max=15" label:"昵称"`
		Status   int    `json:"status" default:"0" validate:"required,numeric" label:"状态"`
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
	c.Success("创建数据成功")
}

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

// Delete 删除用户
func (Controller) Delete(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"required,min=1" label:"用户ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.User.DeleteById(form.Id)
	if err != nil {
		c.Error(errors.Join(errors.New("删除数据失败"), err))
		return
	}
	c.Success("删除数据成功")
}

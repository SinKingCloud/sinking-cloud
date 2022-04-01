package systems

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/system_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/system_user"
	"strconv"
)

// List 接口列表
func (ControllerUser) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id login_time update_time create_time status" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		RoleId          string `form:"role_id" default:"" validate:"omitempty,numeric" label:"角色ID"`
		Status          string `form:"status" default:"" validate:"omitempty,oneof=0 1" label:"状态"`
		Account         string `form:"account" default:"" validate:"omitempty,max=15" label:"用户账号"`
		NickName        string `form:"nick_name" default:"" validate:"omitempty,max=20" label:"用户昵称"`
		Email           string `form:"email" default:"" validate:"omitempty" label:"用户邮箱"`
		Phone           string `form:"phone" default:"" validate:"omitempty,numeric" label:"用户手机"`
		LoginIp         string `form:"login_ip" default:"" validate:"omitempty,max=20" label:"登陆IP"`
		LoginTimeStart  string `form:"login_time_start" default:"" validate:"-" label:"登陆起始时间"`
		LoginTimeEnd    string `form:"login_time_end" default:"" validate:"-" label:"登陆结束时间"`
		CreateTimeStart string `form:"create_time_start" default:"" validate:"-" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" default:"" validate:"-" label:"创建结束时间"`
		UpdateTimeStart string `form:"update_time_start" default:"" validate:"-" label:"更新起始时间"`
		UpdateTimeEnd   string `form:"update_time_end" default:"" validate:"-" label:"更新结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	if form.RoleId != "" {
		where["role_id"] = form.RoleId
	}
	if form.Status != "" {
		where["status"] = form.Status
	}
	if form.Account != "" {
		where["account"] = form.Account
	}
	if form.NickName != "" {
		where["nick_name"] = form.NickName
	}
	if form.Email != "" {
		where["email"] = form.Email
	}
	if form.Phone != "" {
		where["phone"] = form.Phone
	}
	if form.LoginIp != "" {
		where["login_ip"] = form.LoginIp
	}
	if form.LoginTimeStart != "" {
		where["login_time_start"] = form.LoginTimeStart
	}
	if form.LoginTimeEnd != "" {
		where["login_time_end"] = form.LoginTimeEnd
	}
	if form.CreateTimeStart != "" {
		where["create_time_start"] = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where["create_time_end"] = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where["update_time_start"] = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where["update_time_end"] = form.UpdateTimeEnd
	}
	data, total, err := system_user.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err == nil {
		system_log.GetIns().CreateByQueue(c, 0, "查看用户", "查看用户列表")
		c.SuccessWithData("获取用户数据成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
		return
	}
	c.Error("获取数据失败")
}

// Update 修改信息
func (ControllerUser) Update(c *sinking_server.Context) {
	type Form struct {
		Ids      []int  `json:"ids" default:"" validate:"required,min=1,unique" label:"ID列表"`
		RoleId   string `json:"role_id" default:"" validate:"omitempty,numeric" label:"角色ID"`
		Password string `json:"password" default:"" validate:"omitempty,max=20" label:"用户密码"`
		Avatar   string `json:"avatar" default:"" validate:"omitempty,uri" label:"用户头像"`
		Status   string `json:"status" default:"" validate:"omitempty,oneof=0 1" label:"状态"`
		NickName string `json:"nick_name" default:"" validate:"omitempty,max=20" label:"用户昵称"`
		Email    string `json:"email" default:"" validate:"omitempty,email" label:"用户邮箱"`
		Phone    string `json:"phone" default:"" validate:"omitempty,numeric,max=11" label:"用户手机"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	data := make(map[string]interface{})
	if form.RoleId != "" {
		data["role_id"] = form.RoleId
	}
	if form.Avatar != "" {
		data["avatar"] = form.Avatar
	}
	if form.Password != "" {
		data["password"] = str.NewStringTool(form.Password).GetPassword()
	}
	if form.Status != "" {
		data["status"] = form.Status
	}
	if form.NickName != "" {
		data["nick_name"] = form.NickName
	}
	if form.Email != "" {
		data["email"] = form.Email
	}
	if form.Phone != "" {
		data["phone"] = form.Phone
	}
	err := system_user.GetIns().UpdateByIds(form.Ids, data)
	if err == nil {
		//刷新缓存
		go func() {
			for _, v := range form.Ids {
				service.GetSystemUserServiceIns().ClearCache(v)
			}
		}()
		system_log.GetIns().CreateByQueue(c, 3, "修改用户", "修改用户信息")
		c.Success("修改数据成功")
	} else {
		c.Error("修改数据失败")
	}
}

// Create 添加信息
func (ControllerUser) Create(c *sinking_server.Context) {
	type Form struct {
		RoleId   int    `json:"role_id" default:"" validate:"numeric" label:"角色ID"`
		Status   int    `json:"status" default:"" validate:"oneof=0 1" label:"状态"`
		Account  string `json:"account" default:"" validate:"required,max=15" label:"用户账号"`
		Password string `json:"password" default:"" validate:"omitempty,max=20" label:"用户密码"`
		Avatar   string `json:"avatar" default:"" validate:"omitempty,uri" label:"用户头像"`
		NickName string `json:"nick_name" default:"" validate:"omitempty,max=20" label:"用户昵称"`
		Email    string `json:"email" default:"" validate:"omitempty,email" label:"用户邮箱"`
		Phone    string `json:"phone" default:"" validate:"required,numeric,len=11" label:"用户手机"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	if num, err := system_user.GetIns().CountByAccount(form.Account); err != nil || num > 0 {
		c.Error("该账户名已被占用")
		return
	}
	if num, err := system_user.GetIns().CountByPhone(form.Phone); err != nil || num > 0 {
		c.Error("该手机号已被占用")
		return
	}
	if form.Password == "" {
		form.Password = "123456"
	}
	phone, _ := strconv.ParseInt(form.Phone, 10, 64)
	err := system_user.GetIns().Create(&model.SystemUser{
		RoleId:   form.RoleId,
		Account:  form.Account,
		Password: str.NewStringTool(form.Password).GetPassword(),
		NickName: form.NickName,
		Avatar:   form.Avatar,
		Email:    form.Email,
		Phone:    phone,
		Status:   form.Status,
	})
	if err == nil {
		system_log.GetIns().CreateByQueue(c, 1, "添加用户", "添加用户记录")
		c.Success("添加数据成功")
	} else {
		c.Error("添加数据失败")
	}
}

// Delete 删除信息
func (ControllerUser) Delete(c *sinking_server.Context) {
	type Form struct {
		Ids []int `json:"ids" default:"" validate:"required,min=1,unique" label:"ID列表"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	err := system_user.GetIns().DeleteByIds(form.Ids)
	if err == nil {
		system_log.GetIns().CreateByQueue(c, 2, "删除用户", "删除用户记录")
		c.Success("删除数据成功")
	} else {
		c.Error("删除数据失败")
	}
}

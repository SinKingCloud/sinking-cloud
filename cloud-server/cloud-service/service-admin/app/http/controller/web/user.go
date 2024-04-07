package web

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/web_user"
)

// List 列表
func (ControllerUser) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id login_time update_time create_time status real_auth" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		RealAuth        string `form:"real_auth" json:"real_auth" default:"" validate:"omitempty,oneof=0 1" label:"实名"`
		Status          string `form:"status" json:"status" default:"" validate:"omitempty,oneof=0 1" label:"状态"`
		Account         string `form:"account" json:"account" default:"" validate:"omitempty,max=15" label:"用户账号"`
		NickName        string `form:"nick_name" json:"nick_name" default:"" validate:"omitempty,max=20" label:"用户昵称"`
		Email           string `form:"email" json:"email" default:"" validate:"omitempty" label:"用户邮箱"`
		Phone           string `form:"phone" json:"phone" default:"" validate:"omitempty,numeric" label:"用户手机"`
		LoginIp         string `form:"login_ip" json:"login_ip" default:"" validate:"omitempty,max=20" label:"登陆IP"`
		LoginTimeStart  string `form:"login_time_start" json:"login_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"登陆起始时间"`
		LoginTimeEnd    string `form:"login_time_end" json:"login_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"登陆结束时间"`
		CreateTimeStart string `form:"create_time_start" json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `form:"update_time_start" json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `form:"update_time_end" json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	if form.Status != "" {
		where["status"] = form.Status
	}
	if form.RealAuth != "" {
		where["real_auth"] = form.RealAuth
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
	data, total, err := web_user.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err == nil {
		system_log.GetIns().CreateByQueue(c, system_log.EventShow, "查看用户", "查看会员列表")
		c.SuccessWithData("获取用户数据成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
		return
	}
	c.Error("获取数据失败")
}

// Update 修改信息
func (ControllerUser) Update(c *sinking_server.Context) {
	type Form struct {
		Ids      []int  `form:"ids" json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
		RoleId   string `form:"role_id" json:"role_id" default:"" validate:"omitempty,numeric" label:"角色ID"`
		Password string `form:"password" json:"password" default:"" validate:"omitempty,max=20" label:"用户密码"`
		Avatar   string `form:"avatar" json:"avatar" default:"" validate:"omitempty,uri" label:"用户头像"`
		Status   string `form:"status" json:"status" default:"" validate:"omitempty,oneof=0 1" label:"状态"`
		NickName string `form:"nick_name" json:"nick_name" default:"" validate:"omitempty,max=20" label:"用户昵称"`
		Email    string `form:"email" json:"email" default:"" validate:"omitempty,email" label:"用户邮箱"`
		Phone    string `form:"phone" json:"phone" default:"" validate:"omitempty,numeric,max=11" label:"用户手机"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	data := make(map[string]interface{})
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
	if form.Phone != "" && len(form.Ids) == 1 {
		if num, err := web_user.GetIns().CountByPhone(form.Phone); err == nil && num <= 0 {
			data["phone"] = form.Phone
		}
	}
	err := web_user.GetIns().UpdateByIds(form.Ids, data)
	if err == nil {
		//刷新缓存
		go func() {
			for _, v := range form.Ids {
				service.GetWebUserServiceIns().ClearCache(v)
			}
		}()
		system_log.GetIns().CreateByQueue(c, system_log.EventEdit, "修改用户", "修改会员信息")
		c.Success("修改数据成功")
	} else {
		c.Error("修改数据失败")
	}
}

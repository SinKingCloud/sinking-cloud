package systems

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log"
)

// List 获取数据
func (ControllerLog) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id create_time" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Account         string `form:"account" default:"" validate:"omitempty,max=15" label:"用户账号"`
		NickName        string `form:"nick_name" default:"" validate:"omitempty,max=20" label:"用户昵称"`
		Email           string `form:"email" default:"" validate:"omitempty" label:"用户邮箱"`
		Phone           string `form:"phone" default:"" validate:"omitempty,numeric" label:"用户手机"`
		RequestId       string `form:"request_id" default:"" validate:"omitempty,len=36" label:"请求ID"`
		RequestIp       string `form:"request_ip" default:"" validate:"omitempty,max=20" label:"请求IP"`
		Type            string `form:"type" default:"" validate:"omitempty,numeric" label:"类型"`
		Title           string `form:"title" default:"" validate:"omitempty,max=50" label:"标题"`
		Content         string `form:"content" default:"" validate:"omitempty,max=100" label:"内容"`
		CreateTimeStart string `form:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	if form.Type != "" {
		where["type"] = form.Type
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
	if form.RequestId != "" {
		where["request_id"] = form.RequestId
	}
	if form.RequestIp != "" {
		where["request_ip"] = form.RequestIp
	}
	if form.Title != "" {
		where["title"] = form.Title
	}
	if form.Content != "" {
		where["content"] = form.Content
	}
	if form.CreateTimeStart != "" {
		where["create_time_start"] = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where["create_time_end"] = form.CreateTimeEnd
	}
	data, total, err := system_log.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取日志数据失败")
	} else {
		system_log.GetIns().CreateByQueue(c, system_log.EventShow, "查看日志", "查看用户日志")
		c.SuccessWithData("获取日志数据成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

// Delete 删除记录
func (ControllerLog) Delete(c *sinking_server.Context) {
	type Form struct {
		RequestTimeStart string `json:"request_time_start" default:"" validate:"required" label:"创建起始时间"`
		RequestTimeEnd   string `json:"request_time_end" default:"" validate:"required" label:"创建结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	err := system_log.GetIns().DeleteByRequestTime(form.RequestTimeStart, form.RequestTimeEnd)
	if err == nil {
		system_log.GetIns().CreateByQueue(c, system_log.EventDelete, "清理日志", "清理用户日志")
		c.Success("删除数据成功")
	} else {
		c.Error("删除数据失败")
	}
}

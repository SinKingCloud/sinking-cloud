package systems

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log"
)

// List 获取数据
func (ControllerLog) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `json:"order_by_field" default:"id" validate:"oneof=id create_time" label:"排序字段"`
		OrderByType     string `json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Account         string `json:"account" default:"" validate:"omitempty,max=15" label:"用户账号"`
		NickName        string `json:"nick_name" default:"" validate:"omitempty,max=20" label:"用户昵称"`
		Email           string `json:"email" default:"" validate:"omitempty" label:"用户邮箱"`
		Phone           string `json:"phone" default:"" validate:"omitempty,numeric" label:"用户手机"`
		RequestId       string `json:"request_id" default:"" validate:"omitempty,len=36" label:"请求ID"`
		RequestIp       string `json:"request_ip" default:"" validate:"omitempty,max=20" label:"请求IP"`
		Type            string `json:"type" default:"" validate:"omitempty,numeric" label:"类型"`
		Title           string `json:"title" default:"" validate:"omitempty,max=50" label:"标题"`
		Content         string `json:"content" default:"" validate:"omitempty,max=100" label:"内容"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
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
	data, total, err := service.SystemLog.Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取日志数据失败")
	} else {
		service.SystemLog.CreateByQueue(c, system_log.EventShow, "查看日志", "查看用户日志")
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
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	err := service.SystemLog.DeleteByRequestTime(form.RequestTimeStart, form.RequestTimeEnd)
	if err == nil {
		service.SystemLog.CreateByQueue(c, system_log.EventDelete, "清理日志", "清理用户日志")
		c.Success("删除数据成功")
	} else {
		c.Error("删除数据失败")
	}
}

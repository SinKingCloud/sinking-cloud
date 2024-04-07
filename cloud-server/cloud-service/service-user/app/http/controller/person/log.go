package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/log"
	"strconv"
)

func (ControllerPerson) Log(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id create_time" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		RequestId       string `form:"request_id" json:"request_id" default:"" validate:"omitempty,len=36" label:"请求ID"`
		RequestIp       string `form:"request_ip" json:"request_ip" default:"" validate:"omitempty,max=20" label:"请求IP"`
		Type            string `form:"type" json:"type" default:"" validate:"omitempty,numeric" label:"类型"`
		Title           string `form:"title" json:"title" default:"" validate:"omitempty,max=50" label:"标题"`
		Content         string `form:"content" json:"content" default:"" validate:"omitempty,max=100" label:"内容"`
		CreateTimeStart string `form:"create_time_start" json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	where["user_id"] = strconv.Itoa(c.GetWebUserInfo().Id)
	if form.Type != "" {
		where["type"] = form.Type
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
	data, total, err := log.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取操作记录失败")
	} else {
		log.GetIns().CreateByQueue(c, log.EventShow, "查看日志", "查看操作日志")
		c.SuccessWithData("获取操作记录成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

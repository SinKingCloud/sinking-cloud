package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_app"
	"strconv"
)

func (ControllerApp) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id create_time" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Key             string `form:"key" default:"" validate:"omitempty,max=32" label:"密匙"`
		Name            string `form:"name" default:"" validate:"omitempty,max=50" label:"名称"`
		Describe        string `form:"describe" default:"" validate:"omitempty,max=50" label:"描述"`
		CreateTimeStart string `form:"create_time_start" default:"" validate:"-" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" default:"" validate:"-" label:"创建结束时间"`
		UpdateTimeStart string `form:"update_time_start" default:"" validate:"-" label:"修改起始时间"`
		UpdateTimeEnd   string `form:"update_time_end" default:"" validate:"-" label:"修改结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	where["user_id"] = strconv.Itoa(c.GetWebUserInfo().Id)
	if form.Key != "" {
		where["key"] = form.Key
	}
	if form.Name != "" {
		where["name"] = form.Name
	}
	if form.Describe != "" {
		where["describe"] = form.Describe
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
	data, total, err := web_app.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取应用列表失败")
	} else {
		rpc.WebUserLog.Create(c.GetWebUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), 1, "查看应用", "查看应用列表")
		c.SuccessWithData("获取应用列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

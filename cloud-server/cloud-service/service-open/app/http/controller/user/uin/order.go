package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/order"
	"strconv"
)

type ControllerOrder struct {
}

func (ControllerOrder) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		AppId           string `form:"app_id" default:"" validate:"omitempty,numeric" label:"应用ID"`
		UinId           string `form:"uin_id" default:"" validate:"omitempty,numeric" label:"QQID"`
		Type            string `form:"uin_id" default:"" validate:"omitempty,,oneof=0 1 2" label:"订单类型"`
		Name            string `form:"name" default:"" validate:"omitempty,max=50" label:"标题"`
		Status          string `form:"status" default:"" validate:"omitempty,oneof=0 1" label:"状态"`
		CreateTimeStart string `form:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	where["user_id"] = strconv.Itoa(c.GetWebUserInfo().Id)
	if form.AppId != "" {
		where["app_id"] = form.AppId
	}
	if form.UinId != "" {
		where["uin_id"] = form.UinId
	}
	if form.Status != "" {
		where["status"] = form.Status
	}
	if form.Name != "" {
		where["name"] = form.Name
	}
	if form.CreateTimeStart != "" {
		where["create_time_start"] = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where["create_time_end"] = form.CreateTimeEnd
	}
	data, total, err := order.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取订单列表失败")
	} else {
		c.SuccessWithData("获取订单列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/web_uin/uin"
	"strconv"
)

type ControllerLog struct {
}

func (ControllerLog) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id" label:"排序字段"`
		OrderByType  string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		AppId        string `form:"app_id" json:"app_id" default:"" validate:"omitempty,numeric" label:"应用ID"`
		Uin          string `form:"uin" json:"uin" default:"" json:"uin" validate:"required,numeric,min=5,max=10" label:"QQ账号"`
		Title        string `form:"title" json:"title" default:"" validate:"omitempty,max=50" label:"标题"`
		Content      string `form:"content" json:"content" default:"" validate:"omitempty,max=50" label:"描述"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	uinInfo, err := uin.GetIns().FindByUserIdWithUin(c.GetWebUserInfo().Id, form.Uin)
	if err != nil || uinInfo == nil {
		c.Error("未查询到QQ信息,请先添加QQ")
		return
	}
	where := make(map[string]string)
	where["user_id"] = strconv.Itoa(c.GetWebUserInfo().Id)
	where["uin_id"] = strconv.Itoa(uinInfo.Id)
	if form.AppId != "" {
		where["app_id"] = form.AppId
	}
	if form.Title != "" {
		where["title"] = form.Title
	}
	if form.Content != "" {
		where["content"] = form.Content
	}
	data, total, err := log.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取日志列表失败")
	} else {
		c.SuccessWithData("获取日志列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

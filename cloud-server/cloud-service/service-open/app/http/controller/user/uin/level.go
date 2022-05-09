package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/level"
)

type ControllerLevel struct {
}

func (ControllerLevel) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id" label:"排序字段"`
		OrderByType  string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Name         string `form:"name" default:"" validate:"omitempty,max=50" label:"名称"`
		Describe     string `form:"describe" default:"" validate:"omitempty,max=50" label:"描述"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	where["status"] = "0"
	if form.Name != "" {
		where["name"] = form.Name
	}
	if form.Describe != "" {
		where["describe"] = form.Describe
	}
	data, total, err := level.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取折扣列表失败")
	} else {
		c.SuccessWithData("获取折扣列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

// Buy 开通等级
func (ControllerLevel) Buy(c *sinking_server.Context) {
	type Form struct {
		LevelId int `form:"level_id" json:"level_id" default:"" validate:"gte=0" label:"等级ID"`
		Month   int `form:"month" json:"month" default:"" validate:"gte=1" label:"开通月份"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	e := web_uin.BuyLevel(c.GetWebUserInfo().Id, form.LevelId, form.Month, true)
	if e == nil {
		c.Success("开通折扣等级成功")
	} else {
		c.Error(e.Error())
	}
}

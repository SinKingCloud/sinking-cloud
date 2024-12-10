package auth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service"
)

// ShopList 商品列表
func (ControllerAuth) ShopList(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `json:"order_by_field" default:"id" validate:"oneof=id create_time" label:"排序字段"`
		OrderByType     string `json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Name            string `json:"name" default:"" validate:"omitempty,max=50" label:"名称"`
		Describe        string `json:"describe" default:"" validate:"omitempty,len=100" label:"描述"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"修改起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"修改结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
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
	data, total, err := service.AuthApp.Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取列表失败")
	} else {
		c.SuccessWithData("获取列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

// ShopBuy 购买商品
func (ControllerAuth) ShopBuy(c *sinking_server.Context) {
	type Form struct {
		Id     int    `json:"id" default:"" validate:"min=0" label:"程序ID"`
		Domain string `json:"domain" default:"" validate:"required,max=100" label:"域名"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	app, err := service.AuthApp.FindByIdWithCache(form.Id, false)
	if err != nil || app == nil {
		c.Error("未查询到该应用信息")
		return
	}
	e := service.WebAuthShop.Buy(c.GetWebUserInfo().Id, form.Id, form.Domain)
	if e != nil {
		c.Error(e.Error())
		return
	}
	c.Success("购买成功")
}

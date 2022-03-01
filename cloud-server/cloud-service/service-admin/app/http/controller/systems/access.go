package systems

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/access"
)

func (ControllerAccess) List(c *sinking_server.Context) {
	type Form struct {
		Page            int    `form:"page" default:"1" validate:"required,gte=1" label:"页码编号"`
		PageSize        int    `form:"page_size" default:"20" validate:"required,lte=1000" label:"每页数量"`
		OrderByField    string `form:"order_by_field" default:"id" validate:"oneof=id" label:"排序字段"`
		OrderByType     string `form:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Name            string `form:"name" default:"" validate:"max=20" label:"名称"`
		Describe        string `form:"describe" default:"" validate:"max=100" label:"描述"`
		Url             string `form:"url" default:"" validate:"max=50" label:"请求地址"`
		Method          string `form:"method" default:"ALL" validate:"oneof=GET POST PUT UPDATE ALL" label:"请求类型"`
		Status          string `form:"status" default:"ALL" validate:"oneof=ALL 0 1" label:"状态"`
		CreateTimeStart string `form:"create_time_start" default:"" validate:"-" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" default:"" validate:"-" label:"创建结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
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
	if form.Url != "" {
		where["url"] = form.Url
	}
	if form.Method != "ALL" {
		where["method"] = form.Method
	}
	if form.Status != "ALL" {
		where["status"] = form.Status
	}
	if form.CreateTimeStart != "" {
		where["create_time_start"] = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where["create_time_end"] = form.CreateTimeEnd
	}
	err, total, data := access.GetIns().GetAccessList(where, form.OrderByField, form.OrderByType, form.Page, form.PageSize)
	if err == nil {
		c.SuccessWithData("获取权限数据成功", page.NewPage(total, form.Page, form.PageSize, data))
		return
	}
	c.Error("获取数据失败")
}

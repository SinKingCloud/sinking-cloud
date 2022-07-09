package auth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/auth_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/auth_version"
	"strconv"
)

// AppInfo 程序信息
func (ControllerAuth) AppInfo(c *sinking_server.Context) {
	type Form struct {
		Id int `form:"id" default:"" validate:"min=0" label:"程序ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	app, err := auth_app.GetIns().FindByIdWithCache(form.Id, false)
	if err != nil || app == nil {
		c.Error("未查询到该应用信息")
		return
	}
	c.SuccessWithData("获取成功", app)
}

// AppVersion 程序版本
func (ControllerAuth) AppVersion(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField  string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id status create_time" label:"排序字段"`
		OrderByType   string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		AuthId        int    `form:"auth_id" default:"" validate:"min=0" label:"程序ID"`
		VersionNumber string `form:"version_number" default:"" validate:"omitempty,numeric" label:"程序ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	app, err := auth_app.GetIns().FindByIdWithCache(form.AuthId, false)
	if err != nil || app == nil {
		c.Error("未查询到该应用信息")
		return
	}
	where := map[string]string{"auth_id": strconv.Itoa(form.AuthId)}
	if form.VersionNumber != "" {
		where["version_number"] = form.VersionNumber
	}
	data, total, err := auth_version.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取列表失败")
	} else {
		c.SuccessWithData("获取列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

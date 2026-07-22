package permission

import (
	"service-admin/app/repository/permission"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// List 权限列表
func (Controller) List(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time,update_time", "id,create_time,update_time")
	var form struct {
		Code            string `json:"code" default:"" validate:"omitempty" label:"标识"`
		Name            string `json:"name" default:"" validate:"omitempty" label:"名称"`
		Path            string `json:"path" default:"" validate:"omitempty" label:"路径"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &permission.SelectPermission{}
	if form.Code != "" {
		where.Code = form.Code
	}
	if form.Name != "" {
		where.Name = form.Name
	}
	if form.Path != "" {
		where.Path = form.Path
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where.UpdateTimeStart = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where.UpdateTimeEnd = form.UpdateTimeEnd
	}
	data, err := service.Rbac.SelectPermission(where, query)
	if err != nil {
		c.Error("获取数据失败")
	} else {
		c.SuccessWithData("获取数据成功", data)
	}
}

package menu

import (
	"service-admin/app/repository/menu"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Controls 控件列表
func (Controller) Controls(c *context.Context) {
	query := c.ValidatePage("id", "asc", "id,status,create_time,update_time")
	var form struct {
		Id              int64  `json:"id" default:"" validate:"numeric,min=1" label:"页面ID"`
		Code            string `json:"code" default:"" validate:"omitempty" label:"标识"`
		Name            string `json:"name" default:"" validate:"omitempty" label:"名称"`
		Icon            string `json:"icon" default:"" validate:"omitempty" label:"图标"`
		Status          string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &menu.SelectMenu{}
	if form.Code != "" {
		where.Code = form.Code
	}
	if form.Name != "" {
		where.Name = form.Name
	}
	if form.Icon != "" {
		where.Icon = form.Icon
	}
	if form.Status != "" {
		where.Status = form.Status
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
	list, total, err := service.Rbac.GetPageControlsByMenuId(form.Id, where, query)
	if err != nil {
		c.Error("获取数据失败")
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, query.Page, query.PageSize, list))
}

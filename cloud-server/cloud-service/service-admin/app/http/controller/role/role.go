package role

import (
	"service-admin/app/repository/role"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

type Controller struct {
}

// List 权限列表
func (Controller) List(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,create_time,update_time")
	var form struct {
		Name            string `json:"name" default:"" validate:"omitempty" label:"名称"`
		Status          string `json:"status" default:"" validate:"omitempty" label:"状态"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &role.SelectRole{}
	if form.Name != "" {
		where.Name = form.Name
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
	data, total, err := service.Rbac.SelectRole(where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error("获取数据失败")
	} else {
		c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
	}
}

// Info 角色信息
func (Controller) Info(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
}

// Create 创建角色
func (Controller) Create(c *context.Context) {
	var form struct {
		Name   string `json:"name" default:"" validate:"required,max=30" label:"名称"`
		Status int    `json:"status" default:"0" validate:"numeric" label:"状态"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
}

// Update 修改角色
func (Controller) Update(c *context.Context) {
	var form struct {
		Id int64 `json:"ids" default:"" validate:"required,numeric,min=1" label:"ID列表"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
}

// Delete 删除角色
func (Controller) Delete(c *context.Context) {
	var form struct {
		Ids []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
}

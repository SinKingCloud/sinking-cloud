package permission

import (
	"service-admin/app/model"
	"service-admin/app/repository/permission"
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
	data, total, err := service.Rbac.SelectPermission(where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error("获取数据失败")
	} else {
		c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
	}
}

// Create 创建权限
func (Controller) Create(c *context.Context) {
	var form struct {
		Code string `json:"place" default:"" validate:"required,max=50" label:"标识"`
		Name string `json:"title" default:"" validate:"required,max=20" label:"名称"`
		Path string `json:"describe" default:"" validate:"required,max=100" label:"路径"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Rbac.CreatePermission(&model.Permission{
		Name: form.Name,
		Code: form.Code,
		Path: form.Path,
	})
	if err == nil {
		c.Success("创建数据成功")
	} else {
		c.Error(err.Error())
	}
}

// Update 修改权限
func (Controller) Update(c *context.Context) {
	var form struct {
		Ids  []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
		Name string  `json:"title" default:"" validate:"omitempty,max=100" label:"名称"`
		Path string  `json:"describe" default:"" validate:"omitempty,max=200" label:"描述"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &permission.UpdatePermission{}
	if form.Name != "" {
		data.Name = form.Name
	}
	if form.Path != "" {
		data.Path = form.Path
	}
	err := service.Rbac.UpdatePermissionByIds(form.Ids, data)
	if err != nil {
		c.Error("修改数据失败")
	} else {
		c.Success("修改数据成功")
	}
}

// Delete 删除权限
func (Controller) Delete(c *context.Context) {
	var form struct {
		Ids []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Rbac.DeletePermissionByIds(form.Ids)
	if err != nil {
		c.Error("删除数据失败")
	} else {
		c.Success("删除数据失败")
	}
}

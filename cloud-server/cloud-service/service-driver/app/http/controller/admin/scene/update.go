package scene

import (
	"service-driver/app/repository/scene"
	"service-driver/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Update 更新数据
func (Controller) Update(c *context.Context) {
	var form struct {
		Id       int64  `json:"id" default:"" validate:"required,min=1" label:"记录ID"`
		Name     string `json:"name" default:"" validate:"omitempty" label:"名称"`
		MaxSize  string `json:"max_size" default:"" validate:"omitempty,numeric" label:"最大限制"`
		AllowExt string `json:"allow_ext" default:"" validate:"omitempty" label:"允许类型"`
		Remark   string `json:"remark" default:"" validate:"omitempty" label:"备注"`
		Status   string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &scene.UpdateScene{}
	if form.Name != "" {
		data.Name = form.Name
	}
	if form.MaxSize != "" {
		data.MaxSize, _ = strconv.ParseInt(form.MaxSize, 10, 64)
	}
	if form.AllowExt != "" {
		data.AllowExt = form.AllowExt
	}
	if form.Remark != "" {
		data.Remark = form.Remark
	}
	if form.Status != "" {
		data.Status, _ = strconv.Atoi(form.Status)
	}
	err := service.Scene.UpdateById(form.Id, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}

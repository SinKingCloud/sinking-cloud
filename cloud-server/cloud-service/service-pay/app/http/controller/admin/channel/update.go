package channel

import (
	"service-pay/app/repository/channel"
	"service-pay/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Update 更新数据
func (Controller) Update(c *context.Context) {
	var form struct {
		Id     int64  `json:"id" default:"" validate:"required,min=1" label:"记录ID"`
		Name   string `json:"name" default:"" validate:"omitempty" label:"名称"`
		Config string `json:"config" default:"" validate:"omitempty" label:"配置"`
		Status string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &channel.UpdateChannel{}
	if form.Name != "" {
		data.Name = form.Name
	}
	if form.Config != "" {
		data.Config = form.Config
	}
	if form.Status != "" {
		data.Status, _ = strconv.Atoi(form.Status)
	}
	err := service.Channel.UpdateById(form.Id, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}

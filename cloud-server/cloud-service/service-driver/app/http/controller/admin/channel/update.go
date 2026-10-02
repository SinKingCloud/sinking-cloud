package channel

import (
	"fmt"
	"service-driver/app/repository/channel"
	"service-driver/app/service"
	"strconv"

	commonService "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
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
		data.Name = &form.Name
	}
	if form.Config != "" {
		data.Config = &form.Config
	}
	if form.Status != "" {
		status, err := strconv.Atoi(form.Status)
		if err != nil {
			c.Error("状态参数错误")
			return
		}
		data.Status = &status
	}
	err := service.Channel.UpdateById(form.Id, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = commonService.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventEdit, "修改存储通道", fmt.Sprintf("修改存储通道[%d]", form.Id))
	c.Success("更新数据成功")
}

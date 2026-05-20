package scene

import (
	"service-driver/app/model"
	"service-driver/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Create 创建数据
func (Controller) Create(c *context.Context) {
	var form struct {
		Name      string `json:"name" default:"" validate:"required" label:"名称"`
		MaxSize   int64  `json:"max_size" default:"0" validate:"numeric" label:"最大限制(KB)"`
		AllowExt  string `json:"allow_ext" default:"*" validate:"omitempty" label:"允许类型"`
		ChannelId int64  `json:"channel_id" default:"" validate:"numeric" label:"通道ID"`
		Login     int    `json:"login" default:"0" validate:"numeric" label:"登录限制"`
		Status    int    `json:"status" default:"0" validate:"numeric" label:"状态"`
		Remark    string `json:"remark" default:"" validate:"omitempty" label:"备注"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Scene.Create(&model.Scene{
		Name:      form.Name,
		MaxSize:   form.MaxSize,
		AllowExt:  form.AllowExt,
		Login:     form.Login,
		ChannelId: form.ChannelId,
		Remark:    form.Remark,
		Status:    form.Status,
	})
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("创建数据成功")
}

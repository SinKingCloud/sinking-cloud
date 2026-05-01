package task

import (
	"service-cron/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Run 运行任务
func (Controller) Run(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Task.Run(form.Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("运行任务成功")
}

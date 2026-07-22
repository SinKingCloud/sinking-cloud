package admin

import (
	"service-cron/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Log 任务日志
func (Controller) Log(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time", "id,create_time")
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"任务ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data, err := service.Log.SelectByTaskId(form.Id, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", data)
}

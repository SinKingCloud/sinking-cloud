package task

import (
	"service-cron/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Log 任务日志
func (Controller) Log(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,create_time")
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"任务ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data, total, err := service.Log.SelectByTaskId(form.Id, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
}

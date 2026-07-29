package mind_share

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// List 导图分享列表
func (Controller) List(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time,update_time")
	var form struct {
		DataId int64 `json:"data_id" validate:"required,min=1" label:"导图ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data, err := service.Mind.SelectSharesByUserId(c.GetUserInfo().Id, form.DataId, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", data)
}

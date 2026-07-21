package mind_version

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Restore 恢复导图版本
func (Controller) Restore(c *context.Context) {
	var form struct {
		DataId int64 `json:"data_id" validate:"required,min=1" label:"导图ID"`
		Id     int64 `json:"id" validate:"required,min=1" label:"版本ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data, err := service.Mind.RestoreVersion(c.GetUserInfo().Id, form.DataId, form.Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("恢复版本成功", data)
}

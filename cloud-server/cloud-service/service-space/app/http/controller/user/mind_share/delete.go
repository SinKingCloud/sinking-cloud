package mind_share

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Delete 删除导图分享
func (Controller) Delete(c *context.Context) {
	var form struct {
		DataId int64 `json:"data_id" validate:"required,min=1" label:"导图ID"`
		Id     int64 `json:"id" validate:"required,min=1" label:"分享ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	if err := service.Mind.DeleteShare(c.GetUserInfo().Id, form.DataId, form.Id); err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("删除数据成功")
}

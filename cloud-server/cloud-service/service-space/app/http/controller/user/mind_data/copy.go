package mind_data

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Copy 复制导图数据
func (Controller) Copy(c *context.Context) {
	var form struct {
		Id int64 `json:"id" validate:"required,min=1" label:"导图ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data, err := service.Mind.CopyData(c.GetUserInfo().Id, form.Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("复制导图成功", data)
}

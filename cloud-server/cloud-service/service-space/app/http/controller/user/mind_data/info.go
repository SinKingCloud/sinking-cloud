package mind_data

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Info 导图数据详情
func (Controller) Info(c *context.Context) {
	var form struct {
		Id int64 `json:"id" validate:"required,min=1" label:"导图ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Mind.GetDataInfoByUserId(c.GetUserInfo().Id, form.Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", info)
}

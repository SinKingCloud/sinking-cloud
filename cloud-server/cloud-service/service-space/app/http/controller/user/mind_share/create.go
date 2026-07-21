package mind_share

import (
	"service-space/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Create 创建导图分享
func (Controller) Create(c *context.Context) {
	var form struct {
		DataId     int64  `json:"data_id" validate:"required,min=1" label:"导图ID"`
		Status     int    `json:"status" default:"0" validate:"numeric" label:"分享状态"`
		ExpireTime string `json:"expire_time" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"到期时间"`
		Remark     string `json:"remark" default:"" validate:"omitempty,max=60" label:"备注"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data, err := service.Mind.CreateShare(c.GetUserInfo().Id, form.DataId, form.Status, form.ExpireTime, form.Remark)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("创建数据成功", data)
}

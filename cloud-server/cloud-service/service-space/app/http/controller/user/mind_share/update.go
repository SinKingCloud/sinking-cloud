package mind_share

import (
	mindShareRepository "service-space/app/repository/mind_share"
	"service-space/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Update 更新导图分享
func (Controller) Update(c *context.Context) {
	var form struct {
		DataId int64   `json:"data_id" validate:"required,min=1" label:"导图ID"`
		Id     int64   `json:"id" validate:"required,min=1" label:"分享ID"`
		Status string  `json:"status" default:"" validate:"omitempty,numeric" label:"分享状态"`
		Remark *string `json:"remark" validate:"omitempty,max=60" label:"备注"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &mindShareRepository.UpdateMindShare{}
	if form.Status != "" {
		status, err := strconv.Atoi(form.Status)
		if err != nil {
			c.Error("分享状态不合法")
			return
		}
		data.Status = status
	}
	if form.Remark != nil {
		data.Remark = *form.Remark
	}
	if err := service.Mind.UpdateShare(c.GetUserInfo().Id, form.DataId, form.Id, data); err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}

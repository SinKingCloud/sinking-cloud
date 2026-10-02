package file

import (
	"fmt"
	"service-driver/app/repository/file"
	"service-driver/app/service"
	"strconv"

	commonService "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Update 更新数据
func (Controller) Update(c *context.Context) {
	var form struct {
		Ids    []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
		Status string  `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &file.UpdateFile{}
	if form.Status != "" {
		status, err := strconv.Atoi(form.Status)
		if err != nil {
			c.Error("状态参数错误")
			return
		}
		data.Status = &status
	}
	err := service.File.UpdateByIds(form.Ids, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = commonService.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventEdit, "修改文件状态", fmt.Sprintf("修改文件%v的状态", form.Ids))
	c.Success("更新数据成功")
}

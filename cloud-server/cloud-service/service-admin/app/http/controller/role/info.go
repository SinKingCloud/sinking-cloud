package role

import (
	"service-admin/app/service"
	"strconv"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Info 角色信息
func (Controller) Info(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}

	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data, err := service.Rbac.GetRoleInfoById(form.Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = service2.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventShow, "查看角色详情", "查看角色["+strconv.FormatInt(form.Id, 10)+"]详情")
	c.SuccessWithData("获取数据成功", data)
}

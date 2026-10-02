package scene

import (
	"fmt"
	"service-driver/app/repository/scene"
	"service-driver/app/service"
	"strconv"

	commonService "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Update 更新数据
func (Controller) Update(c *context.Context) {
	var form struct {
		Id       int64  `json:"id" default:"" validate:"required,min=1" label:"记录ID"`
		Name     string `json:"name" default:"" validate:"omitempty" label:"名称"`
		MaxSize  string `json:"max_size" default:"" validate:"omitempty,numeric" label:"最大限制(KB)"`
		AllowExt string `json:"allow_ext" default:"" validate:"omitempty" label:"允许类型"`
		Remark   string `json:"remark" default:"" validate:"omitempty" label:"备注"`
		Login    string `json:"login" default:"" validate:"omitempty,numeric" label:"登录限制"`
		Status   string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &scene.UpdateScene{}
	if form.Name != "" {
		data.Name = &form.Name
	}
	if form.MaxSize != "" {
		maxSize, err := strconv.ParseInt(form.MaxSize, 10, 64)
		if err != nil {
			c.Error("最大限制(KB)参数错误")
			return
		}
		data.MaxSize = &maxSize
	}
	if form.AllowExt != "" {
		data.AllowExt = &form.AllowExt
	}
	if form.Remark != "" {
		data.Remark = &form.Remark
	}
	if form.Login != "" {
		login, err := strconv.Atoi(form.Login)
		if err != nil {
			c.Error("登录限制参数错误")
			return
		}
		data.Login = &login
	}
	if form.Status != "" {
		status, err := strconv.Atoi(form.Status)
		if err != nil {
			c.Error("状态参数错误")
			return
		}
		data.Status = &status
	}
	err := service.Scene.UpdateById(form.Id, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = commonService.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventEdit, "修改上传场景", fmt.Sprintf("修改上传场景[%d]", form.Id))
	c.Success("更新数据成功")
}

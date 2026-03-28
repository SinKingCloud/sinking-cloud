package controller

import (
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Enum 枚举数据
func Enum(c *context.Context) {
	var form struct {
		Name string `json:"name" default:"" validate:"required" label:"枚举名称"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	parts := strings.SplitN(form.Name, ".", 2)
	if len(parts) != 2 {
		c.Error("枚举数据不存在")
		return
	}
	data := service.Enum.Data(strings.ToLower(parts[0]), parts[1])
	if data == nil {
		c.Error("枚举类型不存在")
	} else {
		c.SuccessWithData("获取数据成功", data)
	}
}

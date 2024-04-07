package page

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"strings"
)

// Enum 获取枚举数据
func (ControllerPage) Enum(c *sinking_server.Context) {
	type Form struct {
		Keys string `form:"keys" json:"keys" default:"" validate:"required" label:"标识"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	configs := service.GetSystemConfigServiceIns().GetGroupByCache(config.EnumGroup, false)
	arr := strings.Split(form.Keys, "|")
	data := make(map[string]map[string]string)
	for _, v := range arr {
		k := config.EnumGroup + "." + v
		if configs[k] != "" {
			var temp map[string]string
			e := json.Unmarshal([]byte(configs[k]), &temp)
			if e == nil && temp != nil {
				data[v] = temp
			}
		}
	}
	c.SuccessWithData("获取枚举数据成功", data)
}

package config

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Aliyun 阿里云配置
func (s Controller) Aliyun(c *context.Context) {
	s.config(c, constant.AliGroup, []string{
		constant.AliKey,
		constant.AliSecret,
		constant.AliSmsSign,
		constant.AliSmsTemplate,
		constant.AliSmsVar,
		constant.AliMailAccount,
		constant.AliFaceScene,
	})
}

package config

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Web 网站配置
func (s Controller) Web(c *context.Context) {
	s.config(c, constant.WebGroup, []string{
		constant.WebUrl,
		constant.WebApi,
		constant.WebName,
		constant.WebTitle,
		constant.WebKeyWords,
		constant.WebDescribe,
		constant.WebContact,
		constant.WebLogo,
	})
}

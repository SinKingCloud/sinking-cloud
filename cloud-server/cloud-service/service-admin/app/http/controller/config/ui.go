package config

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Ui 界面配置
func (s Controller) Ui(c *context.Context) {
	s.config(c, constant.UiGroup, []string{
		constant.UiLayout,
		constant.UiWaterMark,
		constant.UiTheme,
		constant.UiCompact,
		constant.UiColor,
		constant.UiRadius,
	})
}

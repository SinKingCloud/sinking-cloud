package config

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Ui 界面配置
func (s Controller) Ui(c *context.Context) {
	s.config(c, constant.UiGroup, "界面配置", []string{
		constant.UiUserLayout,
		constant.UiUserWaterMark,
		constant.UiUserTheme,
		constant.UiUserCompact,
		constant.UiUserColor,
		constant.UiUserRadius,
		constant.UiAdminLayout,
		constant.UiAdminWaterMark,
		constant.UiAdminTheme,
		constant.UiAdminCompact,
		constant.UiAdminColor,
		constant.UiAdminRadius,
	})
}

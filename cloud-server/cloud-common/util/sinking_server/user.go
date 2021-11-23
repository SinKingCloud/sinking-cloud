package sinking_server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
)

// GetWebUserInfo 获取登录网站用户信息
func (c *Context) GetWebUserInfo() *model.WebUser {
	user, exists := c.Get(util.WebUserInfo)
	if exists {
		return user.(*model.WebUser)
	}
	return nil
}

// GetSystemUserInfo 获取登录系统用户信息
func (c *Context) GetSystemUserInfo() *model.SystemUser {
	user, exists := c.Get(util.SystemUserInfo)
	if exists {
		return user.(*model.SystemUser)
	}
	return nil
}

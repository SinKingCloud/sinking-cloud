package sinking_server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
)

// SetWebUserInfo 设置登录网站用户信息
func (c *Context) SetWebUserInfo(user interface{}) {
	c.Set(util.WebUserInfo, user)
}

// GetWebUserInfo 获取登录网站用户信息
func (c *Context) GetWebUserInfo() *model.WebUser {
	user, exists := c.Get(util.WebUserInfo)
	if exists {
		return user.(*model.WebUser)
	}
	return nil
}

// SetSystemUserInfo 设置登录系统用户信息
func (c *Context) SetSystemUserInfo(user interface{}) {
	c.Set(util.SystemUserInfo, user)
}

// GetSystemUserInfo 获取登录系统用户信息
func (c *Context) GetSystemUserInfo() *model.SystemUser {
	user, exists := c.Get(util.SystemUserInfo)
	if exists {
		return user.(*model.SystemUser)
	}
	return nil
}

package sinking_server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
)

// SetWebUserInfo 设置登录网站用户信息
func (c *Context) SetWebUserInfo(user interface{}) {
	c.Set(util.WebUserInfo, user)
}

// GetWebUserInfo 获取登录网站用户信息
func (c *Context) GetWebUserInfo() *jwt.WebUser {
	user, exists := c.Get(util.WebUserInfo)
	if exists {
		return user.(*jwt.WebUser)
	}
	return nil
}

// SetSystemUserInfo 设置登录系统用户信息
func (c *Context) SetSystemUserInfo(user interface{}) {
	c.Set(util.SystemUserInfo, user)
}

// GetSystemUserInfo 获取登录系统用户信息
func (c *Context) GetSystemUserInfo() *jwt.SystemUser {
	user, exists := c.Get(util.SystemUserInfo)
	if exists {
		return user.(*jwt.SystemUser)
	}
	return nil
}

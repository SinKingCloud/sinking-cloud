package server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
)

// SetUserInfo 设置登录网站用户信息
func (c *Context) SetUserInfo(user *jwt.User) {
	c.Set(constant.UserInfo, user)
}

// GetUserInfo 获取登录网站用户信息
func (c *Context) GetUserInfo() *jwt.User {
	value, exists := c.Get(constant.UserInfo)
	if exists {
		return value.(*jwt.User)
	}
	return nil
}

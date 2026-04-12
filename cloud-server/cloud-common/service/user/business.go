package user

import (
	"errors"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// GetSystemUserInfo 获取系统用户信息
func (*Service) GetSystemUserInfo(userId int64, refresh bool) (user *SystemUser, err error) {
	err = rpc.InvokeFunc(constant.ServiceAdmin, "GetUserInfo", &GetSystemUserInfoParam{
		UserId:  userId,
		Refresh: refresh,
	}, &user)
	if err != nil || user == nil {
		return nil, err
	}
	return user, nil
}

// GetWebUserInfo 获取网站用户信息
func (*Service) GetWebUserInfo(userId int64, refresh bool) (user *WebUser, err error) {
	err = rpc.InvokeFunc(constant.ServiceUser, "GetUserInfo", &GetWebUserInfoParam{
		UserId:  userId,
		Refresh: refresh,
	}, &user)
	if err != nil || user == nil {
		return nil, err
	}
	return user, nil
}

// GetJwtUserInfo 获取jwt user信息
func (*Service) GetJwtUserInfo(c *context.Context) (*jwt.User, error) {
	token := c.Request.Header.Get(constant.JwtTokenName)
	if token == "" {
		return nil, errors.New("您还未登陆,请先登陆账户")
	}
	key := jwt.CheckToken(token)
	if key == nil || key.User == nil || key.User.Id <= 0 || key.User.Token == "" {
		return nil, errors.New("登陆超时,请重新登陆")
	}
	if key.User == nil || key.User.Id <= 0 || key.User.Token == "" {
		return nil, errors.New("登陆失效,请重新登陆")
	}
	return key.User, nil
}

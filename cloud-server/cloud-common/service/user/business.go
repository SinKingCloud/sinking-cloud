package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
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

package user

import (
	"errors"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
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

// CreateSystemUserLog 异步创建系统用户操作日志，成功表示消息已入队。
func (*Service) CreateSystemUserLog(userId int64, requestId string, requestIp string, types int, title string, content string) error {
	var result bool
	err := rpc.InvokeFunc(constant.ServiceAdmin, "CreateUserLog", &CreateSystemUserLogParam{
		Id:        str.GetSnowWorkIns().GetId(),
		UserId:    userId,
		RequestId: requestId,
		RequestIp: requestIp,
		Type:      types,
		Title:     title,
		Content:   content,
	}, &result)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("创建用户操作日志失败")
	}
	return nil
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

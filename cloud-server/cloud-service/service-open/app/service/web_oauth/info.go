package web_oauth

import (
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth/oauth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth/oauth_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth/oauth_user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_user"
	"strconv"
)

// UserInfo 个人信息
type UserInfo struct {
	Id        int          `json:"id"`
	Account   string       `json:"account"`
	RealAuth  int          `json:"real_auth"`
	NickName  string       `json:"nick_name"`
	Avatar    string       `json:"avatar"`
	Email     string       `json:"email"`
	Phone     string       `json:"phone"`
	LoginIp   string       `json:"login_ip"`
	LoginTime str.DateTime `json:"login_time"`
}

// GetUserInfo 获取用户信息
func GetUserInfo(oauthId int, ticket string) (*UserInfo, error) {
	//查询ticket
	info, err := oauth_log.GetIns().FindByTicket(ticket)
	if err != nil || info.OauthId != oauthId {
		return nil, errors.New("授权码不存在")
	}
	if info.Status != oauth_log.Success {
		return nil, errors.New("用户未授权或授权码已失效")
	}
	//查询授权用户
	oauthUser, err := oauth_user.GetIns().FindByOauthIdAndUserId(oauthId, info.UserId)
	if err != nil || oauthUser == nil {
		return nil, errors.New("获取用户信息失败")
	}
	if oauthUser.Status != oauth_user.Normal {
		return nil, errors.New("该用户已被禁止登陆")
	}
	//查询系统
	user := web_user.GetIns().GetCurrentUserInfo(oauthUser.UserId, false)
	if user == nil {
		return nil, errors.New("获取用户信息失败")
	}
	if user.Status == web_user.Error {
		return nil, errors.New("该用户已被封禁")
	}
	//更改ticket状态并构建用户信息
	_ = oauth_log.GetIns().UpdateById(info.Id, map[string]interface{}{
		"status": oauth_log.Expire,
	})
	data := &UserInfo{
		Id:        user.Id,
		Account:   user.Account,
		RealAuth:  user.RealAuth,
		NickName:  user.NickName,
		Avatar:    user.Avatar,
		Email:     user.Email,
		Phone:     str.NewStringTool(strconv.FormatInt(user.Phone, 10)).Show(3, 4, "*"),
		LoginIp:   user.LoginIp,
		LoginTime: oauthUser.LoginTime,
	}
	return data, nil
}

// GetOauthInfo 获取接入网站信息
func GetOauthInfo(oauthId int, userId int) (*model.AppOauth, error) {
	oauthInfo, err := oauth.GetIns().FindByIdWithCache(oauthId, false)
	if err != nil || oauthInfo == nil || oauthInfo.UserId != userId {
		return nil, errors.New("该网站不存在")
	}
	if oauthInfo.Status == oauth.Error || oauthInfo.Status == oauth.Stop {
		return nil, errors.New("该网站已被停止接入")
	}
	return oauthInfo, nil
}

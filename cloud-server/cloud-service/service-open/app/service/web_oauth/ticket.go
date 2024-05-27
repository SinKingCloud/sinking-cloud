package web_oauth

import (
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	oauth2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth/oauth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth/oauth_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth/oauth_user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_user"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// AuthTicket 授权ticket
func (*Service) AuthTicket(userId int, ticket string, redirectUrl string) error {
	//防止并发锁
	key := constant.LockOauthAuthTicket + strconv.Itoa(userId)
	if !util.Cache.Lock(key, constant.LockTimeOauthAuthTicket) {
		return errors.New("系统繁忙")
	}
	defer util.Cache.UnLock(key) //解锁
	//查询用户
	user := web_user.GetIns().GetCurrentUserInfo(userId, false)
	if user == nil || user.Status == web_user.Error {
		return errors.New("获取用户信息失败")
	}
	//查询ticket
	info, err := oauth_log.GetIns().FindByTicket(ticket)
	if err != nil || info.Status != oauth_log.Normal {
		return errors.New("授权码已失效,请返回来源地重新发起请求")
	}
	//查询网站
	oauth, err := oauth2.GetIns().FindByIdWithCache(info.OauthId, false)
	if err != nil || oauth.Status != oauth2.Normal {
		return errors.New("此网站已暂停登陆注册功能")
	}
	r, rErr := url.Parse(redirectUrl)
	if rErr != nil || !strings.Contains(oauth.Domain, r.Hostname()) {
		return errors.New("来源地验证失败,请联系管理员")
	}
	//查询是否已授权此网站
	oauthUser, err := oauth_user.GetIns().FindByOauthIdAndUserId(info.OauthId, user.Id)
	if err != nil || oauthUser == nil {
		//写入用户信息
		e := oauth_user.GetIns().Create(&model.AppOauthUser{
			OauthId: info.OauthId,
			UserId:  user.Id,
			Status:  oauth_user.Normal,
		})
		if e != nil {
			return errors.New("创建授权用户失败")
		}
		oauthUser, err = oauth_user.GetIns().FindByOauthIdAndUserId(info.OauthId, user.Id)
	}
	//更新ticket状态
	err = oauth_log.GetIns().UpdateById(info.Id, map[string]interface{}{
		"user_id": user.Id,
		"status":  oauth_log.Success,
	})
	if err != nil {
		return errors.New("授权登陆失败")
	}
	//更新登陆时间
	_ = oauth_user.GetIns().UpdateById(oauthUser.Id, map[string]interface{}{
		"login_time": str.DateTime(time.Now()),
	})
	return nil
}

// GenTicket 生成ticket
func (*Service) GenTicket(oauthId int) (string, error) {
	ticket := str.NewStringTool(str.GetSnowWorkIns().GetUuid()).Md5()
	err := oauth_log.GetIns().Create(&model.AppOauthLog{
		OauthId: oauthId,
		Ticket:  ticket,
		UserId:  0,
		Status:  oauth_log.Normal,
	})
	return ticket, err
}

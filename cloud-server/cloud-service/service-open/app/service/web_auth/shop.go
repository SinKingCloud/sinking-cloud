package web_auth

import (
	"errors"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/auth_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/user_auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_user"
	"strconv"
)

type shop struct {
}

var Shop *shop

// Buy 购买授权
func (*shop) Buy(userId int, authId int, domain string) error {
	key := constant.LockAuthBuy + strconv.Itoa(userId) + "_" + strconv.Itoa(authId) + "_" + domain
	//上锁，避免同时下单太多
	if !util.Cache.Lock(key, constant.LockTimeAuthBuy) {
		return errors.New("系统繁忙,请稍后下单")
	}
	defer util.Cache.UnLock(key) //释放并发锁
	user := web_user.GetIns().GetCurrentUserInfo(userId, true)
	if user == nil {
		return errors.New("获取用户信息失败")
	}
	authApp, e := auth_app.GetIns().FindByIdWithCache(authId, false)
	if e != nil || authApp == nil {
		return errors.New("未查询到应用信息")
	}
	if authApp.Status != auth_app.Normal {
		return errors.New("该应用已下架")
	}
	//判断是否已存在此授权
	userAuth, e := user_auth.GetIns().FindByAuthIdDomain(authApp.Id, domain)
	if e != nil || userAuth != nil {
		return errors.New("系统已存在此授权,无需重复授权")
	}
	//判断余额
	if user.Money < authApp.Price {
		return errors.New("账户余额不足,请充值")
	}
	//扣款
	if authApp.Price > 0 && !rpc.WebUser.Money(user.Id, 0, fmt.Sprintf("%.2f", authApp.Price)) {
		return errors.New("扣款失败，请联系管理员")
	}
	//写入授权
	err := user_auth.GetIns().Create(&model.UserAuth{
		UserId: user.Id,
		AuthId: user.Id,
		Domain: domain,
		Code:   str.GetSnowWorkIns().GetUuid(),
		Status: user_auth.Normal,
	})
	//购买失败退会余额
	if err != nil {
		rpc.WebUser.Money(userId, 1, fmt.Sprintf("%.2f", authApp.Price))
		return errors.New("购买失败,请联系管理员")
	}
	rpc.WebUserPayLog.Create(user.Id, 1, 1, authApp.Price, "在线消费", "购买"+authApp.Name+",授权域名:"+domain)
	return nil
}

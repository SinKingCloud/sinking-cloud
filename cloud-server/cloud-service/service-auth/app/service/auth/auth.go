package auth

import (
	"encoding/json"
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/constant"
	"strconv"
	"time"
)

// GetWebUserInfoByPhone 根据短信验证码获取网站用户信息
func GetWebUserInfoByPhone(phone string, code string) (u *model.WebUser, e error) {
	if !SmsCheck(phone, code) {
		e = errors.New("短信验证码验证失败")
		return
	}
	e = util.Database.Db.Where("`phone` = ?", phone).First(&u).Error
	if e != nil || u == nil {
		//防止并发创建用户
		key := constant.LockWebUserReg + phone
		if !util.Cache.Lock(key, constant.LockTimeWebUserReg) {
			e = errors.New("系统繁忙,请稍后重试")
		} else {
			defer util.Cache.UnLock(key) //解锁
			//创建用户
			tempPhone, _ := strconv.ParseInt(phone, 10, 64)
			e = util.Database.Db.Create(&model.WebUser{
				NickName: str.NewNickTool().GetRand(),
				Phone:    tempPhone,
				Money:    0,
				Status:   0,
			}).Error
			if e != nil {
				e = errors.New("注册账号失败")
				return
			}
			util.Database.Db.Where("`phone` = ?", phone).First(&u)
		}
	}
	if u.Status != 0 {
		e = errors.New("您的账户已被禁止登陆")
		return
	}
	return
}

// GetWebUserInfoByPwd 根据用户名密码获取网站用户信息
func GetWebUserInfoByPwd(account string, pwd string) (u *model.WebUser, e error) {
	e = util.Database.Db.Where("`account` = ?", account).First(&u).Error
	if e != nil || u == nil {
		e = util.Database.Db.Where("`phone` = ?", account).First(&u).Error
		if e != nil || u == nil {
			e = errors.New("用户不存在")
			return
		}
	}
	if u.Status != 0 {
		e = errors.New("您的账户已被禁止登陆")
		return
	}
	if !str.NewStringTool(u.Password).CheckPassword(pwd) {
		e = errors.New("用户名或密码错误")
		return
	}
	return
}

// GetSystemUserInfoByPwd 根据用户名密码获取系统用户信息
func GetSystemUserInfoByPwd(account string, pwd string) (u *model.SystemUser, e error) {
	e = util.Database.Db.Where("`account`= ?", account).First(&u).Error
	if e != nil || u == nil {
		e = errors.New("用户不存在")
		return
	}
	if u.Status != 0 {
		e = errors.New("您的账户已被禁止登陆")
		return
	}
	if !str.NewStringTool(u.Password).CheckPassword(pwd) {
		e = errors.New("用户名或密码错误")
		return
	}
	return
}

// OutLogin 注销登录
func OutLogin(types string, u *model.SystemUser, u2 *model.WebUser) {
	fun := func(types string, newToken string, oldToken string) string {
		temp := make(map[string]string)
		_ = json.Unmarshal([]byte(oldToken), &temp)
		temp[types] = newToken
		//生成新的json
		jsonStr, err := json.Marshal(temp)
		if err != nil {
			return ""
		}
		return string(jsonStr)
	}
	if u != nil {
		tempToken := fun(types, "", u.LoginToken)
		if tempToken != "" {
			tempU := model.SystemUser{Id: u.Id, LoginToken: tempToken}
			_ = util.Database.Db.Updates(&tempU).Error
		}
	}
	if u2 != nil {
		tempToken := fun(types, "", u2.LoginToken)
		if tempToken != "" {
			tempU := model.WebUser{Id: u2.Id, LoginToken: tempToken}
			_ = util.Database.Db.Updates(&tempU).Error
		}
	}
}

// GetJwtToken 生成jwtToken
func GetJwtToken(types string, u *model.SystemUser, u2 *model.WebUser, ip string) (s string, e error) {
	//生成token
	token := str.NewStringTool(strconv.FormatInt(time.Now().UnixMilli(), 10)).Md5()
	loginTime := str.DateTime(time.Now())
	jwtLoginToken, _ := json.Marshal(map[string]string{types: token})
	fun := func(types string, newToken string, oldToken string) string {
		temp := make(map[string]string)
		_ = json.Unmarshal([]byte(oldToken), &temp)
		temp[types] = newToken
		//生成新的json
		jsonStr, err := json.Marshal(temp)
		if err != nil {
			return ""
		}
		return string(jsonStr)
	}
	if u != nil {
		tempToken := fun(types, token, u.LoginToken)
		if tempToken != "" {
			tempU := model.SystemUser{Id: u.Id, LoginToken: tempToken, LoginTime: loginTime, LoginIp: ip}
			e = util.Database.Db.Updates(&tempU).Error
		} else {
			e = errors.New("生成token失败")
			return
		}
		u.LoginToken = string(jwtLoginToken)
		u.LoginTime = loginTime
		service.GetSystemUserServiceIns().ClearCache(u.Id)
	}
	if u2 != nil {
		tempToken := fun(types, token, u2.LoginToken)
		if tempToken != "" {
			tempU := model.WebUser{Id: u2.Id, LoginToken: tempToken, LoginTime: loginTime, LoginIp: ip}
			e = util.Database.Db.Updates(&tempU).Error
		} else {
			e = errors.New("生成token失败")
			return
		}
		u2.LoginToken = string(jwtLoginToken)
		u2.LoginTime = loginTime
		service.GetWebUserServiceIns().ClearCache(u2.Id)
	}
	var webUser *jwt.WebUser
	var sysUser *jwt.SystemUser
	if u2 != nil {
		webUser = &jwt.WebUser{
			Id:         u2.Id,
			LoginToken: u2.LoginToken,
			LoginIp:    u2.LoginIp,
			LoginTime:  loginTime,
			Status:     u2.Status,
			RealAuth:   u2.RealAuth,
		}
	}
	if u != nil {
		sysUser = &jwt.SystemUser{
			Id:         u.Id,
			LoginToken: u.LoginToken,
			LoginIp:    u.LoginIp,
			LoginTime:  loginTime,
			Status:     u.Status,
		}
	}
	s = jwt.GetToken(webUser, sysUser)
	return
}

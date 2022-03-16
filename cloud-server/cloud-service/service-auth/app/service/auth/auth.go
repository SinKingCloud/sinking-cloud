package auth

import (
	"encoding/json"
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"time"
)

// GetSystemUserInfoByPwd 根据用户名密码获取系统用户信息
func GetSystemUserInfoByPwd(account string, pwd string) (u *model.SystemUser, e error) {
	e = util.Database.Db.Where("account=?", account).First(&u).Error
	if e != nil {
		e = errors.New("系统执行错误")
		return
	}
	if u == nil {
		e = errors.New("该用户不存在")
		return
	}
	if u.Status != 0 {
		e = errors.New("该用户已被禁止登陆")
		return
	}
	if !str.NetStringTool(u.Password).CheckPassword(pwd) {
		e = errors.New("用户名或密码错误")
		return
	}
	return
}

// GetJwtToken 生成jwtToken
func GetJwtToken(types string, u *model.SystemUser, u2 *model.WebUser, ip string) (s string, e error) {
	//生成token
	token := str.NetStringTool(time.Now().String()).Md5()
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
	s = jwt.GetToken(u2, u)
	return
}

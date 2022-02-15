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
	util.Database.Db.Where("account=?", account).First(&u)
	if u == nil {
		e = errors.New("该用户不存在")
		return
	}
	if !str.NetStringTool(u.Password).CheckPassword(pwd) {
		e = errors.New("用户名或密码错误")
		return
	}
	return
}

// GetJwtToken 生成jwtToken
func GetJwtToken(types string, u *model.SystemUser, u2 *model.WebUser) (s string, e error) {
	//生成token
	token := str.NetStringTool(time.Now().String()).Md5()
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
			tempU := model.SystemUser{Id: u.Id, LoginToken: tempToken, LoginTime: str.DateTime(time.Now())}
			util.Database.Db.Updates(&tempU)
		} else {
			e = errors.New("生成token失败")
			return
		}
		u.LoginToken = string(jwtLoginToken)
		service.GetSystemUserServiceIns().ClearCache(u.Id)
	}
	if u2 != nil {
		tempToken := fun(types, token, u2.LoginToken)
		if tempToken != "" {
			tempU := model.WebUser{Id: u2.Id, LoginToken: tempToken, LoginTime: str.DateTime(time.Now())}
			util.Database.Db.Updates(&tempU)
		} else {
			e = errors.New("生成token失败")
			return
		}
		u2.LoginToken = string(jwtLoginToken)
		service.GetWebUserServiceIns().ClearCache(u2.Id)
	}
	s = jwt.GetToken(u2, u)
	return
}

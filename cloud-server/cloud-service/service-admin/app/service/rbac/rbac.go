package rbac

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/constant"
	access2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/access"
	role2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/role"
	"strconv"
	"strings"
)

func CheckAccessByCache(roleId int, url string, method string) bool {
	key := constant.CacheNameWithRbac + strconv.Itoa(roleId) + url + ":" + method
	res, exists := util.Cache.Get(key)
	if exists {
		return res == "1"
	} else {
		//写入缓存
		if CheckAccess(roleId, url, method) {
			res = strconv.Itoa(1)
		} else {
			res = strconv.Itoa(0)
		}
		util.Cache.SetWithExpire(key, res, constant.CacheTimeWithRbac)
	}
	return res == "1"
}

// CheckAccess 判断是否含有权限
func CheckAccess(roleId int, url string, method string) bool {
	role := role2.GetRoleServiceIns().FindByIdWithCache(roleId, false)
	if role == nil {
		return true
	}
	//获取权限信息
	access := access2.GetRoleServiceIns().FindByUriMethodWithCache(url, strings.ToUpper(method), false)
	if access == nil || access.Id == 0 {
		return true
	}
	//判断是否含有权限
	var roleAccessIds []int
	err := json.Unmarshal([]byte(role.Accesses), &roleAccessIds)
	if err != nil {
		return true
	}
	for _, v := range roleAccessIds {
		if v == access.Id {
			return true
		}
	}
	return false
}

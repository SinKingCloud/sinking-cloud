package rbac

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/constant"
	access2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/access"
	role2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/role"
	"strconv"
	"strings"
)

// service 单例对象
type service struct {
}

// obj 单例对象
var obj *service

// GetIns 获取单例
func GetIns() *service {
	if obj == nil {
		obj = &service{}
	}
	return obj
}

func (s *service) CheckAccessByCache(roleId int, url string, method string) bool {
	key := constant.CacheNameWithRbac + strconv.Itoa(roleId) + url + ":" + method
	res, exists := util.Cache.Get(key)
	if exists {
		return res == "1"
	} else {
		//写入缓存
		r, _ := s.CheckAccess(roleId, url, method)
		if r {
			res = strconv.Itoa(1)
		} else {
			res = strconv.Itoa(0)
		}
		util.Cache.SetWithExpire(key, res, constant.CacheTimeWithRbac)
	}
	return res == "1"
}

// CheckAccess 判断是否含有权限
func (s *service) CheckAccess(roleId int, url string, method string) (bool, *model.SystemAccess) {
	role, err := role2.GetIns().FindByIdWithCache(roleId, false)
	if err != nil || role == nil {
		return true, nil
	}
	//获取权限信息
	access, err := access2.GetIns().FindByUriMethodWithCache(url, strings.ToUpper(method), false)
	if err != nil || access == nil || access.Id == 0 {
		return true, nil
	}
	if access.Status == 1 {
		return false, access
	}
	//判断是否含有权限
	var roleAccessIds []int
	err = json.Unmarshal([]byte(role.Accesses), &roleAccessIds)
	if err != nil {
		return true, access
	}
	for _, v := range roleAccessIds {
		if v == access.Id {
			return true, access
		}
	}
	return false, access
}

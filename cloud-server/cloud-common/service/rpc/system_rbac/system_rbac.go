package system_rbac

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/name"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"strconv"
	"sync"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj  *Service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

// Check 判断rbac
func (*Service) Check(userId int, roleId int, path string, method string) bool {
	key := CacheNameWithRbac + "Res_" + strconv.Itoa(roleId) + path + ":" + method
	res, exists := util.Cache.Get(key)
	if exists {
		return res == "1"
	} else {
		//写入缓存
		client, err := rpc.New(group.CloudService, name.ServiceAdmin).Call("/rpc/rbac/check", &rpc.Param{
			"user_id": userId,
			"role_id": roleId,
			"path":    path,
			"method":  method,
		}).Response()
		if err != nil {
			return false
		}
		if client.Data.(bool) {
			res = strconv.Itoa(1)
		} else {
			res = strconv.Itoa(0)
		}
		util.Cache.SetWithExpire(key, res, int64(CacheTimeWithRbac))
	}
	return res == "1"
}

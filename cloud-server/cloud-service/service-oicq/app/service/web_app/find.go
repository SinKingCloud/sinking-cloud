package web_app

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/user_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/constant"
	"github.com/pkg/errors"
	"strconv"
)

// FindById 通过ID查询
func (s *Service) FindById(id int) (data *user_app.UserApp, err error) {
	data = rpc.UserApp.FindById(id)
	if data == nil {
		err = errors.New("未查询到数据")
	}
	return
}

// FindByIdWithCache 通过ID查询
func (s *Service) FindByIdWithCache(id int, refresh bool) (data *user_app.UserApp, err error) {
	key := constant.CacheNameWithAppKey + strconv.Itoa(id)
	if refresh {
		data, err = s.FindById(id)
		if err != nil {
			return
		}
		str, err := json.Marshal(data)
		if err == nil {
			util.Cache.SetWithExpire(key, string(str), constant.CacheTimeWithAppKey)
		}
	} else {
		str, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(str), &data)
		} else {
			data, err = s.FindById(id)
			if err != nil {
				return
			}
			str, err := json.Marshal(data)
			if err == nil {
				util.Cache.SetWithExpire(key, string(str), constant.CacheTimeWithAppKey)
			}
		}
	}
	return
}

// Clear 清除缓存
func (s *Service) Clear(id int) bool {
	key := constant.CacheNameWithAppKey + strconv.Itoa(id)
	return util.Cache.Delete(key)
}

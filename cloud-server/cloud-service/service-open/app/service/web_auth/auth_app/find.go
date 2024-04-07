package auth_app

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"strconv"
)

// FindById 通过ID查询
func (s *Service) FindById(id int) (data *model.AuthApp, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&data).Error
	return
}

// FindByIdWithCache 通过ID查询
func (s *Service) FindByIdWithCache(id int, refresh bool) (data *model.AuthApp, err error) {
	key := constant.CacheNameWithAuth + strconv.Itoa(id)
	if refresh {
		data, err = s.FindById(id)
		if err != nil {
			return
		}
		str, err := json.Marshal(data)
		if err == nil {
			util.Cache.SetWithExpire(key, string(str), constant.CacheTimeWithAuth)
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
				util.Cache.SetWithExpire(key, string(str), constant.CacheTimeWithAuth)
			}
		}
	}
	return
}

// Clear 清除缓存
func (s *Service) Clear(id int) bool {
	key := constant.CacheNameWithAuth + strconv.Itoa(id)
	return util.Cache.Delete(key)
}

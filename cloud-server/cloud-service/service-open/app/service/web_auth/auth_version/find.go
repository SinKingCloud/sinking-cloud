package auth_version

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"strconv"
)

// FindById 通过ID查询
func (s *service) FindById(id int) (data *model.AuthVersion, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&data).Error
	return
}

// FindByIdWithCache 通过ID查询
func (s *service) FindByIdWithCache(id int, refresh bool) (data *model.AuthVersion, err error) {
	key := constant.CacheNameWithAuthVersion + strconv.Itoa(id)
	if refresh {
		data, err = s.FindById(id)
		if err != nil {
			return
		}
		str, err := json.Marshal(data)
		if err == nil {
			util.Cache.SetWithExpire(key, string(str), constant.CacheTimeWithAuthVersion)
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
				util.Cache.SetWithExpire(key, string(str), constant.CacheTimeWithAuthVersion)
			}
		}
	}
	return
}

// Clear 清除缓存
func (s *service) Clear(id int) bool {
	key := constant.CacheNameWithAuthVersion + strconv.Itoa(id)
	return util.Cache.Delete(key)
}

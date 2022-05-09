package web_app

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"strconv"
)

// FindByUserIdWithId 通过用户ID及ID查询
func (s *service) FindByUserIdWithId(userId int, id int) (data *model.UserApp, err error) {
	err = util.Database.Db.Where("`id` = ? and `user_id` = ?", id, userId).First(&data).Error
	return
}

// FindById 通过ID查询
func (s *service) FindById(id int) (data *model.UserApp, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&data).Error
	return
}

// FindByIdWithCache 通过ID查询
func (s *service) FindByIdWithCache(id int, refresh bool) (data *model.UserApp, err error) {
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
func (s *service) Clear(id int) bool {
	key := constant.CacheNameWithAppKey + strconv.Itoa(id)
	return util.Cache.Delete(key)
}

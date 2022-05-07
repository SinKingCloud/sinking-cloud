package server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
	"time"
)

// UpdateById 通过ID更新
func (s *service) UpdateById(id int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.AppUinServer{}).Where("`id` = ? ", id).Updates(data).Error
	if err == nil {
		s.Clear(id)
	}
	return
}

// UpdateByIds 通过ID列表更新
func (s *service) UpdateByIds(ids []int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.AppUinServer{}).Where("`id` in ? ", ids).Updates(data).Error
	if err == nil {
		go func(ids []int, s *service) {
			for _, v := range ids {
				s.Clear(v)
			}
		}(ids, s)
	}
	return
}

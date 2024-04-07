package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"time"
)

// UpdateByIds 通过ID列表更新
func (s *Service) UpdateByIds(ids []int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.WebUser{}).Where("`id` in ? ", ids).Updates(data).Error
	return
}

// UpdateById 通过ID更新
func (s *Service) UpdateById(id int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.WebUser{}).Where("`id` = ? ", id).Updates(data).Error
	if err == nil {
		service2.GetWebUserServiceIns().ClearCache(id)
	}
	return
}

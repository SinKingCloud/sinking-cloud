package system_access

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"time"
)

// UpdateById 通过ID更新
func (s *Service) UpdateById(id int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.SystemAccess{}).Where("`id` = ? ", id).Updates(data).Error
	return
}

// UpdateByIds 通过ID列表更新
func (s *Service) UpdateByIds(ids []int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.SystemAccess{}).Where("`id` in ? ", ids).Updates(data).Error
	return
}

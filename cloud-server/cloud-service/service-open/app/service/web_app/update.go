package web_app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"time"
)

// UpdateById 通过ID更新
func (s *service) UpdateById(id int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.UserApp{}).Where("`id` = ? ", id).Updates(data).Error
	if err == nil {
		s.Clear(id)
	}
	return
}

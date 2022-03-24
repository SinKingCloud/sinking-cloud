package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"time"
)

// UpdateById 通过ID更新
func (s *service) UpdateById(id int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.WebUserApp{}).Where("`id` = ? ", id).Updates(data).Error
	return
}

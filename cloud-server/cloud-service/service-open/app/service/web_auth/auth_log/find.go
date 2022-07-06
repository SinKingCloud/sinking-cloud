package auth_log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// FindById 通过ID查询
func (s *service) FindById(id int) (data *model.AuthLog, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&data).Error
	return
}

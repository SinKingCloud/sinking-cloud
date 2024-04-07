package level

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
)

// FindById 通过ID查询
func (s *Service) FindById(id int) (data *model.AppUinLevel, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&data).Error
	return
}

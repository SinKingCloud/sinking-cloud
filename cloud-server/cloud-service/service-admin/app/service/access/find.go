package access

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// FindById 通过ID查询
func (s *service) FindById(id int) (access *model.SystemAccess, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&access).Error
	return
}

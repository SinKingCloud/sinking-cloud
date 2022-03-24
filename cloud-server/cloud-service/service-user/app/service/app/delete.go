package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// DeleteById 通过ID删除
func (s *service) DeleteById(id int) (err error) {
	err = util.Database.Db.Where("`id` = ? ", id).Delete(&model.WebUserApp{}).Error
	return
}

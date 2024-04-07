package system_log_request

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// FindById 通过ID查询
func (s *Service) FindById(id string) (access *model.SystemLogRequest, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&access).Error
	return
}

// FindByRequestId 通过请求ID查询
func (s *Service) FindByRequestId(id string) (access *model.SystemLogRequest, err error) {
	err = util.Database.Db.Where("`request_id` = ? ", id).First(&access).Error
	return
}

package system_log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// DeleteByRequestTime 通过请求时间列表删除
func (s *Service) DeleteByRequestTime(createTimeStart string, createTimeEnd string) (err error) {
	err = util.Database.Db.Where("`create_time` between ? and ? ", createTimeStart, createTimeEnd).Delete(&model.SystemLog{}).Error
	return
}

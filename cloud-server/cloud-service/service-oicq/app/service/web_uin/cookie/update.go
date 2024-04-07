package cookie

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
	"time"
)

// UpdateByUinDomain 通过uin domain更新
func (s *Service) UpdateByUinDomain(uin string, domain string, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.AppUinCookie{}).Where("`uin` = ? AND `domain` = ?", uin, domain).Updates(data).Error
	if err == nil {
		s.ClearByUinDomain(uin, domain)
	}
	return
}

package proxy_area

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// Count 查询
func (s *service) Count(where map[string]interface{}) (total int64, err error) {
	query := util.Database.Db.Model(&model.AppUinProxyArea{})
	if where["province"] != nil {
		query.Where("province", where["province"])
	}
	if where["province_key"] != nil {
		query.Where("province_key", where["province_key"])
	}
	if where["city"] != nil {
		query.Where("city", where["city"])
	}
	if where["city_key"] != nil {
		query.Where("city_key", where["city_key"])
	}
	err = query.Count(&total).Error
	return
}

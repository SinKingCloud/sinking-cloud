package proxy_area

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// Find 查询
func (s *service) Find(where map[string]interface{}) (data *model.AppProxyArea, err error) {
	query := util.Database.Db.Model(&model.AppProxyArea{})
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
	err = query.First(&data).Error
	return
}

// FindById 通过ID查询
func (s *service) FindById(id int) (data *model.AppProxyArea, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&data).Error
	return
}

// FindByProvince 通过Province查询
func (s *service) FindByProvince(province string) (data *model.AppProxyArea, err error) {
	err = util.Database.Db.Where("`province` = ? ", province).First(&data).Error
	return
}

// FindByCity 通过City查询
func (s *service) FindByCity(city string) (data *model.AppProxyArea, err error) {
	err = util.Database.Db.Where("`city` = ? ", city).First(&data).Error
	return
}

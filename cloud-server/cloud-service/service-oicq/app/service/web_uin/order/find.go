package order

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
)

// FindById 通过ID查询
func (s *Service) FindById(id int) (data *model.AppUinOrder, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&data).Error
	return
}

// FindByTradeNo 通过trade_no查询
func (s *Service) FindByTradeNo(tradeNo string) (data *model.AppUinOrder, err error) {
	err = util.Database.Db.Where("`trade_no` = ? ", tradeNo).First(&data).Error
	return
}

// FindByUserIdTradeNo 通过user_id trade_no查询
func (s *Service) FindByUserIdTradeNo(userId int, tradeNo string) (data *model.AppUinOrder, err error) {
	err = util.Database.Db.Where("`user_id` = ? AND `trade_no` = ? ", userId, tradeNo).First(&data).Error
	return
}

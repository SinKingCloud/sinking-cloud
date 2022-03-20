package order

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// FindByTradeNo 通过订单号查询
func (s *service) FindByTradeNo(tradeNo string) (order *model.WebUserOrder, err error) {
	err = util.Database.Db.Where("`trade_no` = ? ", tradeNo).First(&order).Error
	return
}

package balance

import (
	"errors"
	"service-pay/app/enum/balance_log_type"
	"service-pay/app/repository/balance_log"
	"strconv"
)

// SelectLogByUserId 查询用户余额明细
func (s *service) SelectLogByUserId(userId int64, where *balance_log.SelectBalanceLog, orderByField string, orderByType string, page int, pageSize int) (list []*Log, total int64, err error) {
	if userId <= 0 {
		return nil, 0, errors.New("用户ID错误")
	}
	if where == nil {
		where = &balance_log.SelectBalanceLog{}
	}
	where.UserId = strconv.FormatInt(userId, 10)
	if where.Type != "" {
		logType, err := strconv.Atoi(where.Type)
		if err != nil {
			return nil, 0, errors.New("类型参数错误")
		}
		if _, ok := balance_log_type.Map()[logType]; !ok {
			return nil, 0, errors.New("类型参数不合法")
		}
	}
	items, total, err := s.balanceLogRepository.Select(where, orderByField, orderByType, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	list = make([]*Log, 0, len(items))
	for _, item := range items {
		list = append(list, &Log{
			Id:           item.Id,
			Type:         item.Type,
			Money:        item.Money,
			MoneyBefore:  item.MoneyBefore,
			MoneyAfter:   item.MoneyAfter,
			FrozenBefore: item.FrozenBefore,
			FrozenAfter:  item.FrozenAfter,
			Describe:     item.Describe,
			CreateTime:   item.CreateTime,
		})
	}
	return list, total, nil
}

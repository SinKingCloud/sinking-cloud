package order

import (
	"errors"
	"service-pay/app/enum/order_log_operator_type"
	"service-pay/app/enum/order_log_type"
	"service-pay/app/enum/order_pay_status"
	"service-pay/app/enum/order_pay_type"
	"service-pay/app/enum/order_status"
	"service-pay/app/enum/order_type"
	"service-pay/app/model"
	orderRepository "service-pay/app/repository/order"
	orderLogRepository "service-pay/app/repository/order_log"
	"strconv"
)

// FindById 根据ID查询订单
func (s *service) FindById(id int64) (*model.Order, error) {
	if id <= 0 {
		return nil, errors.New("订单ID错误")
	}
	info, err := s.orderRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	if info == nil || info.Id <= 0 {
		return nil, errors.New("订单不存在")
	}
	return info, nil
}

// Select 查询订单列表
func (s *service) Select(where *orderRepository.SelectOrder, orderByField string, orderByType string, page int, pageSize int) (list []*orderRepository.Order, total int64, err error) {
	if where != nil {
		if where.UserId != "" {
			userId, err := strconv.ParseInt(where.UserId, 10, 64)
			if err != nil || userId <= 0 {
				return nil, 0, errors.New("用户ID错误")
			}
		}
		if where.ChannelId != "" {
			channelId, err := strconv.ParseInt(where.ChannelId, 10, 64)
			if err != nil || channelId <= 0 {
				return nil, 0, errors.New("通道ID错误")
			}
		}
		if where.Type != "" {
			if err = s.checkEnumValue(where.Type, order_type.Map(), "订单类型"); err != nil {
				return nil, 0, err
			}
		}
		if where.Status != "" {
			if err = s.checkEnumValue(where.Status, order_status.Map(), "订单状态"); err != nil {
				return nil, 0, err
			}
		}
		if where.PayStatus != "" {
			if err = s.checkEnumValue(where.PayStatus, order_pay_status.Map(), "支付状态"); err != nil {
				return nil, 0, err
			}
		}
		if where.PayType != "" {
			if err = s.checkEnumValue(where.PayType, order_pay_type.Map(), "支付类型"); err != nil {
				return nil, 0, err
			}
		}
	}
	return s.orderRepository.Select(where, orderByField, orderByType, page, pageSize)
}

// SelectByUserId 查询用户订单列表
func (s *service) SelectByUserId(userId int64, where *orderRepository.SelectOrder, orderByField string, orderByType string, page int, pageSize int) (list []*UserOrder, total int64, err error) {
	if userId <= 0 {
		return nil, 0, errors.New("用户ID错误")
	}
	if where == nil {
		where = &orderRepository.SelectOrder{}
	}
	where.UserId = strconv.FormatInt(userId, 10)
	items, total, err := s.Select(where, orderByField, orderByType, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	list = make([]*UserOrder, 0, len(items))
	for _, item := range items {
		list = append(list, &UserOrder{
			Id:           item.Id,
			TradeNo:      item.TradeNo,
			Name:         item.Name,
			Type:         item.Type,
			Amount:       item.Amount,
			Status:       item.Status,
			PayAmount:    item.PayAmount,
			RefundAmount: item.RefundAmount,
			PayStatus:    item.PayStatus,
			PayType:      item.PayType,
			PayTime:      item.PayTime,
			ExpiredTime:  item.ExpiredTime,
			CompleteTime: item.CompleteTime,
			CreateTime:   item.CreateTime,
			UpdateTime:   item.UpdateTime,
		})
	}
	return list, total, nil
}

// SelectLog 查询订单日志
func (s *service) SelectLog(where *orderLogRepository.SelectOrderLog, orderByField string, orderByType string, page int, pageSize int) (list []*model.OrderLog, total int64, err error) {
	if where != nil {
		if where.OperatorType != "" {
			if err = s.checkEnumValue(where.OperatorType, order_log_operator_type.Map(), "操作人类型"); err != nil {
				return nil, 0, err
			}
		}
		if where.Type != "" {
			if err = s.checkEnumValue(where.Type, order_log_type.Map(), "日志类型"); err != nil {
				return nil, 0, err
			}
		}
	}
	return s.orderLogRepository.Select(where, orderByField, orderByType, page, pageSize)
}

// checkEnumValue 判断枚举类型
func (s *service) checkEnumValue(value string, data map[int]string, name string) error {
	key, err := strconv.Atoi(value)
	if err != nil {
		return errors.New(name + "参数错误")
	}
	if _, ok := data[key]; !ok {
		return errors.New(name + "参数不合法")
	}
	return nil
}

package order

import (
	"encoding/json"
	"errors"
	"service-pay/app/constant"
	"service-pay/app/enum/balance_log_type"
	"service-pay/app/enum/order_log_operator_type"
	"service-pay/app/enum/order_log_type"
	"service-pay/app/enum/order_pay_status"
	"service-pay/app/enum/order_status"
	"service-pay/app/enum/order_type"
	"service-pay/app/enum/pay_type"
	"service-pay/app/model"
	orderRepository "service-pay/app/repository/order"
	orderLogRepository "service-pay/app/repository/order_log"
	"strconv"
	"time"

	payService "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/pay"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
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
		if where.Id != "" {
			orderId, err := strconv.ParseInt(where.Id, 10, 64)
			if err != nil || orderId <= 0 {
				return nil, 0, errors.New("订单ID错误")
			}
		}
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
			Name:         item.Name,
			Type:         item.Type,
			Amount:       item.Amount,
			Status:       item.Status,
			PayAmount:    item.PayAmount,
			RefundAmount: item.RefundAmount,
			PayStatus:    item.PayStatus,
			PayTime:      item.PayTime,
			PayType:      item.PayType,
			ExpiredTime:  item.ExpiredTime,
			CompleteTime: item.CompleteTime,
			CreateTime:   item.CreateTime,
			UpdateTime:   item.UpdateTime,
		})
	}
	return list, total, nil
}

// SelectLog 查询订单日志
func (s *service) SelectLog(where *orderLogRepository.SelectOrderLog, orderByField string, orderByType string, page int, pageSize int) (list []*orderLogRepository.OrderLog, total int64, err error) {
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

// FindLogById 根据ID查询日志
func (s *service) FindLogById(id int64) (*model.OrderLog, error) {
	if id <= 0 {
		return nil, errors.New("订单ID错误")
	}
	info, err := s.orderLogRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	if info == nil || info.Id <= 0 {
		return nil, errors.New("订单不存在")
	}
	return info, nil
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

// Create 创建订单
func (s *service) Create(userId int64, name string, Type int, payType string, amount float64, param string, callbackRpcService string, callbackRpcFunc string, describe string) (int64, error) {
	if userId <= 0 {
		return 0, errors.New("用户ID错误")
	}
	if name == "" {
		return 0, errors.New("订单名称不能为空")
	}
	if _, ok := order_type.Map()[Type]; !ok {
		return 0, errors.New("订单类型不合法")
	}
	if _, ok := pay_type.Map()[payType]; !ok {
		return 0, errors.New("支付类型不合法")
	}
	channel, err := s.channelService.GetPayChannel(payType)
	if err != nil {
		return 0, err
	}
	if Type == order_type.Recharge && payType == pay_type.Balance {
		return 0, errors.New("支付类型不合法")
	}
	payConfig := s.channelService.GetConfig(false)
	expireTime := str.DateTime(time.Now().Add(time.Duration(payConfig.Expire) * time.Minute))
	if amount <= 0 {
		amount = 0
	}
	extra := Extra{
		RpcService: callbackRpcService,
		RpcFunc:    callbackRpcFunc,
	}
	if param != "" {
		extra.Param = param
	}
	extraBytes, err := json.Marshal(extra)
	if err != nil {
		return 0, errors.Join(errors.New("生成订单扩展信息失败"), err)
	}
	orderId := str.GetSnowWorkIns().GetId()
	data := &model.Order{
		Id:           orderId,
		UserId:       userId,
		ChannelId:    channel.Id,
		Name:         name,
		Type:         Type,
		Extra:        string(extraBytes),
		Amount:       amount,
		Status:       order_status.UnPay,
		PayStatus:    order_pay_status.UnPay,
		PayType:      channel.Type,
		ExpiredTime:  expireTime,
		PayAmount:    0,
		RefundAmount: 0,
		CreateTime:   str.DateTime(time.Now()),
		UpdateTime:   str.DateTime(time.Now()),
	}
	err = s.database.Transaction(func(tx *gorm.DB) error {
		if err2 := s.orderRepository.Create(data, tx); err2 != nil {
			return errors.Join(errors.New("创建订单失败"), err2)
		}
		afterInfo, err2 := json.Marshal(data)
		if err2 != nil {
			return errors.Join(errors.New("生成订单日志失败"), err2)
		}
		err2 = s.orderLogRepository.Create(&model.OrderLog{
			OrderId:      data.Id,
			OperatorId:   data.UserId,
			OperatorType: order_log_operator_type.User,
			Type:         order_log_type.Create,
			Content:      describe,
			BeforeInfo:   "{}",
			AfterInfo:    string(afterInfo),
		}, tx)
		if err2 != nil {
			return errors.Join(errors.New("创建订单日志失败"), err2)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return orderId, nil
}

// Callback 订单回调
func (s *service) Callback(orderId int64, tradeNo string, payAmount float64, payTime str.DateTime) error {
	return s.callback(orderId, tradeNo, payAmount, payTime, 0, order_log_operator_type.System, "在线支付订单成功,等待发货")
}

// Deliver 订单发货
func (s *service) Deliver(id int64, operatorId int64) error {
	if id <= 0 {
		return errors.New("订单ID错误")
	}
	info, err := s.FindById(id)
	if err != nil {
		return errors.Join(errors.New("获取订单信息失败"), err)
	}
	if info.Status != order_status.UnPay && info.Status != order_status.Delivering {
		return errors.New("订单状态异常")
	}
	tradeNo := info.TradeNo
	payAmount := info.PayAmount
	payTime := info.PayTime
	if info.PayStatus != order_pay_status.Payed {
		if tradeNo == "" {
			tradeNo = strconv.FormatInt(str.GetSnowWorkIns().GetId(), 10)
		}
		payAmount = info.Amount
		payTime = str.DateTime(time.Now())
	}
	return s.callback(info.Id, tradeNo, payAmount, payTime, operatorId, order_log_operator_type.Admin, "后台手动发货")
}

// Refound 订单退款
func (s *service) Refound(id int64, money float64, operatorId int64) error {
	if id <= 0 {
		return errors.New("订单ID错误")
	}
	if money < 0 {
		return errors.New("退款金额不能小于0")
	}
	info, err := s.FindById(id)
	if err != nil {
		return errors.Join(errors.New("获取订单信息失败"), err)
	}
	lockKey := constant.LockOrderCallback + strconv.FormatInt(info.Id, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeOrderCallback) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)

	balanceLockKey := constant.LockBalanceOperate + strconv.FormatInt(info.UserId, 10)
	if !s.cache.Lock(balanceLockKey, constant.LockTimeBalanceOperate) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(balanceLockKey)

	info, err = s.FindById(id)
	if err != nil {
		return errors.Join(errors.New("获取订单信息失败"), err)
	}
	if info.Type == order_type.Recharge {
		return errors.New("充值订单不允许退款")
	}
	if info.PayStatus != order_pay_status.Payed {
		return errors.New("该订单未支付,不能退款")
	}
	if info.PayAmount <= 0 {
		return errors.New("该订单实际支付金额为0,无需退款")
	}
	refundableMoney := info.PayAmount - info.RefundAmount
	if refundableMoney <= 0 {
		return errors.New("该订单已全部退款,无需重复退款")
	}
	refundMoney := money
	if refundMoney == 0 || refundMoney > info.Amount {
		refundMoney = refundableMoney
	}
	if refundMoney <= 0 {
		return errors.New("退款金额必须大于0")
	}
	if refundMoney > refundableMoney {
		return errors.New("退款金额不能大于订单可退款金额")
	}

	newRefundAmount := info.RefundAmount + refundMoney
	newStatus := info.Status
	if newRefundAmount >= info.PayAmount {
		newStatus = order_status.Refunded
	}
	update := &orderRepository.UpdateOrder{
		RefundAmount: newRefundAmount,
	}
	if newStatus != info.Status {
		update.Status = newStatus
	}
	after := *info
	after.RefundAmount = newRefundAmount
	after.Status = newStatus
	after.UpdateTime = str.DateTime(time.Now())

	beforeInfo, err := json.Marshal(info)
	if err != nil {
		return errors.Join(errors.New("生成订单日志失败"), err)
	}
	afterInfo, err := json.Marshal(&after)
	if err != nil {
		return errors.Join(errors.New("生成订单日志失败"), err)
	}
	err = s.database.Transaction(func(tx *gorm.DB) error {
		balanceInfo, err := s.balanceRepository.FindOrCreateByUserId(info.UserId, tx)
		if err != nil {
			return errors.Join(errors.New("获取余额信息失败"), err)
		}
		if balanceInfo == nil || balanceInfo.Id <= 0 {
			return errors.New("获取余额信息失败")
		}
		moneyBefore := balanceInfo.Money
		moneyAfter := balanceInfo.Money + refundMoney
		frozenBefore := balanceInfo.Frozen
		frozenAfter := balanceInfo.Frozen
		if err = s.orderRepository.UpdateById(info.Id, update, tx); err != nil {
			return errors.Join(errors.New("更新订单失败"), err)
		}
		if err = s.orderLogRepository.Create(&model.OrderLog{
			OrderId:      info.Id,
			OperatorId:   operatorId,
			OperatorType: order_log_operator_type.Admin,
			Type:         order_log_type.Refunded,
			Content:      "退款金额:" + strconv.FormatFloat(refundMoney, 'f', 2, 64) + "元",
			BeforeInfo:   string(beforeInfo),
			AfterInfo:    string(afterInfo),
		}, tx); err != nil {
			return errors.Join(errors.New("创建订单日志失败"), err)
		}
		if err = s.balanceRepository.Money(info.UserId, refundMoney, tx); err != nil {
			return errors.Join(errors.New("更新余额失败"), err)
		}
		if err = s.balanceLogRepository.Create(&model.BalanceLog{
			Id:           str.GetSnowWorkIns().GetId(),
			UserId:       info.UserId,
			Type:         balance_log_type.Refunded,
			Money:        refundMoney,
			MoneyBefore:  moneyBefore,
			MoneyAfter:   moneyAfter,
			FrozenBefore: frozenBefore,
			FrozenAfter:  frozenAfter,
			Describe:     "订单[" + strconv.FormatInt(info.Id, 10) + "]退款",
		}, tx); err != nil {
			return errors.Join(errors.New("创建余额明细失败"), err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.cache.Delete(constant.CacheNameWithBalance + strconv.FormatInt(info.UserId, 10))
	return nil
}

// CleanExpiredUnPayOrder 清理过期未支付订单
func (s *service) CleanExpiredUnPayOrder(day int) (int64, error) {
	if day < 1 {
		day = 1
	}
	if !s.cache.Lock(constant.LockOrderCleanExpiredUnPay, constant.LockTimeOrderCleanExpiredUnPay) {
		return 0, errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(constant.LockOrderCleanExpiredUnPay)

	expiredTime := time.Now().AddDate(0, 0, -day).Format("2006-01-02 15:04:05")
	ids, err := s.orderRepository.FindExpiredUnPayIds(expiredTime, 1000)
	if err != nil {
		return 0, errors.Join(errors.New("查询过期未支付订单失败"), err)
	}
	if len(ids) <= 0 {
		return 0, nil
	}
	err = s.database.Transaction(func(tx *gorm.DB) error {
		if err := s.orderLogRepository.DeleteByOrderIds(ids, tx); err != nil {
			return errors.Join(errors.New("删除订单日志失败"), err)
		}
		if err := s.orderRepository.DeleteByIds(ids, tx); err != nil {
			return errors.Join(errors.New("删除订单失败"), err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return int64(len(ids)), nil
}

// callback 订单回调
func (s *service) callback(orderId int64, tradeNo string, payAmount float64, payTime str.DateTime, operatorId int64, operatorType int64, content string) error {
	if orderId <= 0 {
		return errors.New("订单ID错误")
	}
	if content == "" {
		content = "订单支付成功,等待发货"
	}
	lockKey := constant.LockOrderCallback + strconv.FormatInt(orderId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeOrderCallback) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	info, err := s.FindById(orderId)
	if err != nil {
		return errors.Join(errors.New("获取订单信息失败"), err)
	}
	if info.Status != order_status.UnPay && info.Status != order_status.Delivering {
		return nil
	}
	if tradeNo == "" {
		tradeNo = info.TradeNo
	}
	if info.PayStatus == order_pay_status.Payed {
		payAmount = info.PayAmount
		payTime = info.PayTime
	}
	if info.PayStatus != order_pay_status.Payed || info.Status == order_status.UnPay {
		after := *info
		after.Status = order_status.Delivering
		after.TradeNo = tradeNo
		after.PayAmount = payAmount
		after.PayStatus = order_pay_status.Payed
		after.PayTime = payTime
		after.UpdateTime = str.DateTime(time.Now())
		beforeInfo, err := json.Marshal(info)
		if err != nil {
			return errors.Join(errors.New("生成订单日志失败"), err)
		}
		afterInfo, err := json.Marshal(&after)
		if err != nil {
			return errors.Join(errors.New("生成订单日志失败"), err)
		}
		update := &orderRepository.UpdateOrder{
			Status:    order_status.Delivering,
			TradeNo:   tradeNo,
			PayAmount: payAmount,
			PayStatus: order_pay_status.Payed,
			PayTime:   payTime,
		}
		if err = s.database.Transaction(func(tx *gorm.DB) error {
			if err := s.orderRepository.UpdateById(info.Id, update, tx); err != nil {
				return errors.Join(errors.New("更新订单失败"), err)
			}
			if err := s.orderLogRepository.Create(&model.OrderLog{
				OrderId:      info.Id,
				OperatorId:   operatorId,
				OperatorType: operatorType,
				Type:         order_log_type.Pay,
				Content:      content,
				BeforeInfo:   string(beforeInfo),
				AfterInfo:    string(afterInfo),
			}, tx); err != nil {
				return errors.Join(errors.New("创建订单日志失败"), err)
			}
			return nil
		}); err != nil {
			return err
		}
		info = &after
	}

	extra := Extra{}
	if info.Extra != "" {
		if err = json.Unmarshal([]byte(info.Extra), &extra); err != nil {
			return errors.Join(errors.New("订单扩展信息错误"), err)
		}
	}
	if extra.RpcService == "" || extra.RpcFunc == "" {
		return errors.New("订单回调RPC配置不能为空")
	}

	notify := &payService.NotifyParam{
		UserId:    info.UserId,
		OrderId:   info.Id,
		Extra:     extra.Param,
		Amount:    info.Amount,
		TradeNo:   tradeNo,
		PayAmount: payAmount,
		PayType:   info.PayType,
		PayTime:   payTime,
	}
	var result bool
	if err := rpc.InvokeFunc(extra.RpcService, extra.RpcFunc, notify, &result); err != nil {
		return errors.Join(errors.New("业务回调失败"), err)
	}
	if !result {
		return errors.New("业务回调返回失败")
	}

	completeTime := str.DateTime(time.Now())
	after := *info
	after.Status = order_status.Credited
	after.CompleteTime = completeTime
	after.UpdateTime = str.DateTime(time.Now())
	beforeInfo, err := json.Marshal(info)
	if err != nil {
		return errors.Join(errors.New("生成订单日志失败"), err)
	}
	afterInfo, err := json.Marshal(&after)
	if err != nil {
		return errors.Join(errors.New("生成订单日志失败"), err)
	}
	return s.database.Transaction(func(tx *gorm.DB) error {
		if err := s.orderRepository.UpdateById(info.Id, &orderRepository.UpdateOrder{
			Status:       order_status.Credited,
			CompleteTime: completeTime,
		}, tx); err != nil {
			return errors.Join(errors.New("更新订单失败"), err)
		}
		if err := s.orderLogRepository.Create(&model.OrderLog{
			OrderId:      info.Id,
			OperatorId:   operatorId,
			OperatorType: operatorType,
			Type:         order_log_type.Pay,
			Content:      "订单已到账",
			BeforeInfo:   string(beforeInfo),
			AfterInfo:    string(afterInfo),
		}, tx); err != nil {
			return errors.Join(errors.New("创建订单日志失败"), err)
		}
		return nil
	})
}

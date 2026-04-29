package order

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
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

// FindByTradeNo 根据traderNo查询订单
func (s *service) FindByTradeNo(tradeNo string) (*model.Order, error) {
	if tradeNo == "" {
		return nil, errors.New("订单ID错误")
	}
	info, err := s.orderRepository.FindByTradeNo(tradeNo)
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
func (s *service) Create(userId int64, name string, Type int, payType string, amount float64, param interface{}, callbackRpcService string, callbackRpcFunc string, describe string) (string, error) {
	if userId <= 0 {
		return "", errors.New("用户ID错误")
	}
	if name == "" {
		return "", errors.New("订单名称不能为空")
	}
	if _, ok := order_type.Map()[Type]; !ok {
		return "", errors.New("订单类型不合法")
	}
	if _, ok := pay_type.Map()[payType]; !ok {
		return "", errors.New("支付类型不合法")
	}
	channelID, err := s.channelService.GetChannelId(payType)
	if err != nil {
		return "", err
	}
	if Type == order_type.Recharge && payType == pay_type.Balance {
		return "", errors.New("支付类型不合法")
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
	if param != nil {
		extra.Param, _ = json.Marshal(param)
	}
	extraBytes, err := json.Marshal(extra)
	if err != nil {
		return "", errors.Join(errors.New("生成订单扩展信息失败"), err)
	}
	tradeNo := func() string {
		var seed string
		formatInt := strconv.FormatInt(str.GetSnowWorkIns().GetId(), 10)
		if len(formatInt) > 6 {
			seed = formatInt[len(formatInt)-6:]
		} else {
			seed = strconv.Itoa(100000 + rand.IntN(899999))
		}
		return time.Now().Format("20060102150405") + seed
	}()
	data := &model.Order{
		Id:           str.GetSnowWorkIns().GetId(),
		UserId:       userId,
		ChannelId:    channelID,
		TradeNo:      tradeNo,
		Name:         name,
		Type:         Type,
		Extra:        string(extraBytes),
		Amount:       amount,
		Status:       order_status.UnPay,
		PayStatus:    order_pay_status.UnPay,
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
		return "", err
	}
	return tradeNo, nil
}

// Callback 订单回调
func (s *service) Callback(tradeNo string, outTradeNo string, payAmount float64, payTime str.DateTime) error {
	return s.callback(tradeNo, outTradeNo, payAmount, payTime, 0, order_log_operator_type.System, "在线支付订单成功,发货成功")
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
	if info.Status != order_status.UnPay {
		return errors.New("订单状态异常")
	}
	return s.callback(info.TradeNo, strconv.FormatInt(str.GetSnowWorkIns().GetId(), 10), info.Amount, str.DateTime(time.Now()), operatorId, order_log_operator_type.Admin, "后台手动发货")
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
	lockKey := constant.LockOrderCallback + info.TradeNo
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
			Describe:     "订单[" + info.TradeNo + "]退款",
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

// callback 订单回调
func (s *service) callback(tradeNo string, outTradeNo string, payAmount float64, payTime str.DateTime, operatorId int64, operatorType int64, content string) error {
	if tradeNo == "" {
		return errors.New("订单号不能为空")
	}
	if content == "" {
		content = "订单发货成功"
	}
	lockKey := constant.LockOrderCallback + tradeNo
	if !s.cache.Lock(lockKey, constant.LockTimeOrderCallback) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	info, err := s.FindByTradeNo(tradeNo)
	if err != nil {
		return errors.Join(errors.New("获取订单信息失败"), err)
	}
	if info.Status == order_status.Payed || info.PayStatus == order_pay_status.Payed {
		return nil
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
		UserId:     info.UserId,
		TradeNo:    info.TradeNo,
		Extra:      string(extra.Param),
		Amount:     info.Amount,
		OutTradeNo: outTradeNo,
		PayAmount:  payAmount,
		PayTime:    payTime,
	}
	completeTime := str.DateTime(time.Now())
	update := &orderRepository.UpdateOrder{
		Status:       order_status.Payed,
		OutTradeNo:   outTradeNo,
		PayAmount:    payAmount,
		PayStatus:    order_pay_status.Payed,
		PayTime:      payTime,
		CompleteTime: completeTime,
	}

	after := *info
	after.Status = order_status.Payed
	after.OutTradeNo = outTradeNo
	after.PayAmount = payAmount
	after.PayStatus = order_pay_status.Payed
	after.PayTime = payTime
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
		var result bool
		if err := rpc.InvokeFunc(extra.RpcService, extra.RpcFunc, notify, &result); err != nil {
			return errors.Join(errors.New("业务回调失败"), err)
		}
		if !result {
			return errors.New("业务回调返回失败")
		}
		return nil
	})
}

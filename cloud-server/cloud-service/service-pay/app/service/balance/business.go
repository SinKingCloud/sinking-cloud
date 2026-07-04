package balance

import (
	"encoding/json"
	"errors"
	"fmt"
	"service-pay/app/constant"
	"service-pay/app/enum/balance_log_type"
	"service-pay/app/enum/order_type"
	"service-pay/app/model"
	"service-pay/app/repository/balance_log"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/pay"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// GetByUserId 通过用户ID查询余额
func (s *service) GetByUserId(userId int64, refresh bool) (*Balance, error) {
	if userId <= 0 {
		return nil, errors.New("用户ID错误")
	}
	cacheKey := constant.CacheNameWithBalance + strconv.FormatInt(userId, 10)
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			var result *Balance
			if err := json.Unmarshal([]byte(data), &result); err == nil {
				return result, nil
			}
		}
	}
	info, err := s.balanceRepository.FindOrCreateByUserId(userId)
	if err != nil {
		return nil, err
	}
	if info == nil || info.Id <= 0 {
		return nil, errors.New("获取余额信息失败")
	}
	data := &Balance{
		Money:         info.Money,
		Frozen:        info.Frozen,
		TotalRecharge: info.TotalRecharge,
		TotalConsume:  info.TotalConsume,
		CreateTime:    info.CreateTime,
		UpdateTime:    info.UpdateTime,
	}
	if bytes, err2 := json.Marshal(data); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(bytes), constant.CacheTimeWithBalance)
	}
	return data, nil
}

// OperateMoney 操作用户余额并写入明细
func (s *service) OperateMoney(userId int64, money float64, isOperateFrozen bool, logType int, describe string) error {
	if userId <= 0 {
		return errors.New("用户ID错误")
	}
	if money == 0 {
		return errors.New("变更金额不能为0")
	}
	if _, ok := balance_log_type.Map()[logType]; !ok {
		return errors.New("类型参数不合法")
	}
	lockKey := constant.LockBalanceOperate + strconv.FormatInt(userId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeBalanceOperate) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	err := s.database.Transaction(func(tx *gorm.DB) error {
		info, err := s.balanceRepository.FindOrCreateByUserId(userId, tx)
		if err != nil {
			return errors.Join(errors.New("获取余额信息失败"), err)
		}
		if info == nil || info.Id <= 0 {
			return errors.New("获取余额信息失败")
		}
		moneyBefore := info.Money
		moneyAfter := info.Money
		frozenBefore := info.Frozen
		frozenAfter := info.Frozen
		logMoney := money
		if isOperateFrozen {
			moneyAfter -= money
			frozenAfter += money
			logMoney = -money
		} else {
			moneyAfter += money
		}
		if moneyAfter < 0 {
			return errors.New("余额不足")
		}
		if frozenAfter < 0 {
			return errors.New("冻结余额不足")
		}
		if isOperateFrozen {
			if err = s.balanceRepository.Money(userId, -money, tx); err != nil {
				return errors.Join(errors.New("更新余额失败"), err)
			}
			if err = s.balanceRepository.Frozen(userId, money, tx); err != nil {
				return errors.Join(errors.New("更新冻结余额失败"), err)
			}
		} else {
			if err = s.balanceRepository.Money(userId, money, tx); err != nil {
				return errors.Join(errors.New("更新余额失败"), err)
			}
		}
		err = s.balanceLogRepository.Create(&model.BalanceLog{
			Id:           str.GetSnowWorkIns().GetId(),
			UserId:       userId,
			Type:         logType,
			Money:        logMoney,
			MoneyBefore:  moneyBefore,
			MoneyAfter:   moneyAfter,
			FrozenBefore: frozenBefore,
			FrozenAfter:  frozenAfter,
			Describe:     describe,
		}, tx)
		if err != nil {
			return errors.Join(errors.New("创建余额明细失败"), err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.cache.Delete(constant.CacheNameWithBalance + strconv.FormatInt(userId, 10))
	return nil
}

// OperateMoneyCount 操作累计充值或消费金额
func (s *service) OperateMoneyCount(userId int64, money float64) error {
	if userId <= 0 {
		return errors.New("用户ID错误")
	}
	if money == 0 {
		return errors.New("变更金额不能为0")
	}
	lockKey := constant.LockBalanceOperate + strconv.FormatInt(userId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeBalanceOperate) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	if money > 0 {
		if err := s.balanceRepository.TotalRecharge(userId, money); err != nil {
			return errors.Join(errors.New("更新累计充值失败"), err)
		}
	} else {
		if err := s.balanceRepository.TotalConsume(userId, -money); err != nil {
			return errors.Join(errors.New("更新累计消费失败"), err)
		}
	}
	s.cache.Delete(constant.CacheNameWithBalance + strconv.FormatInt(userId, 10))
	return nil
}

// SelectLog 查询用户余额明细
func (s *service) SelectLog(where *balance_log.SelectBalanceLog, page *page.Query) (list []*Log, total int64, err error) {
	if where == nil {
		return nil, 0, errors.New("参数错误")
	}
	if where.UserId != "" {
		if _, err = strconv.ParseInt(where.UserId, 10, 64); err != nil {
			return nil, 0, errors.New("用户ID格式错误")
		}
	}
	if where.Type != "" {
		logType, err := strconv.Atoi(where.Type)
		if err != nil {
			return nil, 0, errors.New("类型参数错误")
		}
		if _, ok := balance_log_type.Map()[logType]; !ok {
			return nil, 0, errors.New("类型参数不合法")
		}
	}
	items, total, err := s.balanceLogRepository.Select(where, page)
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

// CreateRechargeOrder 创建充值订单
func (s *service) CreateRechargeOrder(userId int64, money int, payType string) (int64, error) {
	if money <= 0 {
		return 0, errors.New("充值金额必须大于0")
	}
	if s.channelService == nil {
		return 0, errors.New("获取支付配置失败")
	}
	payConfig := s.channelService.GetConfig(false)
	if payConfig == nil {
		return 0, errors.New("获取支付配置失败")
	}
	amount := float64(money)
	if payConfig.Min > 0 && amount < payConfig.Min {
		return 0, errors.New("充值金额不能小于" + strconv.FormatFloat(payConfig.Min, 'f', -1, 64) + "元")
	}
	if payConfig.Max > 0 && amount > payConfig.Max {
		return 0, errors.New("充值金额不能大于" + strconv.FormatFloat(payConfig.Max, 'f', -1, 64) + "元")
	}
	orderId, err := s.orderService.Create(userId,
		"在线充值",
		order_type.Recharge,
		payType,
		amount,
		"",
		s.consul.GetName(),
		"RechargeOrderNotify", "发起在线充值订单")
	if err != nil {
		return 0, err
	}
	return orderId, nil
}

// RechargeOrderNotify 充值订单回调
func (s *service) RechargeOrderNotify(order *pay.NotifyParam) error {
	if order == nil {
		return errors.New("订单信息不能为空")
	}
	if order.UserId <= 0 {
		return errors.New("用户ID错误")
	}
	if order.OrderId <= 0 {
		return errors.New("订单ID错误")
	}
	if order.PayAmount <= 0 {
		return errors.New("支付金额必须大于0")
	}
	lockKey := constant.LockBalanceOperate + strconv.FormatInt(order.UserId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeBalanceOperate) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	logInfo, err := s.balanceLogRepository.FindById(order.OrderId)
	if err != nil {
		return errors.Join(errors.New("查询充值明细失败"), err)
	}
	if logInfo != nil && logInfo.Id > 0 {
		return nil
	}
	err = s.database.Transaction(func(tx *gorm.DB) error {
		info, err := s.balanceRepository.FindOrCreateByUserId(order.UserId, tx)
		if err != nil {
			return errors.Join(errors.New("获取余额信息失败"), err)
		}
		if info == nil || info.Id <= 0 {
			return errors.New("获取余额信息失败")
		}
		moneyBefore := info.Money
		moneyAfter := info.Money + order.PayAmount
		frozenBefore := info.Frozen
		frozenAfter := info.Frozen
		if err = s.balanceRepository.Money(order.UserId, order.PayAmount, tx); err != nil {
			return errors.Join(errors.New("更新余额失败"), err)
		}
		if err = s.balanceRepository.TotalRecharge(order.UserId, order.PayAmount, tx); err != nil {
			return errors.Join(errors.New("更新累计充值失败"), err)
		}
		if err = s.balanceLogRepository.Create(&model.BalanceLog{
			Id:           order.OrderId,
			UserId:       order.UserId,
			Type:         balance_log_type.Recharge,
			Money:        order.PayAmount,
			MoneyBefore:  moneyBefore,
			MoneyAfter:   moneyAfter,
			FrozenBefore: frozenBefore,
			FrozenAfter:  frozenAfter,
			Describe:     fmt.Sprintf("在线充值[%d]", order.OrderId),
		}, tx); err != nil {
			return errors.Join(errors.New("创建余额明细失败"), err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.cache.Delete(constant.CacheNameWithBalance + strconv.FormatInt(order.UserId, 10))
	return nil
}

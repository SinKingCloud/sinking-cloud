package balance

import (
	"encoding/json"
	"errors"
	"service-pay/app/constant"
	"service-pay/app/enum/balance_log_type"
	"service-pay/app/model"
	"service-pay/app/repository/balance_log"
	"strconv"

	"gorm.io/gorm"
)

// FindByUserId 通过用户ID查询余额
func (s *service) FindByUserId(userId int64, refresh bool) (*Balance, error) {
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
func (s *service) OperateMoney(userId int64, money float64, logType int, describe string) error {
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
		moneyAfter := info.Money + money
		if moneyAfter < 0 {
			return errors.New("余额不足")
		}
		if err = s.balanceRepository.Money(userId, money, tx); err != nil {
			return errors.Join(errors.New("更新余额失败"), err)
		}
		err = s.balanceLogRepository.Create(&model.BalanceLog{
			UserId:       userId,
			Type:         logType,
			Money:        money,
			MoneyBefore:  info.Money,
			MoneyAfter:   moneyAfter,
			FrozenBefore: info.Frozen,
			FrozenAfter:  info.Frozen,
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
func (s *service) SelectLog(where *balance_log.SelectBalanceLog, orderByField string, orderByType string, page int, pageSize int) (list []*Log, total int64, err error) {
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

package auth

import (
	"errors"
	constant2 "service-admin/app/constant"
	sessionRepository "service-admin/app/repository/session"

	"gorm.io/gorm"
)

// SelectSessionByUserId 查询用户会话列表
func (s *service) SelectSessionByUserId(userId int64, where *sessionRepository.SelectSession, orderByField string, orderByType string, page int, pageSize int) ([]*Session, int64, error) {
	if userId <= 0 {
		return nil, 0, errors.New("用户ID错误")
	}
	if where == nil {
		where = &sessionRepository.SelectSession{}
	}
	where.UserId = userId
	sessions, total, err := s.repositorySession.Select(where, orderByField, orderByType, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	data := make([]*Session, 0, len(sessions))
	for _, item := range sessions {
		data = append(data, &Session{
			Token:      item.Token,
			Device:     item.Device,
			Ip:         item.Ip,
			Location:   item.Location,
			ExpireTime: item.ExpireTime,
			CreateTime: item.CreateTime,
			UpdateTime: item.UpdateTime,
		})
	}
	return data, total, nil
}

// LogoutByUserIdAndTokens 根据用户ID和token集合注销会话
func (s *service) LogoutByUserIdAndTokens(userId int64, tokens []string) error {
	if userId <= 0 || len(tokens) == 0 {
		return errors.New("参数错误")
	}
	if err := s.repositorySession.DeleteByUserIdAndTokens(userId, tokens); err != nil {
		return errors.Join(errors.New("注销会话失败"), err)
	}
	s.clearSessionCache(tokens...)
	return nil
}

// ClearExpireSessions 清理过期会话
func (s *service) ClearExpireSessions() error {
	tokens := make([]string, 0, 100)
	err := s.database.Transaction(func(tx *gorm.DB) error {
		expireSessions, err := s.repositorySession.SelectExpireSessions(tx)
		if err != nil {
			return errors.Join(errors.New("查询过期会话失败"), err)
		}
		if len(expireSessions) == 0 {
			return nil
		}
		ids := make([]int64, 0, len(expireSessions))
		for id, token := range expireSessions {
			ids = append(ids, id)
			tokens = append(tokens, token)
		}
		if err = s.repositorySession.DeleteByIds(ids, tx); err != nil {
			return errors.Join(errors.New("删除过期会话失败"), err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	s.clearSessionCache(tokens...)
	return nil
}

// clearSessionCache 清理会话缓存
func (s *service) clearSessionCache(tokens ...string) {
	if len(tokens) == 0 {
		return
	}
	keys := make([]string, 0, len(tokens))
	for _, token := range tokens {
		keys = append(keys, constant2.CacheNameWithSession+token)
	}
	s.cache.Delete(keys...)
}

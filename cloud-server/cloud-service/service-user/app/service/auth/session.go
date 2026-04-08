package auth

import (
	"errors"
	constant2 "service-user/app/constant"
	"service-user/app/model"
	sessionRepository "service-user/app/repository/session"
)

// SelectSessionByUserId 查询用户会话列表
func (s *service) SelectSessionByUserId(userId int64, where *sessionRepository.SelectSession, orderByField string, orderByType string, page int, pageSize int) ([]*model.Session, int64, error) {
	if userId <= 0 {
		return nil, 0, errors.New("用户ID错误")
	}
	if where == nil {
		where = &sessionRepository.SelectSession{}
	}
	where.UserId = userId
	list, total, err := s.repositorySession.Select(where, orderByField, orderByType, page, pageSize)
	if err != nil {
		return nil, 0, errors.Join(errors.New("查询会话失败"), err)
	}
	return list, total, nil
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

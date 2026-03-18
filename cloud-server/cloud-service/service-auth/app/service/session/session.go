package session

import "service-auth/app/repository"

// Service Session服务
type Service struct {
	repo repository.Session
}

func NewService(repo repository.Session) *Service {
	return &Service{repo: repo}
}

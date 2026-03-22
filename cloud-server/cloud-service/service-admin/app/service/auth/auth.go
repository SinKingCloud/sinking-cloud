package auth

import (
	"service-admin/app/repository"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Session服务
type service struct {
	authService       *auth.Service
	repositorySession repository.Session
	repositoryUser    repository.User
	cache             cache.Interface
}

// NewService 实例化service
func NewService(authService *auth.Service, repositorySession repository.Session, repositoryUser repository.User, cache cache.Interface) *service {
	return &service{
		authService:       authService,
		repositorySession: repositorySession,
		repositoryUser:    repositoryUser,
		cache:             cache,
	}
}

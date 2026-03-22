package auth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Session服务
type service struct {
	cache         cache.Interface
	configService *config.Service
	authService   *auth.Service
}

// NewService 实例化service
func NewService(cache cache.Interface, configService *config.Service, authService *auth.Service) *service {
	return &service{
		cache:         cache,
		configService: configService,
		authService:   authService,
	}
}

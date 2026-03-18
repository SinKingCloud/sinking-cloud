package service

import (
	"service-auth/app/repository"
	"service-auth/app/service/session"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

var Session *session.Service

func Init() {
	Session = session.NewService(repository.NewSession(global.Database))
}

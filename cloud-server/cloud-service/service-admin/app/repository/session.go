package repository

import (
	"service-admin/app/model"
	"service-admin/app/repository/session"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Session Session仓储接口
type Session interface {
	Create(session *model.Session) error
}

// NewSession 创建Session仓储
func NewSession(db *database.Database) *session.Repository {
	return &session.Repository{Database: db}
}

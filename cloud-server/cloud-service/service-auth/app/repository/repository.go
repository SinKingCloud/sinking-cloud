package repository

import (
	"service-auth/app/model"
	"service-auth/app/repository/session"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Session Session仓储接口
type Session interface {
	FindByToken(token string) (*model.Session, error)
	FindByUserIdAndType(userId int64, userType int) ([]*model.Session, error)
	Create(session *model.Session) error
	Delete(token string) error
	DeleteByUserId(userId int64) error
}

// NewSession 创建Session仓储
func NewSession(db *database.Database) *session.Repository {
	return &session.Repository{Database: db}
}

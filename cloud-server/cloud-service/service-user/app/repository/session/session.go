package session

import (
	"service-user/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Repository Session仓储
type Repository struct {
	Database *database.Database
}

// Create 创建会话
func (r *Repository) Create(session *model.Session) error {
	return r.Database.Db.Create(session).Error
}

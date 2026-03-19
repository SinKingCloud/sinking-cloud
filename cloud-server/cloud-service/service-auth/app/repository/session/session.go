package session

import (
	"service-auth/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Repository Session仓储
type Repository struct {
	Database *database.Database
}

// FindByToken 根据token查询session
func (r *Repository) FindByToken(token string) (*model.Session, error) {
	var session model.Session
	err := r.Database.Db.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

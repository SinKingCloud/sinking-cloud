package session

import (
	"service-auth/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Repository Session仓储
type Repository struct {
	Database *database.Database
}

// FindByToken 根据token查询
func (r *Repository) FindByToken(token string) (*model.Session, error) {
	var session model.Session
	err := r.Database.Db.Where("`token` = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// FindByUserIdAndType 根据用户ID和类型查询
func (r *Repository) FindByUserIdAndType(userId int64, userType int) ([]*model.Session, error) {
	var sessions []*model.Session
	err := r.Database.Db.Where("`user_id` = ? and `user_type` = ?", userId, userType).Find(&sessions).Error
	return sessions, err
}

// Create 创建会话
func (r *Repository) Create(session *model.Session) error {
	return r.Database.Db.Create(session).Error
}

// Delete 删除会话
func (r *Repository) Delete(token string) error {
	return r.Database.Db.Where("`token` = ?", token).Delete(&model.Session{}).Error
}

// DeleteByUserId 根据用户ID删除
func (r *Repository) DeleteByUserId(userId int64) error {
	return r.Database.Db.Where("`user_id` = ?", userId).Delete(&model.Session{}).Error
}

package permission

import (
	"errors"
	"service-admin/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// FindByPath 根据路径查询权限
func (r *Repository) FindByPath(path string) (perm *model.Permission, err error) {
	err = r.Database.Db.Where("path = ?", path).First(&perm).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return perm, nil
}

// FindByCode 根据 code 查询权限
func (r *Repository) FindByCode(code string) (perm *model.Permission, err error) {
	err = r.Database.Db.Where("code = ?", code).First(&perm).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return perm, nil
}

// Create 创建权限
func (r *Repository) Create(perm *model.Permission) error {
	return r.Database.Db.Create(perm).Error
}

// UpdateByCode 根据 code 更新权限
func (r *Repository) UpdateByCode(code string, perm *model.Permission) error {
	if perm == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if perm.Name != "" {
		updates["name"] = perm.Name
	}
	if perm.Code != "" {
		updates["code"] = perm.Code
	}
	if perm.Path != "" {
		updates["path"] = perm.Path
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return r.Database.Db.Model(&model.Permission{}).Where("code = ?", code).Updates(updates).Error
}

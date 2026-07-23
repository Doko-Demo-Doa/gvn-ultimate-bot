package discordrepos

import (
	"doko/gvn-ultimate-bot/models"

	"gorm.io/gorm"
)

type AdminWhitelistedRoleRepo interface {
	List() ([]*models.AdminWhitelistedRole, error)
	GetByRoleNativeIDs(roleNativeIDs []string) ([]*models.AdminWhitelistedRole, error)
	Upsert(role *models.AdminWhitelistedRole) (*models.AdminWhitelistedRole, error)
	Delete(id uint) error
}

type adminWhitelistedRoleRepo struct {
	db *gorm.DB
}

func NewAdminWhitelistedRoleRepo(db *gorm.DB) AdminWhitelistedRoleRepo {
	return &adminWhitelistedRoleRepo{
		db: db,
	}
}

func (r *adminWhitelistedRoleRepo) List() ([]*models.AdminWhitelistedRole, error) {
	var roles []*models.AdminWhitelistedRole
	if err := r.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *adminWhitelistedRoleRepo) GetByRoleNativeIDs(roleNativeIDs []string) ([]*models.AdminWhitelistedRole, error) {
	var roles []*models.AdminWhitelistedRole
	if len(roleNativeIDs) == 0 {
		return roles, nil
	}
	if err := r.db.Where("role_native_id IN ?", roleNativeIDs).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *adminWhitelistedRoleRepo) Upsert(role *models.AdminWhitelistedRole) (*models.AdminWhitelistedRole, error) {
	var existing models.AdminWhitelistedRole
	err := r.db.Where(&models.AdminWhitelistedRole{RoleNativeID: role.RoleNativeID}).First(&existing).Error
	if err != nil {
		if err := r.db.Create(role).Error; err != nil {
			return nil, err
		}
		return role, nil
	}

	existing.AccessLevel = role.AccessLevel
	existing.Label = role.Label
	if err := r.db.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func (r *adminWhitelistedRoleRepo) Delete(id uint) error {
	return r.db.Delete(&models.AdminWhitelistedRole{}, id).Error
}

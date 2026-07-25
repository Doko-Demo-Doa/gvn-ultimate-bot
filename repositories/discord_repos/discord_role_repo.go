package discordrepos

import (
	"doko/gvn-ultimate-bot/models"
	"errors"

	"gorm.io/gorm"
)

type DiscordRoleRepo interface {
	GetByID(id uint) (*models.DiscordRole, error)
	GetByNativeID(nativeId string) (*models.DiscordRole, error)
	AssignRole(user *models.DiscordUser, toRole models.DiscordRole) (*models.DiscordUserRole, error)
	UnassignRole(user *models.DiscordUser, fromRole models.DiscordRole) (*models.DiscordUserRole, error) // For history
	CreateRole(role *models.DiscordRole) (*models.DiscordRole, error)                                    // Actually upsert
	EditRole(role *models.DiscordRole) (*models.DiscordRole, error)
	ListRoles() ([]*models.DiscordRole, error)
	// Upsert updates the Discord-sourced fields (Name, Mentionable, Hoist,
	// Color) of the role matching NativeID, or creates it if absent. Locally
	// managed fields (ImplicitType, Expiry) are left untouched on update.
	Upsert(role *models.DiscordRole) (*models.DiscordRole, error)
	// DeleteNotIn removes all roles whose NativeID is not in nativeIds.
	// If nativeIds is empty, all roles are removed.
	DeleteNotIn(nativeIds []string) (int64, error)
}

type discordRoleRepo struct {
	db *gorm.DB
}

func NewDiscordRoleRepo(db *gorm.DB) DiscordRoleRepo {
	return &discordRoleRepo{
		db: db,
	}
}

// AssignRole implements DiscordRoleRepo
func (*discordRoleRepo) AssignRole(user *models.DiscordUser, toRole models.DiscordRole) (*models.DiscordUserRole, error) {
	panic("unimplemented")
}

// UnassignRole implements DiscordRoleRepo
func (*discordRoleRepo) UnassignRole(user *models.DiscordUser, fromRole models.DiscordRole) (*models.DiscordUserRole, error) {
	panic("unimplemented")
}

func (dr *discordRoleRepo) CreateRole(role *models.DiscordRole) (*models.DiscordRole, error) {
	var r models.DiscordRole
	if err := dr.db.Where(&models.DiscordRole{NativeID: role.NativeID}).First(&r).Error; err != nil {
		dr.db.Create(&role)
		return role, err
	}

	return &r, nil
}

func (dr *discordRoleRepo) EditRole(role *models.DiscordRole) (*models.DiscordRole, error) {
	var r models.DiscordRole

	// Query it first
	if err := dr.db.Where(&models.DiscordRole{NativeID: role.NativeID}).First(&r).Error; err != nil {
		return role, err
	}

	r.Name = role.Name
	r.Mentionable = role.Mentionable
	r.Hoist = role.Hoist
	r.Color = role.Color
	r.Expiry = role.Expiry
	r.ImplicitType = role.ImplicitType

	if err := dr.db.Save(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// Upsert refreshes the Discord-sourced fields of the role matching NativeID
// (creating it if absent), leaving locally managed fields (ImplicitType,
// Expiry) untouched. Used by SyncGuildRoles to reconcile the discord_role
// table against Discord's actual role list.
func (dr *discordRoleRepo) Upsert(role *models.DiscordRole) (*models.DiscordRole, error) {
	var existing models.DiscordRole
	err := dr.db.Where("native_id = ?", role.NativeID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := dr.db.Create(role).Error; err != nil {
				return nil, err
			}
			return role, nil
		}
		return nil, err
	}

	existing.Name = role.Name
	existing.Mentionable = role.Mentionable
	existing.Hoist = role.Hoist
	existing.Color = role.Color
	if err := dr.db.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func (dr *discordRoleRepo) DeleteNotIn(nativeIds []string) (int64, error) {
	q := dr.db
	if len(nativeIds) > 0 {
		q = q.Where("native_id NOT IN ?", nativeIds)
	}
	result := q.Delete(&models.DiscordRole{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (dr *discordRoleRepo) GetByID(id uint) (*models.DiscordRole, error) {
	var r models.DiscordRole
	if err := dr.db.First(&r, id).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (dr *discordRoleRepo) GetByNativeID(nativeId string) (*models.DiscordRole, error) {
	var r models.DiscordRole
	if err := dr.db.Where(&models.DiscordRole{NativeID: nativeId}).First(&r).Error; err != nil {
		return nil, err
	}

	return &r, nil
}

func (dr *discordRoleRepo) ListRoles() ([]*models.DiscordRole, error) {
	var roles []*models.DiscordRole
	if err := dr.db.Find(&roles).Error; err != nil {
		return roles, err
	}

	return roles, nil
}

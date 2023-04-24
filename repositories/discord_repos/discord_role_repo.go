package discordrepos

import (
	"doko/gvn-ultimate-bot/models"

	"gorm.io/gorm"
)

type DiscordRoleRepo interface {
	GetByID(id uint) (*models.DiscordRole, error)
	GetByNativeID(nativeId string) (*models.DiscordRole, error)
	AssignRole(user models.DiscordUser, toRole models.DiscordRole) (*models.DiscordUserRole, error)
	UnassignRole(user models.DiscordUser, fromRole models.DiscordRole) (*models.DiscordUserRole, error) // For history
	CreateRole(role *models.DiscordRole) (*models.DiscordRole, error)                                   // Actually upsert
	EditRole(role *models.DiscordRole) (*models.DiscordRole, error)
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
func (*discordRoleRepo) AssignRole(user models.DiscordUser, toRole models.DiscordRole) (*models.DiscordUserRole, error) {
	panic("unimplemented")
}

// UnassignRole implements DiscordRoleRepo
func (*discordRoleRepo) UnassignRole(user models.DiscordUser, fromRole models.DiscordRole) (*models.DiscordUserRole, error) {
	panic("unimplemented")
}

func (dr *discordRoleRepo) CreateRole(role *models.DiscordRole) (*models.DiscordRole, error) {
	var r models.DiscordRole
	if err := dr.db.Where(&models.DiscordRole{NativeId: role.NativeId}).First(&r).Error; err != nil {
		dr.db.Create(&role)
		return role, err
	}

	return &r, nil
}

func (dr *discordRoleRepo) EditRole(role *models.DiscordRole) (*models.DiscordRole, error) {
	var r models.DiscordRole

	// Query it first
	if err := dr.db.Where(&models.DiscordRole{NativeId: role.NativeId}).First(&r).Error; err != nil {
		return role, err
	}

	dr.db.Save(&r)
	return &r, nil
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
	if err := dr.db.Where(&models.DiscordRole{NativeId: nativeId}).First(&r).Error; err != nil {
		return nil, err
	}

	return &r, nil
}

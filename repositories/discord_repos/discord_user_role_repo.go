package discordrepos

import (
	"doko/gvn-ultimate-bot/models"
	"time"

	"gorm.io/gorm"
)

type DiscordUserRoleRepo interface {
	CreateAssignment(assignment *models.DiscordUserRole) (*models.DiscordUserRole, error)
	GetExpiredAssignments() ([]*models.DiscordUserRole, error)
	GetActiveAssignmentsByUser(nativeUserID string) ([]*models.DiscordUserRole, error)
	RevokeAssignment(id uint) error
}

type discordUserRoleRepo struct {
	db *gorm.DB
}

func NewDiscordUserRoleRepo(db *gorm.DB) DiscordUserRoleRepo {
	return &discordUserRoleRepo{db: db}
}

func (r *discordUserRoleRepo) CreateAssignment(assignment *models.DiscordUserRole) (*models.DiscordUserRole, error) {
	if err := r.db.Create(assignment).Error; err != nil {
		return nil, err
	}
	return assignment, nil
}

func (r *discordUserRoleRepo) GetExpiredAssignments() ([]*models.DiscordUserRole, error) {
	var assignments []*models.DiscordUserRole
	now := time.Now()
	if err := r.db.
		Where("expiration_date <= ?", now).
		Where("deleted_at IS NULL").
		Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

func (r *discordUserRoleRepo) GetActiveAssignmentsByUser(nativeUserID string) ([]*models.DiscordUserRole, error) {
	var assignments []*models.DiscordUserRole
	now := time.Now()
	if err := r.db.
		Where("user_native_id = ?", nativeUserID).
		Where("expiration_date > ?", now).
		Where("deleted_at IS NULL").
		Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

func (r *discordUserRoleRepo) RevokeAssignment(id uint) error {
	return r.db.Delete(&models.DiscordUserRole{}, id).Error
}

package discordservice

import (
	"doko/gvn-ultimate-bot/models"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
)

// Interface for DiscordService
type DiscordService interface {
	// Listing the roles as DiscordRole model array
	ListRoles() ([]*models.DiscordRole, error)
	CreateRole(*models.DiscordRole) (*models.DiscordRole, error)
	EditRole() (*models.DiscordRole, error)

	// This is actually just mark the role as "deleted" (IsDeleted = 1), not actually delete it
	RemoveRole(email uint) (*models.DiscordRole, error)
}

type discordService struct {
	Repo discordrepos.DiscordRoleRepo
}

func NewDiscordService() DiscordService {
	return &discordService{}
}

func (dr *discordService) CreateRole(r *models.DiscordRole) (*models.DiscordRole, error) {
	return dr.Repo.CreateRole(r)
}

// EditRole implements DiscordService
func (*discordService) EditRole() (*models.DiscordRole, error) {
	panic("unimplemented")
}

// ListRoles implements DiscordService
func (*discordService) ListRoles() ([]*models.DiscordRole, error) {
	panic("unimplemented")
}

// RemoveRole implements DiscordService
func (*discordService) RemoveRole(email uint) (*models.DiscordRole, error) {
	panic("unimplemented")
}

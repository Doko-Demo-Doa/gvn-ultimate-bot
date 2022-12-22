package discordservice

import "doko/gvn-ultimate-bot/models"

// Interface for DiscordService
type DiscordService interface {
	// Listing the roles as DiscordRole model array
	ListRoles() ([]*models.DiscordRole, error)
	CreateRole() (*models.DiscordRole, error)
	EditRole() (*models.DiscordRole, error)

	// This is actually just mark the role as "deleted" (IsDeleted = 1), not actually delete it
	RemoveRole(email uint) (*models.DiscordRole, error)
}

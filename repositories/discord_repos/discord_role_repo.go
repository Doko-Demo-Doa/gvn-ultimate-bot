package discordrepos

import "doko/gvn-ultimate-bot/models"

type DiscordRoleRepo interface {
	GetByID(id uint) (*models.DiscordRole, error)
	GetByNativeID(nativeId string) (*models.DiscordRole, error)
	AssignRole(user models.DiscordUser, toRole models.DiscordRole) (*models.DiscordUserRole, error)
	UnassignRole(user models.DiscordUser, fromRole models.DiscordRole) (*models.DiscordUserRole, error) // For history
}

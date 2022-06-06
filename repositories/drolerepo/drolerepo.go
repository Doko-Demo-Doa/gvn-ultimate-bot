package drolerepo

import "doko/gin-sample/models"

type Repo interface {
	GetByID(id uint) (*models.DiscordRole, error)
	GetByNativeID(nativeId string) (*models.DiscordRole, error)
	Assign(user models.DiscordUser, toRole models.DiscordRole) (*models.DiscordUserRole, error)
	Unassign(user models.DiscordUser, fromRole models.DiscordRole) (*models.DiscordUserRole, error) // For history
}

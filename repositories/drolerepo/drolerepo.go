package drolerepo

import "doko/gin-sample/models"

type Repo interface {
	GetByID(id uint) (*models.DiscordRole, error)
	GetByNativeID(nativeID string) (*models.DiscordRole, error)
	Assign(user models.DiscordUser) (*models.DiscordUserRole, error)
}

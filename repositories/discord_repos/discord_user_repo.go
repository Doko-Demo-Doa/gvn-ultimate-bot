package discordrepos

import "doko/gvn-ultimate-bot/models"

type DiscordUserRepo interface {
	GetByID(id uint) (*models.DiscordRole, error)
	GetByNativeID(nativeId string) (*models.DiscordRole, error)
}

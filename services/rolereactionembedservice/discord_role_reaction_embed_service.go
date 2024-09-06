package discordservice

import (
	"doko/gvn-ultimate-bot/models"
)

type DiscordRoleReactionEmbedService interface {
	ListEmbeds() ([]*models.DiscordRoleReactionEmbed, error)
}

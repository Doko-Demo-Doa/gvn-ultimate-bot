package models

import (
	"time"

	"gorm.io/gorm"
)

type DiscordBanlist struct {
	gorm.Model
	DUserId    string `gorm:"NOT NULL;size:255"`
	Duration   uint   `gorm:"DEFAULT:0"` // In minutes, see: https://discord.com/developers/docs/resources/channel#channel-object-channel-types
	LinkedRole DiscordRole
	BanDate    time.Time
}

func (DiscordBanlist) TableName() string {
	return "discord_banlist"
}

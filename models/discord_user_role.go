package models

import (
	"time"

	"gorm.io/gorm"
)

type DiscordUserRole struct {
	gorm.Model
	DiscordUser    DiscordUser
	GrantedRole    DiscordRole
	GrantedDate    time.Time
	ExpirationDate time.Time
}

func (DiscordUserRole) TableName() string {
	return "discord_user_role"
}

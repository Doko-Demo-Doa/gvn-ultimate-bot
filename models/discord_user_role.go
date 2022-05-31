package models

import "time"

type DiscordUserRole struct {
	DiscordUser    DiscordUser
	GrantedRole    DiscordRole
	GrantedDate    time.Time
	ExpirationDate time.Time
}

func (DiscordUserRole) TableName() string {
	return "discord_user_role"
}

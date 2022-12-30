package models

import "time"

type DiscordRole struct {
	NativeId     string `gorm:"NOT NULL;size:255"`
	Name         string
	Mentionable  uint // Boolean
	Hoist        uint // Boolean
	Color        uint // Color mapped from discord API
	Expiry       time.Time
	ImplicitType uint // Type that only Dautovn can understand
}

func (DiscordRole) TableName() string {
	return "discord_role"
}

package models

type DiscordRole struct {
	NativeId    string `gorm:"NOT NULL;size:255"`
	Name        string
	Mentionable uint // Boolean
	Hoist       uint // Boolean
	Color       uint // Color mapped from discord API
}

func (DiscordRole) TableName() string {
	return "discord_role"
}

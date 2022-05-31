package models

type DiscordMessageLog struct {
	FromUser        DiscordUser
	OriginalContent string
	NewContent      string
	Channel         DiscordChannel
}

func (DiscordMessageLog) TableName() string {
	return "discord_message_log"
}

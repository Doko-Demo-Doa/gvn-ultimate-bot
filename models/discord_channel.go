package models

import "gorm.io/gorm"

type DiscordChannel struct {
	gorm.Model
	NativeId    string `gorm:"NOT NULL;unique;size:255"`
	ChannelType uint   `gorm:"NOT NULL;DEFAULT:0"` // https://discord.com/developers/docs/resources/channel#channel-object-channel-types
	Name        string `gorm:"NOT NULL;size:128"`
	IsNsfw      uint   `gorm:"NOT NULL;DEFAULT:0"`
}

func (DiscordChannel) TableName() string {
	return "discord_channel"
}

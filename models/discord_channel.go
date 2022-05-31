package models

import "gorm.io/gorm"

type DiscordChannel struct {
	gorm.Model
	NativeId    string `gorm:"NOT NULL;size:255"`
	ChannelType uint
	Name        string `gorm:"NOT NULL;size:128"`
	IsNsfw      uint   `gorm:"NOT NULL;DEFAULT:0"`
}

func (DiscordChannel) TableName() string {
	return "discord_user"
}

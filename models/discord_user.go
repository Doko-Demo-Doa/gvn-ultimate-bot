package models

import "gorm.io/gorm"

type DiscordUser struct {
	gorm.Model
	Discriminator string `gorm:"size:64"`
	NativeId      string `gorm:"NOT NULL;size:255"`
	Avatar        string `gorm:"NOT NULL"`
	PremiumType   uint
}

func (DiscordUser) TableName() string {
	return "discord_user"
}

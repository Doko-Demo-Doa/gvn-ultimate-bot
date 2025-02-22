package models

import (
	"gorm.io/gorm"
)

type DiscordRoleReactionEmbed struct {
	gorm.Model
	NativeMessageId string `gorm:"NOT NULL;unique;size:255"`
	Name            string `gorm:"NOT NULL;size:128"`
	Payload         string `gorm:"NOT NULL;size:8192"` // All the data from the embed
	Tags            string `gorm:"size:512"`
	Version         uint   `gorm:"NOT NULL;DEFAULT:1"`
}

func (DiscordRoleReactionEmbed) TableName() string {
	return "discord_role_reaction_embed"
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type DiscordRole struct {
	gorm.Model
	ID           uint
	NativeId     string `gorm:"NOT NULL;size:255"`
	Name         string `gorm:"NOT NULL;size:255"`
	Mentionable  uint   `gorm:"NOT NULL"`
	Hoist        uint   `gorm:"NOT NULL"`
	Color        uint   `gorm:"NOT NULL"`
	Expiry       time.Time
	ImplicitType uint `gorm:"NOT NULL;DEFAULT:0"`
}

func (DiscordRole) TableName() string {
	return "discord_role"
}

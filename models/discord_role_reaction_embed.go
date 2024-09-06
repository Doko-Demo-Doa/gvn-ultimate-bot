package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type JSONB map[string]interface{}

// Value Marshal
func (jsonField JSONB) Value() (driver.Value, error) {
	return json.Marshal(jsonField)
}

// Scan Unmarshal
func (jsonField *JSONB) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(data, &jsonField)
}

type DiscordRoleReactionEmbed struct {
	gorm.Model
	NativeMessageId string `gorm:"NOT NULL;size:255"`
	Name            string `gorm:"NOT NULL;size:128"`
	Payload         JSONB  `gorm:"NOT NULL;type:jsonb"` // All the data from the embed
	Tags            string `gorm:"size:512"`
	IsDeleted       uint   `gorm:"NOT NULL;DEFAULT:0"`
	Version         uint   `gorm:"NOT NULL;DEFAULT:1"`
}

func (DiscordRoleReactionEmbed) TableName() string {
	return "discord_role_reaction_embed"
}

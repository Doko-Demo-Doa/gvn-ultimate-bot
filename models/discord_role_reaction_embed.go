package models

import (
	"gorm.io/gorm"
)

// DiscordRoleReactionEmbed stores a composed reaction-role message.
// Payload is kept as a plain string/text column so the model is
// database-agnostic. The JSON document is validated and parsed using
// ReactionRoleMessagePayload. If you need Postgres jsonb indexing/querying
// later, migrate the column to jsonb and change the gorm tag accordingly.
type DiscordRoleReactionEmbed struct {
	gorm.Model
	NativeMessageId string `gorm:"NOT NULL;unique;size:255"`
	Name            string `gorm:"NOT NULL;size:128"`
	Payload         string `gorm:"NOT NULL;size:8192"` // JSON document validated by ReactionRoleMessagePayload
	Mode            string `gorm:"size:32;DEFAULT:'default'"` // "default" | "reverse"
	Tags            string `gorm:"size:512"`
	Version         uint   `gorm:"NOT NULL;DEFAULT:1"`
}

func (DiscordRoleReactionEmbed) TableName() string {
	return "discord_role_reaction_embed"
}

// ParsedPayload returns the structured payload for this embed.
func (e *DiscordRoleReactionEmbed) ParsedPayload() (*ReactionRoleMessagePayload, error) {
	p, err := ParseReactionRolePayload(e.Payload)
	if err != nil {
		return nil, err
	}
	// Keep the stored mode in sync with the payload.
	if e.Mode != "" {
		p.Mode = ReactionMode(e.Mode)
	}
	return p, nil
}

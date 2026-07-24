package models

import "gorm.io/gorm"

// DiscordMessageAuditLog stores deleted or edited message records for audit.
type DiscordMessageAuditLog struct {
	gorm.Model
	NativeMessageId string `gorm:"NOT NULL;size:255;index"`
	ChannelId       string `gorm:"NOT NULL;size:255"`
	GuildId         string `gorm:"size:255"`
	AuthorId        string `gorm:"size:255"`
	AuthorName      string `gorm:"size:255"`
	Action          string `gorm:"NOT NULL;size:32"` // "delete" | "edit"
	BeforeContent   string `gorm:"size:4000"`
	AfterContent    string `gorm:"size:4000"`
	Attachments     string `gorm:"size:4000"` // JSON array of attachment URLs
}

func (DiscordMessageAuditLog) TableName() string {
	return "discord_message_audit_log"
}

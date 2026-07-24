package models

import "gorm.io/gorm"

const (
	SystemEventTypeUserSync = "USER_SYNC"
)

// SystemEventLog is a generic log of background/system events (syncs, jobs, etc),
// distinguished by EventType.
type SystemEventLog struct {
	gorm.Model
	EventType string `gorm:"NOT NULL;size:64;index"`
	Status    string `gorm:"NOT NULL;size:32"` // "success" | "failure"
	Message   string `gorm:"size:1000"`
	Metadata  string `gorm:"size:4000"` // JSON blob, e.g. {"synced_count": 123}
}

func (SystemEventLog) TableName() string {
	return "system_event_log"
}

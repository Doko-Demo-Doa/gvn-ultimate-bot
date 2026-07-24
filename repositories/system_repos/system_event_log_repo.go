package systemrepos

import (
	"doko/gvn-ultimate-bot/models"

	"gorm.io/gorm"
)

type SystemEventLogRepo interface {
	Create(log *models.SystemEventLog) (*models.SystemEventLog, error)
	GetLatestByEventType(eventType string) (*models.SystemEventLog, error)
	ListByEventType(eventType string, limit int) ([]*models.SystemEventLog, error)
}

type systemEventLogRepo struct {
	db *gorm.DB
}

func NewSystemEventLogRepo(db *gorm.DB) SystemEventLogRepo {
	return &systemEventLogRepo{db: db}
}

func (r *systemEventLogRepo) Create(log *models.SystemEventLog) (*models.SystemEventLog, error) {
	if err := r.db.Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

func (r *systemEventLogRepo) GetLatestByEventType(eventType string) (*models.SystemEventLog, error) {
	var log models.SystemEventLog
	err := r.db.Where("event_type = ?", eventType).Order("created_at DESC").First(&log).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

func (r *systemEventLogRepo) ListByEventType(eventType string, limit int) ([]*models.SystemEventLog, error) {
	var logs []*models.SystemEventLog
	q := r.db.Where("event_type = ?", eventType).Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

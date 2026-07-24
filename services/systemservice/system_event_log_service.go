package systemservice

import (
	"doko/gvn-ultimate-bot/models"
	systemrepos "doko/gvn-ultimate-bot/repositories/system_repos"
)

type SystemEventLogService interface {
	LogEvent(eventType, status, message, metadata string) error
	GetLatestByEventType(eventType string) (*models.SystemEventLog, error)
}

type systemEventLogService struct {
	Repo systemrepos.SystemEventLogRepo
}

func NewSystemEventLogService(repo systemrepos.SystemEventLogRepo) SystemEventLogService {
	return &systemEventLogService{Repo: repo}
}

func (s *systemEventLogService) LogEvent(eventType, status, message, metadata string) error {
	_, err := s.Repo.Create(&models.SystemEventLog{
		EventType: eventType,
		Status:    status,
		Message:   message,
		Metadata:  metadata,
	})
	return err
}

func (s *systemEventLogService) GetLatestByEventType(eventType string) (*models.SystemEventLog, error) {
	return s.Repo.GetLatestByEventType(eventType)
}

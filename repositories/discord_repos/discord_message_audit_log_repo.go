package discordrepos

import (
	"doko/gvn-ultimate-bot/models"

	"gorm.io/gorm"
)

type DiscordMessageAuditLogFilter struct {
	FromDate   *string
	ToDate     *string
	ChannelId  *string
	AuthorName *string
}

type DiscordMessageAuditLogRepo interface {
	Create(log *models.DiscordMessageAuditLog) (*models.DiscordMessageAuditLog, error)
	List(filter *DiscordMessageAuditLogFilter, limit int, offset int) ([]*models.DiscordMessageAuditLog, error)
	Count(filter *DiscordMessageAuditLogFilter) (int64, error)
	ClearAll() error
}

type discordMessageAuditLogRepo struct {
	db *gorm.DB
}

func NewDiscordMessageAuditLogRepo(db *gorm.DB) DiscordMessageAuditLogRepo {
	return &discordMessageAuditLogRepo{db: db}
}

func (r *discordMessageAuditLogRepo) Create(log *models.DiscordMessageAuditLog) (*models.DiscordMessageAuditLog, error) {
	if err := r.db.Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

func (r *discordMessageAuditLogRepo) buildQuery(filter *DiscordMessageAuditLogFilter) *gorm.DB {
	q := r.db.Model(&models.DiscordMessageAuditLog{})
	if filter == nil {
		return q
	}
	if filter.FromDate != nil && *filter.FromDate != "" {
		q = q.Where("created_at >= ?", *filter.FromDate)
	}
	if filter.ToDate != nil && *filter.ToDate != "" {
		q = q.Where("created_at <= ?", *filter.ToDate)
	}
	if filter.ChannelId != nil && *filter.ChannelId != "" {
		q = q.Where("channel_id = ?", *filter.ChannelId)
	}
	if filter.AuthorName != nil && *filter.AuthorName != "" {
		q = q.Where("author_name ILIKE ?", "%"+*filter.AuthorName+"%")
	}
	return q
}

func (r *discordMessageAuditLogRepo) List(filter *DiscordMessageAuditLogFilter, limit int, offset int) ([]*models.DiscordMessageAuditLog, error) {
	var logs []*models.DiscordMessageAuditLog
	q := r.buildQuery(filter).Order("created_at DESC").Limit(limit).Offset(offset)
	if err := q.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *discordMessageAuditLogRepo) Count(filter *DiscordMessageAuditLogFilter) (int64, error) {
	var count int64
	if err := r.buildQuery(filter).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *discordMessageAuditLogRepo) ClearAll() error {
	return r.db.Where("1 = 1").Delete(&models.DiscordMessageAuditLog{}).Error
}

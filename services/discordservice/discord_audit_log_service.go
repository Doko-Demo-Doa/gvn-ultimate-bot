package discordservice

import (
	"doko/gvn-ultimate-bot/models"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
	"encoding/json"
	"log"
)

type AuditLogFilter struct {
	FromDate   *string
	ToDate     *string
	ChannelId  *string
	AuthorName *string
}

type DiscordAuditLogService interface {
	LogMessageDelete(messageID, channelID, guildID, authorID, authorName, content string, attachments []string) error
	LogMessageEdit(messageID, channelID, guildID, authorID, authorName, beforeContent, afterContent string, attachments []string) error
	ListLogs(filter *AuditLogFilter, limit, offset int) ([]*models.DiscordMessageAuditLog, int64, error)
	ClearLogs() error
}

type discordAuditLogService struct {
	Repo discordrepos.DiscordMessageAuditLogRepo
}

func NewDiscordAuditLogService(repo discordrepos.DiscordMessageAuditLogRepo) DiscordAuditLogService {
	return &discordAuditLogService{Repo: repo}
}

func (s *discordAuditLogService) attachmentsJSON(attachments []string) string {
	if len(attachments) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(attachments)
	return string(b)
}

func (s *discordAuditLogService) LogMessageDelete(messageID, channelID, guildID, authorID, authorName, content string, attachments []string) error {
	logEntry := &models.DiscordMessageAuditLog{
		NativeMessageId: messageID,
		ChannelId:       channelID,
		GuildId:         guildID,
		AuthorId:        authorID,
		AuthorName:      authorName,
		Action:          "delete",
		BeforeContent:   content,
		Attachments:     s.attachmentsJSON(attachments),
	}
	_, err := s.Repo.Create(logEntry)
	if err != nil {
		log.Printf("[audit_log] failed to log message delete: %v", err)
	}
	return err
}

func (s *discordAuditLogService) LogMessageEdit(messageID, channelID, guildID, authorID, authorName, beforeContent, afterContent string, attachments []string) error {
	logEntry := &models.DiscordMessageAuditLog{
		NativeMessageId: messageID,
		ChannelId:       channelID,
		GuildId:         guildID,
		AuthorId:        authorID,
		AuthorName:      authorName,
		Action:          "edit",
		BeforeContent:   beforeContent,
		AfterContent:    afterContent,
		Attachments:     s.attachmentsJSON(attachments),
	}
	_, err := s.Repo.Create(logEntry)
	if err != nil {
		log.Printf("[audit_log] failed to log message edit: %v", err)
	}
	return err
}

func (s *discordAuditLogService) ListLogs(filter *AuditLogFilter, limit, offset int) ([]*models.DiscordMessageAuditLog, int64, error) {
	repoFilter := &discordrepos.DiscordMessageAuditLogFilter{}
	if filter != nil {
		repoFilter.FromDate = filter.FromDate
		repoFilter.ToDate = filter.ToDate
		repoFilter.ChannelId = filter.ChannelId
		repoFilter.AuthorName = filter.AuthorName
	}
	logs, err := s.Repo.List(repoFilter, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.Repo.Count(repoFilter)
	if err != nil {
		return nil, 0, err
	}
	return logs, count, nil
}

func (s *discordAuditLogService) ClearLogs() error {
	return s.Repo.ClearAll()
}

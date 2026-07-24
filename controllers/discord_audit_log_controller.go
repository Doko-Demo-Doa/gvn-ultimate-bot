package controllers

import (
	"doko/gvn-ultimate-bot/services/discordservice"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DiscordAuditLogController interface {
	ListAuditLogs(*gin.Context)
	ClearAuditLogs(*gin.Context)
}

type discordAuditLogController struct {
	svc discordservice.DiscordAuditLogService
}

func NewDiscordAuditLogController(svc discordservice.DiscordAuditLogService) DiscordAuditLogController {
	return &discordAuditLogController{svc: svc}
}

func (ctl *discordAuditLogController) ListAuditLogs(c *gin.Context) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	channelId := c.Query("channel_id")
	authorName := c.Query("author_name")

	filter := &discordservice.AuditLogFilter{}
	if fromDate != "" {
		filter.FromDate = &fromDate
	}
	if toDate != "" {
		filter.ToDate = &toDate
	}
	if channelId != "" {
		filter.ChannelId = &channelId
	}
	if authorName != "" {
		filter.AuthorName = &authorName
	}

	logs, total, err := ctl.svc.ListLogs(filter, limit, offset)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", gin.H{
		"items":  logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (ctl *discordAuditLogController) ClearAuditLogs(c *gin.Context) {
	if err := ctl.svc.ClearLogs(); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", nil)
}

package controllers

import (
	"doko/gvn-ultimate-bot/models"
	"doko/gvn-ultimate-bot/scheduler"
	"doko/gvn-ultimate-bot/services/discordservice"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DiscordRoleInput struct {
	Name         string    `json:"name"`
	NativeID     string    `json:"native_id"`
	Mentionable  uint      `json:"mentionable"`
	Hoist        uint      `json:"hoist"`
	Color        uint      `json:"color"`
	Expiry       time.Time `json:"expiry"`
	ImplicitType uint      `json:"implicit_type"`
}

type AssignRoleInput struct {
	UserNativeID string `json:"user_native_id" binding:"required"`
	RoleNativeID string `json:"role_native_id" binding:"required"`
	Duration     string `json:"duration"` // e.g. "30m", "2h", "7d"
}

type RoleAssignmentResponse struct {
	ID             uint      `json:"id"`
	UserNativeID   string    `json:"user_native_id"`
	RoleNativeID   string    `json:"role_native_id"`
	GrantedDate    time.Time `json:"granted_date"`
	ExpirationDate time.Time `json:"expiration_date"`
	Status         string    `json:"status"`
	TimeRemaining  string    `json:"time_remaining"`
}

type DiscordRoleReactionEmbedInput struct {
	NativeMessageId string                          `json:"native_message_id"`
	Name            string                          `json:"name"`
	Payload         models.ReactionRoleMessagePayload `json:"payload"`
	Mode            string                          `json:"mode"`
	Tags            string                          `json:"tags"`
	Version         uint                            `json:"version"`
}

type PublishReactionRoleInput struct {
	Payload models.ReactionRoleMessagePayload `json:"payload" binding:"required"`
}

type DiscordController interface {
	ListDiscordRoles(*gin.Context)
	CreateDiscordRole(*gin.Context)

	AssignRoleToUser(*gin.Context)
	RevokeRoleFromUser(*gin.Context)
	ListRoleAssignments(*gin.Context)

	ListDiscordRoleReactions(*gin.Context)
	GetDiscordRoleReaction(*gin.Context)
	UpsertDiscordRoleReaction(*gin.Context)
	PublishDiscordRoleReaction(*gin.Context)
	DeleteDiscordRoleReaction(*gin.Context)

	ListDiscordChannels(*gin.Context)
	ListDiscordEmojis(*gin.Context)
	SearchDiscordMembers(*gin.Context)
}

type discordController struct {
	ds        discordservice.DiscordService
	dre       discordservice.DiscordRoleReactionEmbedService
	scheduler *scheduler.RoleScheduler
}

func NewDiscordController(
	ds discordservice.DiscordService,
	dre discordservice.DiscordRoleReactionEmbedService,
	sch *scheduler.RoleScheduler,
) DiscordController {
	return &discordController{
		ds:        ds,
		dre:       dre,
		scheduler: sch,
	}
}

/* Role-related */

func (ctl *discordController) ListDiscordRoles(c *gin.Context) {
	data, err := ctl.ds.ListRoles()
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *discordController) CreateDiscordRole(c *gin.Context) {
	var dRoleInput DiscordRoleInput
	if err := c.ShouldBindJSON(&dRoleInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	dRole := ctl.inputToDiscordRole(dRoleInput)
	data, err := ctl.ds.CreateRole(&dRole)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *discordController) AssignRoleToUser(c *gin.Context) {
	var input AssignRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var duration time.Duration
	if input.Duration != "" {
		d, err := time.ParseDuration(input.Duration)
		if err != nil {
			HTTPRes(c, http.StatusBadRequest, "Invalid duration format", nil)
			return
		}
		duration = d
	}

	if err := ctl.scheduler.GrantRole(input.UserNativeID, input.RoleNativeID, duration); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	HTTPRes(c, http.StatusOK, "ok", nil)
}

func (ctl *discordController) RevokeRoleFromUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, "Invalid assignment ID", nil)
		return
	}

	if err := ctl.scheduler.RevokeRole(uint(id)); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	HTTPRes(c, http.StatusOK, "ok", nil)
}

func (ctl *discordController) ListRoleAssignments(c *gin.Context) {
	all, err := ctl.ds.GetAllActiveAssignments()
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Also fetch expired ones for completeness
	expired, err := ctl.ds.GetExpiredRoleAssignments()
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	all = append(all, expired...)

	now := time.Now()
	resp := make([]RoleAssignmentResponse, 0, len(all))
	for _, a := range all {
		status := "active"
		remaining := ""
		if a.ExpirationDate.Before(now) || a.ExpirationDate.Equal(now) {
			status = "expired"
			remaining = "0s"
		} else {
			remaining = humanDuration(time.Until(a.ExpirationDate))
		}
		resp = append(resp, RoleAssignmentResponse{
			ID:             a.ID,
			UserNativeID:   a.UserNativeID,
			RoleNativeID:   a.RoleNativeID,
			GrantedDate:    a.GrantedDate,
			ExpirationDate: a.ExpirationDate,
			Status:         status,
			TimeRemaining:  remaining,
		})
	}

	HTTPRes(c, http.StatusOK, "ok", resp)
}

func humanDuration(d time.Duration) string {
	if d < 0 {
		return "0s"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		m := int(d.Minutes())
		s := int(d.Seconds()) % 60
		if s > 0 {
			return fmt.Sprintf("%dm %ds", m, s)
		}
		return fmt.Sprintf("%dm", m)
	}
	if d < 24*time.Hour {
		h := int(d.Hours())
		m := int(d.Minutes()) % 60
		if m > 0 {
			return fmt.Sprintf("%dh %dm", h, m)
		}
		return fmt.Sprintf("%dh", h)
	}
	days := int(d.Hours()) / 24
	h := int(d.Hours()) % 24
	if h > 0 {
		return fmt.Sprintf("%dd %dh", days, h)
	}
	return fmt.Sprintf("%dd", days)
}

func (ctl *discordController) inputToDiscordRole(input DiscordRoleInput) models.DiscordRole {
	return models.DiscordRole{
		Name:         input.NativeID,
		NativeId:     input.Name,
		Mentionable:  input.Mentionable,
		Hoist:        input.Hoist,
		Color:        input.Color,
		Expiry:       &input.Expiry,
		ImplicitType: input.ImplicitType,
	}
}

func (ctl *discordController) inputToDiscordRoleReactionEmbed(input DiscordRoleReactionEmbedInput) (*models.DiscordRoleReactionEmbed, *models.ReactionRoleMessagePayload, error) {
	payload := input.Payload
	if input.Mode != "" {
		payload.Mode = models.ReactionMode(input.Mode)
	}
	if payload.Mode == "" {
		payload.Mode = models.ReactionModeDefault
	}
	if err := payload.Validate(); err != nil {
		return nil, nil, err
	}
	payloadJSON, err := payload.ToJSON()
	if err != nil {
		return nil, nil, err
	}
	return &models.DiscordRoleReactionEmbed{
		NativeMessageId: input.NativeMessageId,
		Name:            input.Name,
		Payload:         payloadJSON,
		Mode:            string(payload.Mode),
		Tags:            input.Tags,
		Version:         input.Version,
	}, &payload, nil
}

/* Role-reaction related */
func (ctl *discordController) ListDiscordRoleReactions(c *gin.Context) {
	data, err := ctl.dre.ListEmbeds()
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *discordController) CreateDiscordRoleReactions(c *gin.Context) {
	var dRoleInput DiscordRoleReactionEmbedInput
	if err := c.ShouldBindJSON(&dRoleInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	dRoleReactionEmbed, payload, err := ctl.inputToDiscordRoleReactionEmbed(dRoleInput)
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	data, err := ctl.dre.UpsertEmbed(dRoleReactionEmbed, payload)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *discordController) GetDiscordRoleReaction(c *gin.Context) {
	roleReactionId := c.Param("id")
	mId, errUint := strconv.ParseUint(roleReactionId, 10, 32)
	if errUint != nil {
		HTTPRes(c, http.StatusBadRequest, "Invalid role reaction ID", nil)
		return
	}

	data, err := ctl.dre.GetSingleEmbed(uint(mId))
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *discordController) UpsertDiscordRoleReaction(c *gin.Context) {
	var dRoleReactionEmbedInput DiscordRoleReactionEmbedInput
	if err := c.ShouldBindJSON(&dRoleReactionEmbedInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	_, payload, err := ctl.inputToDiscordRoleReactionEmbed(dRoleReactionEmbedInput)
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// If native_message_id is provided, edit the existing Discord message.
	if dRoleReactionEmbedInput.NativeMessageId != "" {
		data, err := ctl.dre.EditEmbed(dRoleReactionEmbedInput.NativeMessageId, payload)
		if err != nil {
			HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		HTTPRes(c, http.StatusOK, "ok", data)
		return
	}

	// Otherwise fall back to plain DB upsert.
	dRoleReactionEmbed, _, err := ctl.inputToDiscordRoleReactionEmbed(dRoleReactionEmbedInput)
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	data, err := ctl.dre.UpsertEmbed(dRoleReactionEmbed, payload)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *discordController) PublishDiscordRoleReaction(c *gin.Context) {
	var input PublishReactionRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if input.Payload.Mode == "" {
		input.Payload.Mode = models.ReactionModeDefault
	}
	data, err := ctl.dre.PublishEmbed(&input.Payload)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *discordController) DeleteDiscordRoleReaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, "Invalid role reaction ID", nil)
		return
	}
	if err := ctl.dre.DeleteEmbed(uint(id)); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", nil)
}

func (ctl *discordController) ListDiscordChannels(c *gin.Context) {
	data, err := ctl.dre.ListChannels()
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *discordController) ListDiscordEmojis(c *gin.Context) {
	data, err := ctl.dre.ListEmojis()
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *discordController) SearchDiscordMembers(c *gin.Context) {
	query := c.Query("q")
	data, err := ctl.dre.SearchGuildMembers(query)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", data)
}

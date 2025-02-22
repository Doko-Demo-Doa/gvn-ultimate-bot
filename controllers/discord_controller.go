package controllers

import (
	"doko/gvn-ultimate-bot/models"
	"doko/gvn-ultimate-bot/services/discordservice"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DiscordRoleInput struct {
	Name        string    `json:"name"`
	NativeID    string    `json:"native_id"`
	Mentionable uint      `json:"mentionable"`
	Hoist       uint      `json:"hoist"`
	Color       uint      `json:"color"`
	Expiry      time.Time `json:"expiry"`
	// Use this to indicate what it is in the Admin UI. Maybe it's Discord primodial role, can't be modified (like, "@everyone")
	ImplicitType uint `json:"implicit_type"`
}

type DiscordRoleReactionEmbedInput struct {
	NativeMessageId string `json:"native_message_id"`
	Name            string `json:"name"`
	Payload         string `json:"payload"`
	Tags            string `json:"tags"`
	Version         uint   `json:"version"`
}

type DiscordController interface {
	ListDiscordRoles(*gin.Context)
	CreateDiscordRole(*gin.Context)

	ListDiscordRoleReactions(*gin.Context)
	GetDiscordRoleReaction(*gin.Context)
	UpsertDiscordRoleReaction(*gin.Context)
}

type discordController struct {
	ds  discordservice.DiscordService
	dre discordservice.DiscordRoleReactionEmbedService
}

func NewDiscordController(
	ds discordservice.DiscordService,
	dre discordservice.DiscordRoleReactionEmbedService,
) DiscordController {
	return &discordController{
		ds:  ds,
		dre: dre,
	}
}

/* Role-related */

func (ctl *discordController) ListDiscordRoles(c *gin.Context) {
	data, err := ctl.ds.ListRoles()
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
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
	}

	HTTPRes(c, http.StatusOK, "ok", data)
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

func (ctl *discordController) inputToDiscordRoleReactionEmbed(input DiscordRoleReactionEmbedInput) models.DiscordRoleReactionEmbed {
	data := models.DiscordRoleReactionEmbed{
		NativeMessageId: input.NativeMessageId,
		Name:            input.Name,
		Tags:            input.Tags,
		Version:         input.Version,
		Payload:         input.Payload,
	}

	return data
}

/* Role-reaction related */
func (ctl *discordController) ListDiscordRoleReactions(c *gin.Context) {
	data, err := ctl.dre.ListEmbeds()
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
	}

	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *discordController) CreateDiscordRoleReactions(c *gin.Context) {
	var dRoleInput DiscordRoleReactionEmbedInput

	if err := c.ShouldBindJSON(&dRoleInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	dRoleReactionEmbed := ctl.inputToDiscordRoleReactionEmbed(dRoleInput)
	data, err := ctl.dre.UpsertEmbed(&dRoleReactionEmbed)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
	}

	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *discordController) GetDiscordRoleReaction(c *gin.Context) {
	roleReactionId := c.Param("id")
	mId, errUint := (strconv.ParseUint(roleReactionId, 10, 32))

	if errUint != nil {
		HTTPRes(c, http.StatusBadRequest, "Invalid role reaction ID", nil)
	}

	data, err := ctl.dre.GetSingleEmbed(uint(mId))
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
	}

	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *discordController) UpsertDiscordRoleReaction(c *gin.Context) {
	var dRoleReactionEmbedInput DiscordRoleReactionEmbedInput

	if err := c.ShouldBindJSON(&dRoleReactionEmbedInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	dRoleReactionEmbed := ctl.inputToDiscordRoleReactionEmbed(dRoleReactionEmbedInput)
	data, err := ctl.dre.UpsertEmbed(&dRoleReactionEmbed)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	HTTPRes(c, http.StatusOK, "ok", data)

}

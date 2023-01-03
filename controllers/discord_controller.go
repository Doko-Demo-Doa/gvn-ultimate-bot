package controllers

import (
	"doko/gvn-ultimate-bot/models"
	"doko/gvn-ultimate-bot/services/discordservice"
	"net/http"
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

type DiscordController interface {
	ListDiscordRoles(*gin.Context)
	CreateDiscordRole(*gin.Context)
}

type discordController struct {
	ds discordservice.DiscordService
}

func NewDiscordController(
	ds discordservice.DiscordService,
) DiscordController {
	return &discordController{
		ds: ds,
	}
}

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
		Expiry:       input.Expiry,
		ImplicitType: input.ImplicitType,
	}
}

package controllers

import (
	"doko/gvn-ultimate-bot/services/discordservice"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DiscordController interface {
	ListDiscordRoles(*gin.Context)
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

package controllers

import (
	"doko/gvn-ultimate-bot/services/adminaccessservice"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WhitelistedRoleInput struct {
	RoleNativeID string `json:"role_native_id" binding:"required"`
	AccessLevel  string `json:"access_level"`
	Label        string `json:"label"`
}

type AdminAccessController interface {
	ListWhitelistedRoles(*gin.Context)
	UpsertWhitelistedRole(*gin.Context)
	DeleteWhitelistedRole(*gin.Context)
	CheckAccess(*gin.Context)
}

type adminAccessController struct {
	as adminaccessservice.AdminAccessService
}

func NewAdminAccessController(as adminaccessservice.AdminAccessService) AdminAccessController {
	return &adminAccessController{
		as: as,
	}
}

func (ctl *adminAccessController) ListWhitelistedRoles(c *gin.Context) {
	roles, err := ctl.as.ListWhitelistedRoles()
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", roles)
}

func (ctl *adminAccessController) UpsertWhitelistedRole(c *gin.Context) {
	var input WhitelistedRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	role, err := ctl.as.UpsertWhitelistedRole(input.RoleNativeID, input.AccessLevel, input.Label)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", role)
}

func (ctl *adminAccessController) DeleteWhitelistedRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, "id should be a number", nil)
		return
	}

	if err := ctl.as.DeleteWhitelistedRole(uint(id)); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", nil)
}

// @Summary Check whether a Discord user's roles grant admin dashboard access
// @Produce  json
// @Param discord_user_id query string true "Discord user native ID"
// @Success 200 {object} Response
// @Router /api/admin/access-check [get]
func (ctl *adminAccessController) CheckAccess(c *gin.Context) {
	discordUserID := c.Query("discord_user_id")
	if discordUserID == "" {
		HTTPRes(c, http.StatusBadRequest, "discord_user_id is required", nil)
		return
	}

	result, err := ctl.as.CheckAccess(discordUserID)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", result)
}

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctl *userController) GetDiscordRoles(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	_, err := ctl.us.InitiateResetPassowrd(input.Email)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
	}

	HTTPRes(c, http.StatusOK, "Email sent", nil)
}

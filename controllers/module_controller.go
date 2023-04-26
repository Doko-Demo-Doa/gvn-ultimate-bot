package controllers

import (
	"doko/gvn-ultimate-bot/services/moduleservice"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GVNModuleInput struct {
	ModuleName  string `json:"module_name"`
	IsActivated int8   `json:"is_activated"`
}

type ModuleController interface {
	ActivateOrDisableModule(*gin.Context)
	ListModules(c *gin.Context)
}

type moduleController struct {
	ms moduleservice.ModuleService
}

func NewModuleController(
	ms moduleservice.ModuleService,
) ModuleController {
	return &moduleController{
		ms: ms,
	}
}

func (ctl *moduleController) ListModules(c *gin.Context) {
	data, err := ctl.ms.ListModules()
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
	}

	HTTPRes(c, http.StatusOK, "ok", data)
}

func (ctl *moduleController) ActivateOrDisableModule(c *gin.Context) {

	panic("unimplemented")
}

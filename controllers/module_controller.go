package controllers

import (
	"doko/gvn-ultimate-bot/services/moduleservice"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ActivateDeactivateModuleInput struct {
	ModuleId    uint  `json:"module_id"`
	IsActivated uint8 `json:"is_activated"`
}

type ModuleController interface {
	ActivateOrDisableModule(*gin.Context)
	ListModules(c *gin.Context)
	GetModuleByID(c *gin.Context)
	GetModuleByName(c *gin.Context)
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

func (ctl *moduleController) GetModuleByID(c *gin.Context) {
	moduleId := (c.Param("id"))
	mId, errUint := (strconv.ParseUint(moduleId, 10, 32))

	if errUint != nil {
		HTTPRes(c, http.StatusBadRequest, "Invalid module ID", nil)
	}

	module, err := ctl.ms.GetModuleByID(uint(mId))

	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
	}

	HTTPRes(c, http.StatusOK, "ok", module)
}

func (ctl *moduleController) GetModuleByName(c *gin.Context) {
	moduleName := c.Query("name")
	module, err := ctl.ms.GetModuleByName(moduleName)

	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
	}

	HTTPRes(c, http.StatusOK, "ok", module)
}

func (ctl *moduleController) ActivateOrDisableModule(c *gin.Context) {
	var moduleInput ActivateDeactivateModuleInput

	if err := c.ShouldBindJSON(&moduleInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	module, err := ctl.ms.ActivateOrDisableModule(moduleInput.ModuleId, moduleInput.IsActivated)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
	}

	HTTPRes(c, http.StatusOK, "ok", module)
}

package controllers

import (
	"github.com/gin-gonic/gin"
)

type GVNModuleInput struct {
	ModuleName  string `json:"module_name"`
	IsActivated int8   `json:"is_activated"`
}

type ModuleController interface {
	EnableModule(*gin.Context)
	DisableModule(*gin.Context)
}

type moduleController struct{}

func NewModuleController() ModuleController {
	return &moduleController{}
}

func (*moduleController) DisableModule(*gin.Context) {
	panic("unimplemented")
}

func (*moduleController) EnableModule(*gin.Context) {
	panic("unimplemented")
}

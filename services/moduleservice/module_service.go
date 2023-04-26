package moduleservice

import "doko/gvn-ultimate-bot/models"

type ModuleService interface {
	EnableModule(id uint) (string, error)
	DisableModule(id uint) (string, error)
	ListModules() ([]*models.AppModule, error)
}
